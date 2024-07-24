package main

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"strconv"
)

func rootHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	data := getArtists()
	if data == nil {
		t, err := template.ParseFiles("templates/500.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error loading 500 template:", err)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		t.Execute(w, nil)
		return
	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		t, err = template.ParseFiles("templates/500.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error loading 500 template:", err)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		t.Execute(w, nil)
	}
	t.Execute(w, data)
}

func artistHandler(w http.ResponseWriter, r *http.Request) {
	// Extracting the ID from the URL
	id := path.Base(r.URL.Path[len("/artists/"):])

	// if id is not a number, return nil
	if _, err := strconv.Atoi(id); err != nil {
		t, err := template.ParseFiles("templates/400.html")
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			log.Println("Error loading 400 template:", err)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		t.Execute(w, nil)
		return
	}

	// Fetch the artist by ID
	data := getArtistByID(id)
	if data == nil {
		t, err := template.ParseFiles("templates/404.html")
		if err != nil {
			http.NotFound(w, r)
			log.Println("Error loading 404 template:", err)
			return
		}
		w.WriteHeader(http.StatusNotFound)
		t.Execute(w, nil)
		return
	}

	t, err := template.ParseFiles("templates/artist.html")
	if err != nil {
		t, err = template.ParseFiles("templates/500.html")
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			log.Println("Error loading 500 template:", err)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		t.Execute(w, nil)
		return
	}
	t.Execute(w, data)
}
