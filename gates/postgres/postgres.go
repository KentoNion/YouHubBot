package postgres

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/bool64/sqluct"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	db *sqlx.DB
	sq sq.StatementBuilderType
	sm sqluct.Mapper
}

func NewDB(db *sqlx.DB) *DB {
	return &DB{
		db: db,
		sm: sqluct.Mapper{Dialect: sqluct.DialectPostgres},
	}
}
