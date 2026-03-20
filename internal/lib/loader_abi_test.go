package lib

import "testing"

func TestAlpmFuncResolution(t *testing.T) {
	if err := EnsureAlpmLoaded(); err != nil {
		t.Skipf("libalpm not available: %v", err)
	}
	if AlpmVersion == nil {
		t.Error("expected alpm_version to resolve")
	}
	if AlpmRelease == nil {
		t.Error("expected alpm_release to resolve")
	}
	if AlpmPkgGetName == nil {
		t.Error("expected alpm_pkg_get_name to resolve")
	}
}
