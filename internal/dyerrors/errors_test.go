package dyerrors

import "testing"

func TestErrno_Error(t *testing.T) {
	tests := []struct {
		errno    Errno
		expected string
	}{
		{ErrOK, "No error"},
		{ErrMemory, "Failed to allocate memory"},
		{ErrSystem, "A system error occurred"},
		{ErrBadPerms, "Permission denied"},
		{ErrNotAFile, "Should be a file"},
		{ErrNotADir, "Should be a directory"},
		{ErrWrongArgs, "Function was called with invalid arguments"},
		{ErrDiskSpace, "Insufficient disk space"},
		{ErrHandleNull, "Handle should be null"},
		{ErrHandleNotNull, "Handle should not be null"},
		{ErrHandleLock, "Failed to acquire lock"},
		{ErrDBOpen, "Failed to open database"},
		{ErrDBCreate, "Failed to create database"},
		{ErrDBNull, "Database should not be null"},
		{ErrDBNotNull, "Database should be null"},
		{ErrDBNotFound, "The database could not be found"},
		{ErrDBInvalid, "Database is invalid"},
		{ErrDBInvalidSig, "Database has an invalid signature"},
		{ErrDBVersion, "The localdb is in a newer/older format than libalpm expects"},
		{ErrDBWrite, "Failed to write to the database"},
		{ErrDBRemove, "Failed to remove entry from database"},
		{ErrServerBadURL, "Server URL is in an invalid format"},
		{ErrServerNone, "The database has no configured servers"},
		{ErrTransNotNull, "A transaction is already initialized"},
		{ErrTransNull, "A transaction has not been initialized"},
		{ErrTransDupTarget, "Duplicate target in transaction"},
		{ErrTransDupFilename, "Duplicate filename in transaction"},
		{ErrTransNotInitialized, "A transaction has not been initialized"},
		{ErrTransNotPrepared, "Transaction has not been prepared"},
		{ErrTransAbort, "Transaction was aborted"},
		{ErrTransType, "Failed to interrupt transaction"},
		{ErrTransNotLocked, "Tried to commit transaction without locking the database"},
		{ErrTransHookFailed, "A hook failed to run"},
		{ErrPkgNotFound, "Package not found"},
		{ErrPkgIgnored, "Package is in ignorepkg"},
		{ErrPkgInvalid, "Package is invalid"},
		{ErrPkgInvalidChecksum, "Package has an invalid checksum"},
		{ErrPkgInvalidSig, "Package has an invalid signature"},
		{ErrPkgMissingSig, "Package does not have a signature"},
		{ErrPkgOpen, "Cannot open the package file"},
		{ErrPkgCantRemove, "Failed to remove package files"},
		{ErrPkgInvalidName, "Package has an invalid name"},
		{ErrPkgInvalidArch, "Package has an invalid architecture"},
		{ErrSigMissing, "Signatures are missing"},
		{ErrSigInvalid, "Signatures are invalid"},
		{ErrUnsatisfiedDeps, "Dependencies could not be satisfied"},
		{ErrConflictingDeps, "Conflicting dependencies"},
		{ErrFileConflicts, "Files conflict"},
		{ErrRetrieve, "Download failed"},
		{ErrInvalidRegex, "Invalid Regex"},
		{ErrLibArchive, "Error in libarchive"},
		{ErrLibCurl, "Error in libcurl"},
		{ErrExternalDownload, "Error in external download program"},
		{ErrGpgme, "Error in gpgme"},
		{ErrMissingCapabilitySignatures, "Missing compile-time features"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			t.Parallel()
			got := tt.errno.Error()
			if got != tt.expected {
				t.Errorf("Errno(%d).Error() = %q, want %q", tt.errno, got, tt.expected)
			}
		})
	}
}

func TestErrno_UnknownError(t *testing.T) {
	// Test an unknown error code returns empty string
	unknown := Errno(999)
	got := unknown.Error()
	if got != "" {
		t.Errorf("Unknown Errno.Error() = %q, want empty string", got)
	}
}

func TestNewError(t *testing.T) {
	tests := []struct {
		name     string
		errno    Errno
		msg      string
		expected string
	}{
		{
			name:     "with message",
			errno:    ErrPkgNotFound,
			msg:      "package 'foo' not in database",
			expected: "Package not found: package 'foo' not in database",
		},
		{
			name:     "without message",
			errno:    ErrDiskSpace,
			msg:      "",
			expected: "Insufficient disk space",
		},
		{
			name:     "memory error with context",
			errno:    ErrMemory,
			msg:      "allocating buffer",
			expected: "Failed to allocate memory: allocating buffer",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := NewError(tt.errno, tt.msg)
			if err == nil {
				t.Fatal("NewError returned nil")
			}
			got := err.Error()
			if got != tt.expected {
				t.Errorf("ALPMError.Error() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestALPMError_Fields(t *testing.T) {
	err := NewError(ErrDBNotFound, "custom message")

	if err.Errno != ErrDBNotFound {
		t.Errorf("Errno = %v, want %v", err.Errno, ErrDBNotFound)
	}
	if err.Msg != "custom message" {
		t.Errorf("Msg = %q, want %q", err.Msg, "custom message")
	}
}

func TestErrno_ImplementsError(t *testing.T) {
	var _ error = ErrOK
	var _ error = ErrMemory
}

func TestALPMError_ImplementsError(t *testing.T) {
	var _ error = &ALPMError{}
	var _ error = NewError(ErrOK, "")
}
