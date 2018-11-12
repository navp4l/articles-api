package main

import (
	"bytes"
	"encoding/json"
	. "github.com/palanisn/articles-api/app"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var app App

func TestMain(m *testing.M) {
	app = App{}
	app.InitializeApp("tester", "p@ssword12#", "articles_store")

	checkTableExists()

	returnCode := m.Run()

	emptyTables()

	os.Exit(returnCode)
}

func checkTableExists() {
	if _, err := app.DB.Exec(articlesDDLQuery); err != nil {
		log.Fatal(err)
	}

	if _, err := app.DB.Exec(tagsDDLQuery); err != nil {
		log.Fatal(err)
	}

	if _, err := app.DB.Exec(tagmapDDLQuery); err != nil {
		log.Fatal(err)
	}
}

func emptyTables() {
	app.DB.Exec("DELETE FROM tagmap")
	app.DB.Exec("ALTER TABLE tagmap AUTO_INCREMENT = 1")

	app.DB.Exec("DELETE FROM tags")
	app.DB.Exec("ALTER TABLE tags AUTO_INCREMENT = 1")

	app.DB.Exec("DELETE FROM articles")
	app.DB.Exec("ALTER TABLE articles AUTO_INCREMENT = 1")
}

const articlesDDLQuery = `
CREATE TABLE IF NOT EXISTS articles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    date DATE NOT NULL,
    body VARCHAR(500) NOT NULL,
    constraint articles_id_uindex
		unique (id)
)`

const tagsDDLQuery = `
CREATE TABLE IF NOT EXISTS tags (
    tag_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    constraint tagmap_tags_tag_id_fk
		unique (tag_id),
	  constraint tags_tag_id_uindex
		unique (tag_id),
	  constraint tags_name_uindex
		unique (name)
)`

const tagmapDDLQuery = `
CREATE TABLE IF NOT EXISTS tagmap (
    id INT AUTO_INCREMENT PRIMARY KEY,
    article_id INT NOT NULL,
    tag_id INT NOT NULL,
    CONSTRAINT tagmap_articles_id_fk
		foreign key (article_id) references articles (id),
    constraint tagmap_tags_tag_id_fk
		foreign key (tag_id) references tags (tag_id)
);`

func TestArticleNotAvailable(t *testing.T) {
	emptyTables()

	request, _ := http.NewRequest("GET", "/articles/1", nil)
	response := makeRequest(request)

	validateResponse(t, http.StatusNotFound, response.Code)

	if body := response.Body.String(); body != "{\"error\":\"Article not found\"}" {
		t.Errorf("Expected {\"error\":\"Article not found\"} . Got %s", body)
	}
}

func TestCreateUser(t *testing.T) {
	emptyTables()

	payload := []byte(`{
	"title" : "Test title",
	"date" : "2018-11-12",
	"body" : "Test body for the article",
	"tags" : ["health", "fitness", "gym"]
	}`)

	request, _ := http.NewRequest("POST", "/articles", bytes.NewBuffer(payload))
	response := makeRequest(request)

	validateResponse(t, http.StatusCreated, response.Code)

	var values map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &values)

	titleVal := values["title"]
	if titleVal != "Test title" {
		t.Errorf("Expected title to be 'Test title'. Got '%v'", titleVal)
	}

	dateVal := values["date"]
	if dateVal != "2018-11-12" {
		t.Errorf("Expected date to be '2018-11-12'. Got '%v'", dateVal)
	}

	bodyVal := values["body"]
	if bodyVal != "Test body for the article" {
		t.Errorf("Expected body to be 'Test body for the article'. Got '%v'", bodyVal)
	}

	tagVal := values["tags"]
	tagSliceVal := tagVal.([]interface{})
	if tagSliceVal[0] != "health" || tagSliceVal[1] != "fitness" || tagSliceVal[2] != "gym" {
		t.Errorf("Expected tags to be '[health fitness gym]'. Got '%v'", tagVal)
	}

	idVal := values["id"]
	if idVal != 1.0 {
		t.Errorf("Expected article ID to be '1'. Got '%v'", idVal)
	}
}

func TestArticleRetrieval(t *testing.T) {
	emptyTables()
	populateDB()

	request, _ := http.NewRequest("GET", "/articles/1", nil)
	response := makeRequest(request)

	var values map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &values)

	titleVal := values["title"]
	if titleVal != "Test title" {
		t.Errorf("Expected title to be 'Test title'. Got '%v'", titleVal)
	}

	dateVal := values["date"]
	if dateVal != "2018-11-12" {
		t.Errorf("Expected date to be '2018-11-12'. Got '%v'", dateVal)
	}

	bodyVal := values["body"]
	if bodyVal != "Test body for the article" {
		t.Errorf("Expected body to be 'Test body for the article'. Got '%v'", bodyVal)
	}

	tagVal := values["tags"]
	tagSliceVal := tagVal.([]interface{})
	if tagSliceVal[0] != "health" || tagSliceVal[1] != "fitness" || tagSliceVal[2] != "gym" {
		t.Errorf("Expected tags to be '[health fitness gym]'. Got '%v'", tagVal)
	}

	idVal := values["id"]
	if idVal != 1.0 {
		t.Errorf("Expected article ID to be '1'. Got '%v'", idVal)
	}

	validateResponse(t, http.StatusOK, response.Code)
}

