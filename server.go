package main

import (
	. "github.com/palanisn/articles-api/app"
)

func main() {
	api := &App{}
	api.InitializeApp("tester", "p@ssword12#", "articles_store")
	api.Start(":8080")
}
