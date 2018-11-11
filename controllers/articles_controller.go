package controllers

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"

	. "github.com/palanisn/articles-api/models"
)

func CreateArticle(w http.ResponseWriter, r *http.Request) {
	var art Article
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&art); err != nil {
		log.Print(err)
		handleError(w, http.StatusBadRequest, "Invalid request")
		return
	}

	defer r.Body.Close()

	if err := art.CreateArticle(); err != nil {
		handleError(w, http.StatusInternalServerError, err.Error())
		return
	}

	handleRespAsJSON(w, http.StatusCreated, art)
}

func GetArticle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])

	if err != nil {
		handleError(w, http.StatusBadRequest, "Invalid article ID")
		return
	}

	article := Article{ID: id}

	if err := article.GetArticle(); err != nil {
		switch err {

		case sql.ErrNoRows:
			handleError(w, http.StatusNotFound, "Article not found")

		default:
			handleError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	handleRespAsJSON(w, http.StatusOK, article)
}

func handleError(w http.ResponseWriter, code int, message string) {
	handleRespAsJSON(w, code, map[string]string{"error": message})
}

func handleRespAsJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
