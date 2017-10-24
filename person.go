package main

import (
	"encoding/csv"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type Person struct {
	Id        bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName string        `bson:"first_name" json:"first_name"`
	LastName  string        `bson:"last_name" json:"last_name"`
	Email     string        `bson:"email" json:"email"`
	Phone     string        `bson:"phone_number" json:"phone_number"`
}

type People []Person

type NoSuchPersonErr struct {
	Name string
}

func (err *NoSuchPersonErr) Error() string {
	return "Person, " + err.Name + ", does not exist in the db"
}

func (p *Person) ListDetails(app *AppContext, person string) (per Person, err error) {
	c := app.Mongo.session.DB(app.Mongo.database).C("people")
	err = c.Find(bson.M{"first_name": person}).One(&per)

	if err != nil {
		log.Printf("%#v", err)
		return per, &NoSuchPersonErr{Name: person}
	}

	return per, nil
}

func (p *Person) ListAllEntries(app *AppContext) (people People, err error) {
	c := app.Mongo.session.DB(app.Mongo.database).C("people")
	err = c.Find(nil).All(&people)
	if err != nil {
		return nil, err
	}
	return people, nil
}

type PersonValidationErr struct {
	Item string
}

func (err *PersonValidationErr) Error() string {
	return err.Item
}

func (p *Person) valid() bool {
	return len(p.FirstName) > 0 && len(p.LastName) > 0 && len(p.Email) > 0 && len(p.Phone) > 0
}

func (p *Person) AddPerson(app *AppContext, first_name, last_name, email, phone_number string) (*Person, error) {
	per := Person{
		FirstName: first_name,
		LastName:  last_name,
		Email:     email,
		Phone:     phone_number,
	}
	if per.valid() == false {
		return nil, &PersonValidationErr{Item: "First Name, Last Name, email and Phone Number cannot be blank."}
	}

	c := app.Mongo.session.DB(app.Mongo.database).C("people")
	err := c.Insert(per)

	if mgo.IsDup(err) {
		if err != nil {
			return nil, err
		}
	}

	return &per, err
}

func (p *Person) ModifyPerson(app *AppContext, person *Person) (*Person, error) {
	if person.valid() == false {
		return nil, &PersonValidationErr{Item: "First Name, Last Name, email and Phone Number cannot be blank."}
	}

	c := app.Mongo.session.DB(app.Mongo.database).C("people")
	err := c.Update(bson.M{"_id": person.Id}, bson.M{"$set": bson.M{"first_name": person.FirstName, "last_name": person.LastName, "email": person.Email, "phone_number": person.Phone}})

	if err != nil {
		return nil, err
	}

	return person, nil
}

func (p *Person) RemovePerson(app *AppContext, person *Person) (*Person, error) {
	if person.valid() == false {
		return nil, &PersonValidationErr{Item: "First Name, Last Name, email and Phone Number cannot be blank."}
	}

	per := Person{}
	c := app.Mongo.session.DB(app.Mongo.database).C("people")
	err := c.Find(bson.M{"first_name": person.FirstName}).One(&per)

	if err != nil {
		log.Printf("%#v", err)
		return &per, &NoSuchPersonErr{Name: person.FirstName}
	}
	err = c.Remove(per)

	if err != nil {
		return nil, err
	}

	return person, nil
}

func (p *Person) CSVImport(app *AppContext, people []Person) error {
	c := app.Mongo.session.DB(app.Mongo.database).C("people")
	for _, i := range people {

		i.Id = bson.NewObjectId()
		err := c.Insert(i)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	}

	return nil
}

func (p *Person) CSVExport(app *AppContext, people People, writer *csv.Writer) {
	for _, i := range people {
		slice := []string{i.FirstName, i.LastName, i.Email, i.Phone}
		if err := writer.Write((slice)); err != nil {
			log.Println("error writing records to csv: ", err)
		}
	}
}
