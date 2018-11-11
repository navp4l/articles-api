package routers

import "github.com/gorilla/mux"

func InitializeRoutes(router *mux.Router) {
	InitializeArticleRoutes(router)
	InitializeTagRoutes(router)
}
