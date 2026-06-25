//go:build integration

package dyalpm

import (
	"runtime"
	"testing"
	"unsafe"

	"github.com/ebitengine/purego"
)

// Integration test that dynamically loads alpm_pkg_vercmp from libalpm
// and compares the output against our pure Go implementation.
//
// Run with: go test -tags=integration -v -run TestVerCmp_Integration

const (
	libalpmPath         = "libalpm.so.16"
	libalpmPathFallback = "libalpm.so"
)

// loadVercmp loads alpm_pkg_vercmp from libalpm.so
func loadVercmp() (func(string, string) int, error) {
	// Try primary path first
	lib, err := purego.Dlopen(libalpmPath, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		// Try fallback path
		lib, err = purego.Dlopen(libalpmPathFallback, purego.RTLD_NOW|purego.RTLD_GLOBAL)
		if err != nil {
			return nil, err
		}
	}

	vercmpPtr, err := purego.Dlsym(lib, "alpm_pkg_vercmp")
	if err != nil {
		return nil, err
	}

	return func(a, b string) int {
		aBytes := append([]byte(a), 0)
		bBytes := append([]byte(b), 0)

		result, _, _ := purego.SyscallN(vercmpPtr,
			uintptr(unsafe.Pointer(&aBytes[0])),
			uintptr(unsafe.Pointer(&bBytes[0])))
		runtime.KeepAlive(aBytes)
		runtime.KeepAlive(bBytes)

		return int(int32(result)) // signed int
	}, nil
}

// normalizeIntegrationResult normalizes the result to -1, 0, or 1
func normalizeIntegrationResult(r int) int {
	if r < 0 {
		return -1
	}
	if r > 0 {
		return 1
	}
	return 0
}

func TestVerCmp_Integration(t *testing.T) {
	libalpmVercmp, err := loadVercmp()
	if err != nil {
		t.Skipf("Could not load libalpm: %v", err)
	}

	// Comprehensive test cases covering various version formats
	testCases := []struct {
		name string
		v1   string
		v2   string
	}{
		// Basic comparisons
		{"equal", "1.0", "1.0"},
		{"simple_less", "1.0", "2.0"},
		{"simple_greater", "2.0", "1.0"},

		// Multi-component versions
		{"patch_diff", "1.0.1", "1.0.2"},
		{"patch_vs_no_patch", "1.0.1", "1.0"},
		{"longer_version", "1.0.0.1", "1.0.0"},

		// Releases
		{"release_diff", "1.0-1", "1.0-2"},
		{"release_vs_no_release", "1.0-1", "1.0"},
		{"both_releases", "2.0-1", "2.0-3"},

		// Epochs
		{"epoch_wins", "1:1.0", "2.0"},
		{"higher_epoch", "2:1.0", "1:2.0"},
		{"same_epoch", "1:1.0", "1:2.0"},
		{"epoch_0_implied", "0:1.0", "1.0"},
		{"epoch_complex", "1:5.15.0-1", "5.15.0-1"},

		// Alpha suffixes
		{"alpha_vs_beta", "1.0a", "1.0b"},
		{"alpha_vs_bare", "1.0alpha", "1.0"},
		{"bare_vs_alpha", "1.0", "1.0a"},
		{"beta_vs_rc", "1.0beta", "1.0rc"},
		{"alpha_vs_numeric", "1.0a", "1.0.1"},

		// Leading zeros
		{"leading_zero", "01", "1"},
		{"leading_zeros_segment", "1.01", "1.1"},
		{"many_leading_zeros", "001.002.003", "1.2.3"},

		// Empty strings
		{"empty_both", "", ""},
		{"empty_vs_value", "", "1.0"},
		{"value_vs_empty", "1.0", ""},

		// Special characters/separators
		{"underscore", "1_0", "1.0"},
		{"mixed_separators", "1.0_1", "1.0.1"},
		{"dots_diff", "1.0.0", "1.0"},

		// Real-world versions
		{"kernel", "5.15.0", "5.15.1"},
		{"kernel_major", "5.15.0", "6.0.0"},
		{"gcc", "12.1.0", "12.2.0"},
		{"python", "3.10.0", "3.11.0"},
		{"arch_pkg", "1.2.3-4", "1.2.3-5"},
		{"git_describe", "1.0.0.r10.gabcdef", "1.0.0.r11.g123456"},

		// Edge cases
		{"single_digit", "1", "2"},
		{"single_alpha", "a", "b"},
		{"numeric_vs_alpha", "1", "a"},
		{"long_version", "1.2.3.4.5.6.7.8.9", "1.2.3.4.5.6.7.8.10"},
		{"very_long_number", "99999999999", "99999999998"},

		// More complex scenarios
		{"rc_versions", "1.0rc1", "1.0rc2"},
		{"pre_release", "1.0pre1", "1.0"},
		{"post_release", "1.0.post1", "1.0"},
		{"dev_version", "1.0.dev1", "1.0"},
		{"snapshot", "1.0.20210101", "1.0.20210102"},

		// Arch Linux specific patterns
		{"arch_kernel", "6.6.7.arch1-1", "6.6.8.arch1-1"},
		{"arch_zen", "6.6.7.zen1-1", "6.6.8.zen1-1"},
		{"any_arch", "2024.01.01-1", "2024.01.02-1"},

		// Separator edge cases
		{"dash_in_version", "1.0-beta-1", "1.0-1"},
		{"multiple_dashes", "1-2-3-4", "1-2-3-5"},
	}

	var mismatches []struct {
		name     string
		v1, v2   string
		goResult int
		cResult  int
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			goResult := normalizeIntegrationResult(VerCmp(tc.v1, tc.v2))
			cResult := normalizeIntegrationResult(libalpmVercmp(tc.v1, tc.v2))

			if goResult != cResult {
				t.Errorf("VerCmp(%q, %q): Go=%d, libalpm=%d",
					tc.v1, tc.v2, goResult, cResult)
				mismatches = append(mismatches, struct {
					name     string
					v1, v2   string
					goResult int
					cResult  int
				}{tc.name, tc.v1, tc.v2, goResult, cResult})
			}
		})
	}

	if len(mismatches) > 0 {
		t.Logf("\nSummary of mismatches (%d total):", len(mismatches))
		for _, m := range mismatches {
			t.Logf("  %s: VerCmp(%q, %q) Go=%d vs C=%d",
				m.name, m.v1, m.v2, m.goResult, m.cResult)
		}
	}
}

