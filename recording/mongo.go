package recording

import (
	"gitlab.vailsys.com/jerny/coffer/pkg/mongo"
	"gopkg.in/mgo.v2/bson"
)

type MongoRecordingRepo struct {
	DB              string
	Collection      string
	SessionProvider *mongo.SessionProvider
}

func NewMongoRecordingRepo(opts mongo.MongoConfig, provider *mongo.SessionProvider) RecordingRepo {
	//change this
	return &MongoRecordingRepo{
		DB:              opts.DB,
		Collection:      "recordings",
		SessionProvider: provider,
	}
}

func (mr *MongoRecordingRepo) Get(accountId, recordingId string) (*Recording, error) {
	session, err := mr.SessionProvider.GetSession()
	defer session.Close()

	if err != nil {
		return nil, err
	}

	collection := session.DB(mr.DB).C(mr.Collection)

	var recording *Recording

	query := bson.M{"accountId": accountId, "recordingId": recordingId}
	err = collection.Find(query).One(&recording)

	if err != nil {
		return nil, err

	}

	return recording, nil
}

func (mr *MongoRecordingRepo) List(accountId string) ([]*Recording, string, error) {
	session, err := mr.SessionProvider.GetSession()
	defer session.Close()

	if err != nil {
		return nil, "", err
	}

	collection := session.DB(mr.DB).C(mr.Collection)

	var recordings []*Recording

	query := bson.M{"accountId": accountId}
	err = collection.Find(query).All(&recordings)

	if err != nil {
		return nil, "", err

	}

	return recordings, "", nil
}

func (mr *MongoRecordingRepo) Delete(accountId, recordingId string) error {
	return nil
}
