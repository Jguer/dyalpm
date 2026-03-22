package lib

var (
	AlpmVersion      func() uintptr
	AlpmCapabilities func() uint64
	AlpmErrno        func(handle uintptr) int32
	AlpmStrerror     func(errno int32) uintptr
	AlpmRelease      func(handle uintptr) int32
	AlpmInitialize   func(root string, dbpath string, errno *int32) uintptr
	AlpmGetLocaldb   func(handle uintptr) uintptr
	AlpmGetSyncdbs   func(handle uintptr) uintptr
	AlpmGetLocalDB   func(handle uintptr) uintptr
	AlpmGetSyncDBS   func(handle uintptr) uintptr

	AlpmOptionGetRoot               func(handle uintptr) uintptr
	AlpmOptionGetDbpath             func(handle uintptr) uintptr
	AlpmOptionSetLogfile            func(handle uintptr, path string) int32
	AlpmOptionGetLogfile            func(handle uintptr) uintptr
	AlpmOptionSetGPGDir             func(handle uintptr, path string) int32
	AlpmOptionGetGPGDir             func(handle uintptr) uintptr
	AlpmOptionSetUseSyslog          func(handle uintptr, value int32) int32
	AlpmOptionGetUseSyslog          func(handle uintptr) int32
	AlpmOptionSetCheckspace         func(handle uintptr, value int32) int32
	AlpmOptionGetCheckspace         func(handle uintptr) int32
	AlpmOptionSetDBExt              func(handle uintptr, dbext string) int32
	AlpmOptionGetDBExt              func(handle uintptr) uintptr
	AlpmOptionSetDefaultSigLevel    func(handle uintptr, value int32) int32
	AlpmOptionGetDefaultSigLevel    func(handle uintptr) int32
	AlpmOptionSetLocalFileSigLevel  func(handle uintptr, value int32) int32
	AlpmOptionGetLocalFileSigLevel  func(handle uintptr) int32
	AlpmOptionSetRemoteFileSigLevel func(handle uintptr, value int32) int32
	AlpmOptionGetRemoteFileSigLevel func(handle uintptr) int32
	AlpmOptionSetParallelDownloads  func(handle uintptr, value int32) int32
	AlpmOptionGetParallelDownloads  func(handle uintptr) int32
	AlpmOptionSetCachedirs          func(handle uintptr, list uintptr) int32
	AlpmOptionGetCachedirs          func(handle uintptr) uintptr
	AlpmOptionSetHookdirs           func(handle uintptr, list uintptr) int32
	AlpmOptionGetHookdirs           func(handle uintptr) uintptr
	AlpmOptionSetNoUpgrades         func(handle uintptr, list uintptr) int32
	AlpmOptionGetNoUpgrades         func(handle uintptr) uintptr
	AlpmOptionSetNoextracts         func(handle uintptr, list uintptr) int32
	AlpmOptionGetNoextracts         func(handle uintptr) uintptr
	AlpmOptionSetIgnorepkgs         func(handle uintptr, list uintptr) int32
	AlpmOptionGetIgnorepkgs         func(handle uintptr) uintptr
	AlpmOptionSetIgnoregroups       func(handle uintptr, list uintptr) int32
	AlpmOptionGetIgnoregroups       func(handle uintptr) uintptr
	AlpmOptionSetOverwriteFiles     func(handle uintptr, list uintptr) int32
	AlpmOptionGetOverwriteFiles     func(handle uintptr) uintptr
	AlpmOptionMatchNoUpgrade        func(handle uintptr, path string) int32
	AlpmOptionMatchNoextract        func(handle uintptr, path string) int32
	AlpmOptionSetSandboxuser        func(handle uintptr, user string) int32
	AlpmOptionGetSandboxuser        func(handle uintptr) uintptr
	AlpmOptionSetDisableDlTimeout   func(handle uintptr, value int32) int32
	AlpmOptionSetDisableSandbox     func(handle uintptr, value int32) int32
	AlpmOptionSetArchitectures      func(handle uintptr, list uintptr) int32
	AlpmOptionGetArchitectures      func(handle uintptr) uintptr
	AlpmOptionAddCachedir           func(handle uintptr, path string) int32
	AlpmOptionRemoveCachedir        func(handle uintptr, path string) int32
	AlpmOptionAddHookdir            func(handle uintptr, path string) int32
	AlpmOptionRemoveHookdir         func(handle uintptr, path string) int32
	AlpmOptionAddNoupgrade          func(handle uintptr, path string) int32
	AlpmOptionRemoveNoupgrade       func(handle uintptr, path string) int32
	AlpmOptionAddNoextract          func(handle uintptr, path string) int32
	AlpmOptionRemoveNoextract       func(handle uintptr, path string) int32
	AlpmOptionAddIgnorepkg          func(handle uintptr, pkg string) int32
	AlpmOptionRemoveIgnorepkg       func(handle uintptr, pkg string) int32
	AlpmOptionAddIgnoregroup        func(handle uintptr, group string) int32
	AlpmOptionRemoveIgnoregroup     func(handle uintptr, group string) int32
	AlpmOptionAddOverwriteFile      func(handle uintptr, glob string) int32
	AlpmOptionRemoveOverwriteFile   func(handle uintptr, glob string) int32
	AlpmOptionSetLogcb              func(handle uintptr, cb uintptr, ctx uintptr) int32
	AlpmOptionGetLogcb              func(handle uintptr) uintptr
	AlpmOptionGetLogcbCtx           func(handle uintptr) uintptr
	AlpmOptionSetDlopencb           func(handle uintptr, cb uintptr, ctx uintptr) int32
	AlpmOptionGetDlopencb           func(handle uintptr) uintptr
	AlpmOptionGetDlopencbCtx        func(handle uintptr) uintptr
	AlpmOptionSetDlcb               func(handle uintptr, cb uintptr, ctx uintptr) int32
	AlpmOptionGetDlcb               func(handle uintptr) uintptr
	AlpmOptionGetDlcbCtx            func(handle uintptr) uintptr
	AlpmOptionSetFetchcb            func(handle uintptr, cb uintptr, ctx uintptr) int32
	AlpmOptionGetFetchcb            func(handle uintptr) uintptr
	AlpmOptionGetFetchcbCtx         func(handle uintptr) uintptr
	AlpmOptionSetEventcb            func(handle uintptr, cb uintptr, ctx uintptr) int32
	AlpmOptionGetEventcb            func(handle uintptr) uintptr
	AlpmOptionGetEventcbCtx         func(handle uintptr) uintptr
	AlpmOptionSetQuestioncb         func(handle uintptr, cb uintptr, ctx uintptr) int32
	AlpmOptionGetQuestioncb         func(handle uintptr) uintptr
	AlpmOptionGetQuestioncbCtx      func(handle uintptr) uintptr
	AlpmOptionSetProgresscb         func(handle uintptr, cb uintptr, ctx uintptr) int32
	AlpmOptionGetProgresscb         func(handle uintptr) uintptr
	AlpmOptionGetProgresscbCtx      func(handle uintptr) uintptr
	AlpmOptionSetNoassumeInstalled  func(handle uintptr, list uintptr) int32
	AlpmOptionGetAssumeInstalled    func(handle uintptr) uintptr
	AlpmOptionAddAssumeInstalled    func(handle uintptr, dep uintptr) int32
	AlpmOptionRemoveAssumeInstalled func(handle uintptr, dep uintptr) int32
	AlpmOptionAddArchitecture       func(handle uintptr, arch string) int32
	AlpmOptionRemoveArchitecture    func(handle uintptr, arch string) int32

	AlpmTransInit            func(handle uintptr, flags int32) int32
	AlpmTransPrepare         func(handle uintptr, list *uintptr) int32
	AlpmTransCommit          func(handle uintptr, list *uintptr) int32
	AlpmTransRelease         func(handle uintptr) int32
	AlpmAddPkg               func(handle uintptr, pkg uintptr) int32
	AlpmRemovePkg            func(handle uintptr, pkg uintptr) int32
	AlpmSyncSysupgrade       func(handle uintptr, enable int32) int32
	AlpmTransGetAdd          func(handle uintptr) uintptr
	AlpmTransGetRemove       func(handle uintptr) uintptr
	AlpmTransGetFlags        func(handle uintptr) int32
	AlpmRegisterSyncDB       func(handle uintptr, name string, siglevel int32) uintptr
	AlpmUnregisterAllSyncDBs func(handle uintptr) int32
	AlpmFetchPkgurl          func(handle uintptr, urlsList uintptr, fetchedList *uintptr) int32
	AlpmTransInterrupt       func(handle uintptr) int32
	AlpmUnlock               func(handle uintptr) int32
	AlpmFindGroupPkgs        func(list uintptr, name string) uintptr
	AlpmSandboxSetupChild    func(handle uintptr, user string, dir string) int32

	AlpmLogactionSym uintptr

	AlpmDBGetName            func(db uintptr) uintptr
	AlpmDBGetPkg             func(db uintptr, name string) uintptr
	AlpmDBGetPkgcache        func(db uintptr) uintptr
	AlpmDBSearch             func(db uintptr, needles uintptr, result *uintptr) int32
	AlpmDBGetGroup           func(db uintptr, name string) uintptr
	AlpmDBGetGroupcache      func(db uintptr) uintptr
	AlpmDBUpdate             func(handle uintptr, list uintptr, force int32) int32
	AlpmDBUnregister         func(db uintptr) int32
	AlpmDBGetServers         func(db uintptr) uintptr
	AlpmDBSetServers         func(db uintptr, servers uintptr) int32
	AlpmDBAddServer          func(db uintptr, url string) int32
	AlpmDBRemoveServer       func(db uintptr, url string) int32
	AlpmDBGetCacheServers    func(db uintptr) uintptr
	AlpmDBSetCacheServers    func(db uintptr, servers uintptr) int32
	AlpmDBAddCacheServer     func(db uintptr, url string) int32
	AlpmDBRemoveCacheServer  func(db uintptr, url string) int32
	AlpmDBSetUsage           func(db uintptr, usage int32) int32
	AlpmDBGetUsage           func(db uintptr, usage *int32) int32
	AlpmDBGetValid           func(db uintptr) int32
	AlpmDBGetSiglevel        func(db uintptr) int32
	AlpmDBGetHandle          func(db uintptr) uintptr
	AlpmDBCheckPGPSignature  func(db uintptr, siglist uintptr) int32
	AlpmPkgCheckPGPSignature func(pkg uintptr, siglist uintptr) int32

	AlpmDepComputeString func(dep uintptr) uintptr
	AlpmDepFromString    func(text string) uintptr
	AlpmDepFree          func(dep uintptr) int32
	AlpmDepmissingFree   func(ptr uintptr) int32
	AlpmConflictFree     func(ptr uintptr) int32
	AlpmFileConflictFree func(ptr uintptr) int32
	AlpmCheckDeps        func(handle uintptr, pkgList uintptr, remPkgList uintptr, upgradePkgList uintptr, reverse int32) uintptr
	AlpmCheckConflicts   func(handle uintptr, pkgList uintptr) uintptr
	AlpmFindSatisfier    func(pkgs uintptr, dep string) uintptr
	AlpmFindDBSatisfier  func(handle uintptr, dbList uintptr, dep string) uintptr
	AlpmComputeMd5sum    func(filename string) uintptr
	AlpmComputeSha256sum func(filename string) uintptr
	AlpmExtractKeyID     func(handle uintptr, identifier string, sig uintptr, sigLen int32, keys *uintptr) int32

	AlpmPkgFind               func(list uintptr, name string) uintptr
	AlpmPkgGetName            func(pkg uintptr) uintptr
	AlpmPkgGetVersion         func(pkg uintptr) uintptr
	AlpmPkgGetDesc            func(pkg uintptr) uintptr
	AlpmPkgGetArch            func(pkg uintptr) uintptr
	AlpmPkgGetSize            func(pkg uintptr) int64
	AlpmPkgGetISize           func(pkg uintptr) int64
	AlpmPkgGetDB              func(pkg uintptr) uintptr
	AlpmPkgGetDepends         func(pkg uintptr) uintptr
	AlpmPkgGetConflicts       func(pkg uintptr) uintptr
	AlpmPkgGetProvides        func(pkg uintptr) uintptr
	AlpmPkgGetOptdepends      func(pkg uintptr) uintptr
	AlpmPkgGetReplaces        func(pkg uintptr) uintptr
	AlpmPkgGetGroups          func(pkg uintptr) uintptr
	AlpmPkgGetLicenses        func(pkg uintptr) uintptr
	AlpmPkgGetFilename        func(pkg uintptr) uintptr
	AlpmPkgGetReason          func(pkg uintptr) int32
	AlpmPkgGetOrigin          func(pkg uintptr) int32
	AlpmPkgGetBase            func(pkg uintptr) uintptr
	AlpmPkgGetPackager        func(pkg uintptr) uintptr
	AlpmPkgGetSha256sum       func(pkg uintptr) uintptr
	AlpmPkgGetValidation      func(pkg uintptr) int32
	AlpmPkgGetURL             func(pkg uintptr) uintptr
	AlpmPkgHasScriptlet       func(pkg uintptr) int32
	AlpmPkgDownloadSize       func(pkg uintptr) int64
	AlpmPkgGetBackup          func(pkg uintptr) uintptr
	AlpmPkgGetFiles           func(pkg uintptr) uintptr
	AlpmPkgGetInstallDate     func(pkg uintptr) int32
	AlpmPkgGetBuildDate       func(pkg uintptr) int32
	AlpmPkgGetHandle          func(pkg uintptr) uintptr
	AlpmPkgShouldIgnore       func(handle uintptr, pkg uintptr) int32
	AlpmPkgCheckmd5sum        func(pkg uintptr) int32
	AlpmPkgGetSig             func(pkg uintptr, sig *uintptr, sigLen *uintptr) int32
	AlpmPkgGetBase64Sig       func(pkg uintptr) uintptr
	AlpmPkgChangelogOpen      func(pkg uintptr) uintptr
	AlpmPkgChangelogRead      func(buf uintptr, size uintptr, pkg uintptr, fp uintptr) int
	AlpmPkgChangelogClose     func(pkg uintptr, fp uintptr)
	AlpmPkgSyncGetNewVersion  func(pkg uintptr, dbList uintptr) uintptr
	AlpmPkgGetFilesContains   func(fileList uintptr, path string) int32
	AlpmPkgFree               func(pkg uintptr) int32
	AlpmPkgGetXdata           func(pkg uintptr) uintptr
	AlpmPkgComputeRequiredBy  func(pkg uintptr) uintptr
	AlpmPkgComputeOptionalFor func(pkg uintptr) uintptr
	AlpmPkgLoad               func(handle uintptr, filename string, full int32, siglevel int32, pkg *uintptr) int32

	AlpmSiglistCleanup func(listPtr uintptr, handle uintptr) int32

	AlpmListAdd   func(list uintptr, data uintptr) uintptr
	AlpmListCount func(list uintptr, result *uintptr) int32
	AlpmListFree  func(list uintptr)

	AlpmSetLockFile func(handle uintptr, path string) int32
)

