package dyalpm

import (
	stderrors "errors"
	"runtime"
	"unsafe"

	"github.com/ebitengine/purego"

	"github.com/Jguer/dyalpm/internal/dyerrors"
	"github.com/Jguer/dyalpm/internal/lib"
	alpmlist "github.com/Jguer/dyalpm/internal/list"
)

// Handle represents an ALPM handle
type Handle interface {
	// Release releases the handle and cleans up resources
	Release() error

	LocalDB() (Database, error)
	SyncDBs() ([]Database, error)
	SyncDBListByDBName(name string) ([]Database, error)
	SyncDBByName(name string) (Database, error)

	// Transaction methods
	TransInit(flags TransactionFlag) error
	TransRelease() error
	SyncSysupgrade(enableDowngrade bool) error
	TransGetAdd() PackageIterator

	// RegisterSyncDB registers a sync database
	RegisterSyncDB(name string, siglevel int) (Database, error)

	// Root returns the root path
	Root() string

	// DBPath returns the database path
	DBPath() string

	// Errno returns the current error code
	Errno() dyerrors.Errno

	// StrError returns the error string for an error code
	StrError(errno dyerrors.Errno) string

	// Options
	SetLogFile(path string) error
	LogFile() string
	SetGPGDir(path string) error
	GPGDir() string
	SetUseSyslog(enable bool) error
	UseSyslog() bool
	SetCheckSpace(enable bool) error
	CheckSpace() bool
	SetDBExt(ext string) error
	DBExt() string
	SetDefaultSigLevel(level int) error
	DefaultSigLevel() int
	SetLocalFileSigLevel(level int) error
	LocalFileSigLevel() int
	SetRemoteFileSigLevel(level int) error
	RemoteFileSigLevel() int
	SetParallelDownloads(num int) error
	ParallelDownloads() int

	// Architecture
	Architectures() ([]string, error)
	SetArchitectures(archs []string) error
	AddArchitecture(arch string) error
	RemoveArchitecture(arch string) error

	// List Options
	CacheDirs() ([]string, error)
	SetCacheDirs(dirs []string) error
	AddCacheDir(dir string) error
	RemoveCacheDir(dir string) error

	HookDirs() ([]string, error)
	SetHookDirs(dirs []string) error
	AddHookDir(dir string) error
	RemoveHookDir(dir string) error

	NoUpgrades() ([]string, error)
	SetNoUpgrades(paths []string) error
	AddNoUpgrade(path string) error
	RemoveNoUpgrade(path string) error
	MatchNoUpgrade(path string) (int, error)

	NoExtracts() ([]string, error)
	SetNoExtracts(paths []string) error
	AddNoExtract(path string) error
	RemoveNoExtract(path string) error
	MatchNoExtract(path string) (int, error)

	IgnorePkgs() ([]string, error)
	SetIgnorePkgs(pkgs []string) error
	AddIgnorePkg(pkg string) error
	RemoveIgnorePkg(pkg string) error

	IgnoreGroups() ([]string, error)
	SetIgnoreGroups(groups []string) error
	AddIgnoreGroup(group string) error
	RemoveIgnoreGroup(group string) error

	OverwriteFiles() ([]string, error)
	SetOverwriteFiles(globs []string) error
	AddOverwriteFile(glob string) error
	RemoveOverwriteFile(glob string) error

	SetSandboxUser(user string) error
	SandboxUser() string
	SetDisableDLTimeout(disable bool) error
	SetDisableSandbox(disable bool) error

	// LoadPackage loads a package from a file
	LoadPackage(filename string, full bool, siglevel int) (Package, error)

	// DB Management
	UnregisterAllSyncDBs() error

	// Utils
	LogAction(prefix, message string) error
	Unlock() error
	FetchPkgURL(url string) (string, error)
	FindGroupPkgs(dbs []Database, name string) ([]Package, error)
	ExtractKeyID(identifier string, sig []byte) ([]string, error)
	InterruptTransaction() error
	SandboxSetupChild(user, dir string) error

	// Assume Installed
	AssumeInstalled() ([]Dependency, error)
	SetAssumeInstalled(deps []Dependency) error
	AddAssumeInstalled(dep Dependency) error
	RemoveAssumeInstalled(dep Dependency) error

	// Dependency Resolution
	CheckDeps(pkgs []Package, remPkgs []Package, upgradePkgs []Package, reverseDeps bool) ([]DepMissing, error)
	CheckConflicts(pkgs []Package) ([]Conflict, error)
	FindDBSatisfier(dbs []Database, depstring string) Package

	// Callbacks (raw pointers)
	//
	// NOTE: libalpm's log callback (`alpm_cb_log`) uses a `va_list`, which cannot be
	// safely bridged to a Go function without a C shim. This wrapper therefore
	// exposes logcb as raw pointers only.
	LogCallback() (cb uintptr, ctx uintptr)
	SetLogCallback(cb uintptr, ctx uintptr) error

	DownloadCallback() (cb uintptr, ctx uintptr)
	SetDownloadCallback(cb uintptr, ctx uintptr) error
	SetDownloadCallbackFunc(cb DownloadCallback) error

	FetchCallback() (cb uintptr, ctx uintptr)
	SetFetchCallback(cb uintptr, ctx uintptr) error
	SetFetchCallbackFunc(cb FetchCallback) error

	EventCallback() (cb uintptr, ctx uintptr)
	SetEventCallback(cb uintptr, ctx uintptr) error
	SetEventCallbackFunc(cb EventCallback) error

	QuestionCallback() (cb uintptr, ctx uintptr)
	SetQuestionCallback(cb uintptr, ctx uintptr) error
	SetQuestionCallbackFunc(cb QuestionCallback) error

	ProgressCallback() (cb uintptr, ctx uintptr)
	SetProgressCallback(cb uintptr, ctx uintptr) error
	SetProgressCallbackFunc(cb ProgressCallback) error
}

