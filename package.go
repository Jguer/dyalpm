package dyalpm

import (
	stderrors "errors"
	"io"
	"time"
	"unsafe"

	"github.com/Jguer/dyalpm/internal/lib"
	alpmlist "github.com/Jguer/dyalpm/internal/list"
)

// Package represents an ALPM package (minimal surface used by yay).
type Package interface {
	Name() string
	Version() string
	Description() string
	Architecture() string
	Size() int64
	ISize() int64
	DB() Database
	Depends() []Depend
	Provides() []Depend
	OptionalDepends() []Depend
	Groups() []string
	BuildDate() time.Time
	Reason() PkgReason
	Base() string
	ShouldIgnore() bool
}

// Validation represents package validation status
type Validation int

// PkgValidation represents package validation status
type PkgValidation int

const (
	PkgValidationUnknown   PkgValidation = 0
	PkgValidationNone      PkgValidation = (1 << 0)
	PkgValidationMD5Sum    PkgValidation = (1 << 1)
	PkgValidationSHA256Sum PkgValidation = (1 << 2)
	PkgValidationSignature PkgValidation = (1 << 3)
)

//nolint:revive // legacy name kept for compatibility
type package_ struct {
	ptr    uintptr
	handle *handle
}

func newPackage(ptr uintptr, h *handle) *package_ {
	return &package_{
		ptr:    ptr,
		handle: h,
	}
}

func (p *package_) Name() string {
	if p.ptr == 0 {
		return ""
	}
	if lib.AlpmPkgGetName == nil {
		return ""
	}
	return lib.PtrToString(lib.AlpmPkgGetName(p.ptr))
}

func (p *package_) Version() string {
	if p.ptr == 0 {
		return ""
	}
	if lib.AlpmPkgGetVersion == nil {
		return ""
	}
	return lib.PtrToString(lib.AlpmPkgGetVersion(p.ptr))
}

func (p *package_) Description() string {
	if p.ptr == 0 {
		return ""
	}
	if lib.AlpmPkgGetDesc == nil {
		return ""
	}
	return lib.PtrToString(lib.AlpmPkgGetDesc(p.ptr))
}

func (p *package_) Architecture() string {
	if p.ptr == 0 {
		return ""
	}
	if lib.AlpmPkgGetArch == nil {
		return ""
	}
	return lib.PtrToString(lib.AlpmPkgGetArch(p.ptr))
}

func (p *package_) Size() int64 {
	if p.ptr == 0 {
		return 0
	}
	if lib.AlpmPkgGetSize == nil {
		return 0
	}
	return lib.AlpmPkgGetSize(p.ptr)
}

func (p *package_) ISize() int64 {
	if p.ptr == 0 {
		return 0
	}
	if lib.AlpmPkgGetISize == nil {
		return 0
	}
	return lib.AlpmPkgGetISize(p.ptr)
}

func (p *package_) DB() Database {
	if p.ptr == 0 {
		return nil
	}
	if lib.AlpmPkgGetDB == nil {
		return nil
	}
	dbPtr := lib.AlpmPkgGetDB(p.ptr)
	if dbPtr == 0 {
		return nil
	}

	return newDatabase(dbPtr, p.handle)
}

func (p *package_) Depends() []Depend {
	deps, _ := p.getDependencyList("alpm_pkg_get_depends")
	return toDependList(deps)
}

func (p *package_) Conflicts() []Depend {
	deps, _ := p.getDependencyList("alpm_pkg_get_conflicts")
	return toDependList(deps)
}

func (p *package_) Provides() []Depend {
	deps, _ := p.getDependencyList("alpm_pkg_get_provides")
	return toDependList(deps)
}

