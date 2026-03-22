package dyalpm

import (
	stderrors "errors"
	"runtime"
	"unsafe"

	"github.com/Jguer/dyalpm/internal/lib"
	alpmlist "github.com/Jguer/dyalpm/internal/list"
)

// Usage type for database usage flags
type Usage int

const (
	UsageSync    Usage = 1 << 0
	UsageSearch  Usage = 1 << 1
	UsageInstall Usage = 1 << 2
	UsageUpgrade Usage = 1 << 3
	UsageAll     Usage = (1 << 4) - 1
)

// Database represents an ALPM database (minimal surface used by yay).
type Database interface {
	Name() string
	Pkg(name string) Package
	PkgCache() PackageIterator
	Search(needles []string) PackageIterator
	SetServers(servers []string) error
	SetUsage(usage int) error
}

type database struct {
	ptr    uintptr
	handle *handle
}

func newDatabase(ptr uintptr, h *handle) *database {
	return &database{
		ptr:    ptr,
		handle: h,
	}
}

func (d *database) getList(funcName string) (*alpmlist.List, error) {
	var listPtr uintptr
	switch funcName {
	case "alpm_db_get_pkgcache":
		if lib.AlpmDBGetPkgcache == nil {
			return nil, stderrors.New("missing function: alpm_db_get_pkgcache")
		}
		listPtr = lib.AlpmDBGetPkgcache(d.ptr)
	case "alpm_db_get_groupcache":
		if lib.AlpmDBGetGroupcache == nil {
			return nil, stderrors.New("missing function: alpm_db_get_groupcache")
		}
		listPtr = lib.AlpmDBGetGroupcache(d.ptr)
	case "alpm_db_get_servers":
		if lib.AlpmDBGetServers == nil {
			return nil, stderrors.New("missing function: alpm_db_get_servers")
		}
		listPtr = lib.AlpmDBGetServers(d.ptr)
	case "alpm_db_get_cache_servers":
		if lib.AlpmDBGetCacheServers == nil {
			return nil, stderrors.New("missing function: alpm_db_get_cache_servers")
		}
		listPtr = lib.AlpmDBGetCacheServers(d.ptr)
	default:
		return nil, stderrors.New("missing function: " + funcName)
	}
	if listPtr == 0 {
		return nil, nil
	}

	return alpmlist.NewList(listPtr), nil
}

func (d *database) getServers(funcName string) []string {
	if d.ptr == 0 {
		return nil
	}

	listPtr, err := d.getList(funcName)
	if err != nil || listPtr == nil {
		return nil
	}

	return collectList(listPtr, func(ptr uintptr) string {
		return lib.PtrToString(ptr)
	})
}

func (d *database) GetName() string {
	if d.ptr == 0 {
		return ""
	}
	if lib.AlpmDBGetName == nil {
		return ""
	}
	return lib.PtrToString(lib.AlpmDBGetName(d.ptr))
}

func (d *database) GetHandle() Handle {
	return d.handle
}

func (d *database) GetPkg(name string) (Package, error) {
	if d.ptr == 0 {
		return nil, ErrInvalidDatabase
	}
	if lib.AlpmDBGetPkg == nil {
		return nil, stderrors.New("missing function: alpm_db_get_pkg")
	}

	pkgPtr := lib.AlpmDBGetPkg(d.ptr, name)
	if pkgPtr == 0 {
		return nil, ErrPackageNotFound
	}

	return newPackage(pkgPtr, d.handle), nil
}

func (d *database) GetPkgCache() ([]Package, error) {
	if d.ptr == 0 {
		return nil, ErrInvalidDatabase
	}

	alpmList, err := d.getList("alpm_db_get_pkgcache")
	if err != nil || alpmList == nil {
		return []Package{}, err
	}

	pkgs := collectList(alpmList, func(ptr uintptr) Package {
		return newPackage(ptr, d.handle)
	})

	return pkgs, nil
}

// go-alpm/v2 compatible methods
func (d *database) Name() string { return d.GetName() }

func (d *database) Pkgs() []Package {
	pkgs, _ := d.GetPkgCache()
	return pkgs
}

func (d *database) Pkg(name string) Package {
	pkg, err := d.GetPkg(name)
	if err != nil {
		return nil
	}
	return pkg
}

func (d *database) PkgCache() PackageIterator {
	if d.ptr == 0 {
		return PackageIterator{}
	}

	if lib.AlpmDBGetPkgcache == nil {
		return PackageIterator{}
	}

	listPtr := lib.AlpmDBGetPkgcache(d.ptr)
	if listPtr == 0 {
		return PackageIterator{}
	}

	return newPackageIterator(listPtr, d.handle, false)
}

