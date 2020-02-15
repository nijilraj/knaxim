// basing off of database/file.go

package database

import (
	"bytes"
	"compress/gzip"
	"io"

	"git.maxset.io/web/knaxim/internal/database/filehash"
	"git.maxset.io/web/knaxim/pkg/srverror"
)

type ViewStore struct {
	ID      filehash.StoreID `json:"id" bson:"id"`
	Content []byte           `json:"content" bson:"-"`
}

func NewViewStore(id filehash.StoreID, r io.Reader) (*ViewStore, error) {
	store := new(ViewStore)

	contentBuf := new(bytes.Buffer)
	gzWrite := gzip.NewWriter(contentBuf)

	var err error
	if _, err = io.Copy(gzWrite, r); err != nil {
		return nil, err
	}

	if err = gzWrite.Close(); err != nil {
		return nil, srverror.New(err, 500, "Database Error V2")
	}

	store.Content = contentBuf.Bytes()
	store.ID = id
	return store, nil
}

func (vs *ViewStore) Reader() (io.Reader, error) {
	buf := bytes.NewReader(vs.Content)
	out, err := gzip.NewReader(buf)
	if err != nil {
		srverror.New(err, 500, "Database Error V3", "file reading error")
	}
	return out, err
}
