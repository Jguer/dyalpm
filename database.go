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

// list returns the libalpm-owned list from fn, or an error if fn is not bound.
func (d *database) list(fn func(uintptr) uintptr, name string) (*alpmlist.List, error) {
	if fn == nil {
		return nil, stderrors.New("missing function: " + name)
	}
	return alpmlist.NewList(fn(d.ptr)), nil
}

func (d *database) getServers(fn func(uintptr) uintptr, name string) []string {
	if d.ptr == 0 {
		return nil
	}
	l, err := d.list(fn, name)
	if err != nil {
		return nil
	}
	return collectList(l, lib.PtrToString)
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

	alpmList, err := d.list(lib.AlpmDBGetPkgcache, "alpm_db_get_pkgcache")
	if err != nil || alpmList == nil {
		return []Package{}, err
	}

	return collectList(alpmList, func(ptr uintptr) Package {
		return newPackage(ptr, d.handle)
	}), nil
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

	return newPackageIterator(listPtr, d.handle)
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

	alpmList, bufs := cStringList(needles)
	defer alpmList.Free()

	var resultListPtr uintptr
	ret := lib.AlpmDBSearch(d.ptr, alpmList.Ptr(), &resultListPtr)
	runtime.KeepAlive(bufs)

	if ret != 0 || resultListPtr == 0 {
		return PackageIterator{}
	}

	// The caller owns the result list: copy it out and free it once, here.
	resultList := alpmlist.NewList(resultListPtr)
	defer resultList.Free()
	return PackageIterator{pkgs: newPackageIterator(resultListPtr, d.handle).Collect(), handle: d.handle}
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

	alpmList, err := d.list(lib.AlpmDBGetGroupcache, "alpm_db_get_groupcache")
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
	return d.getServers(lib.AlpmDBGetServers, "alpm_db_get_servers")
}

func (d *database) SetServers(servers []string) error {
	return d.setServers(lib.AlpmDBSetServers, "alpm_db_set_servers", servers)
}

func (d *database) setServers(fn func(uintptr, uintptr) int32, name string, servers []string) error {
	if d.ptr == 0 {
		return ErrInvalidDatabase
	}
	if fn == nil {
		return stderrors.New("missing function: " + name)
	}

	alpmList, bufs := cStringList(servers)
	defer alpmList.Free()

	result := fn(d.ptr, alpmList.Ptr())
	runtime.KeepAlive(bufs)

	if result != 0 {
		return ErrDatabaseUpdateFailed
	}
	return nil
}

func (d *database) serverOp(fn func(uintptr, string) int32, name, url string) error {
	if d.ptr == 0 {
		return ErrInvalidDatabase
	}
	if fn == nil {
		return stderrors.New("missing function: " + name)
	}
	if fn(d.ptr, url) != 0 {
		return ErrDatabaseUpdateFailed
	}
	return nil
}

func (d *database) AddServer(url string) error {
	return d.serverOp(lib.AlpmDBAddServer, "alpm_db_add_server", url)
}

func (d *database) RemoveServer(url string) error {
	return d.serverOp(lib.AlpmDBRemoveServer, "alpm_db_remove_server", url)
}

func (d *database) GetCacheServers() []string {
	return d.getServers(lib.AlpmDBGetCacheServers, "alpm_db_get_cache_servers")
}

func (d *database) SetCacheServers(servers []string) error {
	return d.setServers(lib.AlpmDBSetCacheServers, "alpm_db_set_cache_servers", servers)
}

func (d *database) AddCacheServer(url string) error {
	return d.serverOp(lib.AlpmDBAddCacheServer, "alpm_db_add_cache_server", url)
}

func (d *database) RemoveCacheServer(url string) error {
	return d.serverOp(lib.AlpmDBRemoveCacheServer, "alpm_db_remove_cache_server", url)
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

	return checkPGPSignature(d.ptr, lib.AlpmDBCheckPGPSignature, "alpm_db_check_pgp_signature")
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
