
# Articles REST API

The *Articles* API provides the below endpoints,
* POST /articles - Create new articles
* GET /articles/{id} - Fetch & return article matching provided id
* GET /tags/{tagName}/{date} - Get tag related information for the date provided

## Application Design



## Database Design

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


