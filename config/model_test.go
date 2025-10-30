package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestDatabaseConfigIsValid_Success valida configuraciones válidas de base de datos
func TestDatabaseConfigIsValid_Success(t *testing.T) {
	tests := []struct {
		name   string
		config DatabaseConfig
	}{
		{
			name: "Configuración válida con postgres",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     80,
				User:     "admin",
				Password: "secret123",
				Name:     "mydb",
				Typo:     DatabaseTypePostgres,
			},
		},
		{
			name: "Configuración válida con mssql",
			config: DatabaseConfig{
				Host:     "192.168.1.100",
				Port:     443,
				User:     "sa",
				Password: "StrongP@ss",
				Name:     "production",
				Typo:     DatabaseTypeMssql,
			},
		},
		{
			name: "Puerto mínimo permitido (1)",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     1,
				User:     "user",
				Password: "pass",
				Name:     "db",
				Typo:     DatabaseTypePostgres,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.IsValid()
			if err != nil {
				t.Errorf("Se esperaba configuración válida, pero obtuvo error: %v", err)
			}
		})
	}
}

// TestDatabaseConfigIsValid_Failures valida que se detecten configuraciones inválidas
func TestDatabaseConfigIsValid_Failures(t *testing.T) {
	tests := []struct {
		name          string
		config        DatabaseConfig
		expectedError string
	}{
		{
			name: "Host vacío",
			config: DatabaseConfig{
				Host:     "",
				Port:     5432,
				User:     "admin",
				Password: "secret",
				Name:     "mydb",
				Typo:     DatabaseTypePostgres,
			},
			expectedError: "database.*: todos los campos son obligatorios",
		},
		{
			name: "Puerto cero",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     0,
				User:     "admin",
				Password: "secret",
				Name:     "mydb",
				Typo:     DatabaseTypePostgres,
			},
			expectedError: "database.*: todos los campos son obligatorios",
		},
		{
			name: "Usuario vacío",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "",
				Password: "secret",
				Name:     "mydb",
				Typo:     DatabaseTypePostgres,
			},
			expectedError: "database.*: todos los campos son obligatorios",
		},
		{
			name: "Password vacío",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "admin",
				Password: "",
				Name:     "mydb",
				Typo:     DatabaseTypePostgres,
			},
			expectedError: "database.*: todos los campos son obligatorios",
		},
		{
			name: "Nombre de base de datos vacío",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     5432,
				User:     "admin",
				Password: "secret",
				Name:     "",
				Typo:     DatabaseTypePostgres,
			},
			expectedError: "database.*: todos los campos son obligatorios",
		},
		{
			name: "Puerto mayor o igual a 1024",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     1024,
				User:     "admin",
				Password: "secret",
				Name:     "mydb",
				Typo:     DatabaseTypePostgres,
			},
			expectedError: "database.port: el valor no puede ser menor a 1024",
		},
		{
			name: "Puerto mayor a 1024",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     8080,
				User:     "admin",
				Password: "secret",
				Name:     "mydb",
				Typo:     DatabaseTypePostgres,
			},
			expectedError: "database.port: el valor no puede ser menor a 1024",
		},
		{
			name: "Tipo de base de datos inválido",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     80,
				User:     "admin",
				Password: "secret",
				Name:     "mydb",
				Typo:     "mysql",
			},
			expectedError: "database.typo: el valor `mysql` no es una base de datos válida",
		},
		{
			name: "Tipo de base de datos vacío",
			config: DatabaseConfig{
				Host:     "localhost",
				Port:     80,
				User:     "admin",
				Password: "secret",
				Name:     "mydb",
				Typo:     "",
			},
			expectedError: "database.typo: el valor `` no es una base de datos válida",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.IsValid()
			if err == nil {
				t.Error("Se esperaba un error pero no se obtuvo ninguno")
				return
			}

			if tt.expectedError != "" && !contains(err.Error(), tt.expectedError) {
				t.Errorf("Error esperado que contenga %q, pero obtuvo: %q", tt.expectedError, err.Error())
			}
		})
	}
}

