/*************************************************************************
 *
 * MAXSET CONFIDENTIAL
 * __________________
 *
 *  [2019] - [2020] Maxset WorldWide Inc.
 *  All Rights Reserved.
 *
 * NOTICE:  All information contained herein is, and remains
 * the property of Maxset WorldWide Inc. and its suppliers,
 * if any.  The intellectual and technical concepts contained
 * herein are proprietary to Maxset WorldWide Inc.
 * and its suppliers and may be covered by U.S. and Foreign Patents,
 * patents in process, and are protected by trade secret or copyright law.
 * Dissemination of this information or reproduction of this material
 * is strictly forbidden unless prior written permission is obtained
 * from Maxset WorldWide Inc.
 */

package handlers

import (
	"net/http"

	"git.maxset.io/web/knaxim/internal/database"
	"git.maxset.io/web/knaxim/internal/database/types"
	"git.maxset.io/web/knaxim/pkg/srverror"
	"git.maxset.io/web/knaxim/pkg/srvjson"
	"github.com/gorilla/mux"
)

// AttachRecord adds paths for basic file actions TODO: combine with AttachFile
func AttachRecord(r *mux.Router) {

	r.Use(srvjson.JSONResponse)
	r.Use(ConnectDatabase)
	r.Use(UserCookie)
	r.Use(ParseBody)
	r.Use(groupMiddleware)
	r.HandleFunc("", getOwnedRecords).Methods("GET")
	r.HandleFunc("/view", getPermissionRecords("view")).Methods("GET")
	r.HandleFunc("/{id}/name", changeRecordName).Methods("POST")
}

func changeRecordName(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(USER).(types.Owner)
	filebase := r.Context().Value(types.FILE).(database.Filebase)
	vals := mux.Vars(r)
	fid, err := types.DecodeFileID(vals["id"])
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

func sendMatchedRecords(out http.ResponseWriter, r *http.Request, matches []types.FileI) {
	w := out.(*srvjson.ResponseWriter)
	output := make(map[string]FileInfo)
	for _, match := range matches {
		count, err := r.Context().Value(types.CONTENT).(database.Contentbase).Len(match.GetID().StoreID)
		if err != nil {
			panic(err)
		}
		store, err := r.Context().Value(types.STORE).(database.Storebase).Get(match.GetID().StoreID)
		if err != nil {
			panic(err)
		}
		output[match.GetID().String()] = FileInfo{match, count, store.FileSize}
	}
	w.Set("files", output)
}

func getOwnedRecords(w http.ResponseWriter, r *http.Request) {
	var owner types.Owner
	if group := r.Context().Value(GROUP); group != nil {
		owner = group.(types.Owner)
	} else {
		owner = r.Context().Value(USER).(types.Owner)
	}
	filebase := r.Context().Value(types.FILE).(database.Filebase)
	recs, err := filebase.GetOwned(owner.GetID())
	if err != nil {
		panic(err)
	}
	sendMatchedRecords(w, r, recs)
}

func getPermissionRecords(key string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var owner types.Owner
		if group := r.Context().Value(GROUP); group != nil {
			owner = group.(types.Owner)
		} else {
			owner = r.Context().Value(USER).(types.Owner)
		}
		filebase := r.Context().Value(types.FILE).(database.Filebase)
		recs, err := filebase.GetPermKey(owner.GetID(), key)
		if err != nil {
			panic(err)
		}
		sendMatchedRecords(w, r, recs)
	}
}
