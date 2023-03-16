package dotdb

import (
	"database/sql"
	"fmt"

	"github.com/blockloop/scan"
	"github.com/jmoiron/sqlx"
	"github.com/qustavo/dotsql"
)

type DotDB struct {
	db     *sqlx.DB
	dotSql *dotsql.DotSql
}

func NewDotDB(db *sqlx.DB, queries string) (*DotDB, error) {
	dot, err := dotsql.LoadFromString(queries)
	if err != nil {
		return nil, err
	}

	return &DotDB{
		db:     db,
		dotSql: dot,
	}, nil
}

func (db *DotDB) TxBegin() (*sqlx.Tx, error) {
	return db.db.Beginx()
}

// rows fetching with ? params

func (db *DotDB) GetStructsSlice(name string, v interface{}, args ...interface{}) error {
	rows, err := db.getRows(db.db, name, args)
	if err != nil {
		return fmt.Errorf("error running query `%s`: %w", name, err)
	}

	err = scan.RowsStrict(v, rows)
	if err != nil {
		return fmt.Errorf("error scanning query `%s` result: %w", name, err)
	}

	return nil
}
func (db *DotDB) GetRows(name string, args ...interface{}) (*sqlx.Rows, error) {
	rows, err := db.getRows(db.db, name, args)
	if err != nil {
		return nil, fmt.Errorf("error running query `%s`: %w", name, err)
	}
	return rows, err
}
func (db *DotDB) TxGetRows(tx *sqlx.Tx, name string, args ...interface{}) (*sqlx.Rows, error) {
	rows, err := db.getRows(tx, name, args)
	if err != nil {
		return nil, fmt.Errorf("error running query `%s`: %w", name, err)
	}
	return rows, err
}
func (db *DotDB) getRows(querier sqlx.Queryer, name string, args []interface{}) (*sqlx.Rows, error) {
	query, args, err := db.getQueryAndArgsWithIn(name, args)
	if err != nil {
		return nil, err
	}
	return querier.Queryx(query, args...)
}

// rows fetching with named params

func (db *DotDB) NamedGetRows(name string, args interface{}) (*sqlx.Rows, error) {
	return db.namedGetRows(db.db, name, args)
}

func (db *DotDB) TxNamedGetRows(tx *sqlx.Tx, name string, args interface{}) (*sqlx.Rows, error) {
	return db.namedGetRows(tx, name, args)
}
func (db *DotDB) namedGetRows(queryer sqlx.Queryer, name string, arg interface{}) (*sqlx.Rows, error) {
	query, args, err := db.getNamedQueryAndArgsWithIn(name, arg)
	if err != nil {
		return nil, err
	}

	return queryer.Queryx(query, args...)
}

// row fetching

func (db *DotDB) GetRow(name string, args ...interface{}) (*sqlx.Row, error) {
	return db.getRow(db.db, name, args)
}
func (db *DotDB) TxGetRow(tx *sqlx.Tx, name string, args ...interface{}) (*sqlx.Row, error) {
	return db.getRow(tx, name, args)
}
func (db *DotDB) getRow(queryer sqlx.Queryer, name string, args []interface{}) (*sqlx.Row, error) {
	query, args, err := db.getQueryAndArgsWithIn(name, args)
	if err != nil {
		return nil, err
	}
	return queryer.QueryRowx(query, args...), nil
}

// row fetching with named parameters

func (db *DotDB) NamedGetRow(name string, arg interface{}) (*sqlx.Row, error) {
	return db.namedGetRow(db.db, name, arg)
}
func (db *DotDB) TxNamedGetRow(tx *sqlx.Tx, name string, arg interface{}) (*sqlx.Row, error) {
	return db.namedGetRow(tx, name, arg)
}
func (db *DotDB) namedGetRow(queryer sqlx.Queryer, name string, arg interface{}) (*sqlx.Row, error) {
	query, args, err := db.getNamedQueryAndArgsWithIn(name, arg)
	if err != nil {
		return nil, err
	}
	return queryer.QueryRowx(query, args...), nil
}

