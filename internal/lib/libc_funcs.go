package lib

import "github.com/ebitengine/purego"

var (
	LibcFree func(ptr uintptr)
)

func registerLibcFuncs() {
	libc, err := purego.Dlopen("libc.so.6", purego.RTLD_NOW)
	if err != nil {
		libc, err = purego.Dlopen("libc.so", purego.RTLD_NOW)
		if err != nil {
			return
		}
	}
	tryRegister(&LibcFree, libc, "free")
}
