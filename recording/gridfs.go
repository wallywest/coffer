package recording

import (
	"bytes"
	"io"
	"time"

	"gitlab.vailsys.com/jerny/coffer/pkg/logger"
	"gitlab.vailsys.com/jerny/coffer/storage/driver/mongo"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	DEFAULT_BATCH = 10
)

type GridFSRepo struct {
	SessionProvider *mongo.SessionProvider
	DB              string
	GFSPrefix       string
}

type GFSFile struct {
	Id          bson.ObjectId `bson:"_id"`
	ChunkSize   int           `bson:"chunkSize"`
	Name        string        `bson:"filename"`
	Length      int64         `bson:"length"`
	Md5         string        `bson:"md5"`
	UploadDate  time.Time     `bson:"uploadDate"`
	ContentType string        `bson:"contentType,omitempty"`
}

func NewGridFSRepo(opts mongo.MongoConfig, provider *mongo.SessionProvider) *GridFSRepo {

	return &GridFSRepo{
		SessionProvider: provider,
		GFSPrefix:       opts.GridFSPrefix,
		DB:              opts.DB,
	}
}

func (repo *GridFSRepo) ListFiles() ([]*GFSFile, error) {
	sess, err := repo.SessionProvider.GetSession()
	if err != nil {
		return nil, err
	}

	defer sess.Close()

	gfs := sess.DB(repo.DB).GridFS(repo.GFSPrefix)

	var files []*GFSFile
	cursor := gfs.Find(nil).Batch(DEFAULT_BATCH).Iter()
	cursor.All(&files)

	return files, nil
}

func (repo *GridFSRepo) GetFile(accountId, recordingId string) (*GFSFile, error) {
	sess, err := repo.SessionProvider.GetSession()
	if err != nil {
		return nil, err
	}

	defer sess.Close()

	gfs := sess.DB(repo.DB).GridFS(repo.GFSPrefix)

	var file GFSFile

	err = gfs.Find(bson.M{"metadata.accountId": accountId, "metadata.fileId": recordingId}).One(&file)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (repo *GridFSRepo) GetFileByName(name string) (*GFSFile, error) {
	sess, err := repo.SessionProvider.GetSession()
	if err != nil {
		return nil, err
	}

	defer sess.Close()

	gfs := sess.DB(repo.DB).GridFS(repo.GFSPrefix)

	var file GFSFile

	err = gfs.Find(bson.M{"filename": name}).One(&file)
	if err != nil {
		return nil, err
	}

	return &file, nil

}

func (repo *GridFSRepo) GetFileById(objectId bson.ObjectId) (*GFSFile, error) {
	sess, err := repo.SessionProvider.GetSession()
	if err != nil {
		return nil, err
	}

	defer sess.Close()

	gfs := sess.DB(repo.DB).GridFS(repo.GFSPrefix)

	var file GFSFile

	err = gfs.Find(bson.M{"_id": objectId}).One(&file)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (repo *GridFSRepo) OpenByName(name string) (*mgo.GridFile, error) {
	sess, err := repo.SessionProvider.GetSession()
	if err != nil {
		return nil, err
	}

	defer sess.Close()

	gfs := sess.DB(repo.DB).GridFS(repo.GFSPrefix)

	file, err := gfs.Open(name)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (repo *GridFSRepo) OpenById(id bson.ObjectId) (io.ReadSeeker, error) {
	sess, err := repo.SessionProvider.GetSession()
	if err != nil {
		logger.Logger.Debugf("error retrieving session: %v", err)
		return nil, err
	}

	defer sess.Close()

	gfs := sess.DB(repo.DB).GridFS(repo.GFSPrefix)

	file, err := gfs.OpenId(id)
	defer file.Close()

	if err != nil {
		logger.Logger.Debugf("error fetching file: %v", err)
		return nil, err
	}

	b := make([]byte, file.Size())
	n, err := file.Read(b)
	if err != nil {
		logger.Logger.Debugf("error writing mgo file to buffer: %v", err)
		return nil, err
	}
	logger.Logger.Debugf("wrote %v bytes", n)

	fileBytes := bytes.NewReader(b)

	return fileBytes, nil
}
