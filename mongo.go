package mongodb

import (
	"github.com/globalsign/mgo"
)

type MgoSession struct {
	Session *mgo.Session
}

var session *MgoSession

// ex. mongodb://user:pass@localhost:27017/data
func InitMongoDB(url string) (*MgoSession, error) {
	nSess, err := NewMgoSession(url)
	if err != nil {
		return nil, err
	}
	session = nSess
	return session, err
}

func NewMgoSession(url string) (*MgoSession, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	return &MgoSession{session}, err
}

func GetMgoSession() *MgoSession {
	return session
}

func (s *MgoSession) Copy() *MgoSession {
	return &MgoSession{s.Session.Copy()}
}

func GetMgoCollection(db, col string) *mgo.Collection {
	return session.Session.DB(db).C(col)
}

func (s *MgoSession) Close() {
	if s.Session != nil {
		s.Session.Close()
	}
}