func TestVerCmp_Integration_Symmetry(t *testing.T) {
	libalpmVercmp, err := loadVercmp()
	if err != nil {
		t.Skipf("Could not load libalpm: %v", err)
	}

	pairs := []struct {
		v1, v2 string
	}{
		{"1.0", "2.0"},
		{"1:1.0", "2.0"},
		{"1.0-1", "1.0-2"},
		{"1.0a", "1.0b"},
		{"1.0.1", "1.0"},
		{"1.0alpha", "1.0"},
	}

	for _, p := range pairs {
		t.Run(p.v1+"_vs_"+p.v2, func(t *testing.T) {
			goForward := normalizeIntegrationResult(VerCmp(p.v1, p.v2))
			goReverse := normalizeIntegrationResult(VerCmp(p.v2, p.v1))
			cForward := normalizeIntegrationResult(libalpmVercmp(p.v1, p.v2))
			cReverse := normalizeIntegrationResult(libalpmVercmp(p.v2, p.v1))

			// Go should be symmetric
			if goForward != -goReverse {
				t.Errorf("Go not symmetric: VerCmp(%q,%q)=%d, VerCmp(%q,%q)=%d",
					p.v1, p.v2, goForward, p.v2, p.v1, goReverse)
			}

			// C should be symmetric
			if cForward != -cReverse {
				t.Errorf("C not symmetric: vercmp(%q,%q)=%d, vercmp(%q,%q)=%d",
					p.v1, p.v2, cForward, p.v2, p.v1, cReverse)
			}

			// Both should match
			if goForward != cForward {
				t.Errorf("Go/C mismatch for %q vs %q: Go=%d, C=%d",
					p.v1, p.v2, goForward, cForward)
			}
		})
	}
}

func TestVerCmp_Integration_Fuzz(t *testing.T) {
	libalpmVercmp, err := loadVercmp()
	if err != nil {
		t.Skipf("Could not load libalpm: %v", err)
	}

	// Generate various version-like strings for fuzzing
	components := []string{
		"0", "1", "2", "10", "99", "123",
		"a", "b", "z", "alpha", "beta", "rc",
		"", "0", "00", "01",
	}
	separators := []string{".", "-", "_", ":"}

	var mismatches int
	total := 0

	for _, c1 := range components {
		for _, s1 := range separators {
			for _, c2 := range components {
				v1 := c1 + s1 + c2
				for _, c3 := range components[:6] { // limit to avoid too many tests
					v2 := c3 + s1 + c2
					total++

					goResult := normalizeIntegrationResult(VerCmp(v1, v2))
					cResult := normalizeIntegrationResult(libalpmVercmp(v1, v2))

					if goResult != cResult {
						mismatches++
						if mismatches <= 10 { // Only report first 10
							t.Errorf("VerCmp(%q, %q): Go=%d, C=%d",
								v1, v2, goResult, cResult)
						}
					}
				}
			}
		}
	}

	if mismatches > 0 {
		t.Errorf("Total mismatches: %d out of %d tests", mismatches, total)
	} else {
		t.Logf("All %d fuzz tests passed", total)
	}
}

// Benchmark comparing pure Go vs C implementation
func BenchmarkVerCmp_Integration_Go(b *testing.B) {
	versions := []struct{ v1, v2 string }{
		{"1.0", "2.0"},
		{"1:5.15.0.arch1-1", "1:5.15.1.arch1-1"},
		{"1.2.3.4.5-6", "1.2.3.4.5-6"},
	}

	for b.Loop() {
		for _, v := range versions {
			VerCmp(v.v1, v.v2)
		}
	}
}

func BenchmarkVerCmp_Integration_C(b *testing.B) {
	libalpmVercmp, err := loadVercmp()
	if err != nil {
		b.Skipf("Could not load libalpm: %v", err)
	}

	versions := []struct{ v1, v2 string }{
		{"1.0", "2.0"},
		{"1:5.15.0.arch1-1", "1:5.15.1.arch1-1"},
		{"1.2.3.4.5-6", "1.2.3.4.5-6"},
	}

	for b.Loop() {
		for _, v := range versions {
			libalpmVercmp(v.v1, v.v2)
		}
	}
}
