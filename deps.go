package dyalpm

import (
	stderrors "errors"
	"unsafe"

	"github.com/Jguer/dyalpm/internal/dyerrors"
	"github.com/Jguer/dyalpm/internal/lib"
	alpmlist "github.com/Jguer/dyalpm/internal/list"
)

// DepMod represents dependency version comparison mode
type DepMod int

const (
	DepModAny DepMod = iota + 1
	DepModEQ
	DepModGE
	DepModLE
	DepModGT
	DepModLT
)

// Dependency represents a package dependency
type Dependency interface {
	GetName() string
	GetVersion() string
	GetMod() DepMod
	ComputeString() string
	Free()
}

// Depend is a value type representation of a dependency for consumer code.
type Depend struct {
	Name        string
	Version     string
	Description string
	NameHash    uint64
	Mod         DepMod
}

// String returns the computed string representation of the dependency.
func (d Depend) String() string {
	if d.Version == "" {
		return d.Name
	}
	modStr := ""
	switch d.Mod {
	case DepModEQ:
		modStr = "="
	case DepModGE:
		modStr = ">="
	case DepModLE:
		modStr = "<="
	case DepModGT:
		modStr = ">"
	case DepModLT:
		modStr = "<"
	}
	return d.Name + modStr + d.Version
}

type dependency struct {
	ptr   uintptr
	owned bool
}

func newDependency(ptr uintptr) *dependency {
	return &dependency{
		ptr:   ptr,
		owned: false,
	}
}

// cDepend mirrors alpm_depend_t.
type cDepend struct {
	Name     uintptr
	Version  uintptr
	Desc     uintptr
	NameHash uint64
	Mod      int32
}

func (d *dependency) c() *cDepend {
	return (*cDepend)(unsafe.Pointer(d.ptr))
}

func (d *dependency) GetName() string {
	if d.ptr == 0 {
		return ""
	}
	return lib.PtrToString(d.c().Name)
}

func (d *dependency) GetVersion() string {
	if d.ptr == 0 {
		return ""
	}
	return lib.PtrToString(d.c().Version)
}

func (d *dependency) GetMod() DepMod {
	if d.ptr == 0 {
		return DepModAny
	}
	return DepMod(d.c().Mod)
}

func (d *dependency) toDepend() Depend {
	if d.ptr == 0 {
		return Depend{}
	}
	c := d.c()
	return Depend{
		Name:        lib.PtrToString(c.Name),
		Version:     lib.PtrToString(c.Version),
		Description: lib.PtrToString(c.Desc),
		NameHash:    c.NameHash,
		Mod:         DepMod(c.Mod),
	}
}

func (d *dependency) ComputeString() string {
	if d.ptr == 0 {
		return ""
	}
	if lib.AlpmDepComputeString == nil {
		return ""
	}

	r1 := lib.AlpmDepComputeString(d.ptr)
	if r1 == 0 {
		return ""
	}

	// The string returned by alpm_dep_compute_string MUST be freed by the caller.
	res := lib.PtrToString(r1)
	lib.Free(r1)
	return res
}

func (d *dependency) Free() {
	if d.ptr == 0 || !d.owned {
		return
	}
	if lib.AlpmDepFree == nil {
		return
	}
	lib.AlpmDepFree(d.ptr)
	d.ptr = 0
}

// DepFromString creates a dependency from a string
func DepFromString(depstring string) (Dependency, error) {
	if lib.AlpmDepFromString == nil {
		return nil, stderrors.New("missing function: alpm_dep_from_string")
	}

	r1 := lib.AlpmDepFromString(depstring)

	if r1 == 0 {
		return nil, ErrInvalidDependency
	}

	dep := newDependency(r1)
	dep.owned = true // alpm_dep_from_string returns a newly allocated struct
	return dep, nil
}

// DepMissing represents a missing dependency
type DepMissing interface {
	GetTarget() string
	GetDepend() Dependency
	GetCausingPkg() string
	Free()
}

type depMissing struct {
	ptr uintptr
}

func newDepMissing(ptr uintptr) *depMissing {
	return &depMissing{
		ptr: ptr,
	}
}

func (d *depMissing) GetTarget() string {
	if d.ptr == 0 {
		return ""
	}
	// struct _alpm_depmissing_t { char *target; alpm_depend_t *depend; char *causingpkg; }
	base := unsafe.Pointer(d.ptr)
	targetPtr := *(*uintptr)(base)
	return lib.PtrToString(targetPtr)
}

