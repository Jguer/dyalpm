//go:build integration

package dyalpm

import (
	"strings"
	"testing"
)

// TestLogCallbackFuncForwardsVaList verifies that the log callback bridges
// libalpm's va_list argument correctly. alpm_register_syncdb logs
// "registering sync database '%s'\n", so a correctly forwarded va_list must
// produce a message containing the database name.
func TestLogCallbackFuncForwardsVaList(t *testing.T) {
	h := mustInitializeTestHandle(t)

	const dbName = "dyalpm-vararg-probe"

	var (
		gotName bool
		gotAny  bool
	)
	if err := h.SetLogCallbackFunc(func(level LogLevel, msg string) {
		gotAny = true
		if strings.Contains(msg, dbName) {
			gotName = true
		}
	}); err != nil {
		t.Fatalf("SetLogCallbackFunc: %v", err)
	}

	if _, err := h.RegisterSyncDB(dbName, 0); err != nil {
		t.Fatalf("RegisterSyncDB: %v", err)
	}

	if !gotAny {
		t.Fatal("log callback was never invoked")
	}
	if !gotName {
		t.Errorf("log message did not contain formatted db name %q (va_list not forwarded correctly)", dbName)
	}

	// Clearing the callback must not panic and must stop delivery.
	if err := h.SetLogCallbackFunc(nil); err != nil {
		t.Fatalf("clearing log callback: %v", err)
	}
}
