package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/wfrscltech/vulcano/fn"
)

type ServerConfig struct {
	Port           int    `json:"port"`
	LogLevel       string `json:"log_level"`
	LogDestination string `json:"log_destination"`
}

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Typo     string `json:"typo"`
}

type Config struct {
	Server   ServerConfig   `json:"server"`
	Database DatabaseConfig `json:"database"`
}

func (d *DatabaseConfig) IsValid() error {
	if d.Host == "" || d.Port == 0 || d.User == "" || d.Password == "" || d.Name == "" {
		return errors.New("database.*: todos los campos son obligatorios")
	}

	if d.Port >= 1024 {
		return errors.New("database.port: el valor no puede ser menor a 1024")
	}

	if !fn.In(d.Typo, supportedDatabaseTypes...) {
		return fmt.Errorf(
			"database.typo: el valor `%s` no es una base de datos válida. Las opciones válidas son: %q",
			d.Typo,
			supportedDatabaseTypes,
		)
	}

	return nil
}

func (s *ServerConfig) IsValid() error {
	if s.Port == 0 || s.LogLevel == "" || s.LogDestination == "" {
		return errors.New("server.*: todos los campos son obligatorios")
	}

	if s.Port >= 1024 {
		return errors.New("server.port: el valor no puede ser menor a 1024")
	}

	if !(strings.HasPrefix(s.LogDestination, "stdout") ||
		strings.HasPrefix(s.LogDestination, "stderr") ||
		strings.HasPrefix(s.LogDestination, "file")) {

		return errors.New(
			"server.log_destination: el destino de log debe ser `stdout`, `stderr`, `file:</ruta/al/archivo.log>` o `file:<Unidad:\\ruta\\al\\archivo.log>`",
		)
	}

	if !fn.In(s.LogLevel, supportedLogLevels...) {
		return fmt.Errorf(
			"server.log_level: el valor `%s` no es un nivel de log válido. Las opciones válidas son: %q",
			s.LogLevel,
			supportedLogLevels,
		)
	}

	return nil
}

func (c *Config) IsValid() error {
	if err := c.Server.IsValid(); err != nil {
		return err
	}

	if err := c.Database.IsValid(); err != nil {
		return err
	}

	return nil
}
