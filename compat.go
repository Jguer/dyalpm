package dyalpm

import (
	stderrors "errors"
	"unsafe"
)

// LogLevel represents logging level
type LogLevel uint

const (
	LogError   LogLevel = 1
	LogWarning LogLevel = 1 << 1
	LogDebug   LogLevel = 1 << 2
	LogFunc    LogLevel = 1 << 3
)

// Question types (alpm_question_type_t). Match these against Question.Type.
const (
	QuestionTypeInstallIgnorepkg int32 = 1 << iota
	QuestionTypeReplacePkg
	QuestionTypeConflictPkg
	QuestionTypeCorruptedPkg
	QuestionTypeRemovePkgs
	QuestionTypeSelectProvider
	QuestionTypeImportKey
)

// QuestionAny is a wrapper for question callback data (go-alpm/v2 compatibility)
type QuestionAny struct {
	Question Question
}

// QuestionInstallIgnorepkg returns a QuestionInstallIgnorepkg if this is an install-ignorepkg question
func (qa QuestionAny) QuestionInstallIgnorepkg() (*QuestionInstallIgnorepkg, error) {
	if qa.Question.Type == QuestionTypeInstallIgnorepkg {
		return &QuestionInstallIgnorepkg{q: qa.Question}, nil
	}
	return nil, stderrors.New("not an install ignorepkg question")
}

// QuestionSelectProvider returns a QuestionSelectProvider if this is a select-provider question.
// The library only decodes the question; choosing a provider is up to the application.
func (qa QuestionAny) QuestionSelectProvider() (*QuestionSelectProvider, error) {
	if qa.Question.Type == QuestionTypeSelectProvider {
		return &QuestionSelectProvider{q: qa.Question}, nil
	}
	return nil, stderrors.New("not a select provider question")
}

// QuestionInstallIgnorepkg is for handling install-ignorepkg questions
type QuestionInstallIgnorepkg struct {
	q Question
}

// SetInstall sets whether to install the ignored package
func (qi *QuestionInstallIgnorepkg) SetInstall(install bool) {
	answer := 0
	if install {
		answer = 1
	}
	qi.q.SetAnswerInt(answer)
}

// QuestionSelectProvider is for handling select-provider questions.
// It mirrors alpm_question_select_provider_t { type; use_index; alpm_list_t *providers; alpm_depend_t *depend; }.
type QuestionSelectProvider struct {
	q Question
}

func (qp *QuestionSelectProvider) field(offset uintptr) uintptr {
	if qp.q.Ptr == 0 {
		return 0
	}
	return *(*uintptr)(unsafe.Add(unsafe.Pointer(qp.q.Ptr), offset))
}

// Dep returns the dependency string being resolved, e.g. "java-runtime>=17".
func (qp *QuestionSelectProvider) Dep() string {
	return newDependency(qp.field(16)).ComputeString()
}

// Providers returns the candidate packages, in the order SetUseIndex refers to.
func (qp *QuestionSelectProvider) Providers(h Handle) PackageIterator {
	hImpl, _ := h.(*handle)
	listPtr := qp.field(8)
	if listPtr == 0 {
		return PackageIterator{}
	}
	return newPackageIterator(listPtr, hImpl)
}

// SetUseIndex sets which provider to use (0-based index into Providers).
func (qp *QuestionSelectProvider) SetUseIndex(index int) {
	qp.q.SetAnswerInt(index)
}
