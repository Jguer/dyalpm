package main

import (
	"fmt"
	"log"
	"os"

	alpm "github.com/Jguer/dyalpm"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: example <command>")
		fmt.Println("Commands:")
		fmt.Println("  list-local          List all installed packages")
		fmt.Println("  info <package>      Show package information")
		fmt.Println("  list-sync           List all sync databases")
		os.Exit(1)
	}

	// Initialize ALPM handle
	// Use absolute paths
	root := "/"
	dbpath := "/var/lib/pacman"

	handle, err := alpm.Initialize(root, dbpath)
	if err != nil {
		log.Fatalf("Failed to initialize ALPM (root=%s, dbpath=%s): %v", root, dbpath, err)
	}
	defer func() { _ = handle.Release() }()

	command := os.Args[1]

	var cmdErr error
	switch command {
	case "list-local":
		cmdErr = listLocal(handle)
	case "info":
		if len(os.Args) < 3 {
			log.Fatal("Usage: example info <package>")
		}
		cmdErr = showInfo(handle, os.Args[2])
	case "list-sync":
		cmdErr = listSync(handle)
	default:
		cmdErr = fmt.Errorf("unknown command: %s", command)
	}

	if cmdErr != nil {
		log.Fatal(cmdErr)
	}
}

func listLocal(handle alpm.Handle) error {
	localDB, err := handle.LocalDB()
	if err != nil {
		return err
	}

	count := 0
	err = localDB.PkgCache().ForEach(func(pkg alpm.Package) error {
		fmt.Printf("  %s %s\n", pkg.Name(), pkg.Version())
		count++
		return nil
	})
	if err != nil {
		return err
	}

	fmt.Printf("Total installed packages: %d\n", count)
	return nil
}

func showInfo(handle alpm.Handle, pkgName string) error {
	localDB, err := handle.LocalDB()
	if err != nil {
		return err
	}

	pkg := localDB.Pkg(pkgName)
	if pkg == nil {
		return fmt.Errorf("package not found: %s", pkgName)
	}

	fmt.Printf("Name: %s\n", pkg.Name())
	fmt.Printf("Version: %s\n", pkg.Version())
	fmt.Printf("Description: %s\n", pkg.Description())
	fmt.Printf("Architecture: %s\n", pkg.Architecture())
	fmt.Printf("Installed Size: %d bytes\n", pkg.ISize())

	deps := pkg.Depends()
	if len(deps) > 0 {
		fmt.Println("Dependencies:")
		for _, dep := range deps {
			fmt.Printf("  %s", dep.Name)
			if dep.Version != "" {
				fmt.Printf(" %s", dep.Version)
			}
			fmt.Println()
		}
	}
	return nil
}

func listSync(handle alpm.Handle) error {
	syncDBs, err := handle.SyncDBs()
	if err != nil {
		return err
	}

	fmt.Printf("Sync databases (%d):\n", len(syncDBs))
	for _, db := range syncDBs {
		fmt.Printf("  %s\n", db.Name())
	}
	return nil
}
