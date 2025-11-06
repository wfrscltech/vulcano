# Vulcano

Biblioteca de infraestructura reutilizable para construir servicios HTTP en Go con componentes preconfigurados para bases de datos, servidor web, logging y gestión de configuración.

**Módulo:** `github.com/wfrscltech/vulcano`

## Tabla de Contenidos

- [Características](#características)
- [Instalación](#instalación)
- [Requisitos](#requisitos)
- [Comandos de Desarrollo](#comandos-de-desarrollo)
- [Uso](#uso)
- [Estructura del Proyecto](#estructura-del-proyecto)
- [Arquitectura](#arquitectura)
- [Documentación API](#documentación-api)
- [Endpoints Integrados](#endpoints-integrados)
- [Bases de Datos Soportadas](#bases-de-datos-soportadas)
- [Middleware Incluido](#middleware-incluido)
- [Funciones de Utilidad](#funciones-de-utilidad)
- [Logging](#logging)
- [Cierre Graceful](#cierre-graceful)
- [Extender Vulcano](#extender-vulcano)
- [Dependencias Clave](#dependencias-clave)
- [Notas Técnicas](#notas-técnicas)

## Características

- **Abstracción de Base de Datos**: Interfaz unificada para múltiples motores de base de datos
  - PostgreSQL (usando `jackc/pgx/v5`)
  - Microsoft SQL Server (usando `microsoft/go-mssqldb`)
  - Soporte genérico para cualquier driver compatible con `database/sql`
  - Gestión de transacciones
  - Connection pooling

- **Servidor HTTP**: Configuración predeterminada de Echo Framework
  - Middleware de logging estructurado (slog)
  - Middleware de manejo de errores (RFC 7807 Problem Details)
  - CORS preconfigurado
  - Endpoint de health check
  - Documentación Swagger/OpenAPI integrada

- **Logging Estructurado**: Sistema de logs basado en `log/slog`
  - Formato JSON
  - Rotación automática de archivos
  - Salida dual (consola + archivo)

- **Gestión de Configuración**: Carga y validación de configuración desde JSON
  - Validación estricta de parámetros
  - Mensajes de error descriptivos

- **Cierre Graceful**: Manejo de señales del sistema operativo
  - Timeout configurable
  - Limpieza ordenada de recursos

- **Funciones de Utilidad**: Helpers genéricos para tareas comunes
  - Conversión de texto (camelCase, PascalCase, snake_case)
  - Operador ternario con tipos genéricos
  - Funciones criptográficas
  - Validaciones

## Instalación

```bash
go get github.com/wfrscltech/vulcano
```

## Requisitos

- Go 1.24.0 o superior

## Comandos de Desarrollo

### Testing

```bash
# Ejecutar todos los tests
go test ./...

# Tests con cobertura
go test -cover ./...

# Tests de un paquete específico
go test ./fn
go test ./config
go test ./infra/database

# Ejecutar un test específico
go test ./fn -run TestFunctionName

# Tests con salida detallada
go test -v ./...
```

### Building

```bash
# Descargar dependencias
go mod download

# Ordenar dependencias
go mod tidy

# Verificar dependencias
go mod verify
```

### Calidad de Código

```bash
# Formatear código
go fmt ./...

# Análisis estático
go vet ./...
```

## Uso

### Ejemplo Básico: Servidor HTTP con Echo

```go
package main

import (
    "log/slog"

    "github.com/labstack/echo/v4"
    "github.com/wfrscltech/vulcano/config"
    vulcanoEcho "github.com/wfrscltech/vulcano/infra/echo"
    "github.com/wfrscltech/vulcano/logger"
    "github.com/wfrscltech/vulcano/service"
)

func main() {
    // Cargar configuración
    cfg, err := config.Read("config.json")
    if err != nil {
        panic(err)
    }

    // Inicializar logger
    logger.Init(slog.LevelInfo, cfg.Server.LogDestination)

    // Crear instancia de Echo preconfigurada
    e := vulcanoEcho.NewEchoInstance(
        logger.Log,
        "1.0.0",           // version
        "2024-10-31",      // buildTime
        "abc123",          // commitHash
    )

    // Registrar rutas personalizadas
    e.GET("/api/hello", func(c echo.Context) error {
        return c.JSON(200, map[string]string{"message": "Hola mundo"})
    })

    // Ejecutar servidor con cierre graceful
    if err := service.RunGracefully(logger.Log, e.Server); err != nil {
        logger.Log.Error("Error al ejecutar servidor", "error", err)
    }
}
```

> **Nota:** Se usa `vulcanoEcho` como alias para evitar conflictos con el paquete `echo` de labstack.

### Ejemplo: Conexión a Base de Datos

```go
package main

import (
    "context"
    "log"

    "github.com/wfrscltech/vulcano/config"
    "github.com/wfrscltech/vulcano/infra/database"
)

func main() {
    cfg, _ := config.Read("config.json")

    // Conectar a PostgreSQL
    db, err := database.NewPostgresCnx(&cfg.Database)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Ejecutar consulta
    ctx := context.Background()
    rows, err := db.Query(ctx, "SELECT id, nombre FROM usuarios WHERE activo = $1", true)
    if err != nil {
        log.Fatal(err)
    }
    defer rows.Close()

    // Iterar resultados
    for rows.Next() {
        var id int
        var nombre string
        if err := rows.Scan(&id, &nombre); err != nil {
            log.Fatal(err)
        }
        log.Printf("Usuario: %d - %s\n", id, nombre)
    }
}
```

### Ejemplo: Uso de Transacciones

```go
func transferirSaldo(db database.Database, desde, hacia int, monto float64) error {
    ctx := context.Background()

    // Iniciar transacción
    tx, err := db.BeginTx(ctx)
    if err != nil {
        return err
    }

    // Asegurar rollback en caso de error
    defer tx.Rollback(ctx)

    // Debitar cuenta origen
    _, err = tx.Exec(ctx,
        "UPDATE cuentas SET saldo = saldo - $1 WHERE id = $2",
        monto, desde)
    if err != nil {
        return err
    }

    // Acreditar cuenta destino
    _, err = tx.Exec(ctx,
        "UPDATE cuentas SET saldo = saldo + $1 WHERE id = $2",
        monto, hacia)
    if err != nil {
        return err
    }

    // Confirmar transacción
    return tx.Commit(ctx)
}
```

### Archivo de Configuración (config.json)

```json
{
  "server": {
    "port": 8080,
    "log_level": "info",
    "log_destination": "file:/var/log/miapp"
  },
  "database": {
    "host": "localhost",
    "port": 5432,
    "user": "usuario",
    "password": "contraseña",
    "name": "mi_base_datos",
    "type": "postgres"
  }
}
```

## Estructura del Proyecto

```
vulcano/
├── config/          # Gestión de configuración y validación
├── fn/              # Funciones de utilidad (texto, validaciones, criptografía)
├── infra/           # Implementaciones de infraestructura
│   ├── database/    # Adaptadores de bases de datos
│   └── echo/        # Configuración de Echo Framework
│       ├── apidocs/ # Documentación Swagger/OpenAPI
│       └── middleware/ # Middlewares personalizados
├── logger/          # Sistema de logging estructurado
├── server/          # Abstracción de servidor HTTP
└── service/         # Gestión de ciclo de vida de servicios
```

## Arquitectura

### Patrón de Diseño: Adapter + Dependency Injection

Vulcano sigue una arquitectura por capas con adaptadores e interfaces para portabilidad:

#### 1. Capa de Servicio (`service/`)

Gestiona el ciclo de vida del servicio con cierre graceful:

- **Interfaz `Runner`**: Los servicios deben implementar `Start()` y `Shutdown(ctx)`
- **`RunGracefully()`**: Wrapper que maneja señales del OS (SIGTERM, SIGINT) y activa el cierre graceful con timeout de 5 segundos
- **Build Tags**: Implementaciones específicas por plataforma
  - `service_linux.go`: Implementación para Linux
  - `service_windows.go`: Implementación para Windows

#### 2. Capa de Infraestructura (`infra/`)

Implementaciones concretas de dependencias externas:

**Base de Datos** (`infra/database/`):
- Interfaz `Database`: API unificada para Query/QueryRow/Exec/BeginTx
- `postgres.go`: Adaptador PostgreSQL usando `jackc/pgx/v5` con connection pooling
- `mssql.go`: Adaptador MS SQL Server usando `microsoft/go-mssqldb`
- `sqlbase.go`: Adaptador genérico para cualquier driver compatible con `database/sql`
- Todas las implementaciones envuelven tipos específicos (pgx.Rows, sql.Rows) en interfaces comunes `Rows`/`Row`/`Tx`
- Timeout de conexión: 5 segundos

**Servidor HTTP** (`infra/echo/`):
- `NewEchoInstance()`: Factory que crea instancia Echo preconfigurada
- Stack de middleware: Slog logging → Problem Details (RFC 7807) → CORS
- Endpoint `/health` integrado
- Documentación Swagger/OpenAPI en `/doc/api` (usando Redoc UI)

#### 3. Configuración (`config/`)

Sistema de configuración basado en JSON con validación:

- Struct `Config`: Configuración de servidor + base de datos
- `Read(path)`: Carga y valida archivos JSON
- Validación forzada en métodos `IsValid()` con mensajes de error en español
- Tipos de base de datos y niveles de log soportados definidos en `constants.go`

#### 4. Logging (`logger/`)

Logging estructurado JSON con rotación:

- Usa `log/slog` de Go para logging estructurado
- `lumberjack` para rotación de logs (10MB max, 5 backups, retención 30 días, compresión)
- Escribe simultáneamente a stdout y archivo

#### 5. Patrón de Configuración del Servicio

Los servicios que usan esta biblioteca deberían:

1. Cargar configuración con `config.Read(path)`
2. Inicializar logger con `logger.Init(level, dir)`
3. Crear conexión de base de datos con `database.NewPostgresCnx()` o similar
4. Crear instancia Echo con `echo.NewEchoInstance(logger, version, buildTime, commitHash)`
5. Registrar rutas y handlers
6. Ejecutar con `service.RunGracefully(logger, runner)`

## Documentación API

Vulcano incluye soporte integrado para documentación Swagger/OpenAPI. Una vez iniciado el servidor, la documentación interactiva está disponible en:

```
http://localhost:8080/doc/api
```

La especificación OpenAPI en formato JSON está disponible en:

```
http://localhost:8080/doc/spec/swagger.json
```

## Endpoints Integrados

### Health Check

```
GET /health
```

Retorna información sobre el estado del servicio:

```json
{
  "status": "ok",
  "version": "1.0.0",
  "build_time": "2024-10-31",
  "commit_hash": "abc123"
}
```

## Bases de Datos Soportadas

| Base de Datos | Identificador en Config | Driver | Características |
|---------------|------------------------|---------|-----------------|
| PostgreSQL | `postgres` | `jackc/pgx/v5` | Connection pooling automático, recomendado para apps modernas |
| MS SQL Server | `mssql` | `microsoft/go-mssqldb` | Compatible con SQL Server 2012+ |
| Genérico | `sqlbase` | `database/sql` | Para cualquier driver compatible con `database/sql` |

### Configuración por Tipo de Base de Datos

**PostgreSQL:**
```json
{
  "database": {
    "host": "localhost",
    "port": 5432,
    "user": "usuario",
    "password": "contraseña",
    "name": "mi_base_datos",
    "type": "postgres"
  }
}
```

**MS SQL Server:**
```json
{
  "database": {
    "host": "localhost",
    "port": 1433,
    "user": "sa",
    "password": "contraseña",
    "name": "mi_base_datos",
    "type": "mssql"
  }
}
```

## Middleware Incluido

1. **SlogMiddleware**: Logging estructurado de todas las peticiones HTTP
   - Registra método, ruta, código de estado, latencia e IP
   - Niveles automáticos: ERROR (5xx), WARN (4xx), INFO (otros)

2. **ProblemMiddleware**: Manejo estandarizado de errores según RFC 7807
   - Convierte errores en formato Problem Details JSON
   - Facilita debugging y manejo de errores en clientes

3. **CORS**: Configurado por defecto para permitir todas las origines
   - Permite métodos: GET, POST, PUT, OPTIONS
   - Permite todos los headers

## Funciones de Utilidad

### Operador Ternario

```go
import "github.com/wfrscltech/vulcano/fn"

resultado := fn.TernaryIf(edad >= 18, "adulto", "menor")
```

### Conversión de Texto

```go
import "github.com/wfrscltech/vulcano/fn"

// CamelCase a snake_case
snake := fn.ToSnakeCase("MiVariable") // "mi_variable"

// snake_case a PascalCase
pascal := fn.ToPascalCase("mi_variable") // "MiVariable"
```

## Logging

El sistema de logging usa `log/slog` para logs estructurados:

```go
logger.Log.Info("Operación completada",
    "usuario_id", 123,
    "operacion", "transferencia",
    "monto", 1000.50)

logger.Log.Error("Error en la operación",
    "error", err.Error(),
    "contexto", "procesamiento_pago")
```

### Niveles de Log Soportados

Configure el nivel de log en `config.json`:

- `debug`: Información detallada para debugging
- `info`: Información general del flujo de la aplicación (por defecto)
- `warn`: Advertencias sobre situaciones no ideales
- `error`: Errores que requieren atención

**Ejemplo de configuración:**
```json
{
  "server": {
    "log_level": "info",
    "log_destination": "file:/var/log/miapp"
  }
}
```

### Rotación de Logs

Los logs se guardan en formato JSON con rotación automática:
- Tamaño máximo: 10 MB por archivo
- Backups: 5 archivos
- Retención: 30 días
- Compresión: Habilitada
- Salida dual: stdout + archivo simultáneo

## Cierre Graceful

Los servicios que implementan la interfaz `Runner` pueden usar `service.RunGracefully()` para manejar el cierre ordenado:

```go
type Runner interface {
    Start() error
    Shutdown(ctx context.Context) error
}
```

El cierre graceful:
1. Captura señales del sistema (SIGTERM, SIGINT)
2. Llama al método `Shutdown()` con timeout de 5 segundos
3. Permite limpieza ordenada de recursos

## Extender Vulcano

### Agregar Soporte para Nueva Base de Datos

El patrón de adaptador permite agregar nuevos motores de base de datos fácilmente:

1. Crear archivo en `infra/database/` (ej: `mysql.go`)
2. Implementar la interfaz `Database` con tipos concretos:
   ```go
   type Database interface {
       Query(ctx context.Context, query string, args ...interface{}) (Rows, error)
       QueryRow(ctx context.Context, query string, args ...interface{}) Row
       Exec(ctx context.Context, query string, args ...interface{}) (Result, error)
       BeginTx(ctx context.Context) (Tx, error)
       Close() error
   }
   ```
3. Crear tipos wrapper para envolver tipos específicos del vendor:
   - Wrapper para `Row` (ej: envuelve `sql.Row`)
   - Wrapper para `Rows` (ej: envuelve `sql.Rows`)
   - Wrapper para `Tx` (ej: envuelve `sql.Tx`)
4. Implementar función factory `New<DatabaseName>Cnx(*config.DatabaseConfig)` que:
   - Construye la cadena de conexión
   - Establece connection pooling si está disponible
   - Prueba la conexión con timeout de 5 segundos
   - Retorna instancia que implementa `Database`
5. Agregar el tipo de base de datos a `supportedDatabaseTypes` en `config/constants.go`

**Ejemplo:** PostgreSQL usa connection pooling (`pgxpool`), mientras que MSSQL usa `database/sql` estándar. Ambos implementan la misma interfaz `Database`, haciéndolos intercambiables en el código cliente.

### Agregar Middleware Personalizado

```go
func MiMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        // Lógica antes del handler
        err := next(c)
        // Lógica después del handler
        return err
    }
}

// Usar en Echo
e.Use(MiMiddleware)
```

## Dependencias Clave

Vulcano está construido sobre las siguientes dependencias principales:

- **Echo v4**: Framework HTTP de alto rendimiento
- **jackc/pgx/v5**: Driver PostgreSQL con connection pooling automático
- **microsoft/go-mssqldb**: Driver oficial para Microsoft SQL Server
- **swaggo/swag**: Generación automática de documentación Swagger/OpenAPI
- **lumberjack.v2**: Rotación de archivos de log
- **log/slog**: Biblioteca estándar de Go para logging estructurado

## Notas Técnicas

- Todos los mensajes de error y logs están en español
- La validación de configuración es estricta; campos faltantes o inválidos causan errores inmediatos
- El timeout de cierre graceful está fijado en 5 segundos en `service/runner.go`
- Las conexiones de base de datos se prueban con timeout de 5 segundos en la inicialización
- CORS está configurado por defecto para permitir todas las origenes (`*`)
- El endpoint de health retorna versión, build time y commit hash como metadatos
- Los build tags aseguran que solo se compile el archivo apropiado por plataforma (`service_linux.go` o `service_windows.go`)

## Licencia

[Información de licencia pendiente]

## Soporte

Para reportar problemas o solicitar características, contacte al equipo de desarrollo de CLTech.
