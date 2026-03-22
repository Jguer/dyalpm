package dyalpm

import (
	stderrors "errors"

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

func (t *transaction) getTransactionList(funcName string, failErr error) (*alpmlist.List, error) {
	var fn func(uintptr, *uintptr) int32
	switch funcName {
	case "alpm_trans_prepare":
		fn = lib.AlpmTransPrepare
	case "alpm_trans_commit":
		fn = lib.AlpmTransCommit
	default:
		return nil, stderrors.New("missing function: " + funcName)
	}

	if fn == nil {
		return nil, stderrors.New("missing function: " + funcName)
	}

	var dataListPtr uintptr
	result := fn(t.handle.ptr, &dataListPtr)
	if result != 0 {
		return nil, failErr
	}

	if dataListPtr == 0 {
		return nil, nil
	}

	alpmList := alpmlist.NewList(dataListPtr)
	if alpmList == nil {
		return nil, nil
	}

	return alpmList, nil
}

// NewTransaction creates a new transaction for the given handle
func NewTransaction(h Handle) Transaction {
	handle := h.(*handle)
	return &transaction{
		handle: handle,
	}
}

func (t *transaction) Init(flags TransactionFlag) error {
	if t.handle.ptr == 0 {
		return ErrInvalidHandle
	}

	if lib.AlpmTransInit == nil {
		return stderrors.New("missing function: alpm_trans_init")
	}

	result := lib.AlpmTransInit(t.handle.ptr, clampIntToInt32(int(flags)))
	if result != 0 {
		return ErrTransactionInitFailed
	}

	return nil
}

func (t *transaction) Prepare() ([]DepMissing, error) {
	if t.handle.ptr == 0 {
		return nil, ErrInvalidHandle
	}

	alpmList, err := t.getTransactionList("alpm_trans_prepare", ErrTransactionPrepareFailed)
	if err != nil || alpmList == nil {
		return []DepMissing{}, err
	}
	defer alpmList.Free()

	missing := collectList(alpmList, func(ptr uintptr) DepMissing {
		return newDepMissing(ptr)
	})

	return missing, nil
}

func (t *transaction) Commit() ([]FileConflict, error) {
	if t.handle.ptr == 0 {
		return nil, ErrInvalidHandle
	}

	alpmList, err := t.getTransactionList("alpm_trans_commit", ErrTransactionCommitFailed)
	if err != nil || alpmList == nil {
		return []FileConflict{}, err
	}
	defer alpmList.Free()

	conflicts := collectList(alpmList, func(ptr uintptr) FileConflict {
		return newFileConflict(ptr)
	})

	return conflicts, nil
}

func (t *transaction) Release() error {
	if t.handle.ptr == 0 {
		return ErrInvalidHandle
	}

	if lib.AlpmTransRelease == nil {
		return stderrors.New("missing function: alpm_trans_release")
	}

	result := lib.AlpmTransRelease(t.handle.ptr)
	if result != 0 {
		return ErrTransactionReleaseFailed
	}

	return nil
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
		return ErrAddPackageFailed
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
		return ErrRemovePackageFailed
	}

	return nil
}

func (t *transaction) SyncSysupgrade(enableDowngrade bool) error {
	if t.handle.ptr == 0 {
		return ErrInvalidHandle
	}

	if lib.AlpmSyncSysupgrade == nil {
		return stderrors.New("missing function: alpm_sync_sysupgrade")
	}

	downgrade := int32(0)
	if enableDowngrade {
		downgrade = 1
	}
	result := lib.AlpmSyncSysupgrade(t.handle.ptr, downgrade)
	if result != 0 {
		return ErrSysupgradeFailed
	}

	return nil
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
	return t.getPackageList("alpm_trans_get_add")
}

func (t *transaction) GetRemove() ([]Package, error) {
	return t.getPackageList("alpm_trans_get_remove")
}

func (t *transaction) getPackageList(funcName string) ([]Package, error) {
	if t.handle.ptr == 0 {
		return nil, ErrInvalidHandle
	}

	switch funcName {
	case "alpm_trans_get_add", "alpm_trans_get_remove":
	default:
		return nil, stderrors.New("missing function: " + funcName)
	}

	var getFn func(uintptr) uintptr
	if funcName == "alpm_trans_get_add" {
		getFn = lib.AlpmTransGetAdd
	} else {
		getFn = lib.AlpmTransGetRemove
	}
	if getFn == nil {
		return nil, stderrors.New("missing function: " + funcName)
	}

	listPtr := getFn(t.handle.ptr)
	if listPtr == 0 {
		return []Package{}, nil
	}

	alpmList := alpmlist.NewList(listPtr)
	if alpmList == nil {
		return []Package{}, nil
	}

	var pkgs []Package
	for item := alpmList; item != nil && item.Ptr() != 0; item = item.Next() {
		pkgPtr := item.Data()
		if pkgPtr != 0 {
			pkgs = append(pkgs, newPackage(pkgPtr, t.handle))
		}
	}

	return pkgs, nil
}