// exec query

func (db *DotDB) ShouldExec(name string, args ...interface{}) sql.Result {
	if res, err := db.exec(db.db, name, args); err != nil {
		panic(err.Error())
	} else {
		return res
	}
}
func (db *DotDB) Exec(name string, args ...interface{}) (sql.Result, error) {
	return db.exec(db.db, name, args)
}
func (db *DotDB) TxExec(tx *sqlx.Tx, name string, args ...interface{}) (sql.Result, error) {
	return db.exec(tx, name, args)
}
func (db *DotDB) exec(execer sqlx.Execer, name string, args []interface{}) (sql.Result, error) {
	query, args, err := db.getQueryAndArgsWithIn(name, args)
	if err != nil {
		return nil, err
	}

	return execer.Exec(query, args...)
}

// exec query with named params

func (db *DotDB) NamedExec(name string, arg interface{}) (sql.Result, error) {
	return db.namedExec(db.db, name, arg)
}
func (db *DotDB) TxNamedExec(tx *sqlx.Tx, name string, arg interface{}) (sql.Result, error) {
	return db.namedExec(tx, name, arg)
}
func (db *DotDB) namedExec(execer sqlx.Execer, name string, arg interface{}) (sql.Result, error) {
	query, args, err := db.getNamedQueryAndArgsWithIn(name, arg)
	if err != nil {
		return nil, err
	}

	return execer.Exec(query, args...)
}

// exec query with batch support

func (db *DotDB) ExecWithBatch(name string, args ...interface{}) (sql.Result, error) {
	return db.execWithBatch(db.db, name, args)
}
func (db *DotDB) TxExecWithBatch(tx *sqlx.Tx, name string, args ...interface{}) (sql.Result, error) {
	return db.execWithBatch(tx, name, args)
}
func (db *DotDB) execWithBatch(executor sqlx.Execer, name string, args []interface{}) (sql.Result, error) {
	query, err := db.dotSql.Raw(name)
	if err != nil {
		return nil, err
	}

	return executor.Exec(query, args...)
}

// exec query with named params and batch support

func (db *DotDB) NamedExecWithBatch(name string, args interface{}) (sql.Result, error) {
	return db.namedExecWithBatch(db.db, name, args)
}
func (db *DotDB) TxNamedExecWithBatch(tx *sqlx.Tx, name string, args interface{}) (sql.Result, error) {
	return db.namedExecWithBatch(tx, name, args)
}
func (db *DotDB) namedExecWithBatch(executor sqlx.Execer, name string, arg interface{}) (sql.Result, error) {
	query, err := db.dotSql.Raw(name)
	if err != nil {
		return nil, err
	}

	if named, linearArgs, err := sqlx.Named(query, arg); err != nil {
		return nil, err
	} else {
		return executor.Exec(named, linearArgs...)
	}
}

// utils

func (db *DotDB) getNamedQueryAndArgsWithIn(name string, arg interface{}) (string, []interface{}, error) {
	query, err := db.dotSql.Raw(name)
	if err != nil {
		return "", nil, err
	}

	query, args, err := sqlx.Named(query, arg)
	if err != nil {
		return "", nil, err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return "", nil, err
	}
	query = db.db.Rebind(query)

	return query, args, nil
}

func (db *DotDB) getQueryAndArgsWithIn(name string, args []interface{}) (string, []interface{}, error) {
	query, err := db.dotSql.Raw(name)
	if err != nil {
		return "", nil, err
	}

	query, args, err = sqlx.In(query, args...)
	if err != nil {
		return "", nil, err
	}
	query = db.db.Rebind(query)

	return query, args, nil
}

func (db *DotDB) GetRawDB() *sqlx.DB {
	return db.db
}
