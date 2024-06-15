package db

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/enolgor/golang-webservice-template/models"
	"github.com/stephenafamo/scan"
	"github.com/stephenafamo/scan/stdscan"
)

type DB struct {
	conn *sql.DB
}

func (db *DB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.conn.QueryContext(ctx, query, args...)
}

func (db *DB) CreateUser(ctx context.Context, user *models.User) error {
	query := sq.Insert("users").Columns("id", "email", "roles").Values(user.ID, user.Email, user.Roles)
	sql, args, err := query.ToSql()
	if err != nil {
		return err
	}
	_, err = db.conn.ExecContext(ctx, sql, args...)
	return err
}

func (db *DB) ReadUser(ctx context.Context, id string) (*models.User, error) {
	query := sq.Select("*").From("users").Where(sq.Eq{"id": id}).Limit(1)
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}
	return stdscan.One(ctx, db, scan.StructMapper[*models.User](), sql, args...)
}
