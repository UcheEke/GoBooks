package main

import (
	"fmt"
	"net/http"
	"html/template"
	"log"

	"database/sql"
	_"github.com/mattn/go-sqlite3"
	"encoding/json"
)

type Page struct {
	Name string
	DBStatus bool
}

type SearchResult struct {
	Title string
	Author string
	Year string
	ID string
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

	// Define a handler for search results
	http.HandleFunc("/search", func(w http.ResponseWriter, req *http.Request){
		results := []SearchResult{
			SearchResult{"Moby-Dick", "Herman Melville", "1851", "222222"},
			SearchResult{"The Adventures of Huckleberry Finn", "Mark Twain", "1884", "333333"},
			SearchResult{"A Catcher in the Rye", "J D Salinger", "1951", "444444"},
		}

		encoder := json.NewEncoder(w)
		if err := encoder.Encode(results);err != nil {
			panic(err)
		}
	})

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	port := ":8080"

	fmt.Printf("Starting webserver on port %s\n", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
