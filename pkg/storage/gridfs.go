package storage

import (
	"time"

	"gitlab.vailsys.com/jerny/coffer/cmd/coffer/options"
	"gitlab.vailsys.com/jerny/coffer/pkg/storage/mongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type GridFSProvider struct {
	SessionProvider *mongo.SessionProvider
	GFS             *mgo.GridFS
}

type GFSFile struct {
	Id          bson.ObjectId `bson:"_id"`
	ChunkSize   int           `bson:"chunkSize"`
	Name        string        `bson:"filename"`
	Length      int64         `bson:"length"`
	Md5         string        `bson:"md5"`
	UploadDate  time.Time     `bson:"uploateDate"`
	ContentType string        `bson:"contentType,omitempty"`
}

func NewGridFSProvider(opts *options.CofferConfig) (*GridFSProvider, error) {

	sp, err := mongo.NewSessionProvider(opts.MongoConfig)
	if err != nil {
		return nil, err
	}

	session, err := sp.GetSession()

	if err != nil {
		return nil, err
	}

	gfs := session.DB(opts.MongoConfig.DB).GridFS(opts.MongoConfig.GridFSPrefix)

	provider := &GridFSProvider{
		SessionProvider: sp,
		GFS:             gfs,
	}

	return provider, nil
}
