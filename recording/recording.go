package recording

import (
	"io"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