func (d *database) Search(needles []string) PackageIterator {
	if d.ptr == 0 {
		return PackageIterator{}
	}

	if len(needles) == 0 {
		return PackageIterator{}
	}

	if lib.AlpmDBSearch == nil {
		return PackageIterator{}
	}

	var alpmList *alpmlist.List
	var cStrings [][]byte

	for _, s := range needles {
		cs := lib.CString(s)
		cStrings = append(cStrings, cs)
		alpmList = alpmlist.Add(alpmList, uintptr(unsafe.Pointer(&cs[0])))
	}

	if alpmList == nil {
		return PackageIterator{}
	}
	defer alpmList.Free()

	var resultListPtr uintptr
	ret := int(lib.AlpmDBSearch(d.ptr, alpmList.Ptr(), &resultListPtr))
	runtime.KeepAlive(cStrings)
	runtime.KeepAlive(alpmList)

	if ret != 0 || resultListPtr == 0 {
		return PackageIterator{}
	}

	// alpm_db_search returns a list that must be freed by the caller.
	return newPackageIterator(resultListPtr, d.handle, true)
}

func (d *database) GetGroup(name string) (Group, error) {
	if d.ptr == 0 {
		return nil, ErrInvalidDatabase
	}
	if lib.AlpmDBGetGroup == nil {
		return nil, stderrors.New("missing function: alpm_db_get_group")
	}

	groupPtr := lib.AlpmDBGetGroup(d.ptr, name)
	if groupPtr == 0 {
		return nil, ErrGroupNotFound
	}

	return newGroup(groupPtr, d.handle), nil
}

func (d *database) GetGroupCache() ([]Group, error) {
	if d.ptr == 0 {
		return nil, ErrInvalidDatabase
	}

	alpmList, err := d.getList("alpm_db_get_groupcache")
	if err != nil || alpmList == nil {
		return []Group{}, err
	}

	groups := collectList(alpmList, func(ptr uintptr) Group {
		return newGroup(ptr, d.handle)
	})

	return groups, nil
}

func (d *database) Update(force bool) error {
	if d.ptr == 0 {
		return ErrInvalidDatabase
	}

	if lib.AlpmDBUpdate == nil {
		return stderrors.New("missing function: alpm_db_update")
	}

	// alpm_db_update expects a list of databases
	dbList := alpmlist.Add(nil, d.ptr)
	if dbList == nil {
		return ErrDatabaseUpdateFailed
	}
	defer dbList.Free()

	forceInt := int32(0)
	if force {
		forceInt = 1
	}
	result := lib.AlpmDBUpdate(d.handle.ptr, dbList.Ptr(), forceInt)
	if result != 0 {
		return ErrDatabaseUpdateFailed
	}

	return nil
}

func (d *database) Unregister() error {
	if d.ptr == 0 {
		return ErrInvalidDatabase
	}

	if lib.AlpmDBUnregister == nil {
		return stderrors.New("missing function: alpm_db_unregister")
	}

	result := lib.AlpmDBUnregister(d.ptr)
	if result != 0 {
		return ErrDatabaseUnregisterFailed
	}

	d.ptr = 0
	return nil
}

func (d *database) GetServers() []string {
	return d.getServers("alpm_db_get_servers")
}

func (d *database) SetServers(servers []string) error {
	return d.setServers("alpm_db_set_servers", servers)
}

func (d *database) setServers(funcName string, servers []string) error {
	if d.ptr == 0 {
		return ErrInvalidDatabase
	}

	var setServersFn func(uintptr, uintptr) int32
	switch funcName {
	case "alpm_db_set_servers":
		setServersFn = lib.AlpmDBSetServers
	case "alpm_db_set_cache_servers":
		setServersFn = lib.AlpmDBSetCacheServers
	default:
		return stderrors.New("missing function: " + funcName)
	}

	if setServersFn == nil {
		return stderrors.New("missing function: " + funcName)
	}

	var alpmList *alpmlist.List
	var cStrings [][]byte

	for _, s := range servers {
		cs := lib.CString(s)
		cStrings = append(cStrings, cs)
		alpmList = alpmlist.Add(alpmList, uintptr(unsafe.Pointer(&cs[0])))
	}
	if alpmList != nil {
		defer alpmList.Free()
	}

	result := setServersFn(d.ptr, alpmList.Ptr())

	// Keep strings alive during the call
	runtime.KeepAlive(cStrings)
	runtime.KeepAlive(alpmList)

	if result != 0 {
		return ErrDatabaseUpdateFailed
	}

	return nil
}

func (d *database) AddServer(url string) error {
	if d.ptr == 0 {
		return ErrInvalidDatabase
	}

	if lib.AlpmDBAddServer == nil {
		return stderrors.New("missing function: alpm_db_add_server")
	}

	result := lib.AlpmDBAddServer(d.ptr, url)

	if result != 0 {
		return ErrDatabaseUpdateFailed
	}

	return nil
}

