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
	"log/slog"
	"time"
)

type DB struct {
	db  *sqlx.DB
	sq  sq.StatementBuilderType
	sm  sqluct.Mapper
	log *slog.Logger
}

func NewDB(db *sqlx.DB, log *slog.Logger) *DB {
	return &DB{
		db:  db,
		sm:  sqluct.Mapper{Dialect: sqluct.DialectPostgres},
		sq:  sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		log: log,
	}
}

// Метдоды Admins--------------------------------------------------------------------------------------------------------
// Пытается получить админа по userID
func (p *DB) GetAdmin(ctx context.Context, userID domain.UserID) (domain.Admin, error) {
	const op = "gates.postgres.GetAdmin"
	p.log.Debug(fmt.Sprintf("%s Trying to get admin with userID - %s", op, userID))
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
	p.log.Debug(fmt.Sprintf("%v: sucsess, found admin with userID %v", op, userID))
	return admin, nil
}

// Методы источников----------------------------------------------------------------------------------------------------
// пытается добавить источник фида по структуре Source
func (p *DB) AddTGSource(ctx context.Context, source source) error {
	const op = "gates.postgres.AddTGSource"
	p.log.Debug(fmt.Sprintf("%v: trying to add source %v", op, source))
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
	p.log.Debug(fmt.Sprintf("%v: success, source added", op))
	return nil
}

// Пытается удалить источник по имени канала
func (p *DB) DeleteTGSource(ctx context.Context, name string) error {
	const op = "gates.postgres.DeleteTGSource"
	p.log.Debug(fmt.Sprintf("%v: trying to delete source %v", op, name))
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
	p.log.Debug(fmt.Sprintf("%v: success, source deleted, source name: %v", op, name))
	return nil
}

// Отдаёт список источников сохранённых в бд
func (p *DB) GetSources(ctx context.Context) ([]source, error) {
	const op = "gates.postgres.GetSources"
	p.log.Debug(fmt.Sprintf("%v: trying to get all sources", op))
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
	p.log.Debug(fmt.Sprintf("%v: sucsess, all sources retrieved", op))
	return sources, nil
}

// Методы posts----------------------------------------------------------------------------------------------------------
// Сохраняет пост в бд
func (p *DB) StorePost(ctx context.Context, post post) error {
	const op = "gates.postgres.StorePost"
	p.log.Debug(fmt.Sprintf("%v: trying to store post %v", op, post))
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
	p.log.Debug(fmt.Sprintf("%v: success, post stored", op))
	return nil
}

// возвращает все статьи которые ещё не были опубликованны начиная с определённого времени
func (p *DB) AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]domain.Post, error) {
	const op = "gates.postgres.AllNotPosted"
	p.log.Debug(fmt.Sprintf("%v: trying to get all not posted posts"))
	query := p.sm.Select(p.sq.Select(), &post{}).From("posts").Where(sq.Eq{"posted_at": nil})
	qry, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("#{op}: #{err}")
	}
	var unposted []domain.Post
	err = p.db.SelectContext(ctx, &unposted, qry, args...)
	if err != nil {
		return nil, fmt.Errorf("#{op}: #{err}")
	}
	p.log.Debug(fmt.Sprintf("%v: sucsess, all not posted posts retrieved", op))
	return unposted, nil
}

// Отмечает запощенную статью как запощенную (проставляет time.Now в параметр posted_at)
func (p *DB) MarkPosted(ctx context.Context, postID int) error {
	const op = "gates.postgres.MarkPosted"
	p.log.Debug(fmt.Sprintf("%v: trying to mark posted post %v", op, postID))
	query := p.sm.Update("posts").Set("posted_ad", time.Now).Where(sq.Eq{"post_id": postID})
	qry, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("#{op}: #{err}")
	}
	_, err = p.db.ExecContext(ctx, qry, args...)
	if err != nil {
		return fmt.Errorf("#{op}: #{err}")
	}
	p.log.Debug(fmt.Sprintf("%v: success, posted post marked", op))
	return nil
}
