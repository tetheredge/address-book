package main

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"testing"
)

func setupTestDb() *AppContext {
	testDb := true
	dbSession := SetupDb(testDb)
	mongo := NewSession(&dbSession)
	app := new(AppContext)
	app.Mongo = mongo
	return app
}

func insertPerson(app *AppContext) {
	person := Person{
		FirstName: "Taylor",
		LastName:  "Etheredge",
		Email:     "taylor.etheredge@gmail.com",
		Phone:     "972-885-9135",
	}

	c := app.Mongo.session.DB(app.Mongo.database).C("people")
	c.Insert(person)
}

func removePerson(app *AppContext, person *Person) {
	c := app.Mongo.session.DB(app.Mongo.database).C("people")
	c.Remove(person)
}

func TestListDetails(t *testing.T) {
	t.Run("A=successfulLookup", func(t *testing.T) {
		app := setupTestDb()
		insertPerson(app)
		person := Person{}

		per, err := person.ListDetails(app, "Taylor")
		if err != nil {
			t.Errorf("Expected there to not be an error, but got %v", err)
		}
		removePerson(app, &per)
		CloseDb(app.Mongo)
	})

	t.Run("A=failedLookup", func(t *testing.T) {
		app := setupTestDb()
		insertPerson(app)
		person := Person{}

		_, err := person.ListDetails(app, "Zack")
		if err == nil {
			t.Errorf("Expected to have an error when a person does not exist, but got %v", err)
		}
		per := Person{
			FirstName: "Taylor",
			LastName:  "Etheredge",
			Email:     "taylor.etheredge@gmail.com",
			Phone:     "972-885-9135",
		}
		removePerson(app, &per)
		CloseDb(app.Mongo)
	})
}

func TestAddPerson(t *testing.T) {
	t.Run("A=successfulInsert", func(t *testing.T) {

		app := setupTestDb()
		person := Person{}

		per, err := person.AddPerson(app, "Zack", "Pire", "ZPire@insourcegroup.com", "972-455-1016")

		if err != nil {
			t.Errorf("Expected to not have an error, but got %v", err)
		}
		removePerson(app, per)
		CloseDb(app.Mongo)
	})

	t.Run("A=failedInsert", func(*testing.T) {
		app := setupTestDb()
		person := Person{}

		per, err := person.AddPerson(app, "", "Pire", "ZPire@insourcegroup.com", "972-455-1016")
		fmt.Println(err)
		if err == nil {

			t.Errorf("Expected to have an error, but got %v", err)
		}
		removePerson(app, per)
		CloseDb(app.Mongo)
	})
}

func TestModifyPerson(t *testing.T) {
	t.Run("A=successfulUpdate", func(t *testing.T) {
		app := setupTestDb()
		person := Person{
			FirstName: "Taylor",
			LastName:  "Etheredge",
			Email:     "taylor.etheredge@gmail.com",
			Phone:     "972-400-5373",
		}
		c := app.Mongo.session.DB(app.Mongo.database).C("people")
		c.Insert(person)

		per := Person{}
		c.Find(bson.M{"first_name": "Taylor"}).One(&per)

		per.Phone = "972-885-9135"

		result, err := per.ModifyPerson(app, &per)

		if err != nil {
			t.Errorf("Expected err to be nil, but got %v", err)
		}

		if result.Phone != "972-885-9135" {
			t.Errorf("Expected the phone number to be updated, but got %v", result.Phone)
		}
		removePerson(app, &per)
		CloseDb(app.Mongo)
	})

	t.Run("A=failedUpdate", func(t *testing.T) {
		app := setupTestDb()
		person := Person{
			FirstName: "Taylor",
			LastName:  "Etheredge",
			Email:     "taylor.etheredge@gmail.com",
			Phone:     "972-400-5373",
		}
		c := app.Mongo.session.DB(app.Mongo.database).C("people")
		c.Insert(person)

		per := Person{}
		c.Find(bson.M{"first_name": "Taylor"}).One(&per)

		per.Phone = ""

		_, err := per.ModifyPerson(app, &per)

		if err == nil {
			t.Errorf("Expected to have an error, but got %v", err)
		}

		removePerson(app, &per)
		CloseDb(app.Mongo)
	})
}

func TestRemove(t *testing.T) {
	t.Run("A=successfulRemoval", func(t *testing.T) {
		app := setupTestDb()
		insertPerson(app)
		person := Person{}

		per := Person{}
		c := app.Mongo.session.DB(app.Mongo.database).C("people")
		c.Find(bson.M{"first_name": "Taylor"}).One(&per)
		_, err := person.RemovePerson(app, &per)

		if err != nil {
			t.Errorf("Expected to not have an error, but got %v", err)
		}

		CloseDb(app.Mongo)
	})

	t.Run("A=failedRemoval", func(t *testing.T) {
		app := setupTestDb()
		insertPerson(app)
		person := Person{}

		per := Person{}
		c := app.Mongo.session.DB(app.Mongo.database).C("people")
		c.Find(bson.M{"first_name": "Zack"}).One(&per)
		_, err := person.RemovePerson(app, &per)

		if err == nil {
			t.Errorf("Expected an error to occur, but got %v", err)
		}
		c.Find(bson.M{"first_name": "Taylor"}).One(&per)
		removePerson(app, &per)
		CloseDb(app.Mongo)
	})
}

// func TestCSVImport(t *testing.T) {
// 	t.Run("A=successfulImport", func(t *testing.T) {
// 		app := setupTestDb()
// 		person := Person{}
// 		person.CSVImport(app)
// 	})
// }
