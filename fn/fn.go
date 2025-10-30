package fn

// NativeData Define los tipos nativos de Go para el uso de la función TernaryIf
type NativeData interface {
	~bool | ~rune | ~int | ~int8 | ~int16 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 |
		~uint64 | ~uintptr | ~float32 | ~float64 | ~complex64 | ~complex128 | ~string
}

// TernaryIf Comparador if ternario de la forma c ? rt : rf, solo para tipos nativos y sus derivados
func TernaryIf[T NativeData](c bool, rt, rf T) T {
	if c {
		return rt
	} else {
		return rf
	}
}
