package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func Read(path string) (*Config, error) {
	var cfg *Config

	if err := ReadJSON(path, &cfg); err != nil {
		return nil, err
	}

	if err := cfg.IsValid(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// ReadJSON decodifica el cuerpo de un documento JSON en el destino.
func ReadJSON(path string, dst any) error {
	// Abre el archivo especificado.
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error al abrir el archivo de configuración `%s`: %w", filepath.Base(path), err)
	}
	defer file.Close()

	// Decodifica el cuerpo de la petición en el destino.
	err = json.NewDecoder(file).Decode(dst)
	if err != nil {
		// Si hay un error durante la decodificación, comienza el triaje...
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		// Utilice la función errors.As() para comprobar si el error tiene el tipo
		// *json.SyntaxError. Si lo tiene, devuelve un mensaje de error en inglés plano
		// que incluye la localización del problema.
		case errors.As(err, &syntaxError):
			return fmt.Errorf("el cuerpo de este documento contiene JSON mal formado (en el carácter %d)", syntaxError.Offset)

		// En algunas circunstancias Decode() también puede devolver un error io.ErrUnexpectedEOF
		// por errores de sintaxis en el JSON. Así que comprobamos esto usando errors.Is() y
		// devuelve un mensaje de error genérico. Hay un tema abierto al respecto en
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("el cuerpo contiene JSON mal formado")

		// Del mismo modo, captura cualquier error *json.UnmarshalTypeError. Estos ocurren cuando el valor
		// valor JSON es del tipo incorrecto para el destino. Si el error se refiere
		// con un campo específico, lo incluimos en nuestro mensaje de error para facilitar
		// para facilitar la depuración por parte del cliente.
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("el cuerpo contiene JSON de tipo incorrecto para el campo %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("el cuerpo contiene JSON mal formado (en el carácter %d)", unmarshalTypeError.Offset)

		// Decode() devolverá un error io.EOF si el cuerpo de la petición está vacío. Nosotros
		// comprobamos esto con errors.Is() y devolvemos un mensaje de error en su lugar.
		case errors.Is(err, io.EOF):
			return errors.New("el cuerpo no debe estar vacío")

		// Para cualquier otra cosa, devuelve el mensaje de error tal cual.
		default:
			return err
		}
	}

	return nil
}
