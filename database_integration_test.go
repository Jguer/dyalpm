//go:build integration

package dyalpm

import "testing"

func TestDatabase_LocalDBAndSearchQueries(t *testing.T) {
	h := mustInitializeTestHandle(t)
	localDB := mustLocalDB(t, h)

	pkg := mustInstalledPkg(t, h, "glibc", "pacman", "bash")

	pkgByName := localDB.Pkg(pkg.Name())
	if pkgByName == nil {
		t.Fatalf("expected local package %q to be found by name", pkg.Name())
	}

	cache := localDB.PkgCache().Collect()
	if len(cache) == 0 {
		t.Fatalf("expected local package cache to have entries")
	}

	var cacheMatch bool
	for _, cached := range cache {
		if cached != nil && cached.Name() == pkg.Name() {
			cacheMatch = true
			break
		}
	}
	if !cacheMatch {
		t.Fatalf("expected local package cache to contain %q", pkg.Name())
	}

	searchResults := localDB.Search([]string{pkg.Name()}).Collect()
	var searchMatch bool
	for _, found := range searchResults {
		if found != nil && found.Name() == pkg.Name() {
			searchMatch = true
			break
		}
	}
	if !searchMatch {
		t.Fatalf("expected search results to contain %q", pkg.Name())
	}
}

func TestDatabase_SyncDBManagement(t *testing.T) {
	h := mustInitializeTestHandle(t)

	syncDB, err := h.RegisterSyncDB("core", 0)
	if err != nil {
		t.Fatalf("failed to register core sync DB: %v", err)
	}

	dbImpl, ok := syncDB.(*database)
	if !ok {
		t.Fatalf("sync DB is not the expected internal type")
	}

	if err := syncDB.SetUsage(int(UsageAll)); err != nil {
		t.Fatalf("failed to set sync DB usage: %v", err)
	}
	usage, err := dbImpl.GetUsage()
	if err != nil {
		t.Fatalf("failed to read sync DB usage: %v", err)
	}
	if usage != int(UsageAll) {
		t.Fatalf("expected sync DB usage %d, got %d", UsageAll, usage)
	}

	testServers := []string{"file:///tmp/dyalpm-test-server"}
	if err := syncDB.SetServers(testServers); err != nil {
		t.Fatalf("failed to set sync DB servers: %v", err)
	}

	servers := dbImpl.GetServers()
	if !containsStringInSlice(servers, testServers[0]) {
		t.Fatalf("expected sync DB servers to include %q", testServers[0])
	}

	if err := h.UnregisterAllSyncDBs(); err != nil {
		t.Fatalf("failed to unregister all sync DBs: %v", err)
	}

	if _, err := h.SyncDBByName("core"); err == nil {
		t.Fatalf("expected core sync DB to be unregistered")
	}
}

func TestDatabase_DirectAccessors(t *testing.T) {
	h := mustInitializeTestHandle(t)
	localDB := mustLocalDB(t, h)
	localDBImpl, ok := localDB.(*database)
	if !ok {
		t.Fatalf("local DB is not internal *database type")
	}

	pkgs, err := localDBImpl.GetPkgCache()
	if err != nil {
		t.Fatalf("GetPkgCache failed: %v", err)
	}
	if len(pkgs) == 0 {
		t.Skip("local DB has no packages")
	}

	pkg := mustInstalledPkg(t, h)
	if pkg == nil {
		pkg = pkgs[0]
	}
	if pkg == nil || pkg.Name() == "" {
		t.Fatalf("invalid package returned from GetPkgCache")
	}

	if localDBImpl.GetHandle() != h {
		t.Fatalf("GetHandle did not return original handle")
	}
	if got := localDBImpl.GetSigLevel(); got == 0 {
		t.Logf("GetSigLevel returned default value 0")
	}
	if valid := localDBImpl.IsValid(); !valid {
		t.Logf("local database IsValid returned false")
	}

	groups, err := localDBImpl.GetGroupCache()
	if err != nil {
		t.Fatalf("GetGroupCache failed: %v", err)
	}
	if len(groups) == 0 {
		t.Logf("GetGroupCache returned no entries")
	}

	pkgImpl, ok := pkg.(*package_)
	if !ok {
		t.Fatalf("package from GetPkgCache is not internal *package_ type")
	}

	groupPackages := pkgImpl
	if len(groupPackages.Groups()) == 0 {
		for _, candidate := range pkgs[1:] {
			candidatePkg, ok := candidate.(*package_)
			if !ok {
				continue
			}
			if len(candidatePkg.Groups()) > 0 {
				groupPackages = candidatePkg
				break
			}
		}
	}

	groupNames := groupPackages.Groups()
	if len(groupNames) == 0 {
		t.Logf("no installed package currently reports group metadata; skipping group content assertions")
		return
	}

	for _, g := range groupNames {
		group, err := localDBImpl.GetGroup(g)
		if err != nil {
			t.Fatalf("GetGroup(%q) failed: %v", g, err)
		}
		if group == nil || group.GetName() != g {
			t.Fatalf("expected group name %q, got %q", g, group.GetName())
		}
		pkgsInGroup, err := group.GetPackages()
		if err != nil {
			t.Fatalf("GetPackages for group %q failed: %v", g, err)
		}
		if len(pkgsInGroup) == 0 {
			t.Logf("group %q has no packages", g)
			continue
		}
		if pkgsInGroup[0] == nil || pkgsInGroup[0].Name() == "" {
			t.Fatalf("group %q returned invalid package entry", g)
		}
		break
	}

	pkgByName, err := localDBImpl.GetPkg(pkg.Name())
	if err != nil || pkgByName == nil {
		t.Fatalf("GetPkg(%q) failed: %v", pkg.Name(), err)
	}
}

