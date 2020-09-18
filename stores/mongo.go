package stores

import (
	"time"

	"github.com/globalsign/mgo"
)

type Mongo struct {
	Name     string
	Hosts    []string
	Username string
	Password string
	//AuthDB     string
	//ReplicaSet string
	Session *mgo.Session
}

func NewMongo(hosts []string) *Mongo {
	return &Mongo{
		Name:  DB_NAME,
		Hosts: hosts,
		//Password: password,
	}
}

func (m *Mongo) Connect() error {
	var err error
	m.Session, err = mgo.DialWithInfo(&mgo.DialInfo{
		Addrs:    m.Hosts,
		Timeout:  60 * time.Second,
		Username: m.Username,
		Password: m.Password,
		//Database:       m.AuthDB,
		//ReplicaSetName: m.ReplicaSet,
	})
	if err != nil {
		return err
	}

	return nil
}
