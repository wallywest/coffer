package mongo

import (
	"fmt"
	"sync"

	"gopkg.in/mgo.v2"
)

//add const for errors

//SessionProvider
type SessionProvider struct {
	connector Connector

	masterSessionLock sync.Mutex

	masterSession *mgo.Session

	readPreference mgo.Mode
}

func NewSessionProvider(opts MongoConfig) (*SessionProvider, error) {
	provider := &SessionProvider{
		readPreference: mgo.Primary,
	}

	provider.connector = connectorFactory(opts)

	err := provider.connector.Configure(opts)
	if err != nil {
		return nil, err
	}
	return provider, nil
}

func (p *SessionProvider) GetSession() (*mgo.Session, error) {
	p.masterSessionLock.Lock()
	defer p.masterSessionLock.Unlock()

	if p.masterSession != nil {
		return p.masterSession.Copy(), nil
	}

	var err error
	p.masterSession, err = p.connector.NewSession()

	if err != nil {
		return nil, fmt.Errorf("unable to connect to server: %v", err)
	}

	p.refresh()

	return p.masterSession.Copy(), nil
}

func (p *SessionProvider) Close() {
	if p.masterSession != nil {
		p.masterSession.Close()
	}
}

func (p *SessionProvider) refresh() {
	p.masterSession.SetMode(p.readPreference, true)
}

func connectorFactory(opts MongoConfig) Connector {
	return &DefaultConnector{}
}
