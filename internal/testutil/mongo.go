package testutil

import (
	"bytes"
	"io"
	"os"
	"os/exec"

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

func SeedWithFile(filePath string) error {
	session := mongoServer.Session()
	defer session.Close()
	servers := session.LiveServers()

	cmd := exec.Command("mongoimport", "-h", servers[0], "-d", "test", "-c", "recordings", "--file", filePath)

	var out bytes.Buffer
	var cmdErr bytes.Buffer
	cmd.Stderr = &cmdErr
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func SeedAsset(fileName, filePath string) error {
	session := mongoServer.Session()
	defer session.Close()
	servers := session.LiveServers()

	cmd := exec.Command("mongofiles", "-h", servers[0], "-d", "test", "--prefix", "testfs", "-l", filePath, "put", fileName)

	var out bytes.Buffer
	var cmdErr bytes.Buffer
	cmd.Stderr = &cmdErr
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil

}

func SeedAssetWithMeta(fileName string, file *os.File, meta interface{}) error {
	session := mongoServer.Session()
	defer session.Close()

	gfs := session.DB("test").GridFS("testfs")

	f, err := gfs.Create(fileName)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, file)
	if err != nil {
		return err
	}

	f.SetMeta(meta)

	err = f.Close()
	if err != nil {
		return err
	}

	return nil
}
