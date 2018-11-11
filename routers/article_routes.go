package routers

import (
	"github.com/gorilla/mux"
	. "github.com/palanisn/articles-api/controllers"
)

func InitializeArticleRoutes(router *mux.Router) {
	router.HandleFunc("/articles", CreateArticle).Methods("POST")
	router.HandleFunc("/articles/{id:[0-9]+}", GetArticle).Methods("GET")
}
