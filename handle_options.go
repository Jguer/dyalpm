package dyalpm

import (
	"errors"
	"runtime"
	"unsafe"

	"github.com/Jguer/dyalpm/internal/dyerrors"
	"github.com/Jguer/dyalpm/internal/lib"
	alpmlist "github.com/Jguer/dyalpm/internal/list"
)

func (h *handle) SetLogFile(path string) error {
	return h.setOptionStr(lib.AlpmOptionSetLogfile, "alpm_option_set_logfile", path)
}

func (h *handle) LogFile() string {
	return h.getOptionStr(lib.AlpmOptionGetLogfile)
}

func (h *handle) SetGPGDir(path string) error {
	return h.setOptionStr(lib.AlpmOptionSetGPGDir, "alpm_option_set_gpgdir", path)
}

func (h *handle) GPGDir() string {
	return h.getOptionStr(lib.AlpmOptionGetGPGDir)
}

func (h *handle) SetUseSyslog(enable bool) error {
	return h.setOptionInt(lib.AlpmOptionSetUseSyslog, "alpm_option_set_usesyslog", lib.BoolToInt(enable))
}

func (h *handle) UseSyslog() bool {
	return lib.IntToBool(h.getOptionInt(lib.AlpmOptionGetUseSyslog))
}

func (h *handle) SetCheckSpace(enable bool) error {
	return h.setOptionInt(lib.AlpmOptionSetCheckspace, "alpm_option_set_checkspace", lib.BoolToInt(enable))
}

func (h *handle) CheckSpace() bool {
	return lib.IntToBool(h.getOptionInt(lib.AlpmOptionGetCheckspace))
}

func (h *handle) SetDBExt(ext string) error {
	return h.setOptionStr(lib.AlpmOptionSetDBExt, "alpm_option_set_dbext", ext)
}

func (h *handle) DBExt() string {
	return h.getOptionStr(lib.AlpmOptionGetDBExt)
}

func (h *handle) SetDefaultSigLevel(level int) error {
	return h.setOptionInt(lib.AlpmOptionSetDefaultSigLevel, "alpm_option_set_default_siglevel", level)
}

func (h *handle) DefaultSigLevel() int {
	return int(h.getOptionInt(lib.AlpmOptionGetDefaultSigLevel))
}

func (h *handle) SetLocalFileSigLevel(level int) error {
	return h.setOptionInt(lib.AlpmOptionSetLocalFileSigLevel, "alpm_option_set_local_file_siglevel", level)
}

func (h *handle) LocalFileSigLevel() int {
	return int(h.getOptionInt(lib.AlpmOptionGetLocalFileSigLevel))
}

func (h *handle) SetRemoteFileSigLevel(level int) error {
	return h.setOptionInt(lib.AlpmOptionSetRemoteFileSigLevel, "alpm_option_set_remote_file_siglevel", level)
}

func (h *handle) RemoteFileSigLevel() int {
	return int(h.getOptionInt(lib.AlpmOptionGetRemoteFileSigLevel))
}

func (h *handle) SetParallelDownloads(num int) error {
	return h.setOptionInt(lib.AlpmOptionSetParallelDownloads, "alpm_option_set_parallel_downloads", num)
}

func (h *handle) ParallelDownloads() int {
	return int(h.getOptionInt(lib.AlpmOptionGetParallelDownloads))
}

func (h *handle) CacheDirs() ([]string, error) {
	return h.getOptionStrList(lib.AlpmOptionGetCachedirs)
}

func (h *handle) SetCacheDirs(dirs []string) error {
	return h.setOptionStrList(lib.AlpmOptionSetCachedirs, "alpm_option_set_cachedirs", dirs)
}

func (h *handle) AddCacheDir(dir string) error {
	return h.setOptionStr(lib.AlpmOptionAddCachedir, "alpm_option_add_cachedir", dir)
}