func (d *database) RemoveServer(url string) error {
	if d.ptr == 0 {
		return ErrInvalidDatabase
	}

	if lib.AlpmDBRemoveServer == nil {
		return stderrors.New("missing function: alpm_db_remove_server")
	}

	result := lib.AlpmDBRemoveServer(d.ptr, url)

	if result != 0 {
		return ErrDatabaseUpdateFailed
	}

	return nil
}

func (d *database) GetCacheServers() []string {
	return d.getServers("alpm_db_get_cache_servers")
}

func (d *database) SetCacheServers(servers []string) error {
	return d.setServers("alpm_db_set_cache_servers", servers)
}

func (d *database) AddCacheServer(url string) error {
	if d.ptr == 0 {
		return ErrInvalidDatabase
	}

	if lib.AlpmDBAddCacheServer == nil {
		return stderrors.New("missing function: alpm_db_add_cache_server")
	}

	result := lib.AlpmDBAddCacheServer(d.ptr, url)

	if result != 0 {
		return ErrDatabaseUpdateFailed
	}

	return nil
}

func (d *database) RemoveCacheServer(url string) error {
	if d.ptr == 0 {
		return ErrInvalidDatabase
	}

	if lib.AlpmDBRemoveCacheServer == nil {
		return stderrors.New("missing function: alpm_db_remove_cache_server")
	}

	result := lib.AlpmDBRemoveCacheServer(d.ptr, url)

	if result != 0 {
		return ErrDatabaseUpdateFailed
	}

	return nil
}

func (d *database) SetUsage(usage int) error {
	if d.ptr == 0 {
		return ErrInvalidDatabase
	}

	if lib.AlpmDBSetUsage == nil {
		return stderrors.New("missing function: alpm_db_set_usage")
	}

	usageInt32 := clampIntToInt32(usage)
	result := lib.AlpmDBSetUsage(d.ptr, usageInt32)
	if result != 0 {
		return ErrDatabaseUpdateFailed
	}

	return nil
}

func (d *database) GetUsage() (int, error) {
	if d.ptr == 0 {
		return 0, ErrInvalidDatabase
	}

	if lib.AlpmDBGetUsage == nil {
		return 0, stderrors.New("missing function: alpm_db_get_usage")
	}

	var usagePtr int32
	result := lib.AlpmDBGetUsage(d.ptr, &usagePtr)
	if result != 0 {
		return 0, ErrDatabaseUpdateFailed
	}

	return int(usagePtr), nil
}

func (d *database) IsValid() bool {
	if d.ptr == 0 {
		return false
	}

	if lib.AlpmDBGetValid == nil {
		return false
	}

	result := lib.AlpmDBGetValid(d.ptr)
	return result == 0
}

func (d *database) GetSigLevel() int {
	if d.ptr == 0 {
		return 0
	}

	if lib.AlpmDBGetSiglevel == nil {
		return 0
	}

	result := lib.AlpmDBGetSiglevel(d.ptr)
	return int(result)
}

func (d *database) GetNativeHandle() Handle {
	if d.ptr == 0 {
		return nil
	}

	if lib.AlpmDBGetHandle == nil {
		return nil
	}

	result := lib.AlpmDBGetHandle(d.ptr)
	if result == 0 {
		return nil
	}

	return d.handle
}

func (d *database) CheckPGPSignature() (SigList, error) {
	if d.ptr == 0 {
		return SigList{}, ErrInvalidDatabase
	}

	return checkPGPSignature(d.ptr, d.handle, "alpm_db_check_pgp_signature")
}

// Group represents a package group
type Group interface {
	GetName() string
	GetPackages() ([]Package, error)
}

type group struct {
	ptr    uintptr
	handle *handle
}

func newGroup(ptr uintptr, h *handle) *group {
	return &group{
		ptr:    ptr,
		handle: h,
	}
}

func (g *group) GetName() string {
	if g.ptr == 0 {
		return ""
	}
	// Group structure: char *name at offset 0
	base := unsafe.Pointer(g.ptr)
	namePtr := *(*uintptr)(base)
	return lib.PtrToString(namePtr)
}

func (g *group) GetPackages() ([]Package, error) {
	if g.ptr == 0 {
		return nil, ErrInvalidGroup
	}
	// Group structure: alpm_list_t *packages at offset of pointer size
	base := unsafe.Pointer(g.ptr)
	packagesPtr := *(*uintptr)(unsafe.Add(base, unsafe.Sizeof(uintptr(0))))
	if packagesPtr == 0 {
		return []Package{}, nil
	}

	alpmList := alpmlist.NewList(packagesPtr)
	if alpmList == nil {
		return []Package{}, nil
	}

	var pkgs []Package
	for item := alpmList; item != nil && item.Ptr() != 0; item = item.Next() {
		pkgPtr := item.Data()
		if pkgPtr != 0 {
			pkgs = append(pkgs, newPackage(pkgPtr, g.handle))
		}
	}

	return pkgs, nil
}
