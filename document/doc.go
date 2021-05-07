package document

import (
	"github.com/rosshpayne/graph-sdl/internal/db"
)

var NoItemFoundErr = db.NoItemFoundErr

func GetDocument() string {
	return db.GetDocument()
}

func SetDocument(doc string) string {
	return db.SetDocument(doc)
}

func DeleteTyp(obj string) error {
	return db.DeleteType(obj)
}