func (h *handle) RemoveCacheDir(dir string) error {
	return h.setOptionStr(lib.AlpmOptionRemoveCachedir, "alpm_option_remove_cachedir", dir)
}

func (h *handle) HookDirs() ([]string, error) {
	return h.getOptionStrList(lib.AlpmOptionGetHookdirs)
}

func (h *handle) SetHookDirs(dirs []string) error {
	return h.setOptionStrList(lib.AlpmOptionSetHookdirs, "alpm_option_set_hookdirs", dirs)
}

func (h *handle) AddHookDir(dir string) error {
	return h.setOptionStr(lib.AlpmOptionAddHookdir, "alpm_option_add_hookdir", dir)
}

func (h *handle) RemoveHookDir(dir string) error {
	return h.setOptionStr(lib.AlpmOptionRemoveHookdir, "alpm_option_remove_hookdir", dir)
}

func (h *handle) NoUpgrades() ([]string, error) {
	return h.getOptionStrList(lib.AlpmOptionGetNoUpgrades)
}

func (h *handle) SetNoUpgrades(paths []string) error {
	return h.setOptionStrList(lib.AlpmOptionSetNoUpgrades, "alpm_option_set_noupgrades", paths)
}

func (h *handle) AddNoUpgrade(path string) error {
	return h.setOptionStr(lib.AlpmOptionAddNoupgrade, "alpm_option_add_noupgrade", path)
}

func (h *handle) RemoveNoUpgrade(path string) error {
	return h.setOptionStr(lib.AlpmOptionRemoveNoupgrade, "alpm_option_remove_noupgrade", path)
}

func (h *handle) NoExtracts() ([]string, error) {
	return h.getOptionStrList(lib.AlpmOptionGetNoextracts)
}

func (h *handle) SetNoExtracts(paths []string) error {
	return h.setOptionStrList(lib.AlpmOptionSetNoextracts, "alpm_option_set_noextracts", paths)
}

func (h *handle) AddNoExtract(path string) error {
	return h.setOptionStr(lib.AlpmOptionAddNoextract, "alpm_option_add_noextract", path)
}

func (h *handle) RemoveNoExtract(path string) error {
	return h.setOptionStr(lib.AlpmOptionRemoveNoextract, "alpm_option_remove_noextract", path)
}

func (h *handle) IgnorePkgs() ([]string, error) {
	return h.getOptionStrList(lib.AlpmOptionGetIgnorepkgs)
}

func (h *handle) SetIgnorePkgs(pkgs []string) error {
	return h.setOptionStrList(lib.AlpmOptionSetIgnorepkgs, "alpm_option_set_ignorepkgs", pkgs)
}

func (h *handle) AddIgnorePkg(pkg string) error {
	return h.setOptionStr(lib.AlpmOptionAddIgnorepkg, "alpm_option_add_ignorepkg", pkg)
}

func (h *handle) RemoveIgnorePkg(pkg string) error {
	return h.setOptionStr(lib.AlpmOptionRemoveIgnorepkg, "alpm_option_remove_ignorepkg", pkg)
}

func (h *handle) IgnoreGroups() ([]string, error) {
	return h.getOptionStrList(lib.AlpmOptionGetIgnoregroups)
}

func (h *handle) SetIgnoreGroups(groups []string) error {
	return h.setOptionStrList(lib.AlpmOptionSetIgnoregroups, "alpm_option_set_ignoregroups", groups)
}

func (h *handle) AddIgnoreGroup(group string) error {
	return h.setOptionStr(lib.AlpmOptionAddIgnoregroup, "alpm_option_add_ignoregroup", group)
}

func (h *handle) RemoveIgnoreGroup(group string) error {
	return h.setOptionStr(lib.AlpmOptionRemoveIgnoregroup, "alpm_option_remove_ignoregroup", group)
}

func (h *handle) OverwriteFiles() ([]string, error) {
	return h.getOptionStrList(lib.AlpmOptionGetOverwriteFiles)
}

