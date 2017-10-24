package main

import (
	"gopkg.in/mgo.v2"
	"os"
)

type Mongo struct {
	host     string
	err      error
	database string
	port     string
	session  *mgo.Session
	collects *mgo.Collection
}

type DbSession struct {
	Host     string
	Port     string
	Database string
}

func SetupDb(testDb bool) DbSession {
	var addr string
	var port string
	var db string

	dbSession := DbSession{}

	if addr = os.Getenv("MONGODB_PORT_27017_TCP_ADDR"); addr != "" {
		dbSession.Host = addr
	} else {
		dbSession.Host = os.Getenv("MONGO_HOST")
	}

	if port = os.Getenv("MONGODB_PORT_27017_TCP_PORT"); port != "" {
		dbSession.Port = port
	} else {
		dbSession.Port = os.Getenv("MONGO_PORT")
	}

	if testDb == true {
		db = os.Getenv("MONGO_DB_ADDRESS_TEST")
	} else {
		db = os.Getenv("MONGO_DB_ADDRESS")
	}

	dbSession.Database = db

	return dbSession
}

func NewSession(dbSession *DbSession) *Mongo {
	session, err := mgo.Dial("mongodb://" + dbSession.Host + ":" + dbSession.Port)

	if err != nil {
		panic(err)
	}

	addIndex("first_name", session.DB(dbSession.Database))
	addIndex("last_name", session.DB(dbSession.Database))
	addIndex("email", session.DB(dbSession.Database))
	addIndex("phone_number", session.DB(dbSession.Database))
	return &Mongo{session: session, host: dbSession.Host, database: dbSession.Database}
}

func addIndex(column string, db *mgo.Database) {
	index := mgo.Index{
		Key:      []string{"column"},
		Unique:   false,
		DropDups: true,
	}

	indexErr := db.C("people").EnsureIndex(index)
	if indexErr != nil {
		panic(indexErr)
	}
}

func CloseDb(app *Mongo) {
	app.session.Close()
}
