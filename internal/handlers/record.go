package handlers

import (
	"encoding/json"
	"net/http"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/filehash"

	"git.maxset.io/web/knaxim/pkg/srverror"
	"github.com/gorilla/mux"
)

func setupRecord(r *mux.Router) {
	r.Use(cookieMiddleware)
	r.Use(groupMiddleware)
	r.HandleFunc("", getOwnedRecords).Methods("GET")
	r.HandleFunc("/view", getPermissionRecords("view")).Methods("GET")
	r.HandleFunc("/{id}/name", changeRecordName).Methods("POST")
}

func changeRecordName(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(USER).(database.Owner)
	filebase := r.Context().Value("filebase").(database.Filebase)
	vals := mux.Vars(r)
	fid, err := filehash.DecodeFileID(vals["id"])
	if err != nil {
		panic(srverror.New(err, 400, "Bad File ID"))
	}
	file, err := filebase.Get(fid)
	if err != nil {
		panic(err)
	}
	if !file.GetOwner().Match(user) {
		panic(srverror.Basic(403, "Permission Denied", user.GetID().String(), file.GetID().String()))
	}
	if name := r.FormValue("name"); len(name) > 0 {
		file.SetName(name)
		err = filebase.Update(file)
		if err != nil {
			panic(err)
		}
		w.Write([]byte("name changed"))
	} else {
		panic(srverror.Basic(400, "No Name Given"))
	}
}

func sendMatchedRecords(w http.ResponseWriter, r *http.Request, matches []database.FileI) {
	output := make(map[string]FileInfo)
	for _, match := range matches {
		count, err := r.Context().Value(database.CONTENT).(database.Contentbase).Len(match.GetID().StoreID)
		if err != nil {
			panic(err)
		}
		store, err := r.Context().Value(database.STORE).(database.Storebase).Get(match.GetID().StoreID)
		if err != nil {
			panic(err)
		}
		output[match.GetID().String()] = FileInfo{match, count, store.FileSize}
	}
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"files": output,
	}); err != nil {
		panic(srverror.New(err, 500, "Server Error", "sendMatchedRecords encode json"))
	}
	w.Header().Add("Content-Type", "application/json")
}

func getOwnedRecords(w http.ResponseWriter, r *http.Request) {
	var owner database.Owner
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(database.Owner)
	} else {
		owner = r.Context().Value(USER).(database.Owner)
	}
	filebase := r.Context().Value(database.FILE).(database.Filebase)
	recs, err := filebase.GetOwned(owner.GetID())
	if err != nil {
		panic(err)
	}
	sendMatchedRecords(w, r, recs)
}

func getPermissionRecords(key string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var owner database.Owner
		if group := r.Context().Value(GROUP); group != nil {
			owner = group.(database.Owner)
		} else {
			owner = r.Context().Value(USER).(database.Owner)
		}
		filebase := r.Context().Value(database.FILE).(database.Filebase)
		recs, err := filebase.GetPermKey(owner.GetID(), key)
		if err != nil {
			panic(err)
		}
		sendMatchedRecords(w, r, recs)
	}
}