func (p *package_) getDependencyList(funcName string) ([]Dependency, error) {
	if p.ptr == 0 {
		return nil, ErrInvalidPackage
	}
	var getFn func(uintptr) uintptr
	switch funcName {
	case "alpm_pkg_get_depends":
		getFn = lib.AlpmPkgGetDepends
	case "alpm_pkg_get_conflicts":
		getFn = lib.AlpmPkgGetConflicts
	case "alpm_pkg_get_provides":
		getFn = lib.AlpmPkgGetProvides
	case "alpm_pkg_get_optdepends":
		getFn = lib.AlpmPkgGetOptdepends
	case "alpm_pkg_get_replaces":
		getFn = lib.AlpmPkgGetReplaces
	case "alpm_pkg_get_groups":
		getFn = lib.AlpmPkgGetGroups
	case "alpm_pkg_get_licenses":
		getFn = lib.AlpmPkgGetLicenses
	case "alpm_pkg_get_files":
		getFn = lib.AlpmPkgGetFiles
	default:
		return nil, stderrors.New("missing function: " + funcName)
	}
	if getFn == nil {
		return nil, nil
	}

	listPtr := getFn(p.ptr)
	if listPtr == 0 {
		return []Dependency{}, nil
	}

	alpmList := alpmlist.NewList(listPtr)
	if alpmList == nil {
		return []Dependency{}, nil
	}

	var deps []Dependency
	for item := alpmList; item != nil && item.Ptr() != 0; item = item.Next() {
		depPtr := item.Data()
		if depPtr != 0 {
			deps = append(deps, newDependency(depPtr))
		}
	}

	return deps, nil
}

// PkgFind finds a package in a list by name
func PkgFind(pkgs []Package, name string) Package {
	if len(pkgs) == 0 {
		return nil
	}
	if lib.AlpmPkgFind == nil {
		return nil
	}

	var alpmList *alpmlist.List
	for _, p := range pkgs {
		pkgImpl, ok := p.(*package_)
		if ok {
			alpmList = alpmlist.Add(alpmList, pkgImpl.ptr)
		}
	}
	defer alpmList.Free()

	pkgPtr := lib.AlpmPkgFind(alpmList.Ptr(), name)

	if pkgPtr == 0 {
		return nil
	}

	// We need a handle to create a Package.
	// This is a bit tricky since PkgFind doesn't have a handle.
	// But the packages in the input list should have handles.
	var h *handle
	if len(pkgs) > 0 {
		pkgImpl, ok := pkgs[0].(*package_)
		if ok {
			h = pkgImpl.handle
		}
	}

	return newPackage(pkgPtr, h)
}

func (p *package_) CheckDepends() []Depend {
	deps, _ := p.getDependencyList("alpm_pkg_get_checkdepends")
	return toDependList(deps)
}

func (p *package_) MakeDepends() []Depend {
	deps, _ := p.getDependencyList("alpm_pkg_get_makedepends")
	return toDependList(deps)
}

func (p *package_) OptionalDepends() []Depend {
	deps, _ := p.getDependencyList("alpm_pkg_get_optdepends")
	return toDependList(deps)
}

func (p *package_) Replaces() []Depend {
	deps, _ := p.getDependencyList("alpm_pkg_get_replaces")
	return toDependList(deps)
}

func (p *package_) Groups() []string {
	groups, _ := p.getStringList("alpm_pkg_get_groups")
	return groups
}

func (p *package_) Licenses() []string {
	licenses, _ := p.getStringList("alpm_pkg_get_licenses")
	return licenses
}

func (p *package_) getStringList(funcName string) ([]string, error) {
	return p.getStringListWithFree(funcName, false)
}

