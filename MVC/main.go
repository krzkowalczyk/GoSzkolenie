package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/krzkowalczyk/GoSzkolenie/MVC/model"
)

func main() {
	log.Println("Starting app...")

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&model.Person{})

	// Create
	// db.Create(&model.Person{Name: "Marek", Location: "Warsaw"})
	// db.Create(&model.Person{Name: "Jola", Location: "Olkusz"})

	// getAll()
	// getPersonByID(9)
	// updatePersonByID(9, "Ania", "Szczecin")
	// deletePersonByID(17)

	router := mux.NewRouter()

	router.HandleFunc("/people", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/people/{id}", GetPersonEndpoint).Methods("GET")
	router.HandleFunc("/people", CreatePersonEndpoint).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePersonEndpoint).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}
