package mistake

import (
	"fmt"
	"net/http"
)

/*
// Errores centinelas
var NotFoundError = errors.New("recurso no encontrado")
var UnauthorizedError = errors.New("no autorizado")
var ForbiddenError = errors.New("no tiene permiso")

// Errores de validación
var RequiredError = errors.New("campo obligatorio")
var InvalidError = errors.New("campo inválido")
var DuplicatedError = errors.New("campo duplicado")

var InternalError = errors.New("error del servidor")
*/

type MistakeCode int

const (
	NotFound MistakeCode = iota
	Unauthorized
	Forbidden
	Required
	Invalid
	Duplicated
	Internal
)

var errorMessages = map[MistakeCode]int{
	NotFound:     http.StatusNotFound,
	Unauthorized: http.StatusUnauthorized,
	Forbidden:    http.StatusForbidden,
	Required:     http.StatusBadRequest,
	Invalid:      http.StatusBadRequest,
	Duplicated:   http.StatusConflict,
	Internal:     http.StatusInternalServerError,
}

type Mistake struct {
	code int
	msg  string
	err  error
	path []string
}

func New(code MistakeCode, msg string, err error, path ...string) *Mistake {
	return &Mistake{
		code: int(code),
		msg:  msg,
		err:  err,
		path: path,
	}
}

func (e Mistake) Error() string {
	return e.msg
}

func (e *Mistake) Code() int {
	if c, ok := errorMessages[MistakeCode(e.code)]; ok {
		return c
	}

	return 500
}

func (e *Mistake) Path(p ...string) {
	e.path = append(e.path, p...)
}

func (e *Mistake) DevError() string {
	return fmt.Sprintf("%s: %s", e.msg, e.err)
}
