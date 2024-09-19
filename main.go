package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/sudofrost/markdown-note-taking-app/internal/db"
)

func main(){
	http.HandleFunc("POST /api/grammar-check", func(w http.ResponseWriter, r *http.Request) {
		content := r.FormValue("content")
		if content == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		options:= map[string][]string{
			"language": {"en-US"},
			"text": {content},
		}

		res, err := http.PostForm("https://api.languagetool.org/v2/check", options)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
    defer res.Body.Close()
		w.WriteHeader(res.StatusCode)
		io.Copy(w, res.Body)
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
		rows, err := db.DB.Query("SELECT * FROM notes")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		type Note struct {
			ID      int    `json:"id"`
			Title   string `json:"title"`
			Content string `json:"content"`
		}

		var notes = []Note{}

		for rows.Next() {
			var id int
			var title string
			var content string
			err := rows.Scan(&id, &title, &content)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			notes = append(notes, Note{
				ID:      id,
				Title:   title,
				Content: content,
			})
		}

		res, err := json.Marshal(notes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	})

  http.HandleFunc("GET /api/notes/{id}", func(w http.ResponseWriter, r *http.Request) {
		paramID := r.PathValue("id")
		row := db.DB.QueryRow("SELECT * FROM notes WHERE id = ?", paramID)

    if row == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		type Note struct {
			ID      int    `json:"id"`
			Title   string `json:"title"`
			Content string `json:"content"`
		}

		var id int
		var title string
		var content string
		err := row.Scan(&id, &title, &content)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		res, err := json.Marshal(Note{
			ID:      id,
			Title:   title,
			Content: content,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	})

	http.HandleFunc("GET /api/notes/{id}/render", func(w http.ResponseWriter, r *http.Request) {
		paramID := r.PathValue("id")
		row := db.DB.QueryRow("SELECT content FROM notes WHERE id = ?", paramID)

    if row == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		var content string
		err := row.Scan(&content)
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		p := parser.New()
		flags := html.CommonFlags | html.HrefTargetBlank
		rOPTS := html.RendererOptions{
			Flags: flags,
		}
		rr := html.NewRenderer(rOPTS)
		ad := markdown.ToHTML([]byte(content), p, rr)
		w.Write(ad)
	})

	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}
