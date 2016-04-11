package recording

import (
	"time"

	"gitlab.vailsys.com/jerny/coffer/pkg/mongo"

	"gopkg.in/mgo.v2/bson"
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

func (repo *GridFSRepo) GetFile(accountId, recordingId string) (*RecordingFile, error) {
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

	return &RecordingFile{GFSFile: &file}, nil
}
