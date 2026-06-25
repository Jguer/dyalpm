//go:build integration

package dyalpm

import (
	"os"
	"slices"
	"strings"
	"testing"
)

const (
	integrationRootEnv    = "DYALPM_TEST_ROOT"
	integrationDBPathEnv  = "DYALPM_TEST_DBPATH"
	integrationPackageEnv = "DYALPM_TEST_PACKAGE"
)

func mustInitializeTestHandle(t *testing.T) Handle {
	t.Helper()

	root := os.Getenv(integrationRootEnv)
	if root == "" {
		root = "/"
	}

	dbpath := os.Getenv(integrationDBPathEnv)
	if dbpath == "" {
		dbpath = "/var/lib/pacman"
	}

	h, err := Initialize(root, dbpath)
	if err != nil {
		t.Fatalf("failed to initialize ALPM handle: %v", err)
	}

	t.Cleanup(func() {
		_ = h.Release()
	})

	return h
}

func mustLocalDB(t *testing.T, h Handle) Database {
	t.Helper()

	db, err := h.LocalDB()
	if err != nil {
		t.Fatalf("failed to get local DB: %v", err)
	}

	return db
}

func mustInstalledPkg(t *testing.T, h Handle, names ...string) Package {
	t.Helper()

	localDB := mustLocalDB(t, h)

	candidates := slices.Clone(names)
	if pkgEnv := os.Getenv(integrationPackageEnv); pkgEnv != "" {
		for _, name := range strings.Split(pkgEnv, ",") {
			candidates = append(candidates, strings.TrimSpace(name))
		}
	}
	candidates = append(candidates, "pacman", "bash", "coreutils", "glibc", "linux")

	seen := map[string]struct{}{}
	for _, name := range candidates {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}
		if _, exists := seen[name]; exists {
			continue
		}
		seen[name] = struct{}{}

		pkg := localDB.Pkg(name)
		if pkg != nil {
			return pkg
		}
	}

	packages := localDB.PkgCache().Collect()
	if len(packages) == 0 {
		t.Skip("local DB has no packages to test against")
	}

	for _, pkg := range packages {
		if pkg != nil && pkg.Name() != "" {
			return pkg
		}
	}

	t.Skip("no installed package could be discovered in local DB")
	return nil
}

func requireDownloaderCapability(t *testing.T) {
	t.Helper()

	caps, err := Capabilities()
	if err != nil {
		t.Skipf("failed to read capabilities: %v", err)
	}

	if caps&CapDownloader == 0 {
		t.Skip("alpm downloader capability is not available")
	}
}
