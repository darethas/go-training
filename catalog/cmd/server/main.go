package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	// in go, blank imports are sometimes used to import a package only for it's side effects
	// a.k.a to run the package's "init" function.
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	// set up our database connection
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/?parseTime=true",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)
	// setup an open connection to the database using the mysql driver
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		// since not being able to connect to our database is a fatal condition for us, we log then panic
		log.Fatalf("could not connect to  database: %v\n", err)
	}

	// set up our store
	s, err := newStore(db)
	if err != nil {
		log.Fatalf("failed to setup store: %v", err)
	}

	// set up our mux
	r := mux.NewRouter()

	// set up our paths and handlers
	// GET products
	r.HandleFunc("/v1/products", getProducts(s)).Methods(http.MethodGet)
	r.HandleFunc("/v1/products/{id:[0-9]+}", getProductByID(s)).Methods(http.MethodGet)
	r.HandleFunc("/v1/products/{id:[0-9]+}/decrement", decrementProductQuantityByID(s)).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":3000", r))
}
