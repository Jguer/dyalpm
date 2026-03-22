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
	return h.setOptionStr("alpm_option_set_logfile", path)
}

func (h *handle) LogFile() string {
	return h.getOptionStr("alpm_option_get_logfile")
}

func (h *handle) SetGPGDir(path string) error {
	return h.setOptionStr("alpm_option_set_gpgdir", path)
}

func (h *handle) GPGDir() string {
	return h.getOptionStr("alpm_option_get_gpgdir")
}

func (h *handle) SetUseSyslog(enable bool) error {
	value := 0
	if enable {
		value = 1
	}
	return h.setOptionInt("alpm_option_set_usesyslog", value)
}

func (h *handle) UseSyslog() bool {
	val := h.getOptionInt("alpm_option_get_usesyslog")
	return lib.IntToBool(val)
}

func (h *handle) SetCheckSpace(enable bool) error {
	value := 0
	if enable {
		value = 1
	}
	return h.setOptionInt("alpm_option_set_checkspace", value)
}

func (h *handle) CheckSpace() bool {
	val := h.getOptionInt("alpm_option_get_checkspace")
	return lib.IntToBool(val)
}

func (h *handle) SetDBExt(ext string) error {
	return h.setOptionStr("alpm_option_set_dbext", ext)
}

func (h *handle) DBExt() string {
	return h.getOptionStr("alpm_option_get_dbext")
}

func (h *handle) SetDefaultSigLevel(level int) error {
	return h.setOptionInt("alpm_option_set_default_siglevel", level)
}

func (h *handle) DefaultSigLevel() int {
	return int(h.getOptionInt("alpm_option_get_default_siglevel"))
}

func (h *handle) SetLocalFileSigLevel(level int) error {
	return h.setOptionInt("alpm_option_set_local_file_siglevel", level)
}

func (h *handle) LocalFileSigLevel() int {
	return int(h.getOptionInt("alpm_option_get_local_file_siglevel"))
}

func (h *handle) SetRemoteFileSigLevel(level int) error {
	return h.setOptionInt("alpm_option_set_remote_file_siglevel", level)
}

func (h *handle) RemoteFileSigLevel() int {
	return int(h.getOptionInt("alpm_option_get_remote_file_siglevel"))
}

func (h *handle) SetParallelDownloads(num int) error {
	return h.setOptionInt("alpm_option_set_parallel_downloads", num)
}

func (h *handle) ParallelDownloads() int {
	return int(h.getOptionInt("alpm_option_get_parallel_downloads"))
}

func (h *handle) CacheDirs() ([]string, error) {
	return h.getOptionStrList("alpm_option_get_cachedirs")
}

func (h *handle) SetCacheDirs(dirs []string) error {
	return h.setOptionStrList("alpm_option_set_cachedirs", dirs)
}

func (h *handle) AddCacheDir(dir string) error {
	return h.setOptionStr("alpm_option_add_cachedir", dir)
}

func (h *handle) RemoveCacheDir(dir string) error {
	return h.setOptionStr("alpm_option_remove_cachedir", dir)
}

func (h *handle) HookDirs() ([]string, error) {
	return h.getOptionStrList("alpm_option_get_hookdirs")
}

func (h *handle) SetHookDirs(dirs []string) error {
	return h.setOptionStrList("alpm_option_set_hookdirs", dirs)
}

func (h *handle) AddHookDir(dir string) error {
	return h.setOptionStr("alpm_option_add_hookdir", dir)
}

func (h *handle) RemoveHookDir(dir string) error {
	return h.setOptionStr("alpm_option_remove_hookdir", dir)
}

func (h *handle) NoUpgrades() ([]string, error) {
	return h.getOptionStrList("alpm_option_get_noupgrades")
}

func (h *handle) SetNoUpgrades(paths []string) error {
	return h.setOptionStrList("alpm_option_set_noupgrades", paths)
}

func (h *handle) AddNoUpgrade(path string) error {
	return h.setOptionStr("alpm_option_add_noupgrade", path)
}

func (h *handle) RemoveNoUpgrade(path string) error {
	return h.setOptionStr("alpm_option_remove_noupgrade", path)
}

