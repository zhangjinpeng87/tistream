package utils

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// DBPool is database connection pool
type DBPool struct {
	*sql.DB
}

// NewDBPool creates a new database connection pool.
func NewDBPool(addr string, port int, user, password, dbName string) (*DBPool, error) {
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, addr, port, dbName)
	db, err := sql.Open("mysql", dbUrl)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)

	return &DBPool{
		db,
	}, nil
}
