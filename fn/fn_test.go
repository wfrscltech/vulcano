package fn

import (
	"testing"
)

func TestTernaryIf(t *testing.T) {
	t.Run("bool type", func(t *testing.T) {
		tests := []struct {
			name      string
			condition bool
			trueVal   bool
			falseVal  bool
			want      bool
		}{
			{
				name:      "condición verdadera",
				condition: true,
				trueVal:   true,
				falseVal:  false,
				want:      true,
			},
			{
				name:      "condición falsa",
				condition: false,
				trueVal:   true,
				falseVal:  false,
				want:      false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := TernaryIf(tt.condition, tt.trueVal, tt.falseVal); got != tt.want {
					t.Errorf("TernaryIf() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("int type", func(t *testing.T) {
		tests := []struct {
			name      string
			condition bool
			trueVal   int
			falseVal  int
			want      int
		}{
			{
				name:      "condición verdadera",
				condition: true,
				trueVal:   42,
				falseVal:  0,
				want:      42,
			},
			{
				name:      "condición falsa",
				condition: false,
				trueVal:   42,
				falseVal:  0,
				want:      0,
			},
			{
				name:      "números negativos",
				condition: true,
				trueVal:   -5,
				falseVal:  10,
				want:      -5,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := TernaryIf(tt.condition, tt.trueVal, tt.falseVal); got != tt.want {
					t.Errorf("TernaryIf() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("string type", func(t *testing.T) {
		tests := []struct {
			name      string
			condition bool
			trueVal   string
			falseVal  string
			want      string
		}{
			{
				name:      "condición verdadera",
				condition: true,
				trueVal:   "yes",
				falseVal:  "no",
				want:      "yes",
			},
			{
				name:      "condición falsa",
				condition: false,
				trueVal:   "yes",
				falseVal:  "no",
				want:      "no",
			},
			{
				name:      "strings vacíos",
				condition: true,
				trueVal:   "",
				falseVal:  "default",
				want:      "",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := TernaryIf(tt.condition, tt.trueVal, tt.falseVal); got != tt.want {
					t.Errorf("TernaryIf() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("float64 type", func(t *testing.T) {
		tests := []struct {
			name      string
			condition bool
			trueVal   float64
			falseVal  float64
			want      float64
		}{
			{
				name:      "condición verdadera",
				condition: true,
				trueVal:   3.14,
				falseVal:  2.71,
				want:      3.14,
			},
			{
				name:      "condición falsa",
				condition: false,
				trueVal:   3.14,
				falseVal:  2.71,
				want:      2.71,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := TernaryIf(tt.condition, tt.trueVal, tt.falseVal); got != tt.want {
					t.Errorf("TernaryIf() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("int64 type", func(t *testing.T) {
		tests := []struct {
			name      string
			condition bool
			trueVal   int64
			falseVal  int64
			want      int64
		}{
			{
				name:      "condición verdadera",
				condition: true,
				trueVal:   9223372036854775807,
				falseVal:  0,
				want:      9223372036854775807,
			},
			{
				name:      "condición falsa",
				condition: false,
				trueVal:   9223372036854775807,
				falseVal:  -1,
				want:      -1,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := TernaryIf(tt.condition, tt.trueVal, tt.falseVal); got != tt.want {
					t.Errorf("TernaryIf() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("uint type", func(t *testing.T) {
		tests := []struct {
			name      string
			condition bool
			trueVal   uint
			falseVal  uint
			want      uint
		}{
			{
				name:      "condición verdadera",
				condition: true,
				trueVal:   100,
				falseVal:  200,
				want:      100,
			},
			{
				name:      "condición falsa",
				condition: false,
				trueVal:   100,
				falseVal:  200,
				want:      200,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := TernaryIf(tt.condition, tt.trueVal, tt.falseVal); got != tt.want {
					t.Errorf("TernaryIf() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("rune type", func(t *testing.T) {
		tests := []struct {
			name      string
			condition bool
			trueVal   rune
			falseVal  rune
			want      rune
		}{
			{
				name:      "condición verdadera",
				condition: true,
				trueVal:   'A',
				falseVal:  'B',
				want:      'A',
			},
			{
				name:      "condición falsa",
				condition: false,
				trueVal:   'X',
				falseVal:  'Y',
				want:      'Y',
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if got := TernaryIf(tt.condition, tt.trueVal, tt.falseVal); got != tt.want {
					t.Errorf("TernaryIf() = %v, want %v", got, tt.want)
				}
			})
		}
	})

	t.Run("custom type based on native", func(t *testing.T) {
		type MyInt int
		type MyString string

		intTests := []struct {
			name      string
			condition bool
			trueVal   MyInt
			falseVal  MyInt
			want      MyInt
		}{
			{
				name:      "custom int true",
				condition: true,
				trueVal:   MyInt(10),
				falseVal:  MyInt(20),
				want:      MyInt(10),
			},
			{
				name:      "custom int false",
				condition: false,
				trueVal:   MyInt(10),
				falseVal:  MyInt(20),
				want:      MyInt(20),
			},
		}

		for _, tt := range intTests {
			t.Run(tt.name, func(t *testing.T) {
				if got := TernaryIf(tt.condition, tt.trueVal, tt.falseVal); got != tt.want {
					t.Errorf("TernaryIf() = %v, want %v", got, tt.want)
				}
			})
		}

		stringTests := []struct {
			name      string
			condition bool
			trueVal   MyString
			falseVal  MyString
			want      MyString
		}{
			{
				name:      "custom string true",
				condition: true,
				trueVal:   MyString("hello"),
				falseVal:  MyString("world"),
				want:      MyString("hello"),
			},
			{
				name:      "custom string false",
				condition: false,
				trueVal:   MyString("hello"),
				falseVal:  MyString("world"),
				want:      MyString("world"),
			},
		}

		for _, tt := range stringTests {
			t.Run(tt.name, func(t *testing.T) {
				if got := TernaryIf(tt.condition, tt.trueVal, tt.falseVal); got != tt.want {
					t.Errorf("TernaryIf() = %v, want %v", got, tt.want)
				}
			})
		}
	})
}
