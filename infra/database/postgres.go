package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wfrscltech/vulcano/config"
)

// Adaptador de PostgreSQL siguiendo la especificación de la base de datos y la liberia jackc/pgx
type Postgres struct {
	pool *pgxpool.Pool
}

// Representación de una consulta de varias filas en PostgreSQL
type PostgresRows struct {
	pgx.Rows
}

// Representación de una fila en PostgreSQL
type PostgresRow struct {
	pgx.Row
}

// Representación de una transacción en PostgreSQL
type PostgresTx struct {
	pgx.Tx
}

func NewPostgresCnx(dcfg config.DatabaseConfig) (Database, error) {
	cnx, err := pgxpool.New(context.Background(), psqldsn(dcfg))
	if err != nil {
		return nil, err
	}

	// Control de conexión
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = cnx.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return &Postgres{pool: cnx}, nil
}

func (db *Postgres) Close() {
	if db.pool != nil {
		db.pool.Close()
	}
}

func (db *Postgres) Query(ctx context.Context, query string, args ...any) (Rows, error) {
	rows, err := db.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &PostgresRows{rows}, nil
}

func (db *Postgres) QueryRow(ctx context.Context, query string, args ...any) Row {
	return &PostgresRow{db.pool.QueryRow(ctx, query, args...)}
}

func (db *Postgres) Exec(ctx context.Context, query string, args ...any) (int64, error) {
	cmd, err := db.pool.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return cmd.RowsAffected(), nil
}

func (db *Postgres) BeginTx(ctx context.Context) (Tx, error) {
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	return &PostgresTx{tx}, nil
}

func (db *Postgres) RawConnection() any {
	return db.pool
}

// --- Adaptadores de Rows/Row ---

func (r *PostgresRows) Next() bool {
	return r.Rows.Next()
}

func (r *PostgresRows) Scan(dest ...any) error {
	return r.Rows.Scan(dest...)
}

func (r *PostgresRows) Close() {
	r.Rows.Close()
}

func (r *PostgresRow) Scan(dest ...any) error {
	return r.Row.Scan(dest...)
}

// --- Transacción ---

func (tx *PostgresTx) Query(ctx context.Context, query string, args ...any) (Rows, error) {
	rows, err := tx.Tx.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	return &PostgresRows{rows}, nil
}

func (tx *PostgresTx) QueryRow(ctx context.Context, query string, args ...any) Row {
	return &PostgresRow{tx.Tx.QueryRow(ctx, query, args...)}
}

func (tx *PostgresTx) Exec(ctx context.Context, query string, args ...any) (int64, error) {
	cmd, err := tx.Tx.Exec(ctx, query, args...)
	if err != nil {
		return 0, err
	}
	return cmd.RowsAffected(), nil
}

func (tx *PostgresTx) Commit(ctx context.Context) error {
	return tx.Tx.Commit(ctx)
}

func (tx *PostgresTx) Rollback(ctx context.Context) error {
	return tx.Tx.Rollback(ctx)
}

// --- Adaptador de conexión ---

func psqldsn(dcfg config.DatabaseConfig) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		dcfg.User,
		dcfg.Password,
		dcfg.Host,
		dcfg.Port,
		dcfg.Name,
	)
}