func TestDatabase_ServerAndUsageMutators(t *testing.T) {
	h := mustInitializeTestHandle(t)
	syncDB, err := h.RegisterSyncDB("core", 0)
	if err != nil {
		t.Fatalf("failed to register core sync DB: %v", err)
	}
	dbImpl, ok := syncDB.(*database)
	if !ok {
		t.Fatalf("sync DB is not internal database type")
	}

	servers := []string{"file:///tmp/dyalpm-integration-cache"}
	if err := syncDB.SetServers(servers); err != nil {
		t.Fatalf("failed to set servers: %v", err)
	}
	if got := dbImpl.GetServers(); len(got) == 0 {
		t.Fatalf("GetServers returned no servers")
	} else if !containsStringInSlice(got, servers[0]) {
		t.Fatalf("SetServers value %q not found in %v", servers[0], got)
	}

	if err := dbImpl.AddServer("file:///tmp/dyalpm-extra-cache"); err != nil {
		t.Fatalf("failed to add server: %v", err)
	}
	afterAdd := dbImpl.GetServers()
	if !containsStringInSlice(afterAdd, "file:///tmp/dyalpm-extra-cache") {
		t.Fatalf("expected added server in list, got %v", afterAdd)
	}
	if err := dbImpl.RemoveServer("file:///tmp/dyalpm-extra-cache"); err != nil {
		t.Fatalf("failed to remove server: %v", err)
	}
	afterRemove := dbImpl.GetServers()
	if containsStringInSlice(afterRemove, "file:///tmp/dyalpm-extra-cache") {
		t.Fatalf("expected added server to be removed, still present in %v", afterRemove)
	}

	cacheServers := []string{"file:///tmp/dyalpm-cache-a", "file:///tmp/dyalpm-cache-b"}
	if err := dbImpl.SetCacheServers(cacheServers); err != nil {
		t.Fatalf("failed to set cache servers: %v", err)
	}
	if got := dbImpl.GetCacheServers(); len(got) == 0 {
		t.Fatalf("GetCacheServers returned no values")
	}
	if err := dbImpl.AddCacheServer("file:///tmp/dyalpm-cache-extra"); err != nil {
		t.Fatalf("failed to add cache server: %v", err)
	}
	cacheAfterAdd := dbImpl.GetCacheServers()
	if !containsStringInSlice(cacheAfterAdd, "file:///tmp/dyalpm-cache-extra") {
		t.Fatalf("expected cache server to be added, got %v", cacheAfterAdd)
	}
	if err := dbImpl.RemoveCacheServer("file:///tmp/dyalpm-cache-extra"); err != nil {
		t.Fatalf("failed to remove cache server: %v", err)
	}
	cacheAfterRemove := dbImpl.GetCacheServers()
	if containsStringInSlice(cacheAfterRemove, "file:///tmp/dyalpm-cache-extra") {
		t.Fatalf("expected cache server to be removed, still present in %v", cacheAfterRemove)
	}

	if err := dbImpl.SetUsage(int(UsageSync | UsageSearch | UsageInstall)); err != nil {
		t.Fatalf("SetUsage failed: %v", err)
	}
	gotUsage, err := dbImpl.GetUsage()
	if err != nil {
		t.Fatalf("GetUsage failed: %v", err)
	}
	if gotUsage == 0 {
		t.Fatalf("GetUsage returned zero usage")
	}

	if err := syncDB.SetUsage(0); err != nil {
		t.Fatalf("SetUsage(0) failed: %v", err)
	}
	usageAfter, err := dbImpl.GetUsage()
	if err != nil {
		t.Fatalf("GetUsage after reset failed: %v", err)
	}
	if usageAfter != 0 {
		t.Fatalf("expected usage to be reset to 0, got %d", usageAfter)
	}

	if err := dbImpl.Update(false); err == nil {
		t.Logf("sync DB update succeeded unexpectedly in test environment")
	}

	if err := dbImpl.Unregister(); err != nil {
		t.Fatalf("Unregister failed: %v", err)
	}
	if _, err := h.SyncDBByName("core"); err == nil {
		t.Fatalf("expected core sync DB to be removed after Unregister")
	}
}
