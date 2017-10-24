package main

import (
	//"bytes"
	//"gopkg.in/mgo.v2"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListPerson(t *testing.T) {
	t.Run("A=successfulListing", func(t *testing.T) {
		app := setupTestDb()
		person := &Person{
			Id:        bson.ObjectIdHex("757d0838d89d509c93164d9a"),
			FirstName: "Taylor",
			LastName:  "Etheredge",
			Email:     "taylor.etheredge@gmail.com",
			Phone:     "972-885-9135",
		}
		c := app.Mongo.session.DB(app.Mongo.database).C("people")
		c.Insert(person)
		//requestB := []byte(`{"Id": bson.ObjectIdHex("757d0838d89d509c93164d9a"), "FirstName: "Taylor", "LastName": "Etheredge", "Email": "taylor.etheredge@gmail.com", "Phone": "972-885-9135"}`)
		request, _ := http.NewRequest("GET", "/person/Taylor", nil)
		request.Header.Set("Content-Type", "application/json; charset=UTF-8")

		w := httptest.NewRecorder()
		NewRouter(app).ServeHTTP(w, request)

		per := Person{}
		c.Find(bson.M{"first_name": "Taylor"}).One(&per)

		responseBody :=
			`{ 
	"id": "757d0838d89d509c93164d9a",
	"first_name": "Taylor",
	"last_name": "Etheredge", 
	"email": "taylor.etheredge@gmail.com",
	"phone_number": "972-885-9135"
}`
		fmt.Printf("%T", w.Body.String())
		fmt.Printf("%T", responseBody)
		if w.Body.String() != responseBody {
			t.Errorf("Expected the body to be %v, but got %v", responseBody, w.Body.String())
		}
		removePerson(app, &per)
	})
}