func (p *package_) getStringListWithFree(funcName string, freeList bool) ([]string, error) {
	if p.ptr == 0 {
		return nil, ErrInvalidPackage
	}

	var getter func(uintptr) uintptr
	switch funcName {
	case "alpm_pkg_get_groups":
		getter = lib.AlpmPkgGetGroups
	case "alpm_pkg_get_licenses":
		getter = lib.AlpmPkgGetLicenses
	case "alpm_pkg_get_xdata":
		getter = lib.AlpmPkgGetXdata
	case "alpm_pkg_compute_requiredby":
		getter = lib.AlpmPkgComputeRequiredBy
	case "alpm_pkg_compute_optionalfor":
		getter = lib.AlpmPkgComputeOptionalFor
	default:
		return nil, stderrors.New("missing function: " + funcName)
	}

	if getter == nil {
		return nil, stderrors.New("missing function: " + funcName)
	}

	r1 := getter(p.ptr)
	if r1 == 0 {
		return []string{}, nil
	}

	alpmList := alpmlist.NewList(r1)
	if alpmList == nil {
		return []string{}, nil
	}
	if freeList {
		defer alpmList.Free()
	}

	var items []string
	for item := alpmList; item != nil && item.Ptr() != 0; item = item.Next() {
		strPtr := item.Data()
		if strPtr != 0 {
			items = append(items, lib.PtrToString(strPtr))
		}
	}

	return items, nil
}

// Backup represents a backup file
type Backup interface {
	Name() string
	Hash() string
}

type backup struct {
	ptr uintptr
}

func newBackup(ptr uintptr) *backup {
	return &backup{ptr: ptr}
}

func (b *backup) Name() string {
	if b.ptr == 0 {
		return ""
	}
	base := unsafe.Pointer(b.ptr)
	namePtr := *(*uintptr)(base)
	return lib.PtrToString(namePtr)
}

func (b *backup) Hash() string {
	if b.ptr == 0 {
		return ""
	}
	// hash is second pointer
	base := unsafe.Pointer(b.ptr)
	hashPtr := *(*uintptr)(unsafe.Add(base, unsafe.Sizeof(uintptr(0))))
	return lib.PtrToString(hashPtr)
}

func (p *package_) Backup() []Backup {
	if p.ptr == 0 {
		return []Backup{}
	}

	if lib.AlpmPkgGetBackup == nil {
		return []Backup{}
	}

	r1 := lib.AlpmPkgGetBackup(p.ptr)
	if r1 == 0 {
		return []Backup{}
	}

	alpmList := alpmlist.NewList(r1)
	if alpmList == nil {
		return []Backup{}
	}

	var backups []Backup
	for item := alpmList; item != nil && item.Ptr() != 0; item = item.Next() {
		dataPtr := item.Data()
		if dataPtr != 0 {
			backups = append(backups, newBackup(dataPtr))
		}
	}

	return backups
}

// File represents a file in a package
type File interface {
	Name() string
	Size() int64
	Mode() uint32
}

type file struct {
	name string
	size int64
	mode uint32
}

func (f *file) Name() string { return f.name }
func (f *file) Size() int64  { return f.size }
func (f *file) Mode() uint32 { return f.mode }

func (p *package_) Files() []File {
	if p.ptr == 0 {
		return []File{}
	}

	if lib.AlpmPkgGetFiles == nil {
		return []File{}
	}

	// alpm_filelist_t* returned
	r1 := lib.AlpmPkgGetFiles(p.ptr)
	if r1 == 0 {
		return []File{}
	}

	// Read count (size_t)
	// assuming size_t is uintptr size (safe assumption usually)
	base := unsafe.Pointer(r1)
	count := *(*uintptr)(base)
	if count == 0 {
		return []File{}
	}

	// Read files pointer (alpm_file_t*)
	filesPtr := *(*uintptr)(unsafe.Add(base, unsafe.Sizeof(uintptr(0))))
	if filesPtr == 0 {
		return []File{}
	}

	// Need alpm_file_t size
	// struct { char *name; off_t size; mode_t mode; }
	ptrSize := unsafe.Sizeof(uintptr(0))
	offSize := unsafe.Sizeof(int64(0))   // assuming off_t is 64bit
	modeSize := unsafe.Sizeof(uint32(0)) // assuming mode_t is 32bit

	structSize := ptrSize + offSize + modeSize
	// Add padding for alignment if needed
	if structSize%ptrSize != 0 {
		structSize += ptrSize - (structSize % ptrSize)
	}

	var files []File
	filesBase := unsafe.Pointer(filesPtr)
	for i := uintptr(0); i < count; i++ {
		current := unsafe.Add(filesBase, i*structSize)

		namePtr := *(*uintptr)(current)
		name := lib.PtrToString(namePtr)

		size := *(*int64)(unsafe.Add(current, ptrSize))
		mode := *(*uint32)(unsafe.Add(current, ptrSize+offSize))

		files = append(files, &file{
			name: name,
			size: size,
			mode: mode,
		})
	}

	return files
}

