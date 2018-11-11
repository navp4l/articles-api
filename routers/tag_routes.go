package routers

import (
	"github.com/gorilla/mux"
	. "github.com/palanisn/articles-api/controllers"
)

func InitializeTagRoutes(router *mux.Router) {
	router.HandleFunc("/tags/{tagName:[a-z]+}/{date:[0-9]+}", GetTagInfo).Methods("GET")
}
