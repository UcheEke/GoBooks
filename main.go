package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"database/sql"
	"encoding/json"
	"encoding/xml"
	_ "github.com/mattn/go-sqlite3"
	"io/ioutil"
	"net/url"
)

type Page struct {
	Name     string
	DBStatus bool
}

type SearchResult struct {
	Title  string `xml:"title,attr"`
	Author string `xml:"author,attr"`
	Year   string `xml:"hyr,attr"`
	ID     string `xml:"owi,attr"`
}

type ClassifySearchResponse struct {
	Results []SearchResult `xml:"works>work"`
}

type ClassifyBookResponse struct {
	BookData struct {
		Title  string `xml:"title,attr"`
		Author string `xml:"author,attr"`
		ID     string `xml:"owi,attr"`
	} `xml:"work"`
	Classification struct {
		MostPopular string `xml:"sfa,attr"`
	} `xml:"recommendations>ddc>mostPopular"`
}

func search(query string) ([]SearchResult, error) {
	// Unmarshal the received xml into a ClassifySearchResponse struct
	var c ClassifySearchResponse

	body, err := classifyAPI("http://classify.oclc.org/classify2/Classify?sumary=true&title=" + url.QueryEscape(query))
	if err != nil {
		return []SearchResult{}, err
	}

	if err = xml.Unmarshal(body, &c); err != nil {
		return []SearchResult{}, err
	} else {
		return c.Results, nil
	}
}

func findBook(id string) (ClassifyBookResponse, error) {
	var cbr ClassifyBookResponse
	var err error
	var body []byte

	if body, err = classifyAPI("http://classify.oclc.org/classify2/Classify?summary=true&owi=" + url.QueryEscape(id)); err != nil {
		return ClassifyBookResponse{}, err
	}

	err = xml.Unmarshal(body, &cbr)
	if err != nil {
		return ClassifyBookResponse{}, err
	} else {
		return cbr, nil
	}
}

func classifyAPI(url string) ([]byte, error) {
	var resp *http.Response
	var err error

	if resp, err = http.Get(url); err != nil {
		return []byte{}, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func main() {
	var err error

	// Ensure that template is employed using 'template.Must'
	templates := template.Must(template.ParseFiles("templates/index.html"))

	// Open a connection to the local sqlite database
	db, _ := sql.Open("sqlite3", "assets/db/dev.db")

	// Handle root connections
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := Page{Name: "Gopher"}
		// FormValue checks if there is a query string "name", and we'll update
		// the instance p(Page)'s name parameter with it if it's available
		if name := r.FormValue("name"); name != "" {
			p.Name = name
		}

		// Check whether the database is connected
		p.DBStatus = db.Ping() == nil
		fmt.Fprintf(w, `<p id="status">%s:%v</p>`, "Database connected", p.DBStatus)

		// Execute template (can also define the template and run tmpl.Execute instead)
		if err := templates.ExecuteTemplate(w, "index.html", p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Define a handler for search results
	http.HandleFunc("/search", func(w http.ResponseWriter, req *http.Request) {
		var results []SearchResult
		var err error

		// The data from the form element with 'name' attribute = "search" is parsed and used as the search query string
		if results, err = search(req.FormValue("search")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		encoder := json.NewEncoder(w)
		if err = encoder.Encode(results); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Handle the addition of books to the database
	http.HandleFunc("/books/add", func(w http.ResponseWriter, req *http.Request) {
		var book ClassifyBookResponse
		var err error

		if book, err = findBook(req.FormValue("id")); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if err = db.Ping(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// Use the db.Exec function to execute an SQL command on the database
		_, err = db.Exec("insert into books(pk,title,author,id,classification) values (?,?,?,?,?)",
			nil, book.BookData.Title, book.BookData.Author, book.BookData.ID,
			book.Classification.MostPopular)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Handle the static files (*.js,*.css, images etc)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	// Start the listening on the user specified port
	port := ":8080"
	fmt.Printf("Starting webserver on port %s\n", port)
	if err = http.ListenAndServe(port, nil); err != nil {
		log.Fatal(err)
	}
}
