package recording

import (
	"fmt"
	"io"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	ErrorNotFound = newError("record not found in repository")
)

type Recording struct {
	Id          bson.ObjectId `bson:"_id"`
	RecordingId string        `bson:"recordingId"`
	AccountId   string        `bson:"accountId"`
	FileId      string        `bson:"fileId"`
	FileName    string        `bson:"fileName"`
	FileSize    int64         `bson:"fileSize"`
	MimeType    string        `bson:"mimeType"`
	DownloadUrl string        `bson:"downloadUrl"`
	URI         string        `bson:"uri"`
	Revision    int           `bson:"revision"`
	DateCreated time.Time     `bson:"dateCreated"`
	DateUpdated time.Time     `bson:"dateUpdated"`
	CallId      string        `bson:"callId"`
	Duration    int64         `bson:"duration"`
}

type RecordingRepo interface {
	//get from mongodb
	Get(accountId, recordingId string) (*Recording, error)

	//add filter, token, max results
	List(accountId string) ([]*Recording, string, error)

	Delete(accountId, recordingId string) error
}

//need to think about this
type AssetRepo interface {
	ListFiles() ([]*GFSFile, error)
	GetFile(accountId, recodingId string) (*GFSFile, error)
	OpenById(id bson.ObjectId) (io.ReadSeeker, error)
	OpenByName(name string) (*mgo.GridFile, error)
}

type Error interface {
	error
	Repo() bool
}

type RepoError struct {
	error
}

//func (e RepoError) Error() string {
//return e.Error()
//}

func (e RepoError) Repo() bool {
	return true
}

func newError(text string) RepoError {
	return RepoError{fmt.Errorf(text)}
	//return RepoError{
	//Type: t,
	//Desc: desc,
	//Code: code,
	//}
}

func mapError(e error) error {
	switch e.Error() {
	case "not found":
		return ErrorNotFound
	}

	//if mapped, ok := errorMap[e]; ok {
	//return mapped
	//}
	//return internalError(e)
	return e
}

func internalError(internal error) error {
	return RepoError{fmt.Errorf("server error")}
}
