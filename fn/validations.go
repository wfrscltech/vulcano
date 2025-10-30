package fn

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"
)

// ValidateIntegerList Valida si una lista de números enteros es válida siguiendo los siguientes criterios:
// 1. No puede estar vacío
// 2. Cada valor debe ser un número entero
// 3. Cada valor debe estar separado por comas
func ValidateIntegerList(value string) error {
	if strings.TrimSpace(value) == "" {
		return errors.New("no puede estar vacío")
	}

	parts := strings.Split(value, ",")
	for i, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			return fmt.Errorf("valor vacío en posición %d", i+1)
		}

		if _, err := strconv.Atoi(part); err != nil {
			return fmt.Errorf("poscición %d: `%s` no es un número entero", i+1, part)
		}
	}

	return nil
}

// ValidateStringList Valida si una lista de strings es válida siguiendo los siguientes criterios:
// 1. No puede estar vacío
// 2. Cada valor debe estar separado por comas (no puede tener espacios en blanco)
func ValidateStringList(value string) error {
	value = strings.ToLower(strings.TrimSpace(value))
	if value == "" {
		return errors.New("no puede estar vacío")
	}

	parts := strings.Split(value, ",")
	for i, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			return fmt.Errorf("valor vacío en posición %d", i+1)
		}
	}

	return nil
}

// In devuelve true si el valor `v` se encuentra dentro del conjunto `set`.
// Admite cualquier tipo comparable (string, int, float64, etc.).
func In[T comparable](v T, set ...T) bool {
	return slices.Contains(set, v)
}

// IsNumeric Devuelve true si el valor `v` es numérico.
func IsNumeric(v string) bool {
	v = strings.TrimSpace(v)
	if v == "" {
		return false
	}
	for _, r := range v {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}
