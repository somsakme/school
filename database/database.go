package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func GetDBConn() (*sql.DB, error) {
	db, err := sql.Open("postgres", "postgres://kpplnkjj:owliREOn5WEm14dcPOipnHv33pUn6U_p@john.db.elephantsql.com:5432/kpplnkjj")
	if err != nil {
		return nil, err
	}
	return db, nil
}
