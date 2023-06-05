package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type book struct {
	Title string
	Image string
}

func main() {
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.Dir("./books/images"))
	mux.HandleFunc("/", home)
	mux.HandleFunc("/books2/", books2)
	mux.Handle("/books2/books/images/", http.StripPrefix("/books2/books/images/", fileServer))
	mux.Handle("/books/", books())

	//mux.HandleFunc("/url", file)

	log.Println("Starting server on Port :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}

func home(w http.ResponseWriter, r *http.Request) {

	//w.Write([]byte("Look upon my works ye mighty."))

	ts, err := template.ParseFiles("./static/home.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

}

func books() http.Handler {

	// FromSlash makes the path platform-independent
	path := http.Dir(filepath.FromSlash("./books/"))
	fs := http.StripPrefix("/books/", http.FileServer(path))

	return fs
}

func books2(w http.ResponseWriter, r *http.Request) {

	files, err := os.ReadDir("./books/")
	if err != nil {
		log.Fatal(err)
	}

	books := make([]book, 0)

	for _, file := range files {
		if !file.IsDir() {
			book := book{
				Title: file.Name(),
				Image: filepath.Join("./books/images/", strings.ReplaceAll(file.Name(), ".epub", ".jpg")),
			}

			books = append(books, book)
		}
	}

	ts := template.Must(template.ParseFiles("static/books.html"))
	err = ts.Execute(w, books)
	if err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
