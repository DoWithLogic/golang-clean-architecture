package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"

	"github.com/DoWithLogic/golang-clean-architecture/pkg/custom"
)

type (
	BeginTx interface {
		BeginTx(ctx context.Context, opts *sql.TxOptions) (tx *sql.Tx, err error)
	}
	ExecContext interface {
		ExecContext(ctx context.Context, query string, args ...interface{}) (res sql.Result, err error)
	}
	PingContext interface {
		PingContext(ctx context.Context) (err error)
	}
	PrepareContext interface {
		PrepareContext(ctx context.Context, query string) (stmt *sql.Stmt, err error)
	}
	QueryContext interface {
		QueryContext(ctx context.Context, query string, args ...interface{}) (rows *sql.Rows, err error)
	}
	QueryRowContext interface {
		QueryRowContext(ctx context.Context, query string, args ...interface{}) (row *sql.Row)
	}

	BoxExec interface {
		Scan(rowsAffected, lastInsertID *int64) (err error)
	}

	BoxQuery interface {
		// Scan accept do, a func that accept `i int` as index and returns a List
		// of pointer.
		//  List == nil   // break the loop
		//  len(List) < 1 // skip the current loop
		//  len(List) > 0 // assign the pointer, must be same as the length of columns
		Scan(row func(i int) custom.Array) (err error)
	}

	exec struct {
		sqlResult sql.Result
		err       error
	}

	query struct {
		sqlRows *sql.Rows
		err     error
	}

	SQLConn interface {
		BeginTx
		io.Closer
		PingContext
		SQLTxConn
	}

	SQLTxConn interface {
		ExecContext
		PrepareContext
		QueryContext
		QueryRowContext
	}

	SQL struct{}
)

var (
	_ SQLConn   = (*sql.Conn)(nil)
	_ SQLConn   = (*sql.DB)(nil)
	_ SQLTxConn = (*sql.Tx)(nil)
)

func (x exec) Scan(rowsAffected, lastInsertID *int64) error {
	if x.err != nil {
		return fmt.Errorf("database: BoxExec 1: %w", x.err)
	}

	if x.sqlResult == nil {
		return fmt.Errorf("database: BoxExec2 : %w", errors.New("invalid arguments for scan"))
	}

	if rowsAffected != nil {
		n, err := x.sqlResult.RowsAffected()
		if err != nil {
			return fmt.Errorf("database: BoxExec3: %w", err)
		}
		if n < 1 {
			return fmt.Errorf("database: BoxExec4: %w", sql.ErrNoRows)
		}
		*rowsAffected = int64(n)
	}

	if lastInsertID != nil {
		n, err := x.sqlResult.LastInsertId()
		if err != nil {
			// Print the error to see why lastInsertID failed to scan
			fmt.Println("Error scanning lastInsertID:", err)
		} else {
			*lastInsertID = int64(n)
		}
	}

	return nil
}

func (x query) Scan(row func(i int) custom.Array) error {
	if x.err != nil {
		return x.err
	}

	if x.sqlRows == nil {
		return fmt.Errorf("database: query: %w", sql.ErrNoRows)
	}

	if err := x.sqlRows.Err(); err != nil {
		return err
	}

	defer x.sqlRows.Close()

	columns, err := x.sqlRows.Columns()
	if err != nil {
		return fmt.Errorf("database: query: %w", err)
	}

	if len(columns) < 1 {
		return fmt.Errorf("database: query: %w", errors.New("no columns returned"))
	}

	var idx int = 0
	for x.sqlRows.Next() {
		if x.sqlRows.Err() != nil {
			return fmt.Errorf("database: query: %w", err)
		}

		if row(idx) == nil {
			break
		}

		if len(row(idx)) < 1 {
			continue
		}

		if len(row(idx)) != len(columns) {
			return fmt.Errorf("database: query: %w: [%d] columns on [%d] destinations", errors.New("invalid arguments for scan"), len(columns), len(row(idx)))
		}

		if err = x.sqlRows.Scan(row(idx)...); err != nil {
			return fmt.Errorf("database: query: %w", err)
		}

		idx++
	}

	return err
}

func (SQL) Exec(sqlResult sql.Result, err error) BoxExec { return exec{sqlResult, err} }

func (SQL) Query(sqlRows *sql.Rows, err error) BoxQuery { return query{sqlRows, err} }

// EndTx will end transaction with provided *sql.Tx and error. The tx argument
// should be valid, and then will check the err, if any error occurred, will
// commencing the ROLLBACK else will COMMIT the transaction.
//
//	txc := XSQLTxConn(db) // shared between *sql.Tx, *sql.DB and *sql.Conn
//	if tx, err := db.BeginTx(ctx, nil); err == nil && tx != nil {
//	  defer func() { err = xsql.EndTx(tx, err) }()
//	  txc = tx
//	}
func (SQL) EndTx(tx *sql.Tx, err error) error {
	if tx == nil {
		return fmt.Errorf("database: %w", errors.New("invalid transaction"))
	}

	// if any error occurred, we try to rollback
	if msg := "rollback"; err != nil {
		if errR := tx.Rollback(); errR != nil {
			msg = fmt.Sprintf("%s failed: (%s)", msg, errR.Error())
		}

		return fmt.Errorf("database: %s because: %w", msg, err)
	}

	// we try to commit here
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("database: %w", err)
	}

	return nil
}
