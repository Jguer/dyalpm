//go:build integration

package dyalpm

import (
	"io"
	"testing"
)

func TestPackage_MetadataAndFileAccess(t *testing.T) {
	h := mustInitializeTestHandle(t)
	localDB := mustLocalDB(t, h)
	pkg := mustInstalledPkg(t, h, "go", "glibc", "pacman", "bash", "coreutils")
	pkgImpl, ok := pkg.(*package_)
	if !ok || pkgImpl == nil {
		t.Fatalf("expected package implementation to be *package_, got %T", pkg)
	}

	if pkg.Name() == "" {
		t.Fatalf("package name is empty")
	}
	if pkg.Version() == "" {
		t.Fatalf("package version is empty")
	}
	if pkg.Architecture() == "" {
		t.Fatalf("package architecture is empty")
	}
	if pkg.ISize() <= 0 {
		t.Fatalf("package installed size is not positive: %d", pkg.ISize())
	}
	if pkg.Reason() != PkgReasonExplicit && pkg.Reason() != PkgReasonDepend && pkg.Reason() != PkgReasonUnknown {
		t.Fatalf("unexpected package reason: %v", pkg.Reason())
	}
	if db := pkg.DB(); db == nil || db.Name() != localDB.Name() {
		t.Fatalf("package database %q does not match expected local database", pkg.Name())
	}

	_ = pkg.Description()
	_ = pkg.Base()
	_ = pkgImpl.Licenses()
	_ = pkg.Groups()
	depends := pkg.Depends()
	if len(depends) == 0 {
		t.Logf("package %q has no declared dependencies", pkg.Name())
	}

	provides := pkg.Provides()
	if len(provides) == 0 {
		t.Fatalf("expected package %q to expose provides", pkg.Name())
	}
	_ = pkg.OptionalDepends()
	_ = pkgImpl.XData()
	_ = pkgImpl.Base()

	files := pkgImpl.Files()
	if len(files) == 0 {
		t.Skipf("package %q has no files for filelist assertions", pkg.Name())
	}

	firstFile := files[0].Name()
	if !pkgImpl.Contains(firstFile) {
		t.Fatalf("package does not contain expected file: %s", firstFile)
	}
	if fileInfo, err := pkgImpl.ContainsFile(firstFile); err != nil {
		t.Fatalf("ContainsFile(%q) failed: %v", firstFile, err)
	} else if fileInfo == nil {
		t.Fatalf("ContainsFile(%q) returned nil file info", firstFile)
	}

	// Ignore behavior should be updated through handle option helpers.
	if err := h.SetIgnorePkgs([]string{pkg.Name()}); err != nil {
		t.Fatalf("failed to set ignore packages: %v", err)
	}
	if !pkg.ShouldIgnore() {
		t.Fatalf("package %q should be marked ignored after setting ignore list", pkg.Name())
	}

	required, err := pkgImpl.ComputeRequiredBy()
	if err != nil {
		t.Fatalf("ComputeRequiredBy() failed: %v", err)
	}
	optional, err := pkgImpl.ComputeOptionalFor()
	if err != nil {
		t.Fatalf("ComputeOptionalFor() failed: %v", err)
	}
	_ = required
	_ = optional
}