// PkgFrom represents where the package came from
type PkgFrom int

const (
	PkgFromFile    PkgFrom = 1
	PkgFromLocalDB PkgFrom = 2
	PkgFromSyncDB  PkgFrom = 3
)

// PkgReason represents why the package was installed
type PkgReason int

const (
	PkgReasonExplicit PkgReason = 0
	PkgReasonDepend   PkgReason = 1
	PkgReasonUnknown  PkgReason = 2
)

func (p *package_) Origin() PkgFrom {
	if p.ptr == 0 {
		return PkgFromFile // Default or error?
	}

	if lib.AlpmPkgGetOrigin == nil {
		return PkgFromFile
	}

	r1 := lib.AlpmPkgGetOrigin(p.ptr)
	return PkgFrom(r1)
}

func (p *package_) BuildDate() time.Time {
	if p.ptr == 0 {
		return time.Time{}
	}

	if lib.AlpmPkgGetBuildDate == nil {
		return time.Time{}
	}
	return toTime(int64(lib.AlpmPkgGetBuildDate(p.ptr)))
}

func (p *package_) InstallDate() time.Time {
	if p.ptr == 0 {
		return time.Time{}
	}

	if lib.AlpmPkgGetInstallDate == nil {
		return time.Time{}
	}
	return toTime(int64(lib.AlpmPkgGetInstallDate(p.ptr)))
}

func (p *package_) Reason() PkgReason {
	if p.ptr == 0 {
		return PkgReasonExplicit
	}

	if lib.AlpmPkgGetReason == nil {
		return PkgReasonExplicit
	}

	r1 := lib.AlpmPkgGetReason(p.ptr)
	return PkgReason(r1)
}

func (p *package_) HasScriptlet() bool {
	if p.ptr == 0 {
		return false
	}

	if lib.AlpmPkgHasScriptlet == nil {
		return false
	}

	return lib.AlpmPkgHasScriptlet(p.ptr) != 0
}

func (p *package_) DownloadSize() int64 {
	if p.ptr == 0 {
		return 0
	}

	if lib.AlpmPkgDownloadSize == nil {
		return 0
	}

	return lib.AlpmPkgDownloadSize(p.ptr)
}

func (p *package_) Free() error {
	if p.ptr == 0 {
		return nil
	}

	// Only free if origin is FILE
	if p.Origin() != PkgFromFile {
		return nil
	}

	if lib.AlpmPkgFree == nil {
		return stderrors.New("missing function: alpm_pkg_free")
	}

	if lib.AlpmPkgFree(p.ptr) != 0 {
		return ErrPackageFreeFailed
	}

	p.ptr = 0
	return nil
}

func (p *package_) ComputeRequiredBy() ([]string, error) {
	return p.getStringListWithFree("alpm_pkg_compute_requiredby", true)
}

func (p *package_) ComputeOptionalFor() ([]string, error) {
	return p.getStringListWithFree("alpm_pkg_compute_optionalfor", true)
}

func (p *package_) ShouldIgnore() bool {
	if p.ptr == 0 || p.handle == nil || p.handle.ptr == 0 {
		return false
	}

	if lib.AlpmPkgShouldIgnore == nil {
		return false
	}

	return lib.AlpmPkgShouldIgnore(p.handle.ptr, p.ptr) != 0
}

