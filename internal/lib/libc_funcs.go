package lib

import "github.com/ebitengine/purego"

var (
	LibcFree func(ptr uintptr)
	// LibcVsnprintf formats a message into buf using a forwarded va_list (ap).
	// Used to bridge libalpm's log callback, whose final argument is a va_list.
	LibcVsnprintf func(buf uintptr, size uintptr, format uintptr, ap uintptr) int32
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
	tryRegister(&LibcVsnprintf, libc, "vsnprintf")
}
