package dyalpm

import "testing"

func TestVerCmp_Basic(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected int
	}{
		// Basic comparisons
		{"equal versions", "1.0", "1.0", 0},
		{"v1 less than v2", "1.0", "2.0", -1},
		{"v1 greater than v2", "2.0", "1.0", 1},

		// Multi-component versions
		{"patch version difference", "1.0.1", "1.0.2", -1},
		{"patch vs no patch", "1.0.1", "1.0", 1},
		{"longer version", "1.0.0.1", "1.0.0", 1},

		// Release comparisons
		{"release difference", "1.0-1", "1.0-2", -1},
		{"same version different release", "2.0-1", "2.0-3", -1},
		{"release vs no release", "1.0-1", "1.0", 0}, // libalpm only compares releases when both exist
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := VerCmp(tt.v1, tt.v2)
			result = normalizeResult(result)
			if result != tt.expected {
				t.Errorf("VerCmp(%q, %q) = %d, want %d", tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}

func TestVerCmp_Epochs(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected int
	}{
		{"epoch wins over version", "1:1.0", "2.0", 1},
		{"higher epoch wins", "2:1.0", "1:2.0", 1},
		{"same epoch, compare versions", "1:1.0", "1:2.0", -1},
		{"epoch 0 implied", "0:1.0", "1.0", 0},
		{"epoch on both", "1:1.0", "1:1.0", 0},
		{"explicit epoch 0 vs implicit", "0:2.0", "1.0", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := VerCmp(tt.v1, tt.v2)
			result = normalizeResult(result)
			if result != tt.expected {
				t.Errorf("VerCmp(%q, %q) = %d, want %d", tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}

func TestVerCmp_Alpha(t *testing.T) {
	// Note: libalpm version comparison differs from semantic versioning.
	// In libalpm, alphabetic suffixes are compared lexicographically and
	// having additional content (like "alpha") makes a version "greater".
	// Note: In libalpm, when a version segment is exhausted on one side
	// but the other has remaining alpha characters, the empty one wins.
	// This is because alpha suffixes are considered "pre-release" indicators.
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected int
	}{
		{"alpha vs beta", "1.0a", "1.0b", -1},
		{"alpha suffix vs bare", "1.0alpha", "1.0", -1}, // alpha suffix loses to bare
		{"beta vs rc", "1.0beta", "1.0rc", -1},
		{"same alpha", "1.0a", "1.0a", 0},
		{"bare vs alpha suffix", "1.0", "1.0a", 1}, // bare beats alpha suffix
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := VerCmp(tt.v1, tt.v2)
			result = normalizeResult(result)
			if result != tt.expected {
				t.Errorf("VerCmp(%q, %q) = %d, want %d", tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}

func TestVerCmp_LeadingZeros(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected int
	}{
		{"leading zeros equal", "01", "1", 0},
		{"leading zeros in segment", "1.01", "1.1", 0},
		{"many leading zeros", "001.002.003", "1.2.3", 0},
		{"zero vs nonzero", "0", "1", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := VerCmp(tt.v1, tt.v2)
			result = normalizeResult(result)
			if result != tt.expected {
				t.Errorf("VerCmp(%q, %q) = %d, want %d", tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}

func TestVerCmp_SpecialCases(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected int
	}{
		{"empty strings", "", "", 0},
		{"empty vs non-empty", "", "1.0", -1},
		{"non-empty vs empty", "1.0", "", 1},
		{"dots only diff", "1.0.0", "1.0", 1},
		{"underscores ignored", "1_0", "1.0", 0},
		{"mixed separators", "1.0_1", "1.0.1", 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := VerCmp(tt.v1, tt.v2)
			result = normalizeResult(result)
			if result != tt.expected {
				t.Errorf("VerCmp(%q, %q) = %d, want %d", tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}

func TestVerCmp_RealWorldVersions(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected int
	}{
		// Linux kernel versions
		{"kernel versions", "5.15.0", "5.15.1", -1},
		{"kernel major", "5.15.0", "6.0.0", -1},

		// GCC versions
		{"gcc versions", "12.1.0", "12.2.0", -1},

		// Python versions
		{"python versions", "3.10.0", "3.11.0", -1},
		{"python micro", "3.10.5", "3.10.6", -1},

		// Complex arch packages
		{"arch package", "1.2.3-4", "1.2.3-5", -1},
		{"arch epoch", "1:1.0-1", "2.0-1", 1},

		// Git describe style
		{"git describe", "1.0.0.r10.gabcdef", "1.0.0.r11.g123456", -1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := VerCmp(tt.v1, tt.v2)
			result = normalizeResult(result)
			if result != tt.expected {
				t.Errorf("VerCmp(%q, %q) = %d, want %d", tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}

func TestVerCmp_Symmetry(t *testing.T) {
	// Test that VerCmp(a, b) == -VerCmp(b, a) for all cases
	pairs := [][2]string{
		{"1.0", "2.0"},
		{"1:1.0", "2.0"},
		{"1.0-1", "1.0-2"},
		{"1.0a", "1.0b"},
		{"1.0.1", "1.0"},
	}

	for _, pair := range pairs {
		v1, v2 := pair[0], pair[1]
		r1 := normalizeResult(VerCmp(v1, v2))
		r2 := normalizeResult(VerCmp(v2, v1))
		if r1 != -r2 {
			t.Errorf("VerCmp(%q, %q) = %d, but VerCmp(%q, %q) = %d (not symmetric)",
				v1, v2, r1, v2, v1, r2)
		}
	}
}

func TestVerCmp_Transitivity(t *testing.T) {
	// If a < b and b < c, then a < c
	triples := [][3]string{
		{"1.0", "2.0", "3.0"},
		{"1.0-1", "1.0-2", "1.0-3"},
		{"1.0a", "1.0b", "1.0c"},
		{"1:1.0", "2:1.0", "3:1.0"},
	}

	for _, triple := range triples {
		a, b, c := triple[0], triple[1], triple[2]
		ab := normalizeResult(VerCmp(a, b))
		bc := normalizeResult(VerCmp(b, c))
		ac := normalizeResult(VerCmp(a, c))

		if ab == -1 && bc == -1 && ac != -1 {
			t.Errorf("Transitivity failed: %q < %q < %q, but VerCmp(%q, %q) = %d",
				a, b, c, a, c, ac)
		}
	}
}

// Helper functions for testing

func normalizeResult(r int) int {
	if r < 0 {
		return -1
	}
	if r > 0 {
		return 1
	}
	return 0
}

// Test internal helper functions

func TestSplitVersion(t *testing.T) {
	tests := []struct {
		input   string
		epoch   string
		version string
		rel     string
	}{
		{"1.0", "0", "1.0", ""},
		{"1.0-1", "0", "1.0", "1"},
		{"1:1.0", "1", "1.0", ""},
		{"1:1.0-1", "1", "1.0", "1"},
		{"2:3.4.5-6", "2", "3.4.5", "6"},
		{"1.0-1-2", "0", "1.0-1", "2"}, // Only last dash is release
		// Epoch parsing edge cases (only digits before ':' count as epoch)
		{"a:0", "0", "a:0", ""},            // 'a' is not a digit, so no epoch
		{"1a:0", "0", "1a:0", ""},          // '1a' contains non-digit, so no epoch
		{":1.0", "0", "1.0", ""},           // Empty epoch before ':' means "0"
		{"0:1.0", "0", "1.0", ""},          // Explicit epoch "0"
		{"12:3.4.5-6", "12", "3.4.5", "6"}, // Multi-digit epoch
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			epoch, version, rel := splitVersion(tt.input)
			if epoch != tt.epoch {
				t.Errorf("epoch = %q, want %q", epoch, tt.epoch)
			}
			if version != tt.version {
				t.Errorf("version = %q, want %q", version, tt.version)
			}
			if rel != tt.rel {
				t.Errorf("rel = %q, want %q", rel, tt.rel)
			}
		})
	}
}

func TestCompareNumericString(t *testing.T) {
	tests := []struct {
		a, b     string
		expected int
	}{
		{"1", "2", -1},
		{"2", "1", 1},
		{"1", "1", 0},
		{"10", "2", 1},
		{"01", "1", 0},
		{"001", "1", 0},
		{"100", "99", 1},
	}

	for _, tt := range tests {
		t.Run(tt.a+"_"+tt.b, func(t *testing.T) {
			t.Parallel()
			result := compareNumericString(tt.a, tt.b)
			result = normalizeResult(result)
			if result != tt.expected {
				t.Errorf("compareNumericString(%q, %q) = %d, want %d",
					tt.a, tt.b, result, tt.expected)
			}
		})
	}
}

func TestStripLeadingZeros(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"0", "0"},
		{"00", "0"},
		{"01", "1"},
		{"001", "1"},
		{"123", "123"},
		{"0123", "123"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			result := stripLeadingZeros(tt.input)
			if result != tt.expected {
				t.Errorf("stripLeadingZeros(%q) = %q, want %q",
					tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsDigit(t *testing.T) {
	for c := byte('0'); c <= '9'; c++ {
		if !isDigit(c) {
			t.Errorf("isDigit(%q) = false, want true", c)
		}
	}
	for _, c := range []byte{'a', 'z', 'A', 'Z', '-', '.', '_'} {
		if isDigit(c) {
			t.Errorf("isDigit(%q) = true, want false", c)
		}
	}
}

func TestIsAlpha(t *testing.T) {
	for c := byte('a'); c <= 'z'; c++ {
		if !isAlpha(c) {
			t.Errorf("isAlpha(%q) = false, want true", c)
		}
	}
	for c := byte('A'); c <= 'Z'; c++ {
		if !isAlpha(c) {
			t.Errorf("isAlpha(%q) = false, want true", c)
		}
	}
	for _, c := range []byte{'0', '9', '-', '.', '_'} {
		if isAlpha(c) {
			t.Errorf("isAlpha(%q) = true, want false", c)
		}
	}
}

func TestIsAlnum(t *testing.T) {
	// Should be true for digits and letters
	for c := byte('0'); c <= '9'; c++ {
		if !isAlnum(c) {
			t.Errorf("isAlnum(%q) = false, want true", c)
		}
	}
	for c := byte('a'); c <= 'z'; c++ {
		if !isAlnum(c) {
			t.Errorf("isAlnum(%q) = false, want true", c)
		}
	}
	// Should be false for non-alphanumeric
	for _, c := range []byte{'-', '.', '_', ':', ' '} {
		if isAlnum(c) {
			t.Errorf("isAlnum(%q) = true, want false", c)
		}
	}
}

// Benchmarks

func BenchmarkVerCmp_Simple(b *testing.B) {
	for b.Loop() {
		VerCmp("1.0", "2.0")
	}
}

func BenchmarkVerCmp_Complex(b *testing.B) {
	for b.Loop() {
		VerCmp("1:5.15.0.arch1-1", "1:5.15.1.arch1-1")
	}
}

func BenchmarkVerCmp_Equal(b *testing.B) {
	for b.Loop() {
		VerCmp("1.2.3.4.5-6", "1.2.3.4.5-6")
	}
}
