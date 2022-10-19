package mongodb

import (
	"fmt"
	"strings"
	"time"

	"gopkg.in/mgo.v2"
)

type Index struct {
	Key        []string // Index key fields; prefix name with dash (-) for descending order
	Unique     bool     // Prevent two documents from having the same index key
	DropDups   bool     // Drop documents with the same index key as a previously indexed one
	Background bool     // Build index in background and return immediately
	Sparse     bool     // Only index documents containing the Key fields

	// If ExpireAfter is defined the server will periodically delete
	// documents with indexed time.Time older than the provided delta.
	ExpireAfter time.Duration

	// Name holds the stored index name. On creation if this field is unset it is
	// computed by EnsureIndex based on the index key.
	Name string

	// Properties for spatial indexes.
	//
	// Min and Max were improperly typed as int when they should have been
	// floats.  To preserve backwards compatibility they are still typed as
	// int and the following two fields enable reading and writing the same
	// fields as float numbers. In mgo.v3, these fields will be dropped and
	// Min/Max will become floats.
	Min, Max   int
	Minf, Maxf float64
	BucketSize float64
	Bits       int

	// Properties for text indexes.
	DefaultLanguage  string
	LanguageOverride string

	// Weights defines the significance of provided fields relative to other
	// fields in a text index. The score for a given word in a document is derived
	// from the weighted sum of the frequency for each of the indexed fields in
	// that document. The default field weight is 1.
	Weights map[string]int
}

type MgoSession struct {
	Session *mgo.Session
	Options Option
}

type Option struct {
	ReconnectTime int
}

var session *MgoSession
var multiSession map[string]*MgoSession

// ex. mongodb://user:pass@localhost:27017/data
func InitMongoDB(url string) (*MgoSession, error) {
	if strings.Index(url, "mongodb://") != 0 {
		url = fmt.Sprintf("mongodb://%s", url)
	}
	nSess, err := NewMgoSession(url)
	if err != nil {
		return nil, err
	}
	session = nSess
	return session, nil
}

func InitMultiMongoDB(url, name string) (*MgoSession, error) {
	if multiSession == nil {
		multiSession = make(map[string]*MgoSession)
	}
	if multiSession[name] != nil {
		return multiSession[name], nil
	}
	if strings.Index(url, "mongodb://") != 0 {
		url = fmt.Sprintf("mongodb://%s", url)
	}

	nSess, err := NewMgoSession(url)
	if err != nil {
		return nil, err
	}
	multiSession[name] = nSess
	return multiSession[name], nil
}

func NewMgoSession(url string) (*MgoSession, error) {
	session, err := mgo.DialWithTimeout(url, time.Second*3)
	if err != nil {
		return nil, err
	}
	session.SetSocketTimeout(1 * time.Hour)
	return &MgoSession{Session: session, Options: Option{
		ReconnectTime: 15,
	}}, err
}

func GetMgoSession() *MgoSession {
	keepAliveSession(session.Session)
	return session
}

func (s *MgoSession) Copy() *MgoSession {
	return &MgoSession{Session: s.Session.Copy(), Options: s.Options}
}

func GetMgoCollection(db, col string) *mgo.Collection {
	keepAliveSession(session.Session)
	return session.Session.DB(db).C(col)
}

func keepAliveSession(sess *mgo.Session) {
	for i := 0; i < session.Options.ReconnectTime; i++ {
		if err := sess.Ping(); err != nil {
			sess.Refresh()
			time.Sleep(time.Second)
		} else {
			break
		}
	}
}

func MultiGetMgoCollection(name, db, col string) (*mgo.Collection, error) {
	if multiSession[name] != nil {
		sess := multiSession[name].Session
		keepAliveSession(sess)
		return sess.DB(db).C(col), nil
	}
	return nil, fmt.Errorf("No session %s", name)
}

func (s *MgoSession) Close() {
	if s.Session != nil {
		s.Session.Close()
	}
}