func (h *handle) NoExtracts() ([]string, error) {
	return h.getOptionStrList("alpm_option_get_noextracts")
}

func (h *handle) SetNoExtracts(paths []string) error {
	return h.setOptionStrList("alpm_option_set_noextracts", paths)
}

func (h *handle) AddNoExtract(path string) error {
	return h.setOptionStr("alpm_option_add_noextract", path)
}

func (h *handle) RemoveNoExtract(path string) error {
	return h.setOptionStr("alpm_option_remove_noextract", path)
}

func (h *handle) IgnorePkgs() ([]string, error) {
	return h.getOptionStrList("alpm_option_get_ignorepkgs")
}

func (h *handle) SetIgnorePkgs(pkgs []string) error {
	return h.setOptionStrList("alpm_option_set_ignorepkgs", pkgs)
}

func (h *handle) AddIgnorePkg(pkg string) error {
	return h.setOptionStr("alpm_option_add_ignorepkg", pkg)
}

func (h *handle) RemoveIgnorePkg(pkg string) error {
	return h.setOptionStr("alpm_option_remove_ignorepkg", pkg)
}

func (h *handle) IgnoreGroups() ([]string, error) {
	return h.getOptionStrList("alpm_option_get_ignoregroups")
}

func (h *handle) SetIgnoreGroups(groups []string) error {
	return h.setOptionStrList("alpm_option_set_ignoregroups", groups)
}

func (h *handle) AddIgnoreGroup(group string) error {
	return h.setOptionStr("alpm_option_add_ignoregroup", group)
}

func (h *handle) RemoveIgnoreGroup(group string) error {
	return h.setOptionStr("alpm_option_remove_ignoregroup", group)
}

func (h *handle) OverwriteFiles() ([]string, error) {
	return h.getOptionStrList("alpm_option_get_overwrite_files")
}

func (h *handle) SetOverwriteFiles(globs []string) error {
	return h.setOptionStrList("alpm_option_set_overwrite_files", globs)
}

func (h *handle) AddOverwriteFile(glob string) error {
	return h.setOptionStr("alpm_option_add_overwrite_file", glob)
}

func (h *handle) RemoveOverwriteFile(glob string) error {
	return h.setOptionStr("alpm_option_remove_overwrite_file", glob)
}

func (h *handle) MatchNoUpgrade(path string) (int, error) {
	return h.matchOption("alpm_option_match_noupgrade", path)
}

func (h *handle) MatchNoExtract(path string) (int, error) {
	return h.matchOption("alpm_option_match_noextract", path)
}

func (h *handle) SetSandboxUser(user string) error {
	return h.setOptionStr("alpm_option_set_sandboxuser", user)
}

func (h *handle) SandboxUser() string {
	return h.getOptionStr("alpm_option_get_sandboxuser")
}

func (h *handle) SetDisableDLTimeout(disable bool) error {
	value := 0
	if disable {
		value = 1
	}
	return h.setOptionInt("alpm_option_set_disable_dl_timeout", value)
}

func (h *handle) SetDisableSandbox(disable bool) error {
	value := 0
	if disable {
		value = 1
	}
	return h.setOptionInt("alpm_option_set_disable_sandbox", value)
}

