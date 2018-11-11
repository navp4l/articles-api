package models

import (
	"fmt"
	. "github.com/palanisn/articles-api/database"
	"log"
)

type TagInfo struct {
	Tag         string   `json:"tag"`
	Count       int      `json:"count"`
	Articles    []int    `json:"articles"`
	RelatedTags []string `json:"related_tags"`
	Date        string   `json:"-"`
}

func (tagInfo *TagInfo) GetTagInfo() error {

	tagName := tagInfo.Tag
	date := tagInfo.Date

	log.Println("Date before", date)
	date = fmt.Sprintf("%s-%s-%s", string(date[:4]), string(date[4:6]), string(date[6:]))
	log.Println("Date after", date)

	countStatement := fmt.Sprintf("SELECT COUNT(*) FROM articles t1, tags t2, tagmap t3 WHERE t1.date='%s' "+
		"AND t1.id=t3.article_id AND t2.tag_id=t3.tag_id AND t2.name='%s'", date, tagName)

	articlesSelectStatement := fmt.Sprintf("SELECT t1.id FROM articles t1, tags t2, tagmap t3 WHERE t1.date='%s' "+
		"AND t3.article_id=t1.id AND t3.tag_id=t2.tag_id AND t2.name='%s' LIMIT 10", date, tagName)

	relatedTagsStatement := fmt.Sprintf("SELECT DISTINCT tg.name FROM articles ar, tags tg, tagmap tm WHERE tm.article_id IN "+
		"( SELECT t1.id FROM articles t1, tags t2, tagmap t3 WHERE t1.date='%s' AND t1.id=t3.article_id AND t2.tag_id=t3.tag_id AND t2.name='%s' ) "+
		"AND tm.tag_id=tg.tag_id;", date, tagName)

	log.Println("countStatement ", countStatement)
	countErr := DB.QueryRow(countStatement).Scan(&tagInfo.Count)

	log.Println("After count")

	if countErr != nil {
		log.Println("Inside Count Err")
		return countErr
	}

	log.Println("articlesSelectStatement ", articlesSelectStatement)
	articlesSelectResult, articlesSelectErr := DB.Query(articlesSelectStatement)

	if articlesSelectErr != nil {
		log.Println("Inside articlesSelectErr")
		return articlesSelectErr
	}

	articles := make([]int, 0)
	for articlesSelectResult.Next() {
		var article int
		err := articlesSelectResult.Scan(&article)
		if err != nil {
			return err
		}
		articles = append(articles, article)
	}
	tagInfo.Articles = articles

	log.Println("relatedTagsStatement ", relatedTagsStatement)
	relatedTagsResult, relatedTagsErr := DB.Query(relatedTagsStatement)

	if relatedTagsErr != nil {
		log.Println("Inside relatedTagsErr")
		return relatedTagsErr
	}

	tags := make([]string, 0)
	for relatedTagsResult.Next() {
		var tag string
		err := relatedTagsResult.Scan(&tag)
		if err != nil {
			return err
		}

		// Exclude input tag name from related tags
		if tag != tagName {
			tags = append(tags, tag)
		}
	}
	tagInfo.RelatedTags = tags

	return nil

}
