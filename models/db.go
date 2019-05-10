package models

import (
		"database/sql"
		"fmt"
		_ "github.com/go-sql-driver/mysql"
		_ "github.com/golang-migrate/migrate/source/file"
		"github.com/joho/godotenv"
		"os"
)

var db *sql.DB

func init() {
		// Load env files and assign env vars
		e := godotenv.Load()
		if e != nil {
				fmt.Println(e.Error())
		}

		username := os.Getenv("db_user")
		password := os.Getenv("db_pass")
		dbName := os.Getenv("db_name")

		// Connect to db
		conn, err := sql.Open("mysql", username+":"+password+"@/"+dbName)
		if err != nil {
				fmt.Println(err.Error())
		}

		db = conn
		// TODO: CREATE DATABASE MIGRATION SYSTEM OR TRY TO DO THIS BY USING DOCKER

		//defer db.Close()

}

func GetDb() *sql.DB {
		return db
}
