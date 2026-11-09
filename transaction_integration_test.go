//go:build integration

package dyalpm

import (
	"errors"
	"testing"

	"github.com/Jguer/dyalpm/internal/dyerrors"
)

// Removing glibc breaks most installed packages, so prepare must fail and report them.
func TestTransaction_PrepareReturnsMissingDeps(t *testing.T) {
	h := mustInitializeTestHandle(t)
	glibc := mustLocalDB(t, h).Pkg("glibc")
	if glibc == nil {
		t.Skip("glibc must be installed")
	}

	trans := NewTransaction(h)
	if err := trans.Init(TransFlagNoLock); err != nil {
		t.Fatalf("Init failed: %v", err)
	}
	defer func() { _ = trans.Release() }()

	if err := trans.RemovePkg(glibc); err != nil {
		t.Fatalf("RemovePkg failed: %v", err)
	}

	missing, err := trans.Prepare()
	if !errors.Is(err, ErrTransactionPrepareFailed) || !errors.Is(err, dyerrors.ErrUnsatisfiedDeps) {
		t.Fatalf("Prepare error = %v, want prepare failure with unsatisfied deps", err)
	}
	if len(missing) == 0 {
		t.Fatal("Prepare returned no missing dependencies")
	}
	for _, m := range missing {
		if m.GetTarget() == "" || m.GetDepend().GetName() == "" {
			t.Errorf("incomplete DepMissing: target=%q dep=%q", m.GetTarget(), m.GetDepend().GetName())
		}
		m.Free()
	}
}
