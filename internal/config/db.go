package config

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySQLConnection(user, password, host, port, dbName string) (*sql.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, dbName)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// sql.Open belum benar-benar connect, jadi kita tes dulu
	if err := db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	return db, nil
}
