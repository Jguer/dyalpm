package lib

import "github.com/ebitengine/purego"

func tryRegister[T any](fn *T, lib uintptr, name string) {
	symbol, err := purego.Dlsym(lib, name)
	if err != nil {
		return
	}
	purego.RegisterFunc(fn, symbol)
}

func tryDlsym(dest *uintptr, lib uintptr, name string) {
	symbol, err := purego.Dlsym(lib, name)
	if err != nil {
		return
	}
	*dest = symbol
}
