//go:build integration

package dyalpm_test

import (
	"testing"

	alpm "github.com/Jguer/dyalpm"
)

func TestInitializeAndRelease(t *testing.T) {
	h, err := alpm.Initialize("/", "/var/lib/pacman")
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	defer h.Release()

	// Test that handle methods work
	root := h.Root()
	if root != "/" {
		t.Errorf("Expected root '/', got '%s'", root)
	}

	dbpath := h.DBPath()
	if dbpath != "/var/lib/pacman/" {
		t.Errorf("Expected dbpath '/var/lib/pacman/', got '%s'", dbpath)
	}
}

// A failed lookup leaves a sticky errno; Release must still succeed and not run twice.
func TestReleaseAfterFailedLookup(t *testing.T) {
	h, err := alpm.Initialize("/", "/var/lib/pacman")
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	localDB, err := h.LocalDB()
	if err != nil {
		t.Fatalf("Failed to get local DB: %v", err)
	}
	if localDB.Pkg("dyalpm-no-such-package") != nil {
		t.Fatal("expected missing package")
	}
	if err := h.Release(); err != nil {
		t.Fatalf("Release failed: %v", err)
	}
	if err := h.Release(); err == nil {
		t.Fatal("second Release should report already released")
	}
}

func TestLocalDB(t *testing.T) {
	h, err := alpm.Initialize("/", "/var/lib/pacman")
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	defer h.Release()

	localDB, err := h.LocalDB()
	if err != nil {
		t.Fatalf("Failed to get local DB: %v", err)
	}

	name := localDB.Name()
	if name != "local" {
		t.Errorf("Expected 'local', got '%s'", name)
	}
}

func TestPkgCache(t *testing.T) {
	h, err := alpm.Initialize("/", "/var/lib/pacman")
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	defer h.Release()

	localDB, err := h.LocalDB()
	if err != nil {
		t.Fatalf("Failed to get local DB: %v", err)
	}

	count := 0
	err = localDB.PkgCache().ForEach(func(pkg alpm.Package) error {
		if pkg.Name() == "" {
			t.Error("Package name is empty")
		}
		if pkg.Version() == "" {
			t.Error("Package version is empty")
		}
		count++
		return nil
	})
	if err != nil {
		t.Errorf("ForEach failed: %v", err)
	}

	if count == 0 {
		t.Error("Expected at least one package in local DB")
	}
	t.Logf("Found %d packages in local DB", count)
}

func TestSyncDBs(t *testing.T) {
	h, err := alpm.Initialize("/", "/var/lib/pacman")
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	defer h.Release()

	// Register a sync DB
	coreDB, err := h.RegisterSyncDB("core", 0)
	if err != nil {
		t.Fatalf("Failed to register core DB: %v", err)
	}

	if coreDB.Name() != "core" {
		t.Errorf("Expected 'core', got '%s'", coreDB.Name())
	}

	syncDBs, err := h.SyncDBs()
	if err != nil {
		t.Fatalf("Failed to get sync DBs: %v", err)
	}

	if len(syncDBs) == 0 {
		t.Error("Expected at least one sync DB")
	}
}

func TestVersionCompare(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected int
	}{
		{"1.0", "2.0", -1},
		{"2.0", "1.0", 1},
		{"1.0", "1.0", 0},
		{"1.0-1", "1.0-2", -1},
		{"1.0.1", "1.0", 1},
		{"1.0a", "1.0b", -1},
	}

	for _, tt := range tests {
		result := alpm.VerCmp(tt.v1, tt.v2)
		// Normalize to -1, 0, 1
		if result < 0 {
			result = -1
		} else if result > 0 {
			result = 1
		}
		if result != tt.expected {
			t.Errorf("VerCmp(%q, %q) = %d, want %d", tt.v1, tt.v2, result, tt.expected)
		}
	}
}

func TestPackageDetails(t *testing.T) {
	h, err := alpm.Initialize("/", "/var/lib/pacman")
	if err != nil {
		t.Fatalf("Failed to initialize: %v", err)
	}
	defer h.Release()

	localDB, err := h.LocalDB()
	if err != nil {
		t.Fatalf("Failed to get local DB: %v", err)
	}

	// Get glibc package - should always be installed
	pkg := localDB.Pkg("glibc")
	if pkg == nil {
		t.Skip("glibc not installed, skipping detailed package test")
	}

	// Test various package methods
	if pkg.Name() != "glibc" {
		t.Errorf("Expected name 'glibc', got '%s'", pkg.Name())
	}

	if pkg.Version() == "" {
		t.Error("Expected non-empty version")
	}

	if pkg.Architecture() == "" {
		t.Error("Expected non-empty architecture")
	}

	if pkg.ISize() <= 0 {
		t.Error("Expected positive installed size")
	}

	deps := pkg.Depends()
	// glibc should have some dependencies
	t.Logf("glibc has %d dependencies", len(deps))
}
