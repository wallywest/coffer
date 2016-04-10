package recording

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type RecordingManager struct {
	recordingRepo RecordingRepo
	assetRepo     AssetRepo
}

type Recording struct {
	Id          bson.ObjectId `bson:"_id"`
	FileId      string        `bson:"fileId"`
	FileName    string        `bson:"fileName"`
	FileSize    int64         `bson:"fileSize"`
	MimeType    string        `bson:"mimeType"`
	DownloadUrl string        `bson:"downloadUrl"`
	URI         string        `bson:"uri"`
	Revision    int           `bson:"revision"`
	DateCreated time.Time     `bson:"dateCreated"`
	DateUpdated time.Time     `bson:"dateUpdated"`
	RecordingId string        `bson:"recordingId"`
	AccountId   string        `bson:"accountId"`
	CallId      string        `bson:"callId"`
	Duration    int           `bson:"duration"`
}

type RecordingAsset struct {
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
	Find() ([]RecordingAsset, error)
	Open() (RecordingAsset, error)
}

func NewRecordingManager(repo RecordingRepo, assetRepo AssetRepo) *RecordingManager {
	return &RecordingManager{
		recordingRepo: repo,
		assetRepo:     assetRepo,
	}
}
