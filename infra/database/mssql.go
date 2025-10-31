package database

import (
	"fmt"

	_ "github.com/microsoft/go-mssqldb"
	"github.com/wfrscltech/vulcano/config"
)

// Adaptador de Microsoft SQL Server siguiendo la especificación de la base de datos y la liberia microsoft/go-mssqldb
type MSSQL struct {
	*sqlBase
}

func NewMSSQLCnx(dcfg *config.DatabaseConfig) (Database, error) {
	cnx, err := newConnection("sqlserver", mssqldsn(dcfg))
	if err != nil {
		return nil, err
	}

	return &MSSQL{sqlBase: &sqlBase{DB: cnx}}, nil
}

func mssqldsn(dcfg *config.DatabaseConfig) string {
	return fmt.Sprintf(
		"sqlserver://%s:%s@%s:%d?database=%s&encrypt=disable",
		dcfg.User,
		dcfg.Password,
		dcfg.Host,
		dcfg.Port,
		dcfg.Name,
	)
}
