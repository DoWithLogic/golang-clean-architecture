package datasource

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/DoWithLogic/golang-clean-architecture/pkg/otel/zerolog"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/utils"
	"github.com/jmoiron/sqlx"
)

type (
	Conn interface {
		BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
		PingContext(ctx context.Context) (err error)
		io.Closer
		ConnTx
	}

	ConnTx interface {
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
		PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
		QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error)
		QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
	}

	Exec interface {
		Scan(rowsAffected, lastInsertID *int64) (err error)
	}

	Query interface {
		Scan(row func(i int) utils.Array) (err error)
	}

	exec struct {
		sqlResult sql.Result
		err       error
	}

	query struct {
		sqlRows *sqlx.Rows
		err     error
	}

	DataSource struct{}
)

var (
	_   Conn   = (*sqlx.Conn)(nil)
	_   Conn   = (*sqlx.DB)(nil)
	_   ConnTx = (*sqlx.Tx)(nil)
	log        = zerolog.NewZeroLog(context.Background(), os.Stdout)
)

// datasource errors
var (
	ErrNoColumnReturned   = errors.New("no columns returned")
	ErrDataNotFound       = errors.New("data not found")
	ErrInvalidArguments   = errors.New("invalid arguments for scan")
	ErrInvalidTransaction = errors.New("invalid transaction")
)

func (x exec) Scan(rowsAffected, lastInsertID *int64) error {
	if x.err != nil {
		log.Z().Err(x.err).Msg("[database:exec]error not nil")

		return x.err
	}

	if x.sqlResult == nil {
		log.Z().Err(sql.ErrNoRows).Msg("[database:exec]rows is nil")

		return ErrDataNotFound
	}

	if rowsAffected != nil {
		n, err := x.sqlResult.RowsAffected()
		if err != nil {
			log.Z().Err(err).Msg("[database:exec]scan rows affected error")

			return err
		}
		if n < 1 {
			log.Z().Err(ErrDataNotFound).Msg("[database:exec]rows affected")

			return ErrDataNotFound
		}
		*rowsAffected = int64(n)
	}

	if lastInsertID != nil {
		n, err := x.sqlResult.LastInsertId()
		if err != nil {
			log.Z().Err(err).Msg("[database:exec]last inserted id error")
		} else {
			*lastInsertID = int64(n)
		}
	}

	return nil
}

func (x query) Scan(row func(i int) utils.Array) error {
	if x.err != nil {
		log.Z().Err(x.err).Msg("[database:query]error not nil")

		return x.err
	}

	if x.sqlRows == nil {
		log.Z().Err(sql.ErrNoRows).Msg("[database:query]rows is nil")

		return ErrDataNotFound
	}

	if err := x.sqlRows.Err(); err != nil {
		return err
	}

	defer x.sqlRows.Close()

	columns, err := x.sqlRows.Columns()
	if err != nil {
		log.Z().Err(err).Msg("[database:query]columns")

		return err
	}

	if len(columns) < 1 {
		log.Z().Err(ErrNoColumnReturned).Msg("[database:query]count columns length")

		return ErrNoColumnReturned
	}

	var idx int = 0
	for x.sqlRows.Next() {
		if x.sqlRows.Err() != nil {
			log.Z().Err(x.sqlRows.Err()).Msg("[database:query]error to scan sql rows")

			return x.sqlRows.Err()
		}

		if row(idx) == nil {
			break
		}

		if len(row(idx)) < 1 {
			continue
		}

		if len(row(idx)) != len(columns) {
			err := fmt.Errorf("%w: [%d] columns on [%d] destinations", ErrInvalidArguments, len(columns), len(row(idx)))
			log.Z().Err(err).Msg("[database:query]error invalid args to scan")

			return err
		}

		if err = x.sqlRows.Scan(row(idx)...); err != nil {
			log.Z().Err(err).Msg("[database:query] failed to scan row")

			return err
		}

		idx++
	}

	return err
}

func (DataSource) ExecSQL(sqlResult sql.Result, err error) exec {
	return exec{sqlResult, err}
}

func (DataSource) QuerySQL(sqlRows *sqlx.Rows, err error) Query {
	return query{sqlRows, err}
}

func (DataSource) EndTx(tx *sqlx.Tx, err error) error {
	if tx == nil {
		log.Z().Err(ErrInvalidTransaction).Msg("[database:EndTx]")

		return ErrInvalidTransaction
	}

	// if any error occurred, we try to rollback
	if msg := "rollback"; err != nil {
		if errR := tx.Rollback(); errR != nil {
			msg = fmt.Sprintf("%s failed: (%s)", msg, errR.Error())
		}

		log.Z().Err(err).Msg(fmt.Sprintf("[database:EndTx]%s", msg))

		return err
	}

	// we try to commit here
	if err = tx.Commit(); err != nil {
		log.Z().Err(err).Msg("[database:EndTx]Commit")

		return err
	}

	return nil
}
