package main

import (
	"fmt"
	"net/http"
	"html/template"
	"log"

	"database/sql"
	_"github.com/mattn/go-sqlite3"
)

type Page struct {
	Name string
	DBStatus bool
}

func main(){
	templates := template.Must(template.ParseFiles("templates/index.html"))

	db, _ := sql.Open("sqlite3", "db/dev.db")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		p := Page{Name: "Gopher"}
		// FormValue checks if there is a query string "name", and we'll update
		// the instance p(Page)'s name parameter with it if it's available
		if name := r.FormValue("name"); name != "" {
			p.Name = name
		}

		// Check whether the database is connected
		p.DBStatus = db.Ping() == nil

		// Execute template (can also define the template and run tmpl.Execute instead)
		if err := templates.ExecuteTemplate(w, "index.html", p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	port := ":8080"

	fmt.Printf("Starting webserver on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