func (d *depMissing) GetDepend() Dependency {
	if d.ptr == 0 {
		return nil
	}
	base := unsafe.Pointer(d.ptr)
	depPtr := *(*uintptr)(unsafe.Add(base, unsafe.Sizeof(uintptr(0))))
	return newDependency(depPtr)
}

func (d *depMissing) GetCausingPkg() string {
	if d.ptr == 0 {
		return ""
	}
	base := unsafe.Pointer(d.ptr)
	pkgPtr := *(*uintptr)(unsafe.Add(base, 2*unsafe.Sizeof(uintptr(0))))
	return lib.PtrToString(pkgPtr)
}

func (d *depMissing) Free() {
	if d.ptr == 0 {
		return
	}
	if lib.AlpmDepmissingFree == nil {
		return
	}
	lib.AlpmDepmissingFree(d.ptr)
	d.ptr = 0
}

// Conflict represents a package conflict
type Conflict interface {
	GetPackage1() string
	GetPackage2() string
	GetReason() Dependency
	Free()
}

type conflict struct {
	ptr uintptr
}

func newConflict(ptr uintptr) *conflict {
	return &conflict{
		ptr: ptr,
	}
}

// field returns the i-th pointer of
// struct _alpm_conflict_t { alpm_pkg_t *package1; alpm_pkg_t *package2; alpm_depend_t *reason; }
func (c *conflict) field(i uintptr) uintptr {
	return *(*uintptr)(unsafe.Add(unsafe.Pointer(c.ptr), i*unsafe.Sizeof(uintptr(0))))
}

func (c *conflict) GetPackage1() string {
	if c.ptr == 0 {
		return ""
	}
	return newPackage(c.field(0), nil).Name()
}

func (c *conflict) GetPackage2() string {
	if c.ptr == 0 {
		return ""
	}
	return newPackage(c.field(1), nil).Name()
}

func (c *conflict) GetReason() Dependency {
	if c.ptr == 0 {
		return nil
	}
	return newDependency(c.field(2))
}

// toDependList converts a Dependency slice into a Depend slice.
func toDependList(deps []Dependency) []Depend {
	result := make([]Depend, len(deps))
	for i, dep := range deps {
		result[i] = dep.(*dependency).toDepend()
	}
	return result
}

func (c *conflict) Free() {
	if c.ptr == 0 {
		return
	}
	if lib.AlpmConflictFree == nil {
		return
	}
	lib.AlpmConflictFree(c.ptr)
	c.ptr = 0
}

// FileConflictType represents the type of file conflict
type FileConflictType int

const (
	FileConflictTarget     FileConflictType = 1
	FileConflictFilesystem FileConflictType = 2
)

// FileConflict represents a file conflict
type FileConflict interface {
	GetTarget() string
	GetType() FileConflictType
	GetFile() string
	GetCTarget() string
	Free()
}

type fileConflict struct {
	ptr uintptr
}

func newFileConflict(ptr uintptr) *fileConflict {
	return &fileConflict{
		ptr: ptr,
	}
}

func (f *fileConflict) GetTarget() string {
	if f.ptr == 0 {
		return ""
	}
	// struct _alpm_fileconflict_t { char *target; type; char *file; char *ctarget; }
	base := unsafe.Pointer(f.ptr)
	targetPtr := *(*uintptr)(base)
	return lib.PtrToString(targetPtr)
}

func (f *fileConflict) GetType() FileConflictType {
	if f.ptr == 0 {
		return 0
	}
	base := unsafe.Pointer(f.ptr)
	// alpm_fileconflicttype_t is a 4-byte enum followed by padding.
	typeVal := *(*int32)(unsafe.Add(base, unsafe.Sizeof(uintptr(0))))
	return FileConflictType(typeVal)
}

func (f *fileConflict) GetFile() string {
	if f.ptr == 0 {
		return ""
	}
	base := unsafe.Pointer(f.ptr)
	filePtr := *(*uintptr)(unsafe.Add(base, 2*unsafe.Sizeof(uintptr(0))))
	return lib.PtrToString(filePtr)
}

func (f *fileConflict) GetCTarget() string {
	if f.ptr == 0 {
		return ""
	}
	base := unsafe.Pointer(f.ptr)
	ctargetPtr := *(*uintptr)(unsafe.Add(base, 3*unsafe.Sizeof(uintptr(0))))
	return lib.PtrToString(ctargetPtr)
}

