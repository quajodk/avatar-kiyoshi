package database

import (
	"database/sql"
	"log"
	"time"
)

var DB *sql.DB

func IntDB() {
	var error error
	connStr := "postgres://postgres.ypgghjvsfpxeoxgbfcak:_d4XQv-KLU*JiKF@aws-0-eu-west-1.pooler.supabase.com:5432/postgres?sslmode=verify-full"
	DB, error = sql.Open("postgres", connStr)
	if error != nil {
		log.Fatal(error)
	}
	DB.SetConnMaxLifetime(60 * time.Second)
	DB.SetMaxOpenConns(4)
	DB.SetMaxIdleConns(4)
}