func (h *handle) SetOverwriteFiles(globs []string) error {
	return h.setOptionStrList(lib.AlpmOptionSetOverwriteFiles, "alpm_option_set_overwrite_files", globs)
}

func (h *handle) AddOverwriteFile(glob string) error {
	return h.setOptionStr(lib.AlpmOptionAddOverwriteFile, "alpm_option_add_overwrite_file", glob)
}

func (h *handle) RemoveOverwriteFile(glob string) error {
	return h.setOptionStr(lib.AlpmOptionRemoveOverwriteFile, "alpm_option_remove_overwrite_file", glob)
}

func (h *handle) MatchNoUpgrade(path string) (int, error) {
	return h.matchOption(lib.AlpmOptionMatchNoUpgrade, "alpm_option_match_noupgrade", path)
}

func (h *handle) MatchNoExtract(path string) (int, error) {
	return h.matchOption(lib.AlpmOptionMatchNoextract, "alpm_option_match_noextract", path)
}

func (h *handle) SetSandboxUser(user string) error {
	return h.setOptionStr(lib.AlpmOptionSetSandboxuser, "alpm_option_set_sandboxuser", user)
}

func (h *handle) SandboxUser() string {
	return h.getOptionStr(lib.AlpmOptionGetSandboxuser)
}

func (h *handle) SetDisableDLTimeout(disable bool) error {
	return h.setOptionInt(lib.AlpmOptionSetDisableDlTimeout, "alpm_option_set_disable_dl_timeout", lib.BoolToInt(disable))
}

func (h *handle) SetDisableSandbox(disable bool) error {
	return h.setOptionInt(lib.AlpmOptionSetDisableSandbox, "alpm_option_set_disable_sandbox", lib.BoolToInt(disable))
}

func (h *handle) AssumeInstalled() ([]Dependency, error) {
	if h.ptr == 0 {
		return nil, dyerrors.ErrHandleNull
	}
	if lib.AlpmOptionGetAssumeInstalled == nil {
		return []Dependency{}, nil
	}

	deps := collectList(alpmlist.NewList(lib.AlpmOptionGetAssumeInstalled(h.ptr)), func(ptr uintptr) Dependency {
		return newDependency(ptr)
	})
	if deps == nil {
		return []Dependency{}, nil
	}
	return deps, nil
}

func (h *handle) SetAssumeInstalled(deps []Dependency) error {
	if h.ptr == 0 {
		return dyerrors.ErrHandleNull
	}
	if lib.AlpmOptionSetNoassumeInstalled == nil {
		return unavailable("alpm_option_set_assumeinstalled")
	}

	var alpmList *alpmlist.List
	for _, d := range deps {
		if depImpl, ok := d.(*dependency); ok {
			alpmList = alpmlist.Add(alpmList, depImpl.ptr)
		}
	}
	defer alpmList.Free()

	if lib.AlpmOptionSetNoassumeInstalled(h.ptr, alpmList.Ptr()) != 0 {
		return dyerrors.NewError(h.Errno(), "failed to set option dependency list")
	}
	return nil
}

func (h *handle) AddAssumeInstalled(dep Dependency) error {
	return h.assumeInstalledOp(lib.AlpmOptionAddAssumeInstalled, "alpm_option_add_assumeinstalled", dep)
}

func (h *handle) RemoveAssumeInstalled(dep Dependency) error {
	return h.assumeInstalledOp(lib.AlpmOptionRemoveAssumeInstalled, "alpm_option_remove_assumeinstalled", dep)
}

func (h *handle) Architectures() ([]string, error) {
	return h.getOptionStrList(lib.AlpmOptionGetArchitectures)
}

func (h *handle) SetArchitectures(archs []string) error {
	return h.setOptionStrList(lib.AlpmOptionSetArchitectures, "alpm_option_set_architectures", archs)
}