func registerAlpmFuncs(library uintptr) {
	tryRegister(&AlpmVersion, library, "alpm_version")
	tryRegister(&AlpmCapabilities, library, "alpm_capabilities")
	tryRegister(&AlpmErrno, library, "alpm_errno")
	tryRegister(&AlpmStrerror, library, "alpm_strerror")
	tryRegister(&AlpmRelease, library, "alpm_release")
	tryRegister(&AlpmGetLocaldb, library, "alpm_get_localdb")
	tryRegister(&AlpmGetSyncdbs, library, "alpm_get_syncdbs")
	tryRegister(&AlpmGetLocalDB, library, "alpm_get_localdb")
	tryRegister(&AlpmGetSyncDBS, library, "alpm_get_syncdbs")
	tryRegister(&AlpmInitialize, library, "alpm_initialize")

	tryRegister(&AlpmOptionGetRoot, library, "alpm_option_get_root")
	tryRegister(&AlpmOptionGetDbpath, library, "alpm_option_get_dbpath")
	tryRegister(&AlpmOptionSetLogfile, library, "alpm_option_set_logfile")
	tryRegister(&AlpmOptionGetLogfile, library, "alpm_option_get_logfile")
	tryRegister(&AlpmOptionSetGPGDir, library, "alpm_option_set_gpgdir")
	tryRegister(&AlpmOptionGetGPGDir, library, "alpm_option_get_gpgdir")
	tryRegister(&AlpmOptionSetUseSyslog, library, "alpm_option_set_usesyslog")
	tryRegister(&AlpmOptionGetUseSyslog, library, "alpm_option_get_usesyslog")
	tryRegister(&AlpmOptionSetCheckspace, library, "alpm_option_set_checkspace")
	tryRegister(&AlpmOptionGetCheckspace, library, "alpm_option_get_checkspace")
	tryRegister(&AlpmOptionSetDBExt, library, "alpm_option_set_dbext")
	tryRegister(&AlpmOptionGetDBExt, library, "alpm_option_get_dbext")
	tryRegister(&AlpmOptionSetDefaultSigLevel, library, "alpm_option_set_default_siglevel")
	tryRegister(&AlpmOptionGetDefaultSigLevel, library, "alpm_option_get_default_siglevel")
	tryRegister(&AlpmOptionSetLocalFileSigLevel, library, "alpm_option_set_local_file_siglevel")
	tryRegister(&AlpmOptionGetLocalFileSigLevel, library, "alpm_option_get_local_file_siglevel")
	tryRegister(&AlpmOptionSetRemoteFileSigLevel, library, "alpm_option_set_remote_file_siglevel")
	tryRegister(&AlpmOptionGetRemoteFileSigLevel, library, "alpm_option_get_remote_file_siglevel")
	tryRegister(&AlpmOptionSetParallelDownloads, library, "alpm_option_set_parallel_downloads")
	tryRegister(&AlpmOptionGetParallelDownloads, library, "alpm_option_get_parallel_downloads")
	tryRegister(&AlpmOptionSetCachedirs, library, "alpm_option_set_cachedirs")
	tryRegister(&AlpmOptionGetCachedirs, library, "alpm_option_get_cachedirs")
	tryRegister(&AlpmOptionSetHookdirs, library, "alpm_option_set_hookdirs")
	tryRegister(&AlpmOptionGetHookdirs, library, "alpm_option_get_hookdirs")
	tryRegister(&AlpmOptionSetNoUpgrades, library, "alpm_option_set_noupgrades")
	tryRegister(&AlpmOptionGetNoUpgrades, library, "alpm_option_get_noupgrades")
	tryRegister(&AlpmOptionSetNoextracts, library, "alpm_option_set_noextracts")
	tryRegister(&AlpmOptionGetNoextracts, library, "alpm_option_get_noextracts")
	tryRegister(&AlpmOptionSetIgnorepkgs, library, "alpm_option_set_ignorepkgs")
	tryRegister(&AlpmOptionGetIgnorepkgs, library, "alpm_option_get_ignorepkgs")
	tryRegister(&AlpmOptionSetIgnoregroups, library, "alpm_option_set_ignoregroups")
	tryRegister(&AlpmOptionGetIgnoregroups, library, "alpm_option_get_ignoregroups")
	tryRegister(&AlpmOptionSetOverwriteFiles, library, "alpm_option_set_overwrite_files")
	tryRegister(&AlpmOptionGetOverwriteFiles, library, "alpm_option_get_overwrite_files")
	tryRegister(&AlpmOptionMatchNoUpgrade, library, "alpm_option_match_noupgrade")
	tryRegister(&AlpmOptionMatchNoextract, library, "alpm_option_match_noextract")
	tryRegister(&AlpmOptionSetSandboxuser, library, "alpm_option_set_sandboxuser")
	tryRegister(&AlpmOptionGetSandboxuser, library, "alpm_option_get_sandboxuser")
	tryRegister(&AlpmOptionSetDisableDlTimeout, library, "alpm_option_set_disable_dl_timeout")
	tryRegister(&AlpmOptionSetDisableSandbox, library, "alpm_option_set_disable_sandbox")
	tryRegister(&AlpmOptionSetArchitectures, library, "alpm_option_set_architectures")
	tryRegister(&AlpmOptionGetArchitectures, library, "alpm_option_get_architectures")
	tryRegister(&AlpmOptionAddCachedir, library, "alpm_option_add_cachedir")
	tryRegister(&AlpmOptionRemoveCachedir, library, "alpm_option_remove_cachedir")
	tryRegister(&AlpmOptionAddHookdir, library, "alpm_option_add_hookdir")
	tryRegister(&AlpmOptionRemoveHookdir, library, "alpm_option_remove_hookdir")
	tryRegister(&AlpmOptionAddNoupgrade, library, "alpm_option_add_noupgrade")
	tryRegister(&AlpmOptionRemoveNoupgrade, library, "alpm_option_remove_noupgrade")
	tryRegister(&AlpmOptionAddNoextract, library, "alpm_option_add_noextract")
	tryRegister(&AlpmOptionRemoveNoextract, library, "alpm_option_remove_noextract")
	tryRegister(&AlpmOptionAddIgnorepkg, library, "alpm_option_add_ignorepkg")
	tryRegister(&AlpmOptionRemoveIgnorepkg, library, "alpm_option_remove_ignorepkg")
	tryRegister(&AlpmOptionAddIgnoregroup, library, "alpm_option_add_ignoregroup")
	tryRegister(&AlpmOptionRemoveIgnoregroup, library, "alpm_option_remove_ignoregroup")
	tryRegister(&AlpmOptionAddOverwriteFile, library, "alpm_option_add_overwrite_file")
	tryRegister(&AlpmOptionRemoveOverwriteFile, library, "alpm_option_remove_overwrite_file")
	tryRegister(&AlpmOptionSetLogcb, library, "alpm_option_set_logcb")
	tryRegister(&AlpmOptionGetLogcb, library, "alpm_option_get_logcb")
	tryRegister(&AlpmOptionGetLogcbCtx, library, "alpm_option_get_logcb_ctx")
	tryRegister(&AlpmOptionSetDlopencb, library, "alpm_option_set_dlopen")
	tryRegister(&AlpmOptionGetDlopencb, library, "alpm_option_get_dlopen")
	tryRegister(&AlpmOptionGetDlopencbCtx, library, "alpm_option_get_dlopen_ctx")
	tryRegister(&AlpmOptionSetDlcb, library, "alpm_option_set_dlcb")
	tryRegister(&AlpmOptionGetDlcb, library, "alpm_option_get_dlcb")
	tryRegister(&AlpmOptionGetDlcbCtx, library, "alpm_option_get_dlcb_ctx")
	tryRegister(&AlpmOptionSetFetchcb, library, "alpm_option_set_fetchcb")
	tryRegister(&AlpmOptionGetFetchcb, library, "alpm_option_get_fetchcb")
	tryRegister(&AlpmOptionGetFetchcbCtx, library, "alpm_option_get_fetchcb_ctx")
	tryRegister(&AlpmOptionSetEventcb, library, "alpm_option_set_eventcb")
	tryRegister(&AlpmOptionGetEventcb, library, "alpm_option_get_eventcb")
	tryRegister(&AlpmOptionGetEventcbCtx, library, "alpm_option_get_eventcb_ctx")
	tryRegister(&AlpmOptionSetQuestioncb, library, "alpm_option_set_questioncb")
	tryRegister(&AlpmOptionGetQuestioncb, library, "alpm_option_get_questioncb")
	tryRegister(&AlpmOptionGetQuestioncbCtx, library, "alpm_option_get_questioncb_ctx")
	tryRegister(&AlpmOptionSetProgresscb, library, "alpm_option_set_progresscb")
	tryRegister(&AlpmOptionGetProgresscb, library, "alpm_option_get_progresscb")
	tryRegister(&AlpmOptionGetProgresscbCtx, library, "alpm_option_get_progresscb_ctx")
	tryRegister(&AlpmOptionSetNoassumeInstalled, library, "alpm_option_set_assumeinstalled")
	tryRegister(&AlpmOptionGetAssumeInstalled, library, "alpm_option_get_assumeinstalled")
	tryRegister(&AlpmOptionAddAssumeInstalled, library, "alpm_option_add_assumeinstalled")
	tryRegister(&AlpmOptionRemoveAssumeInstalled, library, "alpm_option_remove_assumeinstalled")
	tryRegister(&AlpmOptionAddArchitecture, library, "alpm_option_add_architecture")
	tryRegister(&AlpmOptionRemoveArchitecture, library, "alpm_option_remove_architecture")

	tryRegister(&AlpmTransInit, library, "alpm_trans_init")
	tryRegister(&AlpmTransPrepare, library, "alpm_trans_prepare")
	tryRegister(&AlpmTransCommit, library, "alpm_trans_commit")
	tryRegister(&AlpmTransRelease, library, "alpm_trans_release")
	tryRegister(&AlpmAddPkg, library, "alpm_add_pkg")
	tryRegister(&AlpmRemovePkg, library, "alpm_remove_pkg")
	tryRegister(&AlpmPkgLoad, library, "alpm_pkg_load")
	tryRegister(&AlpmSyncSysupgrade, library, "alpm_sync_sysupgrade")
	tryRegister(&AlpmTransGetAdd, library, "alpm_trans_get_add")
	tryRegister(&AlpmTransGetRemove, library, "alpm_trans_get_remove")
	tryRegister(&AlpmTransGetFlags, library, "alpm_trans_get_flags")
	tryRegister(&AlpmRegisterSyncDB, library, "alpm_register_syncdb")
	tryRegister(&AlpmUnregisterAllSyncDBs, library, "alpm_unregister_all_syncdbs")
	tryRegister(&AlpmFetchPkgurl, library, "alpm_fetch_pkgurl")
	tryRegister(&AlpmTransInterrupt, library, "alpm_trans_interrupt")
	tryRegister(&AlpmUnlock, library, "alpm_unlock")
	tryRegister(&AlpmFindGroupPkgs, library, "alpm_find_group_pkgs")
	tryRegister(&AlpmSandboxSetupChild, library, "alpm_sandbox_setup_child")
	tryDlsym(&AlpmLogactionSym, library, "alpm_logaction")

	tryRegister(&AlpmDBGetName, library, "alpm_db_get_name")
	tryRegister(&AlpmDBGetPkg, library, "alpm_db_get_pkg")
	tryRegister(&AlpmDBGetPkgcache, library, "alpm_db_get_pkgcache")
	tryRegister(&AlpmDBSearch, library, "alpm_db_search")
	tryRegister(&AlpmDBGetGroup, library, "alpm_db_get_group")
	tryRegister(&AlpmDBGetGroupcache, library, "alpm_db_get_groupcache")
	tryRegister(&AlpmDBUpdate, library, "alpm_db_update")
	tryRegister(&AlpmDBUnregister, library, "alpm_db_unregister")
	tryRegister(&AlpmDBGetServers, library, "alpm_db_get_servers")
	tryRegister(&AlpmDBSetServers, library, "alpm_db_set_servers")
	tryRegister(&AlpmDBAddServer, library, "alpm_db_add_server")
	tryRegister(&AlpmDBRemoveServer, library, "alpm_db_remove_server")
	tryRegister(&AlpmDBGetCacheServers, library, "alpm_db_get_cache_servers")
	tryRegister(&AlpmDBSetCacheServers, library, "alpm_db_set_cache_servers")
	tryRegister(&AlpmDBAddCacheServer, library, "alpm_db_add_cache_server")
	tryRegister(&AlpmDBRemoveCacheServer, library, "alpm_db_remove_cache_server")
	tryRegister(&AlpmDBSetUsage, library, "alpm_db_set_usage")
	tryRegister(&AlpmDBGetUsage, library, "alpm_db_get_usage")
	tryRegister(&AlpmDBGetValid, library, "alpm_db_get_valid")
	tryRegister(&AlpmDBGetSiglevel, library, "alpm_db_get_siglevel")
	tryRegister(&AlpmDBGetHandle, library, "alpm_db_get_handle")
	tryRegister(&AlpmDBCheckPGPSignature, library, "alpm_db_check_pgp_signature")
	tryRegister(&AlpmPkgCheckPGPSignature, library, "alpm_pkg_check_pgp_signature")

	tryRegister(&AlpmDepComputeString, library, "alpm_dep_compute_string")
	tryRegister(&AlpmDepFromString, library, "alpm_dep_from_string")
	tryRegister(&AlpmDepFree, library, "alpm_dep_free")
	tryRegister(&AlpmDepmissingFree, library, "alpm_depmissing_free")
	tryRegister(&AlpmConflictFree, library, "alpm_conflict_free")
	tryRegister(&AlpmFileConflictFree, library, "alpm_fileconflict_free")
	tryRegister(&AlpmCheckDeps, library, "alpm_checkdeps")
	tryRegister(&AlpmCheckConflicts, library, "alpm_checkconflicts")
	tryRegister(&AlpmFindSatisfier, library, "alpm_find_satisfier")
	tryRegister(&AlpmFindDBSatisfier, library, "alpm_find_dbs_satisfier")
	tryRegister(&AlpmComputeMd5sum, library, "alpm_compute_md5sum")
	tryRegister(&AlpmComputeSha256sum, library, "alpm_compute_sha256sum")
	tryRegister(&AlpmExtractKeyID, library, "alpm_extract_keyid")

	tryRegister(&AlpmPkgFind, library, "alpm_pkg_find")
	tryRegister(&AlpmPkgGetName, library, "alpm_pkg_get_name")
	tryRegister(&AlpmPkgGetVersion, library, "alpm_pkg_get_version")
	tryRegister(&AlpmPkgGetDesc, library, "alpm_pkg_get_desc")
	tryRegister(&AlpmPkgGetArch, library, "alpm_pkg_get_arch")
	tryRegister(&AlpmPkgGetSize, library, "alpm_pkg_get_size")
	tryRegister(&AlpmPkgGetISize, library, "alpm_pkg_get_isize")
	tryRegister(&AlpmPkgGetDB, library, "alpm_pkg_get_db")
	tryRegister(&AlpmPkgGetDepends, library, "alpm_pkg_get_depends")
	tryRegister(&AlpmPkgGetConflicts, library, "alpm_pkg_get_conflicts")
	tryRegister(&AlpmPkgGetProvides, library, "alpm_pkg_get_provides")
	tryRegister(&AlpmPkgGetOptdepends, library, "alpm_pkg_get_optdepends")
	tryRegister(&AlpmPkgGetReplaces, library, "alpm_pkg_get_replaces")
	tryRegister(&AlpmPkgGetGroups, library, "alpm_pkg_get_groups")
	tryRegister(&AlpmPkgGetLicenses, library, "alpm_pkg_get_licenses")
	tryRegister(&AlpmPkgGetFilename, library, "alpm_pkg_get_filename")
	tryRegister(&AlpmPkgGetReason, library, "alpm_pkg_get_reason")
	tryRegister(&AlpmPkgGetOrigin, library, "alpm_pkg_get_origin")
	tryRegister(&AlpmPkgGetBase, library, "alpm_pkg_get_base")
	tryRegister(&AlpmPkgGetPackager, library, "alpm_pkg_get_packager")
	tryRegister(&AlpmPkgGetSha256sum, library, "alpm_pkg_get_sha256sum")
	tryRegister(&AlpmPkgGetValidation, library, "alpm_pkg_get_validation")
	tryRegister(&AlpmPkgGetURL, library, "alpm_pkg_get_url")
	tryRegister(&AlpmPkgHasScriptlet, library, "alpm_pkg_has_scriptlet")
	tryRegister(&AlpmPkgDownloadSize, library, "alpm_pkg_download_size")
	tryRegister(&AlpmPkgGetBackup, library, "alpm_pkg_get_backup")
	tryRegister(&AlpmPkgGetFiles, library, "alpm_pkg_get_files")
	tryRegister(&AlpmPkgGetInstallDate, library, "alpm_pkg_get_installdate")
	tryRegister(&AlpmPkgGetBuildDate, library, "alpm_pkg_get_builddate")
	tryRegister(&AlpmPkgGetHandle, library, "alpm_pkg_get_handle")
	tryRegister(&AlpmPkgShouldIgnore, library, "alpm_pkg_should_ignore")
	tryRegister(&AlpmPkgCheckmd5sum, library, "alpm_pkg_checkmd5sum")
	tryRegister(&AlpmPkgGetSig, library, "alpm_pkg_get_sig")
	tryRegister(&AlpmPkgGetBase64Sig, library, "alpm_pkg_get_base64_sig")
	tryRegister(&AlpmPkgChangelogOpen, library, "alpm_pkg_changelog_open")
	tryRegister(&AlpmPkgChangelogRead, library, "alpm_pkg_changelog_read")
	tryRegister(&AlpmPkgChangelogClose, library, "alpm_pkg_changelog_close")
	tryRegister(&AlpmPkgSyncGetNewVersion, library, "alpm_sync_get_new_version")
	tryRegister(&AlpmPkgGetFilesContains, library, "alpm_filelist_contains")
	tryRegister(&AlpmPkgFree, library, "alpm_pkg_free")
	tryRegister(&AlpmPkgGetXdata, library, "alpm_pkg_get_xdata")
	tryRegister(&AlpmPkgComputeRequiredBy, library, "alpm_pkg_compute_requiredby")
	tryRegister(&AlpmPkgComputeOptionalFor, library, "alpm_pkg_compute_optionalfor")

	tryRegister(&AlpmSiglistCleanup, library, "alpm_siglist_cleanup")
	tryRegister(&AlpmListAdd, library, "alpm_list_add")
	tryRegister(&AlpmListCount, library, "alpm_list_count")
	tryRegister(&AlpmListFree, library, "alpm_list_free")
	tryRegister(&AlpmSetLockFile, library, "alpm_option_set_lockfile")
}
