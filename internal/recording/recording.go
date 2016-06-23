package recording

import (
	"fmt"
	"net/http"
	"regexp"
	"time"

	"gopkg.in/mgo.v2/bson"
)

var (
	ErrorNotFound = newError("record not found in repository", http.StatusNotFound)
	ErrorObjectID = newError("invalid objectId", http.StatusNotFound)
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
	DurationSec int64         `bson:"durationSec"`
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
	ListFiles() ([]*GFSMeta, error)
	GetFile(accountId, recodingId string) (*GFSMeta, error)
	OpenById(id bson.ObjectId) (*GFSFile, error)
}

//type Error interface {
//error
//Repo() bool
//}

type RepoError struct {
	error
	Status int
}

func (e RepoError) Repo() bool {
	return true
}

func newError(text string, status int) RepoError {
	return RepoError{error: fmt.Errorf(text), Status: status}
}

var objectIdErrorRegex = regexp.MustCompile(`ObjectIDs must be exactly 12 bytes long`)
var notFound = "not found"

func mapError(e error) error {
	s := e.Error()
	switch {
	case s == notFound:
		return ErrorNotFound
	case objectIdErrorRegex.MatchString(s):
		return ErrorObjectID
	default:
		return e
	}
}

func internalError(internal error) error {
	return RepoError{error: fmt.Errorf("server error"), Status: http.StatusInternalServerError}
}