func TestTagInfo(t *testing.T) {
	emptyTables()
	populateDB()

	request, _ := http.NewRequest("GET", "/tags/health/20181112", nil)
	response := makeRequest(request)

	var values map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &values)

	tagVal := values["tag"]
	if tagVal != "health" {
		t.Errorf("Expected tag to be 'health'. Got '%v'", tagVal)
	}

	countVal := values["count"]
	if countVal != 1.0 {
		t.Errorf("Expected article ID to be '1'. Got '%v'", countVal)
	}

	articlesVal := values["articles"]
	articlesSliceVal := articlesVal.([]interface{})
	if articlesSliceVal[0] != 1.0 {
		t.Errorf("Expected articles to be '[1]'. Got '%v'", articlesVal)
	}

	relTagVal := values["related_tags"]
	relTagSliceVal := relTagVal.([]interface{})
	if relTagSliceVal[0] != "fitness" || relTagSliceVal[1] != "gym" {
		t.Errorf("Expected related tags to be '[fitness gym]'. Got '%v'", relTagVal)
	}

	validateResponse(t, http.StatusOK, response.Code)
}

func populateDB() {
	// Populate tags
	tagsStatements := make([]string, 0)
	tagsStatements = append(tagsStatements, "INSERT INTO tags(name) VALUES('health')")   // 1 health
	tagsStatements = append(tagsStatements, "INSERT INTO tags(name) VALUES('fitness')")  // 2 fitness
	tagsStatements = append(tagsStatements, "INSERT INTO tags(name) VALUES('science')")  // 3 science
	tagsStatements = append(tagsStatements, "INSERT INTO tags(name) VALUES('wellness')") // 4 wellness
	tagsStatements = append(tagsStatements, "INSERT INTO tags(name) VALUES('exercise')") // 5 exercise
	tagsStatements = append(tagsStatements, "INSERT INTO tags(name) VALUES('gym')")      // 6 gym

	for _, statement := range tagsStatements {
		app.DB.Exec(statement)
	}

	// Populate articles
	articleStatements := make([]string, 0)
	articleStatements = append(articleStatements, "INSERT INTO articles(title, date, body) VALUES('Test title', '2018-11-12', 'Test body for the article')")
	articleStatements = append(articleStatements, "INSERT INTO articles(title, date, body) VALUES('Title 2', '2018-11-08', 'Test body 2')")
	articleStatements = append(articleStatements, "INSERT INTO articles(title, date, body) VALUES('Title 3', '2018-11-09', 'Test body 3')")
	articleStatements = append(articleStatements, "INSERT INTO articles(title, date, body) VALUES('Title 4', '2018-11-10', 'Test body 4')")
	articleStatements = append(articleStatements, "INSERT INTO articles(title, date, body) VALUES('Title 5', '2018-11-10', 'Test body 5')")

	for _, statement := range articleStatements {
		app.DB.Exec(statement)
	}

	// Populate tag map
	tagMapStatements := make([]string, 0)
	tagMapStatements = append(tagMapStatements, "INSERT INTO tagmap(article_id, tag_id) VALUES(1,1)")
	tagMapStatements = append(tagMapStatements, "INSERT INTO tagmap(article_id, tag_id) VALUES(1,2)")
	tagMapStatements = append(tagMapStatements, "INSERT INTO tagmap(article_id, tag_id) VALUES(1,6)")

	tagMapStatements = append(tagMapStatements, "INSERT INTO tagmap(article_id, tag_id) VALUES(2,3)")
	tagMapStatements = append(tagMapStatements, "INSERT INTO tagmap(article_id, tag_id) VALUES(2,2)")
	tagMapStatements = append(tagMapStatements, "INSERT INTO tagmap(article_id, tag_id) VALUES(2,5)")

	tagMapStatements = append(tagMapStatements, "INSERT INTO tagmap(article_id, tag_id) VALUES(3,5)")
	tagMapStatements = append(tagMapStatements, "INSERT INTO tagmap(article_id, tag_id) VALUES(3,3)")
	tagMapStatements = append(tagMapStatements, "INSERT INTO tagmap(article_id, tag_id) VALUES(3,6)")

	tagMapStatements = append(tagMapStatements, "INSERT INTO tagmap(article_id, tag_id) VALUES(4,2)")
	tagMapStatements = append(tagMapStatements, "INSERT INTO tagmap(article_id, tag_id) VALUES(4,5)")
	tagMapStatements = append(tagMapStatements, "INSERT INTO tagmap(article_id, tag_id) VALUES(4,6)")

	tagMapStatements = append(tagMapStatements, "INSERT INTO tagmap(article_id, tag_id) VALUES(5,4)")
	tagMapStatements = append(tagMapStatements, "INSERT INTO tagmap(article_id, tag_id) VALUES(5,6)")
	tagMapStatements = append(tagMapStatements, "INSERT INTO tagmap(article_id, tag_id) VALUES(5,2)")

	for _, statement := range tagMapStatements {
		app.DB.Exec(statement)
	}
}

func makeRequest(req *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	app.Router.ServeHTTP(recorder, req)

	return recorder
}

func validateResponse(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
