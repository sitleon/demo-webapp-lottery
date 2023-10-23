package repos

import "database/sql"

type DbConn struct {
	*sql.DB
	Querier *Queries
}

func NewDbConn(db *sql.DB) *DbConn {
	return &DbConn{DB: db, Querier: New(db)}
}
