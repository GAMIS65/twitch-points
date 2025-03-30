package db

import "github.com/jackc/pgx/v5/pgxpool"

type DBStore struct {
	*Queries
	connPool *pgxpool.Pool
}

func NewStore(connPool *pgxpool.Pool) *DBStore {
	return &DBStore{
		connPool: connPool,
		Queries:  New(connPool),
	}
}
