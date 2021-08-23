package main

import (
	"database/sql"
	"log"
	"net/http"
	"test_task1/handler"
	"test_task1/store"

	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {

	// Opening a driver typically will not attempt to connect to the database.
	dsn := "postgresql://postgres:@localhost:5432/news?sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		// This will not be a connection error, but a DSN parse error or
		// another initialization error.
		log.Fatal(err)
	}
	defer db.Close()

	tx, _ := db.Begin()
	defer tx.Commit()

	logger := logrus.New()

	s := store.New(tx)

	h := handler.New(*s, logger)

	mux := http.NewServeMux()

	mux.HandleFunc("/news", h.GetStartNewsPage)
	mux.HandleFunc("/fetch", h.GetFetchPage)
	mux.HandleFunc("/likes", h.GetLikes)

	logger.Infoln("server started on :8080 port")

	log.Println(http.ListenAndServe(":8080", mux))

}
