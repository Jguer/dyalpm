//go:build integration

package dyalpm

import (
	"testing"
	"unsafe"

	alpmlist "github.com/Jguer/dyalpm/internal/list"
)

// cQuestionSelectProvider mirrors alpm_question_select_provider_t.
type cQuestionSelectProvider struct {
	Type      int32
	UseIndex  int32
	Providers uintptr
	Depend    uintptr
}

func TestQuestionSelectProvider_Decode(t *testing.T) {
	h := mustInitializeTestHandle(t)
	localDB := mustLocalDB(t, h)
	pkg1, pkg2 := localDB.Pkg("pacman"), localDB.Pkg("bash")
	if pkg1 == nil || pkg2 == nil {
		t.Skip("pacman and bash must be installed")
	}
	dep, err := DepFromString("sh>=5")
	if err != nil {
		t.Fatalf("DepFromString failed: %v", err)
	}
	defer dep.Free()

	providers := alpmlist.Add(nil, pkg1.(*package_).ptr)
	providers = alpmlist.Add(providers, pkg2.(*package_).ptr)
	defer providers.Free()

	raw := cQuestionSelectProvider{
		Type:      QuestionTypeSelectProvider,
		Providers: providers.Ptr(),
		Depend:    dep.(*dependency).ptr,
	}
	q := Question{Type: raw.Type, Ptr: uintptr(unsafe.Pointer(&raw))}

	if _, err := (QuestionAny{Question: Question{Type: QuestionTypeConflictPkg}}).QuestionSelectProvider(); err == nil {
		t.Fatal("a conflict question must not decode as select-provider")
	}
	qp, err := QuestionAny{Question: q}.QuestionSelectProvider()
	if err != nil {
		t.Fatalf("QuestionSelectProvider failed: %v", err)
	}

	if got := qp.Dep(); got != "sh>=5" {
		t.Errorf("Dep() = %q, want sh>=5", got)
	}
	var names []string
	_ = qp.Providers(h).ForEach(func(p Package) error {
		names = append(names, p.Name())
		return nil
	})
	if len(names) != 2 || names[0] != "pacman" || names[1] != "bash" {
		t.Errorf("Providers() = %v, want [pacman bash]", names)
	}

	qp.SetUseIndex(1)
	if raw.UseIndex != 1 {
		t.Errorf("SetUseIndex wrote %d, want 1", raw.UseIndex)
	}
}
