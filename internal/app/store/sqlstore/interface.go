package sqlstore

import "database/sql"

type StoreInterface interface {
	Currency() CurrencyInterface
}

type DBInterface interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
	Exec(query string, args ...interface{}) (sql.Result, error)
}
