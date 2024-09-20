package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
	"github.com/sudofrost/markdown-note-taking-app/internal/db"
	"gorm.io/gorm"
)

var BASE_DIR = os.Getenv("APP_BASE_DIR")

func main() {
	http.HandleFunc("POST /api/grammar-check", func(w http.ResponseWriter, r *http.Request) {
		content := r.FormValue("content")
		if content == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		options := map[string][]string{
			"language": {"en-US"},
			"text":     {content},
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
		title := r.FormValue("title")
		content := r.FormValue("content")
		categoryId := r.FormValue("category")
		if title == "" || content == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var category db.Category
		if categoryId != "" {
			err := db.DB.Where("id = ?", categoryId).First(&category).Error
			if err == gorm.ErrRecordNotFound {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		note := db.NewNote(title, content)
		note.CategoryID = category.ID
		err := db.DB.Create(&note).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	})

	http.HandleFunc("GET /api/notes", func(w http.ResponseWriter, r *http.Request) {
		var notes []db.Note
		query := db.DB
		if r.URL.Query().Has("category") {
			var categoryId int
			categoryId, err := strconv.Atoi(r.URL.Query().Get("category"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			query = query.Where("category_id = ?", categoryId)
		}

		if r.URL.Query().Has("q") {
			query = query.Where("title LIKE ?", "%"+r.URL.Query().Get("q")+"%")
		}

		err := query.Find(&notes).Error

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
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
		id, err := strconv.Atoi(paramID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var note *db.Note
		err = db.DB.First(&note, id).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if note == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		res, err := json.Marshal(*note)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	})

	http.HandleFunc("PUT /api/notes/{id}", func(w http.ResponseWriter, r *http.Request) {
		paramID := r.PathValue("id")
		id, err := strconv.Atoi(paramID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var note *db.Note
		err = db.DB.First(&note, id).Error
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if note == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		title := r.FormValue("title")
		content := r.FormValue("content")
		categoryId := r.FormValue("category")
		if title == "" || content == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var category db.Category
		if categoryId != "" {
			err := db.DB.Where("id = ?", categoryId).First(&category).Error
			if err == gorm.ErrRecordNotFound {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
		note.CategoryID = category.ID
		note.Title = title
		note.Content = content
		err = db.DB.Save(&note).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("DELETE /api/notes/{id}", func(w http.ResponseWriter, r *http.Request) {
		paramID := r.PathValue("id")
		id, err := strconv.Atoi(paramID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = db.DB.Delete(&db.Note{}, id).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("GET /api/notes/{id}/render", func(w http.ResponseWriter, r *http.Request) {
		paramID := r.PathValue("id")
		id, err := strconv.Atoi(paramID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var note *db.Note
		err = db.DB.First(&note, id).Error
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if note == nil {
			w.WriteHeader(http.StatusNotFound)
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
		ad := markdown.ToHTML([]byte(note.Content), p, rr)
		w.Write(ad)
	})

	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		view := template.New("")
		view.Funcs(template.FuncMap{
			"vite": ViteAssetURL,
		})

		view, err := view.ParseFiles(BASE_DIR + "/resources/views/home.html")
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		res := strings.Builder{}
		err = view.ExecuteTemplate(&res, "home.html", nil)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", res.Len()))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(res.String()))
	})

	http.HandleFunc("GET /api/categories", func(w http.ResponseWriter, r *http.Request) {
		var categories []db.Category
		err := db.DB.Find(&categories).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		res, err := json.Marshal(categories)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(res)
	})

	http.HandleFunc("POST /api/categories", func(w http.ResponseWriter, r *http.Request) {
		var name = r.FormValue("name")
		if name == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var category = db.NewCategory(name)
		err := db.DB.Create(&category).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("PUT /api/categories/{id}", func(w http.ResponseWriter, r *http.Request) {
		paramID := r.PathValue("id")
		id, err := strconv.Atoi(paramID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var category *db.Category
		err = db.DB.First(&category, id).Error
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if category == nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		category.Name = r.FormValue("name")
		if category.Name == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = db.DB.Save(&category).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("DELETE /api/categories/{id}", func(w http.ResponseWriter, r *http.Request) {
		paramID := r.PathValue("id")
		id, err := strconv.Atoi(paramID)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = db.DB.Delete(&db.Category{}, id).Error
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}

func ViteAssetURL(source string) string {
	source = strings.TrimPrefix(source, "/")
	return "http://localhost:5173/" + source
}
