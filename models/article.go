package models

import (
	"fmt"
	. "github.com/palanisn/articles-api/database"
	"log"
)

type Article struct {
	ID    int      `json:"id"`
	Title string   `json:"title"`
	Date  string   `json:"date"`
	Body  string   `json:"body"`
	Tags  []string `json:"tags"`
}

func (article *Article) CreateArticle() error {

	insertArticle := fmt.Sprintf("INSERT INTO articles(title, date, body) VALUES('%s', '%s', '%s')", article.Title, article.Date, article.Body)

	txn, err := DB.Begin()

	if err != nil {
		return err
	}

	defer func() {
		_ = txn.Rollback()
	}()

	result, insertArticleErr := txn.Exec(insertArticle)

	if insertArticleErr != nil {
		return insertArticleErr
	}

	articleId, lastIdErr := result.LastInsertId()

	article.ID = int(articleId)

	log.Println("Article Id is", articleId)

	if lastIdErr != nil {
		return lastIdErr
	}

	for index, tag := range article.Tags {

		countStatement := fmt.Sprintf("SELECT COUNT(*) FROM tags WHERE name='%s'", tag)

		fmt.Println("countStatement", countStatement)

		var countResult int
		countErr := DB.QueryRow(countStatement).Scan(&countResult)

		if countErr != nil {
			fmt.Println("Inside error")
			return countErr
		}

		fmt.Println("countResult", countResult)

		var tagId int64
		var tagIdLastErr error

		if countResult > 0 {
			tagsSelectStatement := fmt.Sprintf("SELECT tag_id FROM tags WHERE name='%s'", tag)

			log.Println("In tag statements count > 0 :: ", tagsSelectStatement)

			tagIdErr := DB.QueryRow(tagsSelectStatement).Scan(&tagId)

			if tagIdErr != nil {
				return tagIdErr
			}

		} else {
			tagsInsertStatement := fmt.Sprintf("INSERT INTO tags(name) VALUES('%s')", tag)

			log.Println("In tag statements count < 0 :: ", tagsInsertStatement)

			result, err := txn.Exec(tagsInsertStatement)

			if err != nil {
				return err
			}

			tagId, tagIdLastErr = result.LastInsertId()

			if tagIdLastErr != nil {
				return tagIdLastErr
			}

		}

		log.Println(fmt.Sprintf("Tag id for index %v is %v", index, tagId))

		tagMapStatement := fmt.Sprintf("INSERT INTO tagmap(article_id, tag_id) VALUES('%v','%v')", articleId, tagId)

		tagMapResult, tagMapErr := txn.Exec(tagMapStatement)

		if tagMapErr != nil {
			return tagMapErr
		}

		tagMapId, tagMapLastErr := tagMapResult.LastInsertId()
		if tagMapLastErr != nil {
			return tagMapLastErr
		}

		log.Println("Last Tag map id is", tagMapId)

	}

	// Commit the transaction.
	return txn.Commit()

}

func (article *Article) GetArticle() error {
	selectArticleStatement := fmt.Sprintf("SELECT title, body, date FROM articles WHERE Id=%d", article.ID)
	selectArticleErr := DB.QueryRow(selectArticleStatement).Scan(&article.Title, &article.Body, &article.Date)

	if selectArticleErr != nil {
		log.Println("Inside Get Article - selectArticleStatement")
		return selectArticleErr
	}

	selectTagStatement := fmt.Sprintf("SELECT t1.name FROM tags t1 RIGHT JOIN tagmap t2 ON t1.tag_id=t2.tag_id WHERE t2.article_id=%d", article.ID)
	selectTagResult, selectTagError := DB.Query(selectTagStatement)

	if selectTagError != nil {
		log.Println("Inside Get Article - selectArticleStatement")
		return selectTagError
	}

	tags := make([]string, 0)
	for selectTagResult.Next() {
		var tag string
		err := selectTagResult.Scan(&tag)
		if err != nil {
			return err
		}
		tags = append(tags, tag)
	}
	article.Tags = tags

	return nil
}
