package dyalpm

import (
	"runtime"
	"testing"
	"unsafe"
)

// These tests build libalpm structs in Go memory and check the decoders read the right fields.
// They need no libalpm, so they catch layout drift on every plain `go test`.

func cstr(s string) (uintptr, []byte) {
	b := append([]byte(s), 0)
	return uintptr(unsafe.Pointer(&b[0])), b
}

func TestLayout_Depend(t *testing.T) {
	if size := unsafe.Sizeof(cDepend{}); size != 40 {
		t.Fatalf("alpm_depend_t is 40 bytes on 64-bit, mirror is %d", size)
	}
	name, nb := cstr("python")
	ver, vb := cstr("3.12")
	desc, db := cstr("for scripts")
	// Non-zero padding after mod must not leak into Mod.
	raw := [5]uint64{uint64(name), uint64(ver), uint64(desc), 42, uint64(DepModGE) | 0xdead<<32}

	got := newDependency(uintptr(unsafe.Pointer(&raw))).toDepend()
	runtime.KeepAlive(nb)
	runtime.KeepAlive(vb)
	runtime.KeepAlive(db)

	want := Depend{Name: "python", Version: "3.12", Description: "for scripts", NameHash: 42, Mod: DepModGE}
	if got != want {
		t.Fatalf("toDepend = %+v, want %+v", got, want)
	}
}

func TestLayout_DepMissing(t *testing.T) {
	target, tb := cstr("foo")
	causing, cb := cstr("bar")
	dep := [5]uint64{}
	raw := [3]uintptr{target, uintptr(unsafe.Pointer(&dep)), causing}

	m := newDepMissing(uintptr(unsafe.Pointer(&raw)))
	if m.GetTarget() != "foo" || m.GetCausingPkg() != "bar" {
		t.Fatalf("depmissing decoded wrong: target=%q causing=%q", m.GetTarget(), m.GetCausingPkg())
	}
	if m.GetDepend().(*dependency).ptr != uintptr(unsafe.Pointer(&dep)) {
		t.Fatal("depmissing depend pointer read from wrong offset")
	}
	runtime.KeepAlive(tb)
	runtime.KeepAlive(cb)
}

func TestLayout_FileConflict(t *testing.T) {
	target, tb := cstr("foo")
	file, fb := cstr("/usr/bin/foo")
	ctarget, cb := cstr("bar")
	// type is a 4-byte enum; the next 4 bytes are padding and must be ignored.
	raw := [4]uint64{uint64(target), uint64(FileConflictFilesystem) | 0xdead<<32, uint64(file), uint64(ctarget)}

	f := newFileConflict(uintptr(unsafe.Pointer(&raw)))
	if f.GetTarget() != "foo" || f.GetType() != FileConflictFilesystem || f.GetFile() != "/usr/bin/foo" || f.GetCTarget() != "bar" {
		t.Fatalf("fileconflict decoded wrong: %q %d %q %q", f.GetTarget(), f.GetType(), f.GetFile(), f.GetCTarget())
	}
	runtime.KeepAlive(tb)
	runtime.KeepAlive(fb)
	runtime.KeepAlive(cb)
}
