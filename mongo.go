package mongodb

import (
	"log"

	"github.com/globalsign/mgo"
)

type MgoSession struct {
	Session *mgo.Session
}

var session MgoSession

func InitMongoDB(url string) error {
	nSess, err := NewMgoSession(url)
	if err != nil {
		log.Fatalln("Can't Create Session: ", err)
	}
	session = *nSess
	return err
}

func NewMgoSession(url string) (*MgoSession, error) {
	session, err := mgo.Dial(url)
	if err != nil {
		return nil, err
	}
	return &MgoSession{session}, err
}

func GetMgoSession() *MgoSession {
	return &session
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
