package recording

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"time"

	"gitlab.vailsys.com/jerny/coffer/internal/logger"
	"gitlab.vailsys.com/jerny/coffer/internal/storage/driver/mongo"

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

type GFSMeta struct {
	Id          bson.ObjectId `bson:"_id"`
	ChunkSize   int           `bson:"chunkSize"`
	Name        string        `bson:"filename"`
	Length      int64         `bson:"length"`
	Md5         string        `bson:"md5"`
	UploadDate  time.Time     `bson:"uploadDate"`
	ContentType string        `bson:"contentType,omitempty"`
}

type GFSFile struct {
	FileReader  *bytes.Reader
	Md5         string
	ContentType string
}

func NewGridFSRepo(opts mongo.MongoConfig, provider *mongo.SessionProvider) *GridFSRepo {

	return &GridFSRepo{
		SessionProvider: provider,
		GFSPrefix:       opts.GridFSPrefix,
		DB:              opts.DB,
	}
}

func (repo *GridFSRepo) ListFiles() ([]*GFSMeta, error) {
	sess, err := repo.SessionProvider.GetSession()
	if err != nil {
		return nil, mapError(err)
	}

	defer sess.Close()

	gfs := sess.DB(repo.DB).GridFS(repo.GFSPrefix)

	var files []*GFSMeta
	cursor := gfs.Find(nil).Batch(DEFAULT_BATCH).Iter()
	err = cursor.All(&files)
	if err != nil {
		return nil, mapError(err)
	}

	return files, nil
}

func (repo *GridFSRepo) GetFile(accountId, recordingId string) (*GFSMeta, error) {
	sess, err := repo.SessionProvider.GetSession()
	if err != nil {
		return nil, mapError(err)
	}

	defer sess.Close()

	gfs := sess.DB(repo.DB).GridFS(repo.GFSPrefix)

	var fileMeta GFSMeta

	err = gfs.Find(bson.M{"metadata.accountId": accountId, "metadata.fileId": recordingId}).One(&fileMeta)
	if err != nil {
		return nil, mapError(err)
	}

	return &fileMeta, nil
}

func (repo *GridFSRepo) GetFileByName(name string) (*GFSMeta, error) {
	sess, err := repo.SessionProvider.GetSession()
	if err != nil {
		return nil, mapError(err)
	}

	defer sess.Close()

	gfs := sess.DB(repo.DB).GridFS(repo.GFSPrefix)

	var fileMeta GFSMeta

	err = gfs.Find(bson.M{"filename": name}).One(&fileMeta)
	if err != nil {
		return nil, mapError(err)
	}

	return &fileMeta, nil

}

func (repo *GridFSRepo) GetFileById(objectId bson.ObjectId) (*GFSMeta, error) {
	sess, err := repo.SessionProvider.GetSession()
	if err != nil {
		return nil, err
	}

	defer sess.Close()

	gfs := sess.DB(repo.DB).GridFS(repo.GFSPrefix)

	var file GFSMeta

	err = gfs.Find(bson.M{"_id": objectId}).One(&file)
	if err != nil {
		return nil, err
	}

	return &file, nil
}

func (repo *GridFSRepo) OpenById(id bson.ObjectId) (*GFSFile, error) {
	sess, err := repo.SessionProvider.GetSession()
	if err != nil {
		logger.Logger.Debugf("error retrieving session: %v", err)
		return nil, mapError(err)
	}

	defer sess.Close()

	gfs := sess.DB(repo.DB).GridFS(repo.GFSPrefix)

	file, err := gfs.OpenId(id)

	if err != nil {
		logger.Logger.Debugf("error fetching file: %v", err)
		return nil, mapError(err)
	}

	defer file.Close()

	b := make([]byte, file.Size())
	n, err := file.Read(b)
	if err != nil {
		logger.Logger.Debugf("error writing mgo file to buffer: %v", err)
		return nil, mapError(err)
	}
	logger.Logger.Debugf("wrote %v bytes", n)

	fileBytes := bytes.NewReader(b)

	f := &GFSFile{
		Md5:         file.MD5(),
		FileReader:  fileBytes,
		ContentType: file.ContentType(),
	}

	if f.Md5 == "" {
		f.Md5 = fileMd5(file)
	}

	if f.ContentType == "" {
		f.ContentType = "audio/wav"
	}

	return f, nil
}

func fileMd5(file io.Reader) string {
	h := md5.New()
	io.Copy(h, file)
	return fmt.Sprintf("%x", h.Sum(nil))
}
