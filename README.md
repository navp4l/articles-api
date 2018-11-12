
# Articles REST API

The *Articles* API provides the below endpoints,
* POST /articles - Create new articles
* GET /articles/{id} - Fetch & return article matching provided id
* GET /tags/{tagName}/{date} - Get tag related information for the date provided

## Application Design

This application is designed as a simple application that exposes REST endpoints to be consumed
by clients,

![Application Design](img/appDesign.png)

The application is organized as multiple packages,

![Application Structure](img/appStruct.png)

## Database Design

The database for the application has been modeled as 3 tables,

*   articles
*   tagmap
*   tags

![Database Design](img/dbDesign.png)

## Programming Language - **Go**

## Tools Used
* MySQL DB
* IntelliJ GoLand IDE
* Mac OS
* Postman
* Curl

In addition to the standard packages from the Go library Open source libraries were used,
* [Gorilla Mux](https://github.com/gorilla/mux)- Application routing implementation

## Setup instructions

### Pre-requisites
* Go environment
* MySQL DB

### Database setup scripts
```sql

# Setup DB schema

CREATE DATABASE articles_store;
USE articles_store;

# Setup tables in schema

CREATE TABLE articles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    title VARCHAR(50) NOT NULL,
    date DATE NOT NULL,
    body VARCHAR(500) NOT NULL,
    constraint articles_id_uindex
		unique (id)
);

CREATE TABLE tags (
    tag_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    constraint tagmap_tags_tag_id_fk
		unique (tag_id),
	  constraint tags_tag_id_uindex
		unique (tag_id),
	  constraint tags_name_uindex
		unique (name)
);

CREATE TABLE tagmap (
    id INT AUTO_INCREMENT PRIMARY KEY,
    article_id INT NOT NULL,
    tag_id INT NOT NULL,
    CONSTRAINT tagmap_articles_id_fk
		foreign key (article_id) references articles (id),
    constraint tagmap_tags_tag_id_fk
		foreign key (tag_id) references tags (tag_id)
);


# Drop table

DROP TABLE articles_store.tagmap;

DROP TABLE articles_store.articles;

DROP TABLE articles_store.tags;

```

### Step by step guide
* Clone the repository into local workspace 
`git clone `
* Change to project directory
`cd articles-api`
* Build the application 
`go build`
* Install the application
`go install`
* Run the application by passing in appropriate command line flags
`articles-api -dbUserName=tester -dbUserPwd=testing -dbName=articles_store -port=8080`

If the app has started successfully, you should see a msg for server listening on port on the console / std.out
*Server is up and running, listening to requests on port 8080 ..*

### Interacting with the application

You can interact with the app using cURL commands.

* Create an article
```
curl -H "Content-Type: application/json" \
    -H "Accept: application/json"\
    -X POST \
   --data '{
 	"Title" : "fifth now this again with tags another lorem",
 	"Date" : "2018-11-12",
 	"Body" : "fourth now this again with tags another lorem lipsum, some text, potentially containing simple markup about how potato chips are great",
 	"Tags" : ["global", "home", "well-being"]
 }' \
   http://localhost:8080/articles
```
   
* Fetch article
```
curl -H "Accept: application/json" -H "Content-Type: application/json" -X GET http://localhost:8080/articles/2
```

* Fetch tag info
```
curl -H "Accept: application/json" -H "Content-Type: application/json" -X GET http://localhost:8080/tags/home/20181112
```

Alternatively, you can import the setup files into [Postman](https://www.getpostman.com/) and test the endpoints from there.
The import file is linked at [https://www.getpostman.com/collections/51008fcd78ff0287f853](https://www.getpostman.com/collections/51008fcd78ff0287f853)
 