type handle struct {
	ptr uintptr
}

// Initialize creates a new ALPM handle
func Initialize(root, dbpath string) (Handle, error) {
	if err := lib.EnsureAlpmLoaded(); err != nil {
		return nil, err
	}
	if lib.AlpmInitialize == nil {
		return nil, stderrors.New("missing function: alpm_initialize")
	}

	var errno int32
	handlePtr := lib.AlpmInitialize(root, dbpath, &errno)

	if handlePtr == 0 {
		if errno != 0 {
			return nil, dyerrors.NewError(dyerrors.Errno(errno), "failed to initialize ALPM handle")
		}
		return nil, dyerrors.NewError(dyerrors.ErrSystem, "failed to initialize ALPM handle")
	}

	return &handle{ptr: handlePtr}, nil
}

func (h *handle) Release() error {
	if h.ptr == 0 {
		return stderrors.New("handle already released")
	}

	oldPtr := h.ptr

	if lib.AlpmRelease == nil {
		return stderrors.New("missing function: alpm_release")
	}

	r1 := lib.AlpmRelease(h.ptr)
	errno := lib.AlpmErrno(h.ptr)
	if errno != 0 || r1 != 0 {
		return stderrors.New("failed to release handle")
	}

	unregisterCallbackSet(oldPtr)
	h.ptr = 0
	return nil
}

func (h *handle) Errno() dyerrors.Errno {
	if h.ptr == 0 {
		return dyerrors.ErrHandleNull
	}
	if lib.AlpmErrno == nil {
		return dyerrors.ErrSystem
	}
	return dyerrors.Errno(lib.AlpmErrno(h.ptr))
}

func (h *handle) StrError(errno dyerrors.Errno) string {
	if lib.AlpmStrerror == nil {
		return "unknown error"
	}

	r1 := lib.AlpmStrerror(clampIntToInt32(int(errno)))
	if r1 == 0 {
		return "unknown error"
	}

	return lib.PtrToString(r1)
}

func (h *handle) Root() string {
	if h.ptr == 0 {
		return ""
	}
	if lib.AlpmOptionGetRoot == nil {
		return ""
	}
	return lib.PtrToString(lib.AlpmOptionGetRoot(h.ptr))
}

func (h *handle) DBPath() string {
	if h.ptr == 0 {
		return ""
	}
	if lib.AlpmOptionGetDbpath == nil {
		return ""
	}
	return lib.PtrToString(lib.AlpmOptionGetDbpath(h.ptr))
}