func TestPackage_MoreValueMethods(t *testing.T) {
	h := mustInitializeTestHandle(t)
	pkg := mustInstalledPkg(t, h, "go", "glibc", "pacman", "bash", "coreutils")
	pkgImpl, ok := pkg.(*package_)
	if !ok || pkgImpl == nil {
		t.Fatalf("expected package implementation to be *package_, got %T", pkg)
	}

	if size := pkgImpl.Size(); size < 0 {
		t.Fatalf("package size should be non-negative: %d", size)
	}

	if size := pkgImpl.DownloadSize(); size < 0 {
		t.Fatalf("package download size should be non-negative: %d", size)
	}

	_ = pkgImpl.CheckDepends()
	_ = pkgImpl.MakeDepends()
	replaces := pkgImpl.Replaces()
	if len(replaces) == 0 {
		t.Logf("package %q has no replacements", pkg.Name())
	}
	_ = pkgImpl.Conflicts()
	base := pkgImpl.Base()
	if base == "" {
		t.Fatalf("expected package base to be non-empty for %q", pkg.Name())
	}
	_ = pkgImpl.FileName()
	_ = pkgImpl.SHA256Sum()
	_ = pkgImpl.Packager()
	_ = pkgImpl.URL()
	_ = pkgImpl.Base64Sig()
	_ = pkgImpl.Base64Signature()
	_ = pkgImpl.PkgValidation()
	_ = pkgImpl.Validation()
	if native := pkgImpl.NativeHandle(); native == nil {
		t.Logf("native handle is nil for %s", pkg.Name())
	}
	_ = pkgImpl.Origin()
	_ = pkgImpl.HasScriptlet()
	_ = pkgImpl.Type()
	backups := pkgImpl.Backup()
	for _, b := range backups {
		_ = b.Name()
		_ = b.Hash()
	}
	_ = pkgImpl.BuildDate()
	_ = pkgImpl.InstallDate()

	if err := pkgImpl.CheckMD5Sum(); err != nil {
		t.Logf("CheckMD5Sum returned error for %s: %v", pkg.Name(), err)
	}
	if _, err := pkgImpl.CheckPGPSignature(); err != nil {
		t.Logf("CheckPGPSignature returned error for %s: %v", pkg.Name(), err)
	}

	changelog, err := pkgImpl.Changelog()
	if err == nil && changelog != nil {
		buf := make([]byte, 64)
		n, readErr := changelog.Read(buf)
		if readErr != nil && readErr != io.EOF {
			t.Fatalf("changelog read failed for %s: %v", pkg.Name(), readErr)
		}
		if n > 0 {
			t.Logf("read %d bytes from changelog for %s", n, pkg.Name())
		}
		if err := changelog.Close(); err != nil {
			t.Fatalf("changelog close failed: %v", err)
		}
	}

	syncDB := mustLocalDB(t, h)
	newPkg := pkgImpl.SyncGetNewVersion([]Database{syncDB})
	if newPkg == nil {
		t.Logf("SyncGetNewVersion returned nil package")
	}
	aliasedPkg := pkgImpl.SyncNewVersion([]Database{syncDB})
	if newPkg != nil && aliasedPkg == nil {
		t.Fatalf("SyncNewVersion should mirror SyncGetNewVersion behavior")
	}

	if signatures, err := pkgImpl.CheckPGPSignature(); err != nil {
		t.Logf("CheckPGPSignature returned error: %v", err)
	} else if len(signatures.Results) == 0 {
		t.Logf("CheckPGPSignature returned no signatures")
	}

	sigBytes, sigErr := pkgImpl.Sig()
	if sigErr != nil {
		t.Logf("Sig() returned error for %s: %v", pkg.Name(), sigErr)
	}
	if len(sigBytes) == 0 {
		t.Logf("Sig() returned no bytes for %s", pkg.Name())
	}

}

func TestPackage_NilPointerSafety(t *testing.T) {
	var p package_

	if p.Size() != 0 {
		t.Fatalf("expected zero-value package size to be 0")
	}
	if got := p.Conflicts(); len(got) != 0 {
		t.Fatalf("expected zero conflicts for zero-value package, got %#v", got)
	}
	if p.MakeDepends() == nil {
		// nil slice is acceptable; assert no panic only.
	}
	if p.Replaces() == nil {
		// nil slice is acceptable; assert no panic only.
	}
	if p.Base() != "" {
		t.Fatalf("expected empty base for zero-value package")
	}
	if p.FileName() != "" {
		t.Fatalf("expected empty file name for zero-value package")
	}
	if p.PkgValidation() != PkgValidationUnknown {
		t.Fatalf("expected zero-value validation")
	}
	if err := p.CheckMD5Sum(); err == nil {
		t.Fatalf("expected CheckMD5Sum error for zero-value package")
	}
	if _, err := p.CheckPGPSignature(); err == nil {
		t.Fatalf("expected CheckPGPSignature error for zero-value package")
	}
	if p.NativeHandle() != nil {
		t.Fatalf("expected zero-value native handle to be nil")
	}
	if p.Contains("does-not-exist") {
		t.Fatalf("expected zero-value package to have no files")
	}
}
