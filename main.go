package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
)

func main() {
	//  init database
	connStr := "postgres://postgres:love@localhost/imp?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Database initialized.")
	defer db.Close()

	r := mux.NewRouter()
	r.HandleFunc("/", rootEndpoint)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(`:%s`, "8000"), r))
}

func rootEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "imp-assessment")
}
