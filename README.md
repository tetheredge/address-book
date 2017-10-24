Address Book
============

This tool allows for the ability to store information about people into a mongodb datastore. This is an API which allows you to list all entries, show a specific entry, add an entry, modify and entry and delete an entry.  There is also the ability to import entries via a CSV file, and export to a CSV file.

>Notes:
>You will need to use golang version 1.5+. This repo uses
version 1.9, which by default supports a vendor directory to
store git submodules.  If using version 1.5, you will need
to set the environemnt variable, `GO15VENDOREXPERIMENT=1 `, so
that it will use the dependencies in the vendor directory.

Getting Started
====

Get started by cloning down the repository:
```shell
git clone git@github.com:tetheredge/address-book.git
```
Once cloned, run the following command to pull down the external dependencies.
```shell
go get
```

This application requires that certain environment variables be setup in order to connect to the mongodb database. They are as follows:
```shell
export MONGO_HOST=name_of_mongo_host
export MONGO_PORT=port_mongo_is_running_on
export MONGO_DB_ADDRESS_TEST=test_db_instance
export MONGO_DB_ADDRESS=production_db_instance
```

This repo also uses git submodules, to activate those run the following command:
```shell
git submodule init
```

Some of the curl commands to test the api with are:
```shell
To add a new address entry
curl -v -d '{"first_name":"Taylor", "last_name":"Etheredge", "email":"taylor.etheredge@gmail.comm", "phone_number":"972-885-9135"}' http://localhost:8088/person
```
To update an adddress entry
```shell
curl -v  -X "PUT" -d '{"id": "59ea04f3ae6b7d8d5f418b7a", "first_name":"Taylor", "last_name":"Ethredge", "email":"taylor.etheredge@gmail.com", "phone_number":"972-885-9135"}' http://localhost:8088/person
```
To delete an address entry
```shell
curl -v -X "DELETE" -d '{"id": "59ea04f3ae6b7d8d5f418b7a", "first_name":"Taylor", "last_name":"Ethredge", "email":"taylor.etheredge@gmail.com", "phone_number":"972-885-9135"}' http://localhost:8088/person
```
To upload a csv file for importing bulk address entries
```shell
curl http://localhost:8088/person/upload -vvv -F "people=@people.csv"
```

