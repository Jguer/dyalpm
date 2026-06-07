package lib

import "testing"

func TestAlpmFuncResolution(t *testing.T) {
	if err := EnsureAlpmLoaded(); err != nil {
		t.Skipf("libalpm not available: %v", err)
	}
	if AlpmVersion == nil {
		t.Error("expected alpm_version to resolve")
	}
	VersionPtr := AlpmVersion()
	if VersionPtr == 0 {
		t.Error("expected alpm_version to return a non-empty pointer")
	} else {
		version := PtrToString(VersionPtr)
		if version == "" {
			t.Error("expected alpm_version to return a non-empty version string")
		}
	}
	if AlpmRelease == nil {
		t.Error("expected alpm_release to resolve")
	}
	if AlpmPkgGetName == nil {
		t.Error("expected alpm_pkg_get_name to resolve")
	}
	if AlpmInitialize == nil {
		t.Error("expected alpm_initialize to resolve")
	}
	if AlpmGetLocalDB == nil {
		t.Error("expected alpm_get_localdb to resolve")
	}
	if AlpmGetSyncDBS == nil {
		t.Error("expected alpm_get_syncdbs to resolve")
	}
	if AlpmOptionGetRoot == nil {
		t.Error("expected alpm_option_get_root to resolve")
	}
	if AlpmOptionSetLogfile == nil {
		t.Error("expected alpm_option_set_logfile to resolve")
	}
	if AlpmOptionGetLogfile == nil {
		t.Error("expected alpm_option_get_logfile to resolve")
	}
	if AlpmOptionSetNoUpgrades == nil {
		t.Error("expected alpm_option_set_noupgrades to resolve")
	}
	if AlpmOptionMatchNoUpgrade == nil {
		t.Error("expected alpm_option_match_noupgrade to resolve")
	}
	if AlpmPkgLoad == nil {
		t.Error("expected alpm_pkg_load to resolve")
	}
	if AlpmDBCheckPGPSignature == nil {
		t.Error("expected alpm_db_check_pgp_signature to resolve")
	}
	if AlpmPkgCheckPGPSignature == nil {
		t.Error("expected alpm_pkg_check_pgp_signature to resolve")
	}
	if LibcFree == nil {
		t.Error("expected libc free to resolve")
	}
	if LibcVsnprintf == nil {
		t.Error("expected libc vsnprintf to resolve")
	}
	if AlpmCapabilities == nil {
		t.Error("expected alpm_capabilities to resolve")
	} else {
		caps := AlpmCapabilities()
		if caps == 0 {
			t.Logf("alpm_capabilities returned zero bitmask")
		}
	}
}
