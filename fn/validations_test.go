package fn

import (
	"testing"
)

func TestValidateIntegerList(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "lista válida de enteros",
			value:   "1,2,3,4,5",
			wantErr: false,
		},
		{
			name:    "lista válida con espacios",
			value:   "1, 2, 3, 4, 5",
			wantErr: false,
		},
		{
			name:    "string vacío",
			value:   "",
			wantErr: true,
			errMsg:  "no puede estar vacío",
		},
		{
			name:    "string con solo espacios",
			value:   "   ",
			wantErr: true,
			errMsg:  "no puede estar vacío",
		},
		{
			name:    "valor no numérico",
			value:   "1,abc,3",
			wantErr: true,
			errMsg:  "`abc` no es un número entero",
		},
		{
			name:    "valor decimal",
			value:   "1,2.5,3",
			wantErr: true,
			errMsg:  "`2.5` no es un número entero",
		},
		{
			name:    "valor vacío en posición",
			value:   "1,,3",
			wantErr: true,
			errMsg:  "valor vacío en posición",
		},
		{
			name:    "número negativo válido",
			value:   "-1,-2,3",
			wantErr: false,
		},
		{
			name:    "un solo número",
			value:   "42",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateIntegerList(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateIntegerList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateIntegerList() error message = %v, want to contain %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestValidateStringList(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "lista válida de strings",
			value:   "foo,bar,baz",
			wantErr: false,
		},
		{
			name:    "lista válida con espacios",
			value:   "foo, bar, baz",
			wantErr: false,
		},
		{
			name:    "string vacío",
			value:   "",
			wantErr: true,
			errMsg:  "no puede estar vacío",
		},
		{
			name:    "string con solo espacios",
			value:   "   ",
			wantErr: true,
			errMsg:  "no puede estar vacío",
		},
		{
			name:    "valor vacío en posición",
			value:   "foo,,baz",
			wantErr: true,
			errMsg:  "valor vacío en posición",
		},
		{
			name:    "un solo string",
			value:   "single",
			wantErr: false,
		},
		{
			name:    "strings con mayúsculas (normalizado)",
			value:   "FOO,Bar,BAZ",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStringList(tt.value)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateStringList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("ValidateStringList() error message = %v, want to contain %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

func TestIn(t *testing.T) {
	t.Run("strings", func(t *testing.T) {
		tests := []struct {
			name  string
			value string
			set   []string
			want  bool
		}{
			{
				name:  "valor presente",
				value: "foo",
				set:   []string{"foo", "bar", "baz"},
				want:  true,
			},
			{
				name:  "valor ausente",
				value: "qux",
				set:   []string{"foo", "bar", "baz"},
				want:  false,
			},
			{
				name:  "conjunto vacío",
				value: "foo",
				set:   []string{},
				want:  false,
			},
			{
				name:  "un solo elemento presente",
				value: "foo",
				set:   []string{"foo"},
				want:  true,
			},
			{
				name:  "un solo elemento ausente",
				value: "bar",
				set:   []string{"foo"},
				want:  false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := In(tt.value, tt.set...); got != tt.want {
					t.Errorf("In() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("integers", func(t *testing.T) {
		tests := []struct {
			name  string
			value int
			set   []int
			want  bool
		}{
			{
				name:  "valor presente",
				value: 2,
				set:   []int{1, 2, 3},
				want:  true,
			},
			{
				name:  "valor ausente",
				value: 4,
				set:   []int{1, 2, 3},
				want:  false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := In(tt.value, tt.set...); got != tt.want {
					t.Errorf("In() = %v, want %v", got, tt.want)
				}
			})
		}
	})
}

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  bool
	}{
		{
			name:  "número válido",
			value: "12345",
			want:  true,
		},
		{
			name:  "cero",
			value: "0",
			want:  true,
		},
		{
			name:  "con letras",
			value: "123abc",
			want:  false,
		},
		{
			name:  "con espacios",
			value: "123 456",
			want:  false,
		},
		{
			name:  "string vacío",
			value: "",
			want:  false,
		},
		{
			name:  "solo espacios",
			value: "   ",
			want:  false,
		},
		{
			name:  "número negativo",
			value: "-123",
			want:  false,
		},
		{
			name:  "número decimal",
			value: "123.45",
			want:  false,
		},
		{
			name:  "caracteres especiales",
			value: "123!",
			want:  false,
		},
		{
			name:  "número con espacios al inicio",
			value: "  123",
			want:  true,
		},
		{
			name:  "número con espacios al final",
			value: "123  ",
			want:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsNumeric(tt.value); got != tt.want {
				t.Errorf("IsNumeric() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Función auxiliar para verificar si un string contiene otro
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
