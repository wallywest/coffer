package testutil

import (
	"os"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/dbtest"
)

const (
	MONGO_TEST_DIR string = "/tmp/recording_test_db"
)

var mongoServer *dbtest.DBServer

func SetupMongo() {
	mongoServer = &dbtest.DBServer{}
	os.MkdirAll(MONGO_TEST_DIR, 0777)
	mongoServer.SetPath(MONGO_TEST_DIR)
}

func MongoSession() *mgo.Session {
	// we need to initialize the mongoServer by calling Session()
	// the first time
	return mongoServer.Session()
}

func WipeMongo() {
	mongoServer.Wipe()
}

func StopMongo() {
	mongoServer.Stop()
	os.Remove(MONGO_TEST_DIR)
}

func SeedWithFile(filePath string) {
}
