package storage

import (
	"YoutHubBot/domain"
	"context"
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/bool64/sqluct"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
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

// Метдоды Admins--------------------------------------------------------------------------------------------------------
// Пытается получить админа по userID
func (p *DB) GetAdmin(ctx context.Context, userID domain.UserID) (domain.Admin, error) {
	const op = "gates.postgres.GetAdmin"
	query := p.sm.Select(p.sq.Select(), &admin{}).
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

// Методы источников----------------------------------------------------------------------------------------------------
// пытается добавить источник фида по структуре Source
func (p *DB) AddTGSource(ctx context.Context, source source) error {
	const op = "gates.postgres.AddTGSource"
	query := p.sm.Insert(p.sq.Insert("sources"), source, sqluct.InsertIgnore)
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

// Пытается удалить источник по имени канала
func (p *DB) DeleteTGSource(ctx context.Context, name string) error {
	const op = "gates.postgres.DeleteTGSource"
	query := p.sq.Delete("sources").Where(sq.Eq{"tg_chan_name": name})
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
func (p *DB) GetSources(ctx context.Context) ([]source, error) {
	const op = "gates.postgres.GetSources"
	query := p.sm.Select(p.sq.Select(), &source{}).From("sources")
	qry, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("#{op}: #{err}")
	}
	var sources []source
	err = p.db.SelectContext(ctx, &sources, qry, args...)
	if err != nil {
		return nil, fmt.Errorf("#{op}: #{err}")
	}
	return sources, nil
}

// Методы posts----------------------------------------------------------------------------------------------------------
// Сохраняет пост в бд
func (p *DB) StorePost(ctx context.Context, post post) error {
	const op = "gates.postgres.StorePost"
	query := p.sm.Insert(p.sq.Insert("posts"), post, sqluct.InsertIgnore)
	qry, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("#{op}: #{err}")
	}
	rows, err := p.db.ExecContext(ctx, qry, args...)
	if err != nil {
		return fmt.Errorf("#{op}: #{err}")
	}
	if rows, _ := rows.RowsAffected(); rows == 0 {
		return errors.New("#{op}: affected 0 rows")
	}
	return nil
}

// возвращает все статьи которые ещё не были опубликованны начиная с определённого времени
func (p *DB) AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]domain.Post, error) {

}

// Отмечает запощенную статью как запощенную
func (p *DB) MarkPosted(ctx context.Context, id int) error {

}