func (f *fileConflict) Free() {
	if f.ptr == 0 {
		return
	}
	if lib.AlpmFileConflictFree == nil {
		return
	}
	lib.AlpmFileConflictFree(f.ptr)
	f.ptr = 0
}

// Resolution functions

func (h *handle) CheckDeps(pkgs []Package, remPkgs []Package, upgradePkgs []Package, reverseDeps bool) ([]DepMissing, error) {
	if h.ptr == 0 {
		return nil, dyerrors.ErrHandleNull
	}
	if lib.AlpmCheckDeps == nil {
		return nil, stderrors.New("missing function: alpm_checkdeps")
	}

	var pkgList, remPkgList, upgradePkgList *alpmlist.List
	for _, p := range pkgs {
		pkgImpl, _ := p.(*package_)
		pkgList = alpmlist.Add(pkgList, pkgImpl.ptr)
	}
	defer pkgList.Free()

	for _, p := range remPkgs {
		pkgImpl, _ := p.(*package_)
		remPkgList = alpmlist.Add(remPkgList, pkgImpl.ptr)
	}
	defer remPkgList.Free()

	for _, p := range upgradePkgs {
		pkgImpl, _ := p.(*package_)
		upgradePkgList = alpmlist.Add(upgradePkgList, pkgImpl.ptr)
	}
	defer upgradePkgList.Free()

	rev := int32(0)
	if reverseDeps {
		rev = 1
	}
	r1 := lib.AlpmCheckDeps(h.ptr, pkgList.Ptr(), remPkgList.Ptr(), upgradePkgList.Ptr(), rev)
	if r1 == 0 {
		return []DepMissing{}, nil
	}

	resList := alpmlist.NewList(r1)
	defer resList.Free()

	var missings []DepMissing
	for item := resList; item != nil && item.Ptr() != 0; item = item.Next() {
		ptr := item.Data()
		if ptr != 0 {
			missings = append(missings, newDepMissing(ptr))
		}
	}

	return missings, nil
}

func (h *handle) CheckConflicts(pkgs []Package) ([]Conflict, error) {
	if h.ptr == 0 {
		return nil, dyerrors.ErrHandleNull
	}
	if lib.AlpmCheckConflicts == nil {
		return nil, stderrors.New("missing function: alpm_checkconflicts")
	}

	var pkgList *alpmlist.List
	for _, p := range pkgs {
		pkgImpl, _ := p.(*package_)
		pkgList = alpmlist.Add(pkgList, pkgImpl.ptr)
	}
	defer pkgList.Free()

	r1 := lib.AlpmCheckConflicts(h.ptr, pkgList.Ptr())
	if r1 == 0 {
		return []Conflict{}, nil
	}

	resList := alpmlist.NewList(r1)
	defer resList.Free()

	var conflicts []Conflict
	for item := resList; item != nil && item.Ptr() != 0; item = item.Next() {
		ptr := item.Data()
		if ptr != 0 {
			conflicts = append(conflicts, newConflict(ptr))
		}
	}

	return conflicts, nil
}

// FindSatisfier finds a package that satisfies a dependency in a list of packages
func FindSatisfier(pkgs []Package, depstring string) Package {
	if lib.AlpmFindSatisfier == nil {
		return nil
	}

	var pkgList *alpmlist.List
	var h *handle
	for _, p := range pkgs {
		pkgImpl, ok := p.(*package_)
		if ok {
			pkgList = alpmlist.Add(pkgList, pkgImpl.ptr)
			if h == nil {
				h = pkgImpl.handle
			}
		}
	}
	defer pkgList.Free()

	r1 := lib.AlpmFindSatisfier(pkgList.Ptr(), depstring)

	if r1 == 0 {
		return nil
	}

	return newPackage(r1, h)
}

func (h *handle) FindDBSatisfier(dbs []Database, depstring string) Package {
	if h.ptr == 0 {
		return nil
	}
	if lib.AlpmFindDBSatisfier == nil {
		return nil
	}

	var dbList *alpmlist.List
	for _, db := range dbs {
		dbImpl, ok := db.(*database)
		if ok {
			dbList = alpmlist.Add(dbList, dbImpl.ptr)
		}
	}
	defer dbList.Free()

	r1 := lib.AlpmFindDBSatisfier(h.ptr, dbList.Ptr(), depstring)

	if r1 == 0 {
		return nil
	}

	return newPackage(r1, h)
}
