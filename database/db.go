package database

import (
	"database/sql"
	"fmt"
)

var (
	DB *sql.DB
)

func InitializeDB(uname, pwd, dbname string) error {
	connectionString := fmt.Sprintf("%s:%s@/%s", uname, pwd, dbname)

	var err error
	DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		return err
	}

	return nil
}