func (h *handle) AssumeInstalled() ([]Dependency, error) {
	if h.ptr == 0 {
		return nil, dyerrors.ErrHandleNull
	}
	if lib.AlpmOptionGetAssumeInstalled == nil {
		return []Dependency{}, nil
	}

	r1 := lib.AlpmOptionGetAssumeInstalled(h.ptr)
	if r1 == 0 {
		return []Dependency{}, nil
	}

	alpmList := alpmlist.NewList(r1)
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

func (h *handle) SetAssumeInstalled(deps []Dependency) error {
	return h.setOptionDepList("alpm_option_set_assumeinstalled", deps)
}

func (h *handle) AddAssumeInstalled(dep Dependency) error {
	if h.ptr == 0 {
		return dyerrors.ErrHandleNull
	}
	if lib.AlpmOptionAddAssumeInstalled == nil {
		return errors.New("function unavailable: alpm_option_add_assumeinstalled")
	}

	d, ok := dep.(*dependency)
	if !ok {
		return ErrInvalidPackage
	}

	r1 := lib.AlpmOptionAddAssumeInstalled(h.ptr, d.ptr)
	if r1 != 0 {
		return dyerrors.NewError(h.Errno(), "failed to add assume-installed dependency")
	}

	return nil
}

func (h *handle) RemoveAssumeInstalled(dep Dependency) error {
	if h.ptr == 0 {
		return dyerrors.ErrHandleNull
	}
	if lib.AlpmOptionRemoveAssumeInstalled == nil {
		return errors.New("function unavailable: alpm_option_remove_assumeinstalled")
	}

	d, ok := dep.(*dependency)
	if !ok {
		return ErrInvalidPackage
	}

	r1 := lib.AlpmOptionRemoveAssumeInstalled(h.ptr, d.ptr)
	if r1 != 0 {
		return dyerrors.NewError(h.Errno(), "failed to remove assume-installed dependency")
	}

	return nil
}

func (h *handle) Architectures() ([]string, error) {
	return h.getOptionStrList("alpm_option_get_architectures")
}

func (h *handle) SetArchitectures(archs []string) error {
	return h.setOptionStrList("alpm_option_set_architectures", archs)
}

func (h *handle) AddArchitecture(arch string) error {
	return h.setOptionStr("alpm_option_add_architecture", arch)
}

func (h *handle) RemoveArchitecture(arch string) error {
	return h.setOptionStr("alpm_option_remove_architecture", arch)
}

// Helper methods

func (h *handle) matchOption(funcName, path string) (int, error) {
	if h.ptr == 0 {
		return 0, dyerrors.ErrHandleNull
	}

	switch funcName {
	case "alpm_option_match_noupgrade":
		if lib.AlpmOptionMatchNoUpgrade == nil {
			return 0, errors.New("function unavailable: alpm_option_match_noupgrade")
		}
		return int(lib.AlpmOptionMatchNoUpgrade(h.ptr, path)), nil
	case "alpm_option_match_noextract":
		if lib.AlpmOptionMatchNoextract == nil {
			return 0, errors.New("function unavailable: alpm_option_match_noextract")
		}
		return int(lib.AlpmOptionMatchNoextract(h.ptr, path)), nil
	default:
		return 0, errors.New("unsupported function: " + funcName)
	}
}

func (h *handle) setOptionStr(funcName, value string) error {
	if h.ptr == 0 {
		return dyerrors.ErrHandleNull
	}

	var result int32
	switch funcName {
	case "alpm_option_set_logfile":
		if lib.AlpmOptionSetLogfile == nil {
			return errors.New("function unavailable: alpm_option_set_logfile")
		}
		result = lib.AlpmOptionSetLogfile(h.ptr, value)
	case "alpm_option_set_gpgdir":
		if lib.AlpmOptionSetGPGDir == nil {
			return errors.New("function unavailable: alpm_option_set_gpgdir")
		}
		result = lib.AlpmOptionSetGPGDir(h.ptr, value)
	case "alpm_option_set_dbext":
		if lib.AlpmOptionSetDBExt == nil {
			return errors.New("function unavailable: alpm_option_set_dbext")
		}
		result = lib.AlpmOptionSetDBExt(h.ptr, value)
	case "alpm_option_add_cachedir":
		if lib.AlpmOptionAddCachedir == nil {
			return errors.New("function unavailable: alpm_option_add_cachedir")
		}
		result = lib.AlpmOptionAddCachedir(h.ptr, value)
	case "alpm_option_remove_cachedir":
		if lib.AlpmOptionRemoveCachedir == nil {
			return errors.New("function unavailable: alpm_option_remove_cachedir")
		}
		result = lib.AlpmOptionRemoveCachedir(h.ptr, value)
	case "alpm_option_add_hookdir":
		if lib.AlpmOptionAddHookdir == nil {
			return errors.New("function unavailable: alpm_option_add_hookdir")
		}
		result = lib.AlpmOptionAddHookdir(h.ptr, value)
	case "alpm_option_remove_hookdir":
		if lib.AlpmOptionRemoveHookdir == nil {
			return errors.New("function unavailable: alpm_option_remove_hookdir")
		}
		result = lib.AlpmOptionRemoveHookdir(h.ptr, value)
	case "alpm_option_add_noupgrade":
		if lib.AlpmOptionAddNoupgrade == nil {
			return errors.New("function unavailable: alpm_option_add_noupgrade")
		}
		result = lib.AlpmOptionAddNoupgrade(h.ptr, value)
	case "alpm_option_remove_noupgrade":
		if lib.AlpmOptionRemoveNoupgrade == nil {
			return errors.New("function unavailable: alpm_option_remove_noupgrade")
		}
		result = lib.AlpmOptionRemoveNoupgrade(h.ptr, value)
	case "alpm_option_add_noextract":
		if lib.AlpmOptionAddNoextract == nil {
			return errors.New("function unavailable: alpm_option_add_noextract")
		}
		result = lib.AlpmOptionAddNoextract(h.ptr, value)
	case "alpm_option_remove_noextract":
		if lib.AlpmOptionRemoveNoextract == nil {
			return errors.New("function unavailable: alpm_option_remove_noextract")
		}
		result = lib.AlpmOptionRemoveNoextract(h.ptr, value)
	case "alpm_option_add_ignorepkg":
		if lib.AlpmOptionAddIgnorepkg == nil {
			return errors.New("function unavailable: alpm_option_add_ignorepkg")
		}
		result = lib.AlpmOptionAddIgnorepkg(h.ptr, value)
	case "alpm_option_remove_ignorepkg":
		if lib.AlpmOptionRemoveIgnorepkg == nil {
			return errors.New("function unavailable: alpm_option_remove_ignorepkg")
		}
		result = lib.AlpmOptionRemoveIgnorepkg(h.ptr, value)
	case "alpm_option_add_ignoregroup":
		if lib.AlpmOptionAddIgnoregroup == nil {
			return errors.New("function unavailable: alpm_option_add_ignoregroup")
		}
		result = lib.AlpmOptionAddIgnoregroup(h.ptr, value)
	case "alpm_option_remove_ignoregroup":
		if lib.AlpmOptionRemoveIgnoregroup == nil {
			return errors.New("function unavailable: alpm_option_remove_ignoregroup")
		}
		result = lib.AlpmOptionRemoveIgnoregroup(h.ptr, value)
	case "alpm_option_add_overwrite_file":
		if lib.AlpmOptionAddOverwriteFile == nil {
			return errors.New("function unavailable: alpm_option_add_overwrite_file")
		}
		result = lib.AlpmOptionAddOverwriteFile(h.ptr, value)
	case "alpm_option_remove_overwrite_file":
		if lib.AlpmOptionRemoveOverwriteFile == nil {
			return errors.New("function unavailable: alpm_option_remove_overwrite_file")
		}
		result = lib.AlpmOptionRemoveOverwriteFile(h.ptr, value)
	case "alpm_option_set_sandboxuser":
		if lib.AlpmOptionSetSandboxuser == nil {
			return errors.New("function unavailable: alpm_option_set_sandboxuser")
		}
		result = lib.AlpmOptionSetSandboxuser(h.ptr, value)
	case "alpm_option_add_architecture":
		if lib.AlpmOptionAddArchitecture == nil {
			return errors.New("function unavailable: alpm_option_add_architecture")
		}
		result = lib.AlpmOptionAddArchitecture(h.ptr, value)
	case "alpm_option_remove_architecture":
		if lib.AlpmOptionRemoveArchitecture == nil {
			return errors.New("function unavailable: alpm_option_remove_architecture")
		}
		result = lib.AlpmOptionRemoveArchitecture(h.ptr, value)
	default:
		return errors.New("unsupported function: " + funcName)
	}

	if result != 0 {
		return dyerrors.NewError(h.Errno(), "failed to set option")
	}
	return nil
}

func (h *handle) getOptionStr(funcName string) string {
	if h.ptr == 0 {
		return ""
	}
	switch funcName {
	case "alpm_option_get_logfile":
		if lib.AlpmOptionGetLogfile == nil {
			return ""
		}
		return lib.PtrToString(lib.AlpmOptionGetLogfile(h.ptr))
	case "alpm_option_get_gpgdir":
		if lib.AlpmOptionGetGPGDir == nil {
			return ""
		}
		return lib.PtrToString(lib.AlpmOptionGetGPGDir(h.ptr))
	case "alpm_option_get_dbext":
		if lib.AlpmOptionGetDBExt == nil {
			return ""
		}
		return lib.PtrToString(lib.AlpmOptionGetDBExt(h.ptr))
	case "alpm_option_get_sandboxuser":
		if lib.AlpmOptionGetSandboxuser == nil {
			return ""
		}
		return lib.PtrToString(lib.AlpmOptionGetSandboxuser(h.ptr))
	case "alpm_option_get_root":
		if lib.AlpmOptionGetRoot == nil {
			return ""
		}
		return lib.PtrToString(lib.AlpmOptionGetRoot(h.ptr))
	case "alpm_option_get_dbpath":
		if lib.AlpmOptionGetDbpath == nil {
			return ""
		}
		return lib.PtrToString(lib.AlpmOptionGetDbpath(h.ptr))
	default:
		return ""
	}
}

func (h *handle) setOptionInt(funcName string, value int) error {
	if h.ptr == 0 {
		return dyerrors.ErrHandleNull
	}

	valueInt32 := clampIntToInt32(value)

	var result int32
	switch funcName {
	case "alpm_option_set_usesyslog":
		if lib.AlpmOptionSetUseSyslog == nil {
			return errors.New("function unavailable: alpm_option_set_usesyslog")
		}
		result = lib.AlpmOptionSetUseSyslog(h.ptr, valueInt32)
	case "alpm_option_set_checkspace":
		if lib.AlpmOptionSetCheckspace == nil {
			return errors.New("function unavailable: alpm_option_set_checkspace")
		}
		result = lib.AlpmOptionSetCheckspace(h.ptr, valueInt32)
	case "alpm_option_set_default_siglevel":
		if lib.AlpmOptionSetDefaultSigLevel == nil {
			return errors.New("function unavailable: alpm_option_set_default_siglevel")
		}
		result = lib.AlpmOptionSetDefaultSigLevel(h.ptr, valueInt32)
	case "alpm_option_set_local_file_siglevel":
		if lib.AlpmOptionSetLocalFileSigLevel == nil {
			return errors.New("function unavailable: alpm_option_set_local_file_siglevel")
		}
		result = lib.AlpmOptionSetLocalFileSigLevel(h.ptr, valueInt32)
	case "alpm_option_set_remote_file_siglevel":
		if lib.AlpmOptionSetRemoteFileSigLevel == nil {
			return errors.New("function unavailable: alpm_option_set_remote_file_siglevel")
		}
		result = lib.AlpmOptionSetRemoteFileSigLevel(h.ptr, valueInt32)
	case "alpm_option_set_parallel_downloads":
		if lib.AlpmOptionSetParallelDownloads == nil {
			return errors.New("function unavailable: alpm_option_set_parallel_downloads")
		}
		result = lib.AlpmOptionSetParallelDownloads(h.ptr, valueInt32)
	case "alpm_option_set_disable_dl_timeout":
		if lib.AlpmOptionSetDisableDlTimeout == nil {
			return errors.New("function unavailable: alpm_option_set_disable_dl_timeout")
		}
		result = lib.AlpmOptionSetDisableDlTimeout(h.ptr, valueInt32)
	case "alpm_option_set_disable_sandbox":
		if lib.AlpmOptionSetDisableSandbox == nil {
			return errors.New("function unavailable: alpm_option_set_disable_sandbox")
		}
		result = lib.AlpmOptionSetDisableSandbox(h.ptr, valueInt32)
	default:
		return errors.New("unsupported function: " + funcName)
	}

	if result != 0 {
		return dyerrors.NewError(h.Errno(), "failed to set option")
	}
	return nil
}

func (h *handle) getOptionInt(funcName string) int32 {
	if h.ptr == 0 {
		return 0
	}
	switch funcName {
	case "alpm_option_get_usesyslog":
		if lib.AlpmOptionGetUseSyslog == nil {
			return 0
		}
		return lib.AlpmOptionGetUseSyslog(h.ptr)
	case "alpm_option_get_checkspace":
		if lib.AlpmOptionGetCheckspace == nil {
			return 0
		}
		return lib.AlpmOptionGetCheckspace(h.ptr)
	case "alpm_option_get_default_siglevel":
		if lib.AlpmOptionGetDefaultSigLevel == nil {
			return 0
		}
		return lib.AlpmOptionGetDefaultSigLevel(h.ptr)
	case "alpm_option_get_local_file_siglevel":
		if lib.AlpmOptionGetLocalFileSigLevel == nil {
			return 0
		}
		return lib.AlpmOptionGetLocalFileSigLevel(h.ptr)
	case "alpm_option_get_remote_file_siglevel":
		if lib.AlpmOptionGetRemoteFileSigLevel == nil {
			return 0
		}
		return lib.AlpmOptionGetRemoteFileSigLevel(h.ptr)
	case "alpm_option_get_parallel_downloads":
		if lib.AlpmOptionGetParallelDownloads == nil {
			return 0
		}
		return lib.AlpmOptionGetParallelDownloads(h.ptr)
	default:
		return 0
	}
}

func (h *handle) getOptionStrList(funcName string) ([]string, error) {
	if h.ptr == 0 {
		return nil, dyerrors.ErrHandleNull
	}

	var r1 uintptr
	switch funcName {
	case "alpm_option_get_cachedirs":
		if lib.AlpmOptionGetCachedirs == nil {
			return []string{}, nil
		}
		r1 = lib.AlpmOptionGetCachedirs(h.ptr)
	case "alpm_option_get_hookdirs":
		if lib.AlpmOptionGetHookdirs == nil {
			return []string{}, nil
		}
		r1 = lib.AlpmOptionGetHookdirs(h.ptr)
	case "alpm_option_get_noupgrades":
		if lib.AlpmOptionGetNoUpgrades == nil {
			return []string{}, nil
		}
		r1 = lib.AlpmOptionGetNoUpgrades(h.ptr)
	case "alpm_option_get_noextracts":
		if lib.AlpmOptionGetNoextracts == nil {
			return []string{}, nil
		}
		r1 = lib.AlpmOptionGetNoextracts(h.ptr)
	case "alpm_option_get_ignorepkgs":
		if lib.AlpmOptionGetIgnorepkgs == nil {
			return []string{}, nil
		}
		r1 = lib.AlpmOptionGetIgnorepkgs(h.ptr)
	case "alpm_option_get_ignoregroups":
		if lib.AlpmOptionGetIgnoregroups == nil {
			return []string{}, nil
		}
		r1 = lib.AlpmOptionGetIgnoregroups(h.ptr)
	case "alpm_option_get_overwrite_files":
		if lib.AlpmOptionGetOverwriteFiles == nil {
			return []string{}, nil
		}
		r1 = lib.AlpmOptionGetOverwriteFiles(h.ptr)
	case "alpm_option_get_architectures":
		if lib.AlpmOptionGetArchitectures == nil {
			return []string{}, nil
		}
		r1 = lib.AlpmOptionGetArchitectures(h.ptr)
	default:
		return nil, errors.New("unsupported function: " + funcName)
	}
	if r1 == 0 {
		return []string{}, nil
	}

	alpmList := alpmlist.NewList(r1)
	if alpmList == nil {
		return []string{}, nil
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

func (h *handle) setOptionStrList(funcName string, values []string) error {
	if h.ptr == 0 {
		return dyerrors.ErrHandleNull
	}
	switch funcName {
	case "alpm_option_set_cachedirs":
		if lib.AlpmOptionSetCachedirs == nil {
			return errors.New("function unavailable: alpm_option_set_cachedirs")
		}
	case "alpm_option_set_hookdirs":
		if lib.AlpmOptionSetHookdirs == nil {
			return errors.New("function unavailable: alpm_option_set_hookdirs")
		}
	case "alpm_option_set_noupgrades":
		if lib.AlpmOptionSetNoUpgrades == nil {
			return errors.New("function unavailable: alpm_option_set_noupgrades")
		}
	case "alpm_option_set_noextracts":
		if lib.AlpmOptionSetNoextracts == nil {
			return errors.New("function unavailable: alpm_option_set_noextracts")
		}
	case "alpm_option_set_ignorepkgs":
		if lib.AlpmOptionSetIgnorepkgs == nil {
			return errors.New("function unavailable: alpm_option_set_ignorepkgs")
		}
	case "alpm_option_set_ignoregroups":
		if lib.AlpmOptionSetIgnoregroups == nil {
			return errors.New("function unavailable: alpm_option_set_ignoregroups")
		}
	case "alpm_option_set_overwrite_files":
		if lib.AlpmOptionSetOverwriteFiles == nil {
			return errors.New("function unavailable: alpm_option_set_overwrite_files")
		}
	case "alpm_option_set_architectures":
		if lib.AlpmOptionSetArchitectures == nil {
			return errors.New("function unavailable: alpm_option_set_architectures")
		}
	default:
		return errors.New("unsupported function: " + funcName)
	}

	var alpmList *alpmlist.List
	var cStrings [][]byte

	for _, v := range values {
		cS := lib.CString(v)
		cStrings = append(cStrings, cS)
		alpmList = alpmlist.Add(alpmList, uintptr(unsafe.Pointer(&cS[0])))
	}
	if alpmList != nil {
		defer alpmList.Free()
	}

	var listPtr uintptr
	if alpmList != nil {
		listPtr = alpmList.Ptr()
	}

	var r1 int32
	switch funcName {
	case "alpm_option_set_cachedirs":
		r1 = lib.AlpmOptionSetCachedirs(h.ptr, listPtr)
	case "alpm_option_set_hookdirs":
		r1 = lib.AlpmOptionSetHookdirs(h.ptr, listPtr)
	case "alpm_option_set_noupgrades":
		r1 = lib.AlpmOptionSetNoUpgrades(h.ptr, listPtr)
	case "alpm_option_set_noextracts":
		r1 = lib.AlpmOptionSetNoextracts(h.ptr, listPtr)
	case "alpm_option_set_ignorepkgs":
		r1 = lib.AlpmOptionSetIgnorepkgs(h.ptr, listPtr)
	case "alpm_option_set_ignoregroups":
		r1 = lib.AlpmOptionSetIgnoregroups(h.ptr, listPtr)
	case "alpm_option_set_overwrite_files":
		r1 = lib.AlpmOptionSetOverwriteFiles(h.ptr, listPtr)
	case "alpm_option_set_architectures":
		r1 = lib.AlpmOptionSetArchitectures(h.ptr, listPtr)
	}

	runtime.KeepAlive(cStrings)
	runtime.KeepAlive(alpmList)

	if r1 != 0 {
		return dyerrors.NewError(h.Errno(), "failed to set option list")
	}
	return nil
}

func (h *handle) setOptionDepList(funcName string, deps []Dependency) error {
	if h.ptr == 0 {
		return dyerrors.ErrHandleNull
	}
	switch funcName {
	case "alpm_option_set_assumeinstalled":
		if lib.AlpmOptionSetNoassumeInstalled == nil {
			return errors.New("function unavailable: alpm_option_set_assumeinstalled")
		}
	default:
		return errors.New("unsupported function: " + funcName)
	}

	var alpmList *alpmlist.List
	for _, d := range deps {
		depImpl, ok := d.(*dependency)
		if ok {
			alpmList = alpmlist.Add(alpmList, depImpl.ptr)
		}
	}
	if alpmList != nil {
		defer alpmList.Free()
	}

	var listPtr uintptr
	if alpmList != nil {
		listPtr = alpmList.Ptr()
	}

	r1 := lib.AlpmOptionSetNoassumeInstalled(h.ptr, listPtr)
	runtime.KeepAlive(alpmList)

	if r1 != 0 {
		return dyerrors.NewError(h.Errno(), "failed to set option dependency list")
	}
	return nil
}
