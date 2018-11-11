package main

import (
	"flag"
	"fmt"
	. "github.com/palanisn/articles-api/app"
	"log"
)

func main() {
	api := &App{}

	dbUserName := flag.String("dbUserName", "test", "DB user name")
	dbUserPwd := flag.String("dbUserPwd", "test", "DB user pwd")
	dbName := flag.String("dbName", "articles_store", "Name of DB schema")
	port := flag.Int("port", 8080, "Server listening on port")

	flag.Parse()

	api.InitializeApp(*dbUserName, *dbUserPwd, *dbName)

	addrsString := fmt.Sprintf(":%d", *port)

	log.Println(fmt.Sprintf("Server is up and running, listening to requests on port %d ..", *port))

	api.Start(addrsString)
}
