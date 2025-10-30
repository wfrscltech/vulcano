package fn

import (
	"testing"
)

func TestMaskDSN(t *testing.T) {
	tests := []struct {
		name string
		dsn  string
		want string
	}{
		{
			name: "DSN con usuario y contraseña",
			dsn:  "postgres://user:password@localhost:5432/mydb",
			want: "postgres://user:*****@localhost:5432/mydb",
		},
		{
			name: "DSN sin contraseña",
			dsn:  "postgres://user@localhost:5432/mydb",
			want: "postgres://user@localhost:5432/mydb",
		},
		{
			name: "DSN sin usuario ni contraseña",
			dsn:  "postgres://localhost:5432/mydb",
			want: "postgres://localhost:5432/mydb",
		},
		{
			name: "DSN con query params",
			dsn:  "postgres://user:password@localhost:5432/mydb?sslmode=disable",
			want: "postgres://user:*****@localhost:5432/mydb?sslmode=disable",
		},
		{
			name: "DSN con contraseña vacía",
			dsn:  "postgres://user:@localhost:5432/mydb",
			want: "postgres://user@localhost:5432/mydb",
		},
		{
			name: "DSN con contraseña compleja",
			dsn:  "mysql://root:P@ssw0rd!123@localhost:3306/database",
			want: "mysql://root:*****@localhost:3306/database",
		},
		{
			name: "DSN inválido sin scheme",
			dsn:  "invalid-dsn-format",
			want: "://invalid-dsn-format",
		},
		{
			name: "DSN de MongoDB",
			dsn:  "mongodb://admin:secret123@cluster0.mongodb.net/test?retryWrites=true",
			want: "mongodb://admin:*****@cluster0.mongodb.net/test?retryWrites=true",
		},
		{
			name: "DSN de Redis",
			dsn:  "redis://user:password@localhost:6379/0",
			want: "redis://user:*****@localhost:6379/0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MaskDSN(tt.dsn); got != tt.want {
				t.Errorf("MaskDSN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSlugify(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "texto simple",
			input: "Hello World",
			want:  "hello-world",
		},
		{
			name:  "texto con acentos",
			input: "Café con Leché",
			want:  "cafe-con-leche",
		},
		{
			name:  "texto con ñ",
			input: "Año Nuevo",
			want:  "ano-nuevo",
		},
		{
			name:  "texto con múltiples espacios",
			input: "Hello    World",
			want:  "hello-world",
		},
		{
			name:  "texto con caracteres especiales",
			input: "Hello@World#2024!",
			want:  "helloworld2024",
		},
		{
			name:  "texto con guiones múltiples",
			input: "Hello---World",
			want:  "hello-world",
		},
		{
			name:  "texto con acentos variados mayúsculas",
			input: "ÀÉÍÓÚ ÄËÏÖÜâêîôû",
			want:  "aeiou-aeiouaeiou",
		},
		{
			name:  "texto con acentos variados minúsculas",
			input: "àéíóú äëïöü âêîôû",
			want:  "aeiou-aeiou-aeiou",
		},
		{
			name:  "texto con slash (se preserva)",
			input: "hello/world",
			want:  "hello/world",
		},
		{
			name:  "texto vacío",
			input: "",
			want:  "",
		},
		{
			name:  "solo espacios",
			input: "   ",
			want:  "-",
		},
		{
			name:  "mezcla compleja",
			input: "El Niño comió ñoquis en 2024!",
			want:  "el-nino-comio-noquis-en-2024",
		},
		{
			name:  "números y letras",
			input: "Test123 456",
			want:  "test123-456",
		},
		{
			name:  "URL-like string",
			input: "api/v1/users",
			want:  "api/v1/users",
		},
		{
			name:  "texto con cedilla",
			input: "Façade",
			want:  "facade",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Slugify(tt.input); got != tt.want {
				t.Errorf("Slugify() = %v, want %v", got, tt.want)
			}
		})
	}
}
