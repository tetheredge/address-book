package main

import (
	"fmt"
	//"gopkg.in/mgo.v2"
	"log"
	"net/http"
	//"gopkg.in/mgo.v2/bson"
)

func main() {
	fmt.Println("Welcome to the online address book...")
	var testDb = false
	dbSession := SetupDb(testDb)
	mongo := NewSession(&dbSession)
	app := new(AppContext)
	app.Mongo = mongo
	router := NewRouter(app)

	log.Fatal(http.ListenAndServe(":8088", router))
}