func (h *handle) getLocalDB() (Database, error) {
	if h.ptr == 0 {
		return nil, stderrors.New("handle is invalid")
	}
	if lib.AlpmGetLocalDB == nil {
		return nil, stderrors.New("missing function: alpm_get_localdb")
	}
	r1 := lib.AlpmGetLocalDB(h.ptr)
	if r1 == 0 {
		return nil, dyerrors.NewError(h.Errno(), "failed to get local database")
	}
	return newDatabase(r1, h), nil
}

func (h *handle) getSyncDBs() ([]Database, error) {
	if h.ptr == 0 {
		return nil, stderrors.New("handle is invalid")
	}
	if lib.AlpmGetSyncDBS == nil {
		return nil, stderrors.New("missing function: alpm_get_syncdbs")
	}
	r1 := lib.AlpmGetSyncDBS(h.ptr)
	if r1 == 0 {
		errno := h.Errno()
		if errno != dyerrors.ErrOK {
			return nil, dyerrors.NewError(errno, "failed to get sync databases")
		}
		return []Database{}, nil
	}
	listPtr := r1

	alpmList := alpmlist.NewList(listPtr)
	if alpmList == nil {
		return []Database{}, nil
	}

	var dbs []Database
	for item := alpmList; item != nil && item.Ptr() != 0; item = item.Next() {
		dbPtr := item.Data()
		if dbPtr != 0 {
			dbs = append(dbs, newDatabase(dbPtr, h))
		}
	}

	return dbs, nil
}

// go-alpm/v2 compatible methods
func (h *handle) LocalDB() (Database, error) {
	return h.getLocalDB()
}

func (h *handle) SyncDBs() ([]Database, error) {
	dbs, err := h.getSyncDBs()
	if err != nil {
		return nil, err
	}
	return dbs, nil
}

func (h *handle) SyncDBListByDBName(name string) ([]Database, error) {
	dbs, err := h.getSyncDBs()
	if err != nil {
		return nil, err
	}
	for _, db := range dbs {
		if db.Name() == name {
			return []Database{db}, nil
		}
	}
	return nil, stderrors.New("database not found: " + name)
}

func (h *handle) SyncDBByName(name string) (Database, error) {
	dbs, err := h.getSyncDBs()
	if err != nil {
		return nil, err
	}
	for _, db := range dbs {
		if db.Name() == name {
			return db, nil
		}
	}
	return nil, stderrors.New("database not found: " + name)
}

func (h *handle) TransInit(flags TransactionFlag) error {
	if h.ptr == 0 {
		return stderrors.New("invalid handle")
	}
	if lib.AlpmTransInit == nil {
		return stderrors.New("missing function: alpm_trans_init")
	}

	result := lib.AlpmTransInit(h.ptr, clampIntToInt32(int(flags)))
	if result != 0 {
		return stderrors.New("failed to initialize transaction")
	}

	return nil
}

func (h *handle) TransRelease() error {
	if h.ptr == 0 {
		return stderrors.New("invalid handle")
	}
	if lib.AlpmTransRelease == nil {
		return stderrors.New("missing function: alpm_trans_release")
	}

	result := lib.AlpmTransRelease(h.ptr)
	if result != 0 {
		return stderrors.New("failed to release transaction")
	}

	return nil
}

func (h *handle) SyncSysupgrade(enableDowngrade bool) error {
	if h.ptr == 0 {
		return stderrors.New("invalid handle")
	}
	if lib.AlpmSyncSysupgrade == nil {
		return stderrors.New("missing function: alpm_sync_sysupgrade")
	}

	down := int32(0)
	if enableDowngrade {
		down = 1
	}
	result := lib.AlpmSyncSysupgrade(h.ptr, down)
	if result != 0 {
		return stderrors.New("failed to sync sysupgrade")
	}

	return nil
}

func (h *handle) TransGetAdd() PackageIterator {
	if h.ptr == 0 {
		return PackageIterator{}
	}
	if lib.AlpmTransGetAdd == nil {
		return PackageIterator{}
	}
	listPtr := lib.AlpmTransGetAdd(h.ptr)
	if listPtr == 0 {
		return PackageIterator{}
	}

	return newPackageIterator(listPtr, h, false)
}

