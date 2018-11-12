package models

import (
	"fmt"
	. "github.com/palanisn/articles-api/database"
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

	// Insert into articles table
	result, insertArticleErr := txn.Exec(insertArticle)

	if insertArticleErr != nil {
		return insertArticleErr
	}

	articleId, lastIdErr := result.LastInsertId()

	article.ID = int(articleId)

	if lastIdErr != nil {
		return lastIdErr
	}

	for _, tag := range article.Tags {

		// Check if tags exist in tags table
		countStatement := fmt.Sprintf("SELECT COUNT(*) FROM tags WHERE name='%s'", tag)

		var countResult int
		countErr := DB.QueryRow(countStatement).Scan(&countResult)

		if countErr != nil {
			return countErr
		}

		var tagId int64
		var tagIdLastErr error

		// If tag exists then in tags table then retrieve tag id
		if countResult > 0 {
			tagsSelectStatement := fmt.Sprintf("SELECT tag_id FROM tags WHERE name='%s'", tag)

			tagIdErr := DB.QueryRow(tagsSelectStatement).Scan(&tagId)

			if tagIdErr != nil {
				return tagIdErr
			}

		} else { // If tag doesn't exist then insert into tags table
			tagsInsertStatement := fmt.Sprintf("INSERT INTO tags(name) VALUES('%s')", tag)

			result, err := txn.Exec(tagsInsertStatement)

			if err != nil {
				return err
			}

			tagId, tagIdLastErr = result.LastInsertId()

			if tagIdLastErr != nil {
				return tagIdLastErr
			}

		}

		// Insert tag and article mapping in tag map table
		tagMapStatement := fmt.Sprintf("INSERT INTO tagmap(article_id, tag_id) VALUES('%v','%v')", articleId, tagId)

		_, tagMapErr := txn.Exec(tagMapStatement)

		if tagMapErr != nil {
			return tagMapErr
		}
	}

	// Commit the transaction.
	return txn.Commit()

}

func (article *Article) GetArticle() error {

	// Retrieve article from articles table
	selectArticleStatement := fmt.Sprintf("SELECT title, body, date FROM articles WHERE Id=%d", article.ID)
	selectArticleErr := DB.QueryRow(selectArticleStatement).Scan(&article.Title, &article.Body, &article.Date)

	if selectArticleErr != nil {
		return selectArticleErr
	}

	// Retrieve tags mapped to article
	selectTagStatement := fmt.Sprintf("SELECT t1.name FROM tags t1 RIGHT JOIN tagmap t2 ON t1.tag_id=t2.tag_id WHERE t2.article_id=%d", article.ID)
	selectTagResult, selectTagError := DB.Query(selectTagStatement)

	if selectTagError != nil {
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