func (p *package_) CheckMD5Sum() error {
	if p.ptr == 0 {
		return ErrInvalidPackage
	}

	if lib.AlpmPkgCheckmd5sum == nil {
		return stderrors.New("missing function: alpm_pkg_checkmd5sum")
	}

	if lib.AlpmPkgCheckmd5sum(p.ptr) != 0 {
		return stderrors.New("MD5 sum mismatch")
	}

	return nil
}

func (p *package_) NativeHandle() Handle {
	if p.ptr == 0 {
		return nil
	}

	if lib.AlpmPkgGetHandle == nil {
		return nil
	}

	if lib.AlpmPkgGetHandle(p.ptr) == 0 {
		return nil
	}

	return p.handle
}

func (p *package_) CheckPGPSignature() (SigList, error) {
	if p.ptr == 0 {
		return SigList{}, ErrInvalidPackage
	}

	return checkPGPSignature(p.ptr, p.handle, "alpm_pkg_check_pgp_signature")
}

func (p *package_) Sig() ([]byte, error) {
	if p.ptr == 0 {
		return nil, ErrInvalidPackage
	}

	if lib.AlpmPkgGetSig == nil {
		return nil, stderrors.New("missing function: alpm_pkg_get_sig")
	}

	// alpm_pkg_get_sig signature: int alpm_pkg_get_sig(pkg, &sig, &sig_len)
	// Returns error code and writes signature bytes to output parameters
	var sigPtr uintptr
	var sigLen uintptr
	result := lib.AlpmPkgGetSig(p.ptr, &sigPtr, &sigLen)
	if result != 0 || sigPtr == 0 || sigLen == 0 {
		return nil, nil // No signature or error
	}

	// Copy the signature bytes to a Go slice
	sig := make([]byte, sigLen)
	base := unsafe.Pointer(sigPtr)
	for i := uintptr(0); i < sigLen; i++ {
		sig[i] = *(*byte)(unsafe.Add(base, i))
	}

	return sig, nil
}

// Base64Sig returns the base64-encoded package signature.
func (p *package_) Base64Sig() string {
	if p.ptr == 0 {
		return ""
	}

	if lib.AlpmPkgGetBase64Sig == nil {
		return ""
	}

	result := lib.AlpmPkgGetBase64Sig(p.ptr)
	return lib.PtrToString(result)
}

func (p *package_) PkgValidation() PkgValidation {
	if p.ptr == 0 {
		return PkgValidationUnknown
	}

	if lib.AlpmPkgGetValidation == nil {
		return PkgValidationUnknown
	}

	result := lib.AlpmPkgGetValidation(p.ptr)
	return PkgValidation(result)
}

func (p *package_) XData() []string {
	if p.ptr == 0 {
		return nil
	}

	// alpm_pkg_get_xdata returns a list owned by the package
	result, _ := p.getStringList("alpm_pkg_get_xdata")
	return result
}

func (p *package_) Contains(path string) bool {
	if p.ptr == 0 {
		return false
	}

	if lib.AlpmPkgGetFiles == nil || lib.AlpmPkgGetFilesContains == nil {
		return false
	}

	filelistPtr := lib.AlpmPkgGetFiles(p.ptr)
	if filelistPtr == 0 {
		return false
	}
	result := lib.AlpmPkgGetFilesContains(filelistPtr, path)

	return result != 0
}

func (p *package_) Changelog() (io.ReadCloser, error) {
	if p.ptr == 0 {
		return nil, ErrInvalidPackage
	}

	if lib.AlpmPkgChangelogOpen == nil {
		return nil, stderrors.New("missing function: alpm_pkg_changelog_open")
	}

	fp := lib.AlpmPkgChangelogOpen(p.ptr)
	if fp == 0 {
		return nil, stderrors.New("no changelog found")
	}

	return &changelogReader{
		pkg: p,
		fp:  fp,
	}, nil
}

