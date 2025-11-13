package database

import (
	"context"
	"fmt"

	"github.com/wfrscltech/vulcano/config"
)

var cnx Database

// Row representa un registro de la base de datos
type Row interface {
	Scan(dest ...any) error
}

// Rows representa un conjunto de registros
type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Close()
}

// DB define operaciones básicas que puede usar la capa de negocio
type Database interface {
	Close()
	Query(ctx context.Context, query string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) Row
	Exec(ctx context.Context, query string, args ...any) (int64, error)

	// Transacciones
	BeginTx(ctx context.Context) (Tx, error)

	// Devuelve la conexión 'en crudo' para que pueda ser usada para operaciones no soportadas
	RawConnection() any
}

// Tx representa una transacción
type Tx interface {
	Query(ctx context.Context, query string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, query string, args ...any) Row
	Exec(ctx context.Context, query string, args ...any) (int64, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

func GetDatabase() Database {
	return cnx
}

func New(dcfg config.DatabaseConfig) error {
	var err error = nil
	switch dcfg.Typo {
	case "postgres":
		cnx, err = newPostgresCnx(dcfg)
		return err
	case "mssql":
		cnx, err = newMSSQLCnx(dcfg)
		return err
	default:
		return fmt.Errorf("no se reconoce el tipo de base de datos %s", dcfg.Typo)
	}
}
