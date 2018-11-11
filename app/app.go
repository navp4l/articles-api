package app

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	. "github.com/palanisn/articles-api/database"
	. "github.com/palanisn/articles-api/routers"
	"log"
	"net/http"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (app *App) InitializeApp(uname, pwd, dbname string) {
	err := InitializeDB(uname, pwd, dbname)
	if err != nil {
		log.Fatal(err)
	}
	app.DB = DB

	app.Router = mux.NewRouter()
	InitializeRoutes(app.Router)
}

func (app *App) Start(addr string) {
	log.Fatal(http.ListenAndServe(addr, app.Router))
}
