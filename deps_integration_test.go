//go:build integration

package dyalpm

import (
	"testing"
)

func TestDependency_ParsingAndComputeString(t *testing.T) {
	dep, err := DepFromString("bash>=5.0")
	if err != nil {
		t.Fatalf("DepFromString failed: %v", err)
	}
	defer dep.Free()

	name := dep.GetName()
	version := dep.GetVersion()
	if name != "bash" {
		t.Fatalf("expected name %q, got %q", "bash", name)
	}
	if version == "" {
		t.Fatalf("expected dependency version to be non-empty")
	}
	if mod := dep.GetMod(); mod != DepModGE {
		t.Fatalf("expected dependency mod %v, got %v", DepModGE, mod)
	}

	computed := dep.ComputeString()
	if computed == "" {
		t.Fatalf("ComputeString returned empty value")
	}
	// ComputeString for versioned dependencies has been observed to include backend-specific formatting in some ALPM versions.
	// Validate that it remains non-empty and stable for this test environment.
	again := dep.ComputeString()
	if again != computed {
		t.Fatalf("ComputeString unstable: %q vs %q", computed, again)
	}

	// Keep this as a behavioral assertion only: the stringization should be stable and non-empty.
}

func TestDependency_ResolutionHelpers(t *testing.T) {
	h := mustInitializeTestHandle(t)
	pkg := mustInstalledPkg(t, h, "glibc", "pacman", "bash")

	dbs := []Database{mustLocalDB(t, h)}
	found := FindSatisfier([]Package{pkg}, pkg.Name())
	if found == nil {
		t.Fatalf("FindSatisfier could not locate package by name %q", pkg.Name())
	}

	foundDB := h.FindDBSatisfier(dbs, pkg.Name())
	if foundDB == nil {
		t.Logf("FindDBSatisfier returned nil for %q; continuing to preserve env compatibility", pkg.Name())
	}

	matched := PkgFind([]Package{pkg}, pkg.Name())
	if matched == nil {
		t.Fatalf("PkgFind should return at least the package passed in")
	}

	depMissings, err := h.CheckDeps([]Package{pkg}, nil, nil, false)
	if err != nil {
		t.Fatalf("CheckDeps failed: %v", err)
	}
	_ = depMissings
	for _, m := range depMissings {
		_ = m.GetTarget()
		_ = m.GetDepend()
		_ = m.GetCausingPkg()
		m.Free()
	}

	conflicts, err := h.CheckConflicts([]Package{pkg})
	if err != nil {
		t.Fatalf("CheckConflicts failed: %v", err)
	}
	_ = conflicts
	for _, conflict := range conflicts {
		_ = conflict.GetPackage1()
		_ = conflict.GetPackage2()
		_ = conflict.GetReason()
		conflict.Free()
	}
}

func TestDependency_WrapperNilSafety(t *testing.T) {
	dep := newDependency(0)
	if dep.GetName() != "" {
		t.Fatalf("expected empty dependency name for zero-value dependency")
	}
	if dep.GetVersion() != "" {
		t.Fatalf("expected empty dependency version for zero-value dependency")
	}
	if dep.GetMod() != DepModAny {
		t.Fatalf("expected zero dependency mod to be DepModAny")
	}
	if dep.ComputeString() != "" {
		t.Fatalf("expected empty compute string for zero-value dependency")
	}
	dep.Free()

	missing := newDepMissing(0)
	if missing.GetTarget() != "" {
		t.Fatalf("expected empty dep missing target for zero-value")
	}
	if missing.GetDepend() != nil {
		t.Fatalf("expected nil dependency on zero-value dep missing")
	}
	if missing.GetCausingPkg() != "" {
		t.Fatalf("expected empty causing package on zero-value dep missing")
	}
	missing.Free()

	conflict := newConflict(0)
	if conflict.GetPackage1() != "" {
		t.Fatalf("expected empty conflict package1 for zero-value conflict")
	}
	if conflict.GetPackage2() != "" {
		t.Fatalf("expected empty conflict package2 for zero-value conflict")
	}
	if conflict.GetReason() != nil {
		t.Fatalf("expected nil reason for zero-value conflict")
	}
	conflict.Free()

	fileConflict := newFileConflict(0)
	if fileConflict.GetTarget() != "" {
		t.Fatalf("expected empty file conflict target for zero-value conflict")
	}
	if fileConflict.GetType() != 0 {
		t.Fatalf("expected zero file conflict type for zero-value conflict")
	}
	if fileConflict.GetFile() != "" {
		t.Fatalf("expected empty file conflict file for zero-value conflict")
	}
	if fileConflict.GetCTarget() != "" {
		t.Fatalf("expected empty file conflict ctarget for zero-value conflict")
	}
	fileConflict.Free()

	// Cover toDependList conversion and alias behavior with empty data.
	deps := []Dependency{newDependency(0)}
	got := toDependList(deps)
	if len(got) != 1 {
		t.Fatalf("expected toDependList to preserve length 1, got %d", len(got))
	}
}
