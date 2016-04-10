package storage

import (
	"gitlab.vailsys.com/jerny/coffer/pkg/storage/mongo"
	"gitlab.vailsys.com/jerny/coffer/recording"
)

type MongoRecordingRepo struct {
	SessionProvider *mongo.SessionProvider
}

func NewMongoRecordingRepo(provider *mongo.SessionProvider) recording.RecordingRepo {
	return &MongoRecordingRepo{
		SessionProvider: provider,
	}
}

func (mr *MongoRecordingRepo) Get(accountId, recordingId string) (*recording.Recording, error) {
	return nil, nil
}

func (mr *MongoRecordingRepo) List(accountId string) ([]*recording.Recording, string, error) {
	return nil, "", nil
}

func (mr *MongoRecordingRepo) Delete(accountId, recordingId string) error {
	return nil
}
