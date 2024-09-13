package main

import (
	"fmt"
	"net/http"
)

func main(){
	http.HandleFunc("POST /api/grammar-check", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("POST /api/notes", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
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