// TestServerConfigIsValid_Success valida configuraciones válidas de servidor
func TestServerConfigIsValid_Success(t *testing.T) {
	tests := []struct {
		name   string
		config ServerConfig
	}{
		{
			name: "Configuración válida con stdout y debug",
			config: ServerConfig{
				Port:           80,
				LogLevel:       "debug",
				LogDestination: "stdout",
			},
		},
		{
			name: "Configuración válida con stderr y error",
			config: ServerConfig{
				Port:           443,
				LogLevel:       "error",
				LogDestination: "stderr",
			},
		},
		{
			name: "Configuración válida con archivo en Linux",
			config: ServerConfig{
				Port:           999,
				LogLevel:       "info",
				LogDestination: "file:/var/log/app.log",
			},
		},
		{
			name: "Configuración válida con archivo en Windows",
			config: ServerConfig{
				Port:           800,
				LogLevel:       "warning",
				LogDestination: "file:C:\\logs\\app.log",
			},
		},
		{
			name: "Puerto mínimo permitido (1)",
			config: ServerConfig{
				Port:           1,
				LogLevel:       "info",
				LogDestination: "stdout",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.IsValid()
			if err != nil {
				t.Errorf("Se esperaba configuración válida, pero obtuvo error: %v", err)
			}
		})
	}
}

// TestServerConfigIsValid_Failures valida que se detecten configuraciones inválidas
func TestServerConfigIsValid_Failures(t *testing.T) {
	tests := []struct {
		name          string
		config        ServerConfig
		expectedError string
	}{
		{
			name: "Puerto cero",
			config: ServerConfig{
				Port:           0,
				LogLevel:       "info",
				LogDestination: "stdout",
			},
			expectedError: "server.*: todos los campos son obligatorios",
		},
		{
			name: "LogLevel vacío",
			config: ServerConfig{
				Port:           8080,
				LogLevel:       "",
				LogDestination: "stdout",
			},
			expectedError: "server.*: todos los campos son obligatorios",
		},
		{
			name: "LogDestination vacío",
			config: ServerConfig{
				Port:           8080,
				LogLevel:       "info",
				LogDestination: "",
			},
			expectedError: "server.*: todos los campos son obligatorios",
		},
		{
			name: "Puerto mayor o igual a 1024",
			config: ServerConfig{
				Port:           1024,
				LogLevel:       "info",
				LogDestination: "stdout",
			},
			expectedError: "server.port: el valor no puede ser menor a 1024",
		},
		{
			name: "Puerto mayor a 1024",
			config: ServerConfig{
				Port:           8080,
				LogLevel:       "info",
				LogDestination: "stdout",
			},
			expectedError: "server.port: el valor no puede ser menor a 1024",
		},
		{
			name: "Puerto estándar que es rechazado (1024)",
			config: ServerConfig{
				Port:           1024,
				LogLevel:       "info",
				LogDestination: "stdout",
			},
			expectedError: "server.port: el valor no puede ser menor a 1024",
		},
		{
			name: "LogDestination inválido",
			config: ServerConfig{
				Port:           500,
				LogLevel:       "info",
				LogDestination: "syslog",
			},
			expectedError: "server.log_destination: el destino de log debe ser",
		},
		{
			name: "LogLevel inválido",
			config: ServerConfig{
				Port:           500,
				LogLevel:       "critical",
				LogDestination: "stdout",
			},
			expectedError: "server.log_level: el valor `critical` no es un nivel de log válido",
		},
		{
			name: "Múltiples campos inválidos",
			config: ServerConfig{
				Port:           0,
				LogLevel:       "",
				LogDestination: "",
			},
			expectedError: "server.*: todos los campos son obligatorios",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.IsValid()
			if err == nil {
				t.Error("Se esperaba un error pero no se obtuvo ninguno")
				return
			}

			if tt.expectedError != "" && !contains(err.Error(), tt.expectedError) {
				t.Errorf("Error esperado que contenga %q, pero obtuvo: %q", tt.expectedError, err.Error())
			}
		})
	}
}

// TestConfigIsValid_Success valida configuraciones completas válidas
func TestConfigIsValid_Success(t *testing.T) {
	config := Config{
		Server: ServerConfig{
			Port:           80,
			LogLevel:       "info",
			LogDestination: "stdout",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     443,
			User:     "admin",
			Password: "secret",
			Name:     "mydb",
			Typo:     DatabaseTypePostgres,
		},
	}

	err := config.IsValid()
	if err != nil {
		t.Errorf("Se esperaba configuración válida, pero obtuvo error: %v", err)
	}
}

