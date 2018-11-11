package models

import (
	"fmt"
	. "github.com/palanisn/articles-api/database"
)

type Article struct {
	ID    int    `json:"Id"`
	Title string `json:"Title"`
	Date  string `json:"Date"`
	Body  string `json:"Body"`
	//Tags  []string `json:"Tags"`
}

func (article *Article) CreateArticle() error {

	statement := fmt.Sprintf("INSERT INTO articles(title, date, body) VALUES('%s', '%s', '%s')", article.Title, article.Date, article.Body)
	_, err := DB.Exec(statement)

	if err != nil {
		return err
	}

	err = DB.QueryRow("SELECT LAST_INSERT_ID()").Scan(&article.ID)

	if err != nil {
		return err
	}

	return nil
}

func (article *Article) GetArticle() error {
	statement := fmt.Sprintf("SELECT Title, Body, Date FROM articles WHERE Id=%d", article.ID)
	return DB.QueryRow(statement).Scan(&article.Title, &article.Date, &article.Body)
}