func (h *handle) RegisterSyncDB(name string, siglevel int) (Database, error) {
	if h.ptr == 0 {
		return nil, stderrors.New("handle is invalid")
	}
	if lib.AlpmRegisterSyncDB == nil {
		return nil, stderrors.New("missing function: alpm_register_syncdb")
	}

	siglevelInt32 := clampIntToInt32(siglevel)
	r1 := lib.AlpmRegisterSyncDB(h.ptr, name, siglevelInt32)
	if r1 == 0 {
		return nil, dyerrors.NewError(h.Errno(), "failed to register sync database")
	}
	return newDatabase(r1, h), nil
}

func (h *handle) UnregisterAllSyncDBs() error {
	if h.ptr == 0 {
		return dyerrors.ErrHandleNull
	}
	if lib.AlpmUnregisterAllSyncDBs == nil {
		return stderrors.New("missing function: alpm_unregister_all_syncdbs")
	}

	result := lib.AlpmUnregisterAllSyncDBs(h.ptr)
	if result != 0 {
		return ErrDatabaseUnregisterFailed
	}

	return nil
}

func (h *handle) LogAction(prefix, message string) error {
	if h.ptr == 0 {
		return dyerrors.ErrHandleNull
	}

	if lib.AlpmLogactionSym == 0 {
		return stderrors.New("missing function: alpm_logaction")
	}

	cPrefix := lib.CString(prefix)
	cMsg := lib.CString(message)
	cFmt := lib.CString("%s")

	r1, _, _ := purego.SyscallN(
		lib.AlpmLogactionSym,
		h.ptr,
		uintptr(unsafe.Pointer(&cPrefix[0])),
		uintptr(unsafe.Pointer(&cFmt[0])),
		uintptr(unsafe.Pointer(&cMsg[0])),
	)

	runtime.KeepAlive(cPrefix)
	runtime.KeepAlive(cMsg)
	runtime.KeepAlive(cFmt)

	if r1 != 0 {
		return stderrors.New("failed to log action")
	}
	return nil
}

func (h *handle) Unlock() error {
	if h.ptr == 0 {
		return dyerrors.ErrHandleNull
	}
	if lib.AlpmUnlock == nil {
		return stderrors.New("missing function: alpm_unlock")
	}

	if lib.AlpmUnlock(h.ptr) != 0 {
		return stderrors.New("failed to unlock")
	}

	return nil
}

func (h *handle) FetchPkgURL(url string) (string, error) {
	if h.ptr == 0 {
		return "", dyerrors.ErrHandleNull
	}
	if lib.AlpmFetchPkgurl == nil {
		return "", stderrors.New("missing function: alpm_fetch_pkgurl")
	}

	// alpm_fetch_pkgurl(handle, urls_list, &fetched_list)
	cURL := lib.CString(url)
	urlList := alpmlist.Add(nil, uintptr(unsafe.Pointer(&cURL[0])))
	if urlList == nil {
		return "", stderrors.New("failed to create URL list")
	}
	defer urlList.Free()

	var fetchedListPtr uintptr
	r1 := int(lib.AlpmFetchPkgurl(h.ptr, urlList.Ptr(), &fetchedListPtr))
	runtime.KeepAlive(cURL)

	if r1 != 0 {
		return "", stderrors.New("failed to fetch package URL")
	}

	if fetchedListPtr == 0 {
		return "", nil
	}

	fetchedList := alpmlist.NewList(fetchedListPtr)
	defer fetchedList.Free()

	// Return the first fetched path (since we only requested one URL)
	if fetchedList.Ptr() != 0 {
		ptr := fetchedList.Data()
		if ptr != 0 {
			return lib.PtrToString(ptr), nil
		}
	}

	return "", nil
}

func (h *handle) InterruptTransaction() error {
	if h.ptr == 0 {
		return dyerrors.ErrHandleNull
	}
	if lib.AlpmTransInterrupt == nil {
		return stderrors.New("missing function: alpm_trans_interrupt")
	}
	if lib.AlpmTransInterrupt(h.ptr) != 0 {
		return stderrors.New("failed to interrupt transaction")
	}

	return nil
}

func (h *handle) SandboxSetupChild(user, dir string) error {
	if h.ptr == 0 {
		return dyerrors.ErrHandleNull
	}
	if lib.AlpmSandboxSetupChild == nil {
		return stderrors.New("missing function: alpm_sandbox_setup_child")
	}
	if lib.AlpmSandboxSetupChild(h.ptr, user, dir) != 0 {
		return stderrors.New("failed to setup sandbox child")
	}

	return nil
}
