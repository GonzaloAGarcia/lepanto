package main

import (
	"html/template"
	"log"
	"net/http"
	"path/filepath"
)

func home(w http.ResponseWriter, r *http.Request) {

	//w.Write([]byte("Look upon my works ye mighty."))

	ts, err := template.ParseFiles("./home.html")
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
	path := http.Dir(filepath.FromSlash("/mnt/c/Users/Gonza/Documents/Books/"))
	fs := http.StripPrefix("/books/", http.FileServer(path))

	return fs
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)

	mux.Handle("/books/", books())

	//mux.HandleFunc("/url", file)

	log.Println("Starting server on Port :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)

}
