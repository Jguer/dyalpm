# dyalpm

A Go wrapper for the Arch Linux Package Manager (ALPM) library using [purego](https://github.com/ebitengine/purego).

## Features

- **No cgo**: libalpm is loaded at runtime with purego
- **Type-Safe Interfaces**: Clean Go interfaces for all ALPM operations
- **Maintainable Structure**: Well-organized codebase with clear separation of concerns
- **Error Handling**: Comprehensive error types matching ALPM error codes

## Requirements

- Go 1.26 or later
- libalpm.so.16 (pacman 7; falls back to libalpm.so)
- Linux system with ALPM installed

## Installation

```bash
go get github.com/Jguer/dyalpm
```

## Usage

### Basic Example

```go
package main

import (
	"fmt"
	"log"

	alpm "github.com/Jguer/dyalpm"
)

func main() {
	// Initialize ALPM handle
	handle, err := alpm.Initialize("/", "/var/lib/pacman")
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Release()

	// Get local database
	localDB, err := handle.LocalDB()
	if err != nil {
		log.Fatal(err)
	}

	// Get a package
	pkg := localDB.Pkg("pacman")
	if pkg == nil {
		log.Fatal(err)
	}

	fmt.Printf("Package: %s\n", pkg.Name())
	fmt.Printf("Version: %s\n", pkg.Version())
	fmt.Printf("Description: %s\n", pkg.Description())
}
```

### Working with Sync Databases

```go
// List registered sync databases
syncDBs, err := handle.SyncDBs()
if err != nil {
	log.Fatal(err)
}

for _, db := range syncDBs {
	fmt.Printf("%s\n", db.Name())
}
```

### Transactions

```go
// Create a transaction
trans := alpm.NewTransaction(handle)

// Initialize transaction
err := trans.Init(0)
if err != nil {
	log.Fatal(err)
}
defer trans.Release()

// Add a package to install
pkg := syncDB.Pkg("vim")
err = trans.AddPkg(pkg)
if err != nil {
	log.Fatal(err)
}

// Prepare the transaction. On failure, missing holds the unsatisfied dependencies.
missing, err := trans.Prepare()
for _, dep := range missing {
	fmt.Printf("  %s requires %s\n", dep.GetTarget(), dep.GetDepend().GetName())
	dep.Free()
}
if err != nil {
	log.Fatal(err)
}

// Commit the transaction. On failure, conflicts holds the file conflicts.
conflicts, err := trans.Commit()
for _, c := range conflicts {
	fmt.Printf("  %s: %s\n", c.GetTarget(), c.GetFile())
	c.Free()
}
if err != nil {
	log.Fatal(err)
}
```

## Loading

libalpm is opened on the first `Initialize`, `Version` or `Capabilities` call and all symbols are resolved then. A symbol missing from an older libalpm is left unbound, and the methods that need it return a "missing function" error instead of crashing.

## Error Handling

`handle.Errno()` returns the current libalpm error code, and `handle.StrError` turns it into text:

```go
if errno := handle.Errno(); errno != 0 {
	fmt.Printf("Error: %s\n", handle.StrError(errno))
}
```

Transaction errors wrap a sentinel such as `alpm.ErrTransactionPrepareFailed` together with the libalpm error code, so `errors.Is(err, alpm.ErrTransactionPrepareFailed)` works.
