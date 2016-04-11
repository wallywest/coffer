package testutil

import (
	"fmt"
	"os"

	"github.com/mongodb/mongo-tools/common/db"
	mongoLogger "github.com/mongodb/mongo-tools/common/log"
	"github.com/mongodb/mongo-tools/common/options"
	"github.com/mongodb/mongo-tools/mongofiles"
	"github.com/mongodb/mongo-tools/mongoimport"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/dbtest"
)

const (
	MONGO_TEST_DIR string = "/tmp/recording_test_db"
)

var mongoServer *dbtest.DBServer
var mimport mongoimport.MongoImport

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

func SeedWithFile(filePath string) (uint64, error) {
	session := mongoServer.Session()
	defer session.Close()
	servers := session.LiveServers()

	opts := options.New("mongoimport", "", options.EnabledOptions{Connection: true, Namespace: true})
	opts.Direct = true
	opts.DB = "test"
	opts.Collection = "recordings"
	opts.Host = servers[0]
	opts.Quiet = true

	mongoLogger.SetVerbosity(opts.Verbosity)

	sessionProvider, err := db.NewSessionProvider(*opts)
	if err != nil {
		fmt.Println("error establishing session")
		return 0, err
	}

	inputOpts := &mongoimport.InputOptions{
		File: filePath,
	}
	ingestOptions := &mongoimport.IngestOptions{}

	mimport = mongoimport.MongoImport{
		ToolOptions:     opts,
		InputOptions:    inputOpts,
		IngestOptions:   ingestOptions,
		SessionProvider: sessionProvider,
	}

	err = mimport.ValidateSettings([]string{})
	if err != nil {
		fmt.Println("error validating input")
		return 0, err
	}

	numDocs, err := mimport.ImportDocuments()

	if err != nil {
		return 0, err
	}

	return numDocs, nil
}

func SeedAsset(filePath string) error {
	session := mongoServer.Session()
	defer session.Close()
	servers := session.LiveServers()

	opts := options.New("mongofiles", "", options.EnabledOptions{Connection: true, Namespace: true})
	opts.Direct = true
	opts.DB = "test"
	opts.Host = servers[0]
	opts.Quiet = true

	sessionProvider, err := db.NewSessionProvider(*opts)
	if err != nil {
		return err
	}

	mongoLogger.SetVerbosity(opts.Verbosity)

	inputOpts := &mongofiles.InputOptions{}
	storageOpts := &mongofiles.StorageOptions{
		DB:           opts.DB,
		GridFSPrefix: "testfs",
	}

	mf := mongofiles.MongoFiles{
		ToolOptions:     opts,
		StorageOptions:  storageOpts,
		SessionProvider: sessionProvider,
		InputOptions:    inputOpts,
		Command:         "put",
		FileName:        filePath,
	}

	_, err = mf.Run(true)
	if err != nil {
		return err
	}

	return nil
}
