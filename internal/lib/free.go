package lib

// Free releases memory allocated by libc malloc.
func Free(ptr uintptr) {
	if ptr == 0 {
		return
	}
	if LibcFree == nil {
		return
	}
	LibcFree(ptr)
}