// TestConfigIsValid_Failures valida que se detecten errores en cualquier sección
func TestConfigIsValid_Failures(t *testing.T) {
	tests := []struct {
		name          string
		config        Config
		expectedError string
	}{
		{
			name: "Error en configuración de servidor",
			config: Config{
				Server: ServerConfig{
					Port:           0,
					LogLevel:       "info",
					LogDestination: "stdout",
				},
				Database: DatabaseConfig{
					Host:     "localhost",
					Port:     5432,
					User:     "admin",
					Password: "secret",
					Name:     "mydb",
					Typo:     DatabaseTypePostgres,
				},
			},
			expectedError: "server.*: todos los campos son obligatorios",
		},
		{
			name: "Error en configuración de base de datos",
			config: Config{
				Server: ServerConfig{
					Port:           80,
					LogLevel:       "info",
					LogDestination: "stdout",
				},
				Database: DatabaseConfig{
					Host:     "",
					Port:     443,
					User:     "admin",
					Password: "secret",
					Name:     "mydb",
					Typo:     DatabaseTypePostgres,
				},
			},
			expectedError: "database.*: todos los campos son obligatorios",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.IsValid()
			if err == nil {
				t.Error("Se esperaba un error pero no se obtuvo ninguno")
				return
			}

			if tt.expectedError != "" && !contains(err.Error(), tt.expectedError) {
				t.Errorf("Error esperado que contenga %q, pero obtuvo: %q", tt.expectedError, err.Error())
			}
		})
	}
}

// TestReadJSON_Success valida la lectura correcta de archivos JSON
func TestReadJSON_Success(t *testing.T) {
	// Crear un archivo temporal con JSON válido
	content := `{
		"server": {
			"port": 80,
			"log_level": "info",
			"log_destination": "stdout"
		},
		"database": {
			"host": "localhost",
			"port": 443,
			"user": "admin",
			"password": "secret",
			"name": "mydb",
			"typo": "postgres"
		}
	}`

	tmpFile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatalf("Error al crear archivo temporal: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatalf("Error al escribir archivo temporal: %v", err)
	}
	tmpFile.Close()

	var config Config
	err = ReadJSON(tmpFile.Name(), &config)
	if err != nil {
		t.Errorf("No se esperaba error, pero obtuvo: %v", err)
	}

	// Validar que los datos se cargaron correctamente
	if config.Server.Port != 80 {
		t.Errorf("Se esperaba Port=80, obtuvo: %d", config.Server.Port)
	}
	if config.Database.Host != "localhost" {
		t.Errorf("Se esperaba Host=localhost, obtuvo: %s", config.Database.Host)
	}
}

// TestReadJSON_Failures valida el manejo de errores en la lectura de JSON
func TestReadJSON_Failures(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		expectedError string
	}{
		{
			name:          "Archivo no existe",
			content:       "",
			expectedError: "error al abrir el archivo de configuración",
		},
		{
			name:          "JSON mal formado - sintaxis inválida",
			content:       `{"server": "port": 8080}`,
			expectedError: "el cuerpo de este documento contiene JSON mal formado",
		},
		{
			name:          "JSON vacío",
			content:       "",
			expectedError: "el cuerpo no debe estar vacío",
		},
		{
			name:          "JSON con tipo incorrecto",
			content:       `{"server": {"port": "not-a-number", "log_level": "info", "log_destination": "stdout"}}`,
			expectedError: "el cuerpo contiene JSON de tipo incorrecto",
		},
		{
			name:          "JSON incompleto",
			content:       `{"server": {`,
			expectedError: "el cuerpo contiene JSON mal formado",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var config Config

			if tt.name == "Archivo no existe" {
				err := ReadJSON("/ruta/inexistente/config.json", &config)
				if err == nil {
					t.Error("Se esperaba un error pero no se obtuvo ninguno")
					return
				}
				if !contains(err.Error(), tt.expectedError) {
					t.Errorf("Error esperado que contenga %q, pero obtuvo: %q", tt.expectedError, err.Error())
				}
				return
			}

			// Para otros casos, crear archivo temporal
			tmpFile, err := os.CreateTemp("", "config-*.json")
			if err != nil {
				t.Fatalf("Error al crear archivo temporal: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			if tt.content != "" {
				if _, err := tmpFile.Write([]byte(tt.content)); err != nil {
					t.Fatalf("Error al escribir archivo temporal: %v", err)
				}
			}
			tmpFile.Close()

			err = ReadJSON(tmpFile.Name(), &config)
			if err == nil {
				t.Error("Se esperaba un error pero no se obtuvo ninguno")
				return
			}

			if tt.expectedError != "" && !contains(err.Error(), tt.expectedError) {
				t.Errorf("Error esperado que contenga %q, pero obtuvo: %q", tt.expectedError, err.Error())
			}
		})
	}
}

