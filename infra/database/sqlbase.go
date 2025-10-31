package database

import (
	"context"
	"database/sql"
	"time"
)

// Adaptador base con soporte para las bases de datos compatibles con el driver sql de golang
type sqlBase struct {
	DB *sql.DB
}

// Representación de una transacción de una base de datos compatible con el driver sql de golang
type sqlBaseTx struct {
	*sql.Tx
}

// Representación de una consulta de varias filas de una base de datos compatible con el driver sql de golang
type sqlBaseRows struct {
	*sql.Rows
}

// Representación de una fila de una base de datos compatible con el driver sql de golang
type sqlBaseRow struct {
	*sql.Row
}

func newConnection(driverName, dataSourceName string) (*sql.DB, error) {
	cnx, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	// Control de conexión
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = cnx.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	return cnx, nil
}

func (db *sqlBase) Close() {
	if db.DB != nil {
		db.DB.Close()
	}
}

func (db *sqlBase) Query(ctx context.Context, query string, args ...any) (Rows, error) {
	rows, err := db.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &sqlBaseRows{rows}, nil
}

func (db *sqlBase) QueryRow(ctx context.Context, query string, args ...any) Row {
	return &sqlBaseRow{db.DB.QueryRowContext(ctx, query, args...)}
}

func (db *sqlBase) Exec(ctx context.Context, query string, args ...any) (int64, error) {
	res, err := db.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}

func (db *sqlBase) BeginTx(ctx context.Context) (Tx, error) {
	tx, err := db.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return &sqlBaseTx{tx}, nil
}

func (db *sqlBase) RawConnection() any {
	return db.DB
}

// --- Adaptadores de Rows/Row ---

func (r *sqlBaseRows) Next() bool {
	return r.Rows.Next()
}

func (r *sqlBaseRows) Scan(dest ...any) error {
	return r.Rows.Scan(dest...)
}

func (r *sqlBaseRows) Close() {
	r.Rows.Close()
}

func (r *sqlBaseRow) Scan(dest ...any) error {
	return r.Row.Scan(dest...)
}

// --- Transacción ---

func (tx *sqlBaseTx) Query(ctx context.Context, query string, args ...any) (Rows, error) {
	rows, err := tx.Tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &sqlBaseRows{rows}, nil
}

func (tx *sqlBaseTx) QueryRow(ctx context.Context, query string, args ...any) Row {
	return &sqlBaseRow{tx.Tx.QueryRowContext(ctx, query, args...)}
}

func (tx *sqlBaseTx) Exec(ctx context.Context, query string, args ...any) (int64, error) {
	res, err := tx.Tx.ExecContext(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}

func (tx *sqlBaseTx) Commit(ctx context.Context) error {
	return tx.Tx.Commit()
}

func (tx *sqlBaseTx) Rollback(ctx context.Context) error {
	return tx.Tx.Rollback()
}
