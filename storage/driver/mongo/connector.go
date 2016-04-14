package mongo

import (
	"time"

	"gopkg.in/mgo.v2"
)

const DefaultDialTimeoutSeconds = 3

type MongoConfig struct {
	DB             string
	Username       string
	Password       string
	ServerList     []string
	DialTimeout    int
	GridFSPrefix   string
	WriteConcern   string
	ReadPreference string
}

type Connector interface {
	Configure(MongoConfig) error
	NewSession() (*mgo.Session, error)
}

type DefaultConnector struct {
	dialInfo *mgo.DialInfo
}

func (c *DefaultConnector) Configure(options MongoConfig) error {

	timeout := time.Duration(DefaultDialTimeoutSeconds) * time.Second
	if options.DialTimeout != 0 {
		timeout = time.Duration(options.DialTimeout) * time.Second
	}

	c.dialInfo = &mgo.DialInfo{
		Addrs:    options.ServerList,
		Timeout:  timeout,
		Username: options.Username,
		Password: options.Password,
	}
	return nil
}

func (c *DefaultConnector) NewSession() (*mgo.Session, error) {
	return mgo.DialWithInfo(c.dialInfo)
}