// TestRead_Success valida la función Read completa con archivo válido
func TestRead_Success(t *testing.T) {
	content := `{
		"server": {
			"port": 999,
			"log_level": "debug",
			"log_destination": "file:/var/log/app.log"
		},
		"database": {
			"host": "192.168.1.100",
			"port": 443,
			"user": "dbuser",
			"password": "dbpass123",
			"name": "production",
			"typo": "postgres"
		}
	}`

	tmpFile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatalf("Error al crear archivo temporal: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatalf("Error al escribir archivo temporal: %v", err)
	}
	tmpFile.Close()

	config, err := Read(tmpFile.Name())
	if err != nil {
		t.Errorf("No se esperaba error, pero obtuvo: %v", err)
	}

	if config == nil {
		t.Fatal("La configuración no debería ser nil")
	}

	// Validar algunos valores cargados
	if config.Server.Port != 999 {
		t.Errorf("Se esperaba Port=999, obtuvo: %d", config.Server.Port)
	}
	if config.Database.Name != "production" {
		t.Errorf("Se esperaba Name=production, obtuvo: %s", config.Database.Name)
	}
}

// TestRead_Failures valida el manejo de errores en la función Read
func TestRead_Failures(t *testing.T) {
	tests := []struct {
		name          string
		content       string
		setupFile     bool
		expectedError string
	}{
		{
			name:          "Archivo no existe",
			content:       "",
			setupFile:     false,
			expectedError: "error al abrir el archivo de configuración",
		},
		{
			name: "JSON inválido",
			content: `{
				"server": {
					"port": "invalid"
				}
			}`,
			setupFile:     true,
			expectedError: "el cuerpo contiene JSON de tipo incorrecto",
		},
		{
			name: "Configuración inválida - puerto de servidor inválido",
			content: `{
				"server": {
					"port": 0,
					"log_level": "info",
					"log_destination": "stdout"
				},
				"database": {
					"host": "localhost",
					"port": 5432,
					"user": "admin",
					"password": "secret",
					"name": "mydb",
					"typo": "postgres"
				}
			}`,
			setupFile:     true,
			expectedError: "server.*: todos los campos son obligatorios",
		},
		{
			name: "Configuración inválida - tipo de base de datos incorrecto",
			content: `{
				"server": {
					"port": 80,
					"log_level": "info",
					"log_destination": "stdout"
				},
				"database": {
					"host": "localhost",
					"port": 443,
					"user": "admin",
					"password": "secret",
					"name": "mydb",
					"typo": "oracle"
				}
			}`,
			setupFile:     true,
			expectedError: "database.typo: el valor `oracle` no es una base de datos válida",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filePath string

			if tt.setupFile {
				tmpFile, err := os.CreateTemp("", "config-*.json")
				if err != nil {
					t.Fatalf("Error al crear archivo temporal: %v", err)
				}
				defer os.Remove(tmpFile.Name())

				if _, err := tmpFile.Write([]byte(tt.content)); err != nil {
					t.Fatalf("Error al escribir archivo temporal: %v", err)
				}
				tmpFile.Close()
				filePath = tmpFile.Name()
			} else {
				filePath = filepath.Join(os.TempDir(), "nonexistent-config.json")
			}

			config, err := Read(filePath)
			if err == nil {
				t.Error("Se esperaba un error pero no se obtuvo ninguno")
				return
			}

			if config != nil {
				t.Error("La configuración debería ser nil cuando hay un error")
			}

			if tt.expectedError != "" && !contains(err.Error(), tt.expectedError) {
				t.Errorf("Error esperado que contenga %q, pero obtuvo: %q", tt.expectedError, err.Error())
			}
		})
	}
}

// TestReadJSON_EmptyObject valida el manejo de un objeto JSON vacío
func TestReadJSON_EmptyObject(t *testing.T) {
	content := `{}`

	tmpFile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatalf("Error al crear archivo temporal: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatalf("Error al escribir archivo temporal: %v", err)
	}
	tmpFile.Close()

	var config Config
	err = ReadJSON(tmpFile.Name(), &config)

	// ReadJSON no debería fallar con un objeto vacío, ya que es JSON válido
	if err != nil {
		t.Errorf("No se esperaba error para JSON vacío, pero obtuvo: %v", err)
	}
}

// TestRead_EmptyObject valida que Read falle con configuración vacía
func TestRead_EmptyObject(t *testing.T) {
	content := `{}`

	tmpFile, err := os.CreateTemp("", "config-*.json")
	if err != nil {
		t.Fatalf("Error al crear archivo temporal: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	if _, err := tmpFile.Write([]byte(content)); err != nil {
		t.Fatalf("Error al escribir archivo temporal: %v", err)
	}
	tmpFile.Close()

	config, err := Read(tmpFile.Name())

	// Read debería fallar porque la configuración está vacía y no es válida
	if err == nil {
		t.Error("Se esperaba un error para configuración vacía, pero no se obtuvo ninguno")
	}

	if config != nil {
		t.Error("La configuración debería ser nil cuando hay un error de validación")
	}
}

// contains es una función auxiliar para verificar si una cadena contiene otra
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && len(substr) > 0 && containsHelper(s, substr)))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
