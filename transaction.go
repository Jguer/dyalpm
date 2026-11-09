package dyalpm

import (
	stderrors "errors"

	"github.com/Jguer/dyalpm/internal/dyerrors"
	"github.com/Jguer/dyalpm/internal/lib"
	alpmlist "github.com/Jguer/dyalpm/internal/list"
)

// TransactionFlag represents transaction flags
type TransactionFlag int

const (
	TransFlagNoDeps       TransactionFlag = 1
	TransFlagNoSave       TransactionFlag = 1 << 2
	TransFlagNoDepVersion TransactionFlag = 1 << 3
	TransFlagCascade      TransactionFlag = 1 << 4
	TransFlagRecurse      TransactionFlag = 1 << 5
	TransFlagDBOnly       TransactionFlag = 1 << 6
	TransFlagNoHooks      TransactionFlag = 1 << 7
	TransFlagAllDeps      TransactionFlag = 1 << 8
	TransFlagDownloadOnly TransactionFlag = 1 << 9
	TransFlagNoScriptlet  TransactionFlag = 1 << 10
	TransFlagNoConflicts  TransactionFlag = 1 << 11
	TransFlagNeeded       TransactionFlag = 1 << 13
	TransFlagAllExplicit  TransactionFlag = 1 << 14
	TransFlagUnneeded     TransactionFlag = 1 << 15
	TransFlagRecurseAll   TransactionFlag = 1 << 16
	TransFlagNoLock       TransactionFlag = 1 << 17
)

// Transaction represents an ALPM transaction
type Transaction interface {
	// Init initializes the transaction
	Init(flags TransactionFlag) error

	// Prepare prepares the transaction
	Prepare() ([]DepMissing, error)

	// Commit commits the transaction
	Commit() ([]FileConflict, error)

	// Release releases the transaction
	Release() error

	// AddPkg adds a package to the transaction
	AddPkg(pkg Package) error

	// RemovePkg removes a package from the transaction
	RemovePkg(pkg Package) error

	// SyncSysupgrade adds packages to upgrade to the transaction
	SyncSysupgrade(enableDowngrade bool) error

	// GetFlags returns the transaction flags
	GetFlags() TransactionFlag

	// GetAdd returns the list of packages to be added
	GetAdd() ([]Package, error)

	// GetRemove returns the list of packages to be removed
	GetRemove() ([]Package, error)
}

type transaction struct {
	handle *handle
}

// runWithData calls alpm_trans_prepare or alpm_trans_commit. libalpm fills data
// only on failure, and the item type depends on the errno that was set.
func (t *transaction) runWithData(fn func(uintptr, *uintptr) int32, name string, failErr error) (*alpmlist.List, dyerrors.Errno, error) {
	if fn == nil {
		return nil, dyerrors.ErrOK, stderrors.New("missing function: " + name)
	}

	var data uintptr
	if fn(t.handle.ptr, &data) == 0 {
		return nil, dyerrors.ErrOK, nil
	}
	return alpmlist.NewList(data), t.handle.Errno(), t.handle.transErr(failErr)
}

// freeItems frees each item of a libalpm-owned result list, then the list itself.
func freeItems(l *alpmlist.List, free func(uintptr)) {
	for item := l; item != nil; item = item.Next() {
		free(item.Data())
	}
	l.Free()
}

// NewTransaction creates a new transaction for the given handle
func NewTransaction(h Handle) Transaction {
	handle := h.(*handle)
	return &transaction{
		handle: handle,
	}
}

func (t *transaction) Init(flags TransactionFlag) error {
	return t.handle.TransInit(flags)
}

// Prepare prepares the transaction. On ErrUnsatisfiedDeps it also returns the
// missing dependencies; the caller owns them and should call Free on each.
func (t *transaction) Prepare() ([]DepMissing, error) {
	if t.handle.ptr == 0 {
		return nil, ErrInvalidHandle
	}

	data, errno, err := t.runWithData(lib.AlpmTransPrepare, "alpm_trans_prepare", ErrTransactionPrepareFailed)
	switch errno {
	case dyerrors.ErrUnsatisfiedDeps:
		defer data.Free()
		return collectList(data, func(ptr uintptr) DepMissing { return newDepMissing(ptr) }), err
	case dyerrors.ErrConflictingDeps:
		freeItems(data, func(ptr uintptr) { lib.AlpmConflictFree(ptr) })
	case dyerrors.ErrPkgInvalidArch:
		freeItems(data, lib.Free)
	}
	return []DepMissing{}, err
}

// Commit commits the transaction. On ErrFileConflicts it also returns the
// conflicts; the caller owns them and should call Free on each.
func (t *transaction) Commit() ([]FileConflict, error) {
	if t.handle.ptr == 0 {
		return nil, ErrInvalidHandle
	}

	data, errno, err := t.runWithData(lib.AlpmTransCommit, "alpm_trans_commit", ErrTransactionCommitFailed)
	switch errno {
	case dyerrors.ErrFileConflicts:
		defer data.Free()
		return collectList(data, func(ptr uintptr) FileConflict { return newFileConflict(ptr) }), err
	case dyerrors.ErrPkgInvalid:
		freeItems(data, lib.Free)
	}
	return []FileConflict{}, err
}

func (t *transaction) Release() error {
	return t.handle.TransRelease()
}

func (t *transaction) AddPkg(pkg Package) error {
	if t.handle.ptr == 0 {
		return ErrInvalidHandle
	}

	if lib.AlpmAddPkg == nil {
		return stderrors.New("missing function: alpm_add_pkg")
	}

	pkgImpl := pkg.(*package_)
	result := lib.AlpmAddPkg(t.handle.ptr, pkgImpl.ptr)
	if result != 0 {
		return t.handle.transErr(ErrAddPackageFailed)
	}

	return nil
}

func (t *transaction) RemovePkg(pkg Package) error {
	if t.handle.ptr == 0 {
		return ErrInvalidHandle
	}

	if lib.AlpmRemovePkg == nil {
		return stderrors.New("missing function: alpm_remove_pkg")
	}

	pkgImpl := pkg.(*package_)
	result := lib.AlpmRemovePkg(t.handle.ptr, pkgImpl.ptr)
	if result != 0 {
		return t.handle.transErr(ErrRemovePackageFailed)
	}

	return nil
}

func (t *transaction) SyncSysupgrade(enableDowngrade bool) error {
	return t.handle.SyncSysupgrade(enableDowngrade)
}

func (t *transaction) GetFlags() TransactionFlag {
	if t.handle.ptr == 0 {
		return 0
	}

	if lib.AlpmTransGetFlags == nil {
		return 0
	}

	result := lib.AlpmTransGetFlags(t.handle.ptr)
	return TransactionFlag(result)
}

func (t *transaction) GetAdd() ([]Package, error) {
	return t.getPackageList(lib.AlpmTransGetAdd, "alpm_trans_get_add")
}

func (t *transaction) GetRemove() ([]Package, error) {
	return t.getPackageList(lib.AlpmTransGetRemove, "alpm_trans_get_remove")
}

func (t *transaction) getPackageList(fn func(uintptr) uintptr, name string) ([]Package, error) {
	if t.handle.ptr == 0 {
		return nil, ErrInvalidHandle
	}
	if fn == nil {
		return nil, stderrors.New("missing function: " + name)
	}

	pkgs := collectList(alpmlist.NewList(fn(t.handle.ptr)), func(ptr uintptr) Package {
		return newPackage(ptr, t.handle)
	})
	if pkgs == nil {
		return []Package{}, nil
	}
	return pkgs, nil
}
