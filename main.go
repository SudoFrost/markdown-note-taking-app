package main

import (
	"fmt"
	"net/http"

	"github.com/sudofrost/markdown-note-taking-app/internal/db"
)

func main(){
	http.HandleFunc("POST /api/grammar-check", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("POST /api/notes", func(w http.ResponseWriter, r *http.Request) {
		name := r.FormValue("name")
		content := r.FormValue("content")
		if name == "" || content == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		sql := "INSERT INTO notes (title, content) VALUES (?, ?)"
		_, err := db.Exec(sql, name, content)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})

	http.HandleFunc("GET /api/notes", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

  http.HandleFunc("GET /api/notes/{id}", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("GET /api/notes/{id}/render", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}
