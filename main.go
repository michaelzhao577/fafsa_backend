package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// User model for postgres DB
type User struct {
	gorm.Model
	First_Name    string
	Last_Name     string
	SSN           string
	DOB           string
	Email         string
	FSAIDPassword string
}

// map for checking whether questions are valid
doesExist := map[string]bool {
		"Name": true, 
		"Email": true,
		"FSAIDPassword": true
}

var db *gorm.DB
var err error

func main() {
	// load environment variables
	os.Setenv("HOST", "localhost")
	os.Setenv("DBPORT", "5432")
	os.Setenv("USER", "postgres")
	os.Setenv("NAME", "fafsa_data")
	os.Setenv("PASSWORD", "Wzc@jw0724")
	host := os.Getenv("HOST")
	dbPort := os.Getenv("DBPORT")
	user := os.Getenv("USER")
	dbName := os.Getenv("NAME")
	password := os.Getenv("PASSWORD")

	// database connection
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s port=%s", host, user, dbName, password, dbPort)

	// open connection to db
	db, err = gorm.Open("postgres", dbURI)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("Successfully connected to database")
	}

	// close connection to db when main func terminates
	defer db.Close()

	// make migration to the db if they have not already been created
	db.AutoMigrate(&User{})

	// testUser := User{
	// 	First_Name:    "Jack",
	// 	Last_Name:     "Doe",
	// 	SSN:           "123456789",
	// 	DOB:           "07/21/1999",
	// 	Email:         "jack@gmail.com",
	// 	FSAIDPassword: "password",
	// }

	// db.Create(&testUser)

	router := setupRouters()

	log.Fatal(http.ListenAndServe(":8080", router))
}

//create new user
func storeData(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	createdUser := db.Create(&user)
	err = createdUser.Error
	if err != nil {
		json.NewEncoder(w).Encode(err)
	} else {
		json.NewEncoder(w).Encode(&user)
	}
}

//get a user's answer to a question specified in http request
func getData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	question := string(vars["question"])

	var user User
	// Escaping arguments to avoid SQL injections following https://gorm.io/docs/security.html
	db.Where("id = ?", id).First(&user)

	ref := reflect.ValueOf(user)
	answer := reflect.Indirect(ref).FieldByName(question)
	output := answer.Interface().(string)
	if output == "" {
		json.NewEncoder(w).Encode("answer to question " + question + " has not been set")
	} else {
		json.NewEncoder(w).Encode(output)
	}
}

//get a user's entire dataset
func getAllData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	var user User

	db.First(&user, id)

	json.NewEncoder(w).Encode(&user)
}

// update a user's data
func updateData(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	vars := mux.Vars(r)

	id := vars["id"]

	var user User

	db.First(&user, id)

}

func deleteData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

	var user User

	db.First(&user, id)
	db.Delete(&user)

	json.NewEncoder(w).Encode(&user)
}

func setupRouters() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/data/{id}/{question}", getData).Methods("GET")
	router.HandleFunc("/data/{id}", getAllData).Methods("GET")
	router.HandleFunc("/data/{id}").Methods("DELETE")
	router.HandleFunc("/data", storeData).Methods("POST")

	return router
}
