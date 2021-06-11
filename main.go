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

type User struct {
	gorm.Model

	Name          string
	Email         string
	FSAIDPassword string
}

// var (
// 	user = &User{
// 		Name:          "Jack",
// 		Email:         "jack@gmail.com",
// 		FSAIDPassword: "password",
// 	}
// )

var db *gorm.DB
var err error

func main() {
	// load environment variables
	os.Setenv("HOST", "localhost")
	os.Setenv("DBPORT", "5432")
	os.Setenv("USER", "postgres")
	os.Setenv("NAME", "fafsa_data")
	os.Setenv("PASSWORD", "postgres")
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

	testUser := User{
		Name:          "Jack",
		Email:         "jack@gmail.com",
		FSAIDPassword: "password",
	}

	db.Create(&testUser)

	router := setupRouters()

	log.Fatal(http.ListenAndServe(":8080", router))
}

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

func getData(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	question := vars["question"]

	var user User

	db.First(&user, id)

	ref := reflect.ValueOf(user)
	answer := reflect.Indirect(ref).FieldByName(question)

	json.NewEncoder(w).Encode(answer)
}

func setupRouters() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/data/{id}/{question}", getData).Methods("GET")

	return router
}
