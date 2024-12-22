package postgres

import (
	"YoutHubBot/domain"
	"context"
	"database/sql"
	"fmt"
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
		sq: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Пытается получить админа по userID
func (p *DB) GetAdmin(ctx context.Context, userID domain.UserID) (domain.Admin, error) {
	const op = "gates.postgres.GetAdmin"
	query := p.sm.Select(p.sq.Select(), &Admin{}).
		From("admins").
		Where(sq.Eq{"user_id": userID})
	qry, args, err := query.ToSql()
	if err == sql.ErrNoRows {
		return domain.Admin{}, domain.ErrNotAdmin
	}
	if err != nil {
		return domain.Admin{}, fmt.Errorf("#{op}: #{err}")
	}
	var admin domain.Admin
	err = p.db.GetContext(ctx, &admin, qry, args...)
	if err != nil {
		return domain.Admin{}, fmt.Errorf("#{op}: #{admin}")
	}
	return admin, nil
}

// пытается добавить источник фида по структуре Source
func (p *DB) AddTGSource(ctx context.Context, source Source) error {
	const op = "gates.postgres.AddTGSource"
	query := p.sm.Insert(p.sq.Insert("tg_sources"), source, sqluct.InsertIgnore)
	qry, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("#{op}: #{err}")
	}
	rows, err := p.db.ExecContext(ctx, qry, args...)
	if err != nil {
		return fmt.Errorf("#{op}: #{err}")
	}
	if rows, _ := rows.RowsAffected(); rows == 0 {
		return domain.ErrSourceAlreadyExist
	}
	return nil
}

// Пытается удалить источник по его имени
func (p *DB) DeleteTGSource(ctx context.Context, name string) error {
	const op = "gates.postgres.DeleteTGSource"
	query := p.sq.Delete("tg_sources").Where(sq.Eq{"name": name})
	qry, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("#{op}: #{err}")
	}
	rows, err := p.db.ExecContext(ctx, qry, args...)
	if err != nil {
		return fmt.Errorf("#{op}: #{err}")
	}
	if rows, _ := rows.RowsAffected(); rows == 0 {
		return fmt.Errorf("#{op}: affected 0 rows")
	}
	return nil
}

// Отдаёт список источников сохранённых в бд
func (p *DB) GetSources(ctx context.Context) ([]Source, error) {
	const op = "gates.postgres.GetSources"
	query := p.sm.Select(p.sq.Select(), &Source{}).From("tg_sources")
	qry, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("#{op}: #{err}")
	}
	var sources []Source
	err = p.db.SelectContext(ctx, &sources, qry, args...)
	if err != nil {
		return nil, fmt.Errorf("#{op}: #{err}")
	}
	return sources, nil
}