func (h *handle) AddArchitecture(arch string) error {
	return h.setOptionStr(lib.AlpmOptionAddArchitecture, "alpm_option_add_architecture", arch)
}

func (h *handle) RemoveArchitecture(arch string) error {
	return h.setOptionStr(lib.AlpmOptionRemoveArchitecture, "alpm_option_remove_architecture", arch)
}

// Helper methods. Each takes the bound libalpm function and its C name (for errors).

func unavailable(name string) error {
	return errors.New("function unavailable: " + name)
}

// cStringList builds an alpm_list_t of C strings. The caller frees the list and
// keeps the returned buffers alive until the C call returns.
func cStringList(values []string) (*alpmlist.List, [][]byte) {
	var l *alpmlist.List
	bufs := make([][]byte, 0, len(values))
	for _, v := range values {
		cs := lib.CString(v)
		bufs = append(bufs, cs)
		l = alpmlist.Add(l, uintptr(unsafe.Pointer(&cs[0])))
	}
	return l, bufs
}

func (h *handle) assumeInstalledOp(fn func(uintptr, uintptr) int32, name string, dep Dependency) error {
	if h.ptr == 0 {
		return dyerrors.ErrHandleNull
	}
	if fn == nil {
		return unavailable(name)
	}
	d, ok := dep.(*dependency)
	if !ok {
		return ErrInvalidPackage
	}
	if fn(h.ptr, d.ptr) != 0 {
		return dyerrors.NewError(h.Errno(), "failed to update assume-installed dependency")
	}
	return nil
}

func (h *handle) matchOption(fn func(uintptr, string) int32, name, path string) (int, error) {
	if h.ptr == 0 {
		return 0, dyerrors.ErrHandleNull
	}
	if fn == nil {
		return 0, unavailable(name)
	}
	return int(fn(h.ptr, path)), nil
}

func (h *handle) setOptionStr(fn func(uintptr, string) int32, name, value string) error {
	if h.ptr == 0 {
		return dyerrors.ErrHandleNull
	}
	if fn == nil {
		return unavailable(name)
	}
	if fn(h.ptr, value) != 0 {
		return dyerrors.NewError(h.Errno(), "failed to set option")
	}
	return nil
}

func (h *handle) getOptionStr(fn func(uintptr) uintptr) string {
	if h.ptr == 0 || fn == nil {
		return ""
	}
	return lib.PtrToString(fn(h.ptr))
}

func (h *handle) setOptionInt(fn func(uintptr, int32) int32, name string, value int) error {
	if h.ptr == 0 {
		return dyerrors.ErrHandleNull
	}
	if fn == nil {
		return unavailable(name)
	}
	if fn(h.ptr, clampIntToInt32(value)) != 0 {
		return dyerrors.NewError(h.Errno(), "failed to set option")
	}
	return nil
}

func (h *handle) getOptionInt(fn func(uintptr) int32) int32 {
	if h.ptr == 0 || fn == nil {
		return 0
	}
	return fn(h.ptr)
}

func (h *handle) getOptionStrList(fn func(uintptr) uintptr) ([]string, error) {
	if h.ptr == 0 {
		return nil, dyerrors.ErrHandleNull
	}
	if fn == nil {
		return []string{}, nil
	}
	items := collectList(alpmlist.NewList(fn(h.ptr)), lib.PtrToString)
	if items == nil {
		return []string{}, nil
	}
	return items, nil
}

func (h *handle) setOptionStrList(fn func(uintptr, uintptr) int32, name string, values []string) error {
	if h.ptr == 0 {
		return dyerrors.ErrHandleNull
	}
	if fn == nil {
		return unavailable(name)
	}

	alpmList, bufs := cStringList(values)
	defer alpmList.Free()

	r1 := fn(h.ptr, alpmList.Ptr())
	runtime.KeepAlive(bufs)

	if r1 != 0 {
		return dyerrors.NewError(h.Errno(), "failed to set option list")
	}
	return nil
}
