package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

func ListPerson(app *AppContext, w http.ResponseWriter, r *http.Request) (*jsonData, *appError) {
	var person Person
	params := mux.Vars(r)
	personDetails := params["person"]
	per, err := person.ListDetails(app, personDetails)

	if err != nil {
		log.Printf("%#v", err)
		return nil, &appError{
			Code:    http.StatusNotFound,
			Message: "The person, " + personDetails + ", does not exist in the db.",
			Err:     err.Error(),
		}
	}

	buf, _ := json.MarshalIndent(per, "", "\t")
	return &jsonData{
		Code: http.StatusOK,
		Byte: buf,
	}, nil
}

func ListAllEntries(app *AppContext, w http.ResponseWriter, r *http.Request) (*jsonData, *appError) {
	var person Person
	people, err := person.ListAllEntries(app)
	if err != nil {
		return nil, &appError{
			Code:    http.StatusNotFound,
			Message: "There are no entries stored in the db.",
			Err:     err.Error(),
		}
	}
	buf, _ := json.MarshalIndent(people, "", "\t")
	return &jsonData{
		Code: http.StatusOK,
		Byte: buf,
	}, nil
}

func AddPerson(app *AppContext, w http.ResponseWriter, r *http.Request) (*jsonData, *appError) {
	w.Header().Set("Content-Type", "application/json")

	var person Person
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &person)

	if err != nil {
		log.Printf("%v", err)
		return nil, &appError{
			Code:    http.StatusBadRequest,
			Message: "Bad request, invalid post data must be in JSON.",
			Err:     err.Error(),
		}
	}

	result, err := person.AddPerson(app, person.FirstName, person.LastName, person.Email, person.Phone)

	if err != nil {
		log.Printf("%v", err)
		return nil, &appError{
			Code:    http.StatusConflict,
			Message: "Person already exists in the db.",
			Err:     err.Error(),
		}
	}

	buf, _ := json.MarshalIndent(result, "", "\t")

	return &jsonData{
		Code: http.StatusCreated,
		Byte: buf,
	}, nil
}

func UpdatePerson(app *AppContext, w http.ResponseWriter, r *http.Request) (*jsonData, *appError) {
	w.Header().Set("Content-Type", "application/json")

	var person Person
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &person)

	if err != nil {
		log.Printf("%v", err)
		return nil, &appError{
			Code:    http.StatusBadRequest,
			Message: "Bad request, invalid post data must be in JSON.",
			Err:     err.Error(),
		}
	}

	result, err := person.ModifyPerson(app, &person)

	if err != nil {
		log.Printf("%v", err)
		return nil, &appError{
			Code:    http.StatusNotModified,
			Message: "Person does not exist in the db.",
			Err:     err.Error(),
		}
	}

	buf, _ := json.MarshalIndent(result, "", "\t")

	return &jsonData{
		Code: http.StatusOK,
		Byte: buf,
	}, nil
}

func RemovePerson(app *AppContext, w http.ResponseWriter, r *http.Request) (*jsonData, *appError) {
	w.Header().Set("Content-Type", "application/json")

	var person Person
	body, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(body, &person)

	if err != nil {
		log.Printf("%v", err)
		return nil, &appError{
			Code:    http.StatusBadRequest,
			Message: "Bad request, invalid post data must be in JSON.",
			Err:     err.Error(),
		}
	}

	result, err := person.RemovePerson(app, &person)
	log.Println(result)
	if err != nil {
		log.Printf("%v", err)
		return nil, &appError{
			Code:    http.StatusNotModified,
			Message: "Person does not exist in the db.",
			Err:     err.Error(),
		}
	}

	buf, _ := json.MarshalIndent(result, "", "\t")

	return &jsonData{
		Code: http.StatusOK,
		Byte: buf,
	}, nil
}

func ImportCSV(app *AppContext, w http.ResponseWriter, r *http.Request) (*jsonData, *appError) {
	w.Header().Set("Content-Type", "application/csv")
	var Buf bytes.Buffer
	file, _, err := r.FormFile("people")
	if err != nil {
		return nil, &appError{
			Code:    http.StatusNotFound,
			Message: "File is not found.",
			Err:     err.Error(),
		}
	}
	io.Copy(&Buf, file)
	contents := Buf.String()
	contentsToString := strings.NewReader(contents)
	x := csv.NewReader(bufio.NewReader(contentsToString))
	var people []Person
	person := Person{}
	record, _ := x.ReadAll()
	for _, i := range record {
		people = append(people, Person{FirstName: i[0], LastName: i[1], Email: i[2], Phone: i[3]})
	}
	person.CSVImport(app, people)

	return &jsonData{}, nil
}

var per People

func ExportCSV(app *AppContext, w http.ResponseWriter, r *http.Request) (*jsonData, *appError) {
	w.Header().Set("Content-Type", "application/csv")
	file, err := os.Create("export-people.csv")
	if err != nil {
		return nil, &appError{
			Code:    http.StatusNotFound,
			Message: "File is not found.",
			Err:     err.Error(),
		}
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	var people []Person
	c := app.Mongo.session.DB(app.Mongo.database).C("people")
	err = c.Find(nil).All(&people)
	if err != nil {
		return nil, &appError{
			Code:    http.StatusNotFound,
			Message: "Result does not exist in the db.",
			Err:     err.Error(),
		}
	}
	person := Person{}
	person.CSVExport(app, people, writer)
	return &jsonData{}, nil

}