func (p *package_) SyncGetNewVersion(dbsSync []Database) Package {
	if p.ptr == 0 {
		return nil
	}

	if lib.AlpmPkgSyncGetNewVersion == nil {
		return nil
	}

	var dbList *alpmlist.List
	for _, db := range dbsSync {
		dbImpl, ok := db.(*database)
		if ok {
			dbList = alpmlist.Add(dbList, dbImpl.ptr)
		}
	}
	defer dbList.Free()

	r1 := lib.AlpmPkgSyncGetNewVersion(p.ptr, dbList.Ptr())
	if r1 == 0 {
		return nil
	}

	return newPackage(r1, p.handle)
}

type changelogReader struct {
	pkg *package_
	fp  uintptr
}

func (r *changelogReader) Read(p []byte) (n int, err error) {
	if r.fp == 0 {
		return 0, io.EOF
	}

	if lib.AlpmPkgChangelogRead == nil {
		return 0, stderrors.New("missing function: alpm_pkg_changelog_read")
	}

	// size_t alpm_pkg_changelog_read(void *ptr, size_t size, const alpm_pkg_t *pkg, void *fp);
	res := lib.AlpmPkgChangelogRead(uintptr(unsafe.Pointer(&p[0])), uintptr(len(p)), r.pkg.ptr, r.fp)
	if res == 0 {
		return 0, io.EOF
	}
	return res, nil
}

func (r *changelogReader) Close() error {
	if r.fp == 0 {
		return nil
	}

	if lib.AlpmPkgChangelogClose == nil {
		return stderrors.New("missing function: alpm_pkg_changelog_close")
	}

	lib.AlpmPkgChangelogClose(r.pkg.ptr, r.fp)
	r.fp = 0
	return nil
}

func (p *package_) Validation() Validation {
	return Validation(p.PkgValidation())
}

func (p *package_) Base() string {
	if p.ptr == 0 {
		return ""
	}
	if lib.AlpmPkgGetBase == nil {
		return ""
	}
	result := lib.AlpmPkgGetBase(p.ptr)
	return lib.PtrToString(result)
}

func (p *package_) FileName() string {
	if p.ptr == 0 {
		return ""
	}
	if lib.AlpmPkgGetFilename == nil {
		return ""
	}
	result := lib.AlpmPkgGetFilename(p.ptr)
	return lib.PtrToString(result)
}

func (p *package_) Base64Signature() string {
	return p.Base64Sig()
}

func (p *package_) SHA256Sum() string {
	if p.ptr == 0 {
		return ""
	}
	if lib.AlpmPkgGetSha256sum == nil {
		return ""
	}
	result := lib.AlpmPkgGetSha256sum(p.ptr)
	return lib.PtrToString(result)
}

func (p *package_) Packager() string {
	if p.ptr == 0 {
		return ""
	}
	if lib.AlpmPkgGetPackager == nil {
		return ""
	}
	result := lib.AlpmPkgGetPackager(p.ptr)
	return lib.PtrToString(result)
}

func (p *package_) URL() string {
	if p.ptr == 0 {
		return ""
	}
	if lib.AlpmPkgGetURL == nil {
		return ""
	}
	result := lib.AlpmPkgGetURL(p.ptr)
	return lib.PtrToString(result)
}

func (p *package_) Type() string {
	// Returns "pkg" for regular packages
	return "pkg"
}

func (p *package_) ContainsFile(path string) (File, error) {
	if p.Contains(path) {
		files := p.Files()
		for _, f := range files {
			if f.Name() == path {
				return f, nil
			}
		}
	}
	return nil, stderrors.New("file not found")
}

func (p *package_) SyncNewVersion(dbs []Database) Package {
	return p.SyncGetNewVersion(dbs)
}

// toTime converts Unix timestamp to time.Time
func toTime(ts int64) time.Time {
	if ts == 0 {
		return time.Time{}
	}
	return time.Unix(ts, 0)
}
