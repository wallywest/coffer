package recording

import (
	"fmt"

	"gitlab.vailsys.com/vail-cloud-services/coffer/internal/logger"
	"gitlab.vailsys.com/vail-cloud-services/coffer/internal/storage/driver/mongo"
	"gopkg.in/mgo.v2/bson"
)

const (
	MONGO_BATCH_SIZE = 50
)

type MongoRecordingRepo struct {
	DB              string
	Collection      string
	SessionProvider *mongo.SessionProvider
	CursorMap       map[string]*mongoCursor
}

func NewMongoRecordingRepo(opts mongo.MongoConfig, provider *mongo.SessionProvider) RecordingRepo {
	//change this
	return &MongoRecordingRepo{
		DB:              opts.DB,
		Collection:      "recordings",
		SessionProvider: provider,
		CursorMap:       make(map[string]*mongoCursor, 0),
	}
}

func (mr *MongoRecordingRepo) Get(accountId, recordingId string) (*Recording, error) {
	session, err := mr.SessionProvider.GetSession()
	defer session.Close()

	if err != nil {
		logger.Logger.Errorf("error fetching mongo session: %s", err)
		return nil, mapError(err)
	}

	collection := session.DB(mr.DB).C(mr.Collection)

	var rec *Recording

	logger.Logger.Debugf("fetching recording: %v accountId: %v", recordingId, accountId)
	query := bson.M{"accountId": accountId, "recordingId": recordingId}
	err = collection.Find(query).One(&rec)

	if err != nil {
		logger.Logger.Errorf("error fetching recording: %s", err)
		return nil, mapError(err)

	}

	return rec, nil
}

func (mr *MongoRecordingRepo) List(accountId string) ([]*Recording, *CursorInfo, error) {
	session, err := mr.SessionProvider.GetSession()

	if err != nil {
		return nil, nil, mapError(err)
	}

	cq := &cursorQuery{
		Session:    session,
		DB:         mr.DB,
		Collection: mr.Collection,
		Selector:   bson.M{"accountId": accountId},
	}

	cursor, err := newRecordingCursor(cq)

	if err != nil {
		return nil, nil, mapError(err)
	}

	recordings, info, err := cursor.NextPage()

	if err != nil {
		return nil, nil, nil
	}
	logger.Logger.Debugf("cursor with info: %s", info)

	if cursor.HasNext() {
		mr.CursorMap[cursor.Id] = cursor
	}

	return recordings, info, nil
}

func (mr *MongoRecordingRepo) ListByCursor(cursorId string) ([]*Recording, *CursorInfo, error) {
	cursor, ok := mr.CursorMap[cursorId]

	if !ok {
		return nil, nil, mapError(fmt.Errorf("invalid cursorId"))
	}

	recordings, info, err := cursor.NextPage()

	if err != nil {
		return nil, nil, mapError(err)

	}

	return recordings, info, nil
}

func (mr *MongoRecordingRepo) ListByCall(accountId, callId string) ([]*Recording, *CursorInfo, error) {
	session, err := mr.SessionProvider.GetSession()

	if err != nil {
		return nil, nil, mapError(err)
	}

	cq := &cursorQuery{
		Session:    session,
		DB:         mr.DB,
		Collection: mr.Collection,
		Selector:   bson.M{"accountId": accountId, "callId": callId},
	}

	cursor, err := newRecordingCursor(cq)

	if err != nil {
		return nil, nil, mapError(err)
	}

	recordings, info, err := cursor.NextPage()

	if err != nil {
		return nil, nil, mapError(err)

	}

	return recordings, info, nil
}

func (mr *MongoRecordingRepo) Delete(accountId, recordingId string) error {
	session, err := mr.SessionProvider.GetSession()
	defer session.Close()

	if err != nil {
		logger.Logger.Errorf("error fetching mongo session: %s", err)
		return mapError(err)
	}

	collection := session.DB(mr.DB).C(mr.Collection)

	logger.Logger.Debugf("deleting recording: %v accountId: %v", recordingId, accountId)
	query := bson.M{"accountId": accountId, "recordingId": recordingId}
	err = collection.Remove(query)

	if err != nil {
		logger.Logger.Errorf("error deleting recording: %s", err)
		return mapError(err)
	}

	return nil
}
