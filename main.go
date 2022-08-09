package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
)

type SearchedBook struct {
	Docs [1]struct {
		Title        string `json:"title"`
		Isbn         [1]string
		Publish      int      `json:"first_publish_year"`
		EditCount    int      `json:"edition_count"`
		AuthorKeys   []string `json:"author_key"`
		AuthorsNames []string `json:"author_name"`
	}
}

type Abook struct {
	Title    string `json:"title"`
	Publish  string `json:"first_publish_date"`
	Revision int    `json:"revision"`
}

type AuthorsBooks struct {
	Entries []Abook
}

func TemplatingIndex() *template.Template {

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	templateRef := template.Must(template.ParseFiles(wd + "/public/index.html"))

	return templateRef

}

func main() {

	log.Println("Starting server...")

	http.HandleFunc("/search", ListBooks)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		err := TemplatingIndex().ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			log.Fatalln(err)
		}

	})

	log.Println("Server listening ... localhost port:3000")

	http.ListenAndServe(":3000", nil)

}

func ListBooks(w http.ResponseWriter, r *http.Request) {

	searchValue := r.FormValue("search")
	if searchValue != "" {

		searchedBookData := RequestBook(TrimingSearch(searchValue))

		authorsWorks := ReqWorks(searchedBookData.Docs[0].AuthorKeys)

		mapedData := SortingData(searchedBookData.Docs[0].AuthorsNames, authorsWorks)

		err := TemplatingIndex().ExecuteTemplate(w, "index.html", mapedData)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		err := TemplatingIndex().ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			log.Fatalln(err)
		}
	}
}

func TrimingSearch(s string) string {

	trimedString := strings.ReplaceAll(s, " ", "+")
	return trimedString
}

func RequestBook(b string) SearchedBook {

	apiSearch := fmt.Sprintf("http://openlibrary.org/search.json?title=%s", b)

	res, err := http.Get(apiSearch)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	var searchedBook SearchedBook
	_ = json.NewDecoder(res.Body).Decode(&searchedBook)

	return searchedBook

}

func ReqWorks(authorsslices []string) []AuthorsBooks {

	worksData := make([]AuthorsBooks, 0)

	for _, author := range authorsslices {

		respond, err := http.Get(fmt.Sprintf("https://openlibrary.org/authors/%s/works.json", author))
		if err != nil {
			log.Fatal(err)
		}
		defer respond.Body.Close()

		var authorsBook AuthorsBooks
		_ = json.NewDecoder(respond.Body).Decode(&authorsBook)

		worksData = append(worksData, authorsBook)

	}

	return worksData

}

func SortingData(authorsnames []string, authorsworks []AuthorsBooks) map[string]map[string]Abook {

	sortedWorkss := make(map[string]map[string]Abook)

	if authorsnames != nil {

		for i := 0; i < len(authorsnames); i++ {

			sortedWorkss[authorsnames[i]] = map[string]Abook{}

			for j := 0; j < len(authorsworks[i].Entries); j++ {

				sortedWorkss[authorsnames[i]][authorsworks[i].Entries[j].Title] = authorsworks[i].Entries[j]

			}

		}

	}

	return sortedWorkss
}
