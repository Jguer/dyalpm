package dyerrors

// Errno represents an ALPM error code
type Errno int

const (
	ErrOK Errno = iota
	ErrMemory
	ErrSystem
	ErrBadPerms
	ErrNotAFile
	ErrNotADir
	ErrWrongArgs
	ErrDiskSpace
	ErrHandleNull
	ErrHandleNotNull
	ErrHandleLock
	ErrDBOpen
	ErrDBCreate
	ErrDBNull
	ErrDBNotNull
	ErrDBNotFound
	ErrDBInvalid
	ErrDBInvalidSig
	ErrDBVersion
	ErrDBWrite
	ErrDBRemove
	ErrServerBadURL
	ErrServerNone
	ErrTransNotNull
	ErrTransNull
	ErrTransDupTarget
	ErrTransDupFilename
	ErrTransNotInitialized
	ErrTransNotPrepared
	ErrTransAbort
	ErrTransType
	ErrTransNotLocked
	ErrTransHookFailed
	ErrPkgNotFound
	ErrPkgIgnored
	ErrPkgInvalid
	ErrPkgInvalidChecksum
	ErrPkgInvalidSig
	ErrPkgMissingSig
	ErrPkgOpen
	ErrPkgCantRemove
	ErrPkgInvalidName
	ErrPkgInvalidArch
	ErrSigMissing
	ErrSigInvalid
	ErrUnsatisfiedDeps
	ErrConflictingDeps
	ErrFileConflicts
	ErrRetrieve
	ErrInvalidRegex
	ErrLibArchive
	ErrLibCurl
	ErrExternalDownload
	ErrGpgme
	ErrMissingCapabilitySignatures
)

// Error implements the error interface
func (e Errno) Error() string {
	return errnoStrings[e]
}

var errnoStrings = map[Errno]string{
	ErrOK:                          "No error",
	ErrMemory:                      "Failed to allocate memory",
	ErrSystem:                      "A system error occurred",
	ErrBadPerms:                    "Permission denied",
	ErrNotAFile:                    "Should be a file",
	ErrNotADir:                     "Should be a directory",
	ErrWrongArgs:                   "Function was called with invalid arguments",
	ErrDiskSpace:                   "Insufficient disk space",
	ErrHandleNull:                  "Handle should be null",
	ErrHandleNotNull:               "Handle should not be null",
	ErrHandleLock:                  "Failed to acquire lock",
	ErrDBOpen:                      "Failed to open database",
	ErrDBCreate:                    "Failed to create database",
	ErrDBNull:                      "Database should not be null",
	ErrDBNotNull:                   "Database should be null",
	ErrDBNotFound:                  "The database could not be found",
	ErrDBInvalid:                   "Database is invalid",
	ErrDBInvalidSig:                "Database has an invalid signature",
	ErrDBVersion:                   "The localdb is in a newer/older format than libalpm expects",
	ErrDBWrite:                     "Failed to write to the database",
	ErrDBRemove:                    "Failed to remove entry from database",
	ErrServerBadURL:                "Server URL is in an invalid format",
	ErrServerNone:                  "The database has no configured servers",
	ErrTransNotNull:                "A transaction is already initialized",
	ErrTransNull:                   "A transaction has not been initialized",
	ErrTransDupTarget:              "Duplicate target in transaction",
	ErrTransDupFilename:            "Duplicate filename in transaction",
	ErrTransNotInitialized:         "A transaction has not been initialized",
	ErrTransNotPrepared:            "Transaction has not been prepared",
	ErrTransAbort:                  "Transaction was aborted",
	ErrTransType:                   "Failed to interrupt transaction",
	ErrTransNotLocked:              "Tried to commit transaction without locking the database",
	ErrTransHookFailed:             "A hook failed to run",
	ErrPkgNotFound:                 "Package not found",
	ErrPkgIgnored:                  "Package is in ignorepkg",
	ErrPkgInvalid:                  "Package is invalid",
	ErrPkgInvalidChecksum:          "Package has an invalid checksum",
	ErrPkgInvalidSig:               "Package has an invalid signature",
	ErrPkgMissingSig:               "Package does not have a signature",
	ErrPkgOpen:                     "Cannot open the package file",
	ErrPkgCantRemove:               "Failed to remove package files",
	ErrPkgInvalidName:              "Package has an invalid name",
	ErrPkgInvalidArch:              "Package has an invalid architecture",
	ErrSigMissing:                  "Signatures are missing",
	ErrSigInvalid:                  "Signatures are invalid",
	ErrUnsatisfiedDeps:             "Dependencies could not be satisfied",
	ErrConflictingDeps:             "Conflicting dependencies",
	ErrFileConflicts:               "Files conflict",
	ErrRetrieve:                    "Download failed",
	ErrInvalidRegex:                "Invalid Regex",
	ErrLibArchive:                  "Error in libarchive",
	ErrLibCurl:                     "Error in libcurl",
	ErrExternalDownload:            "Error in external download program",
	ErrGpgme:                       "Error in gpgme",
	ErrMissingCapabilitySignatures: "Missing compile-time features",
}

// ALPMError wraps an ALPM error with context
type ALPMError struct {
	Errno Errno
	Msg   string
}

func (e *ALPMError) Error() string {
	if e.Msg != "" {
		return e.Errno.Error() + ": " + e.Msg
	}
	return e.Errno.Error()
}

// NewError creates a new ALPMError
func NewError(errno Errno, msg string) *ALPMError {
	return &ALPMError{
		Errno: errno,
		Msg:   msg,
	}
}
