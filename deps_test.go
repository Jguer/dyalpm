package dyalpm

import "testing"

func TestDepMod_Constants(t *testing.T) {
	// Verify DepMod constants have expected values (iota + 1)
	tests := []struct {
		mod      DepMod
		expected DepMod
	}{
		{DepModAny, 1},
		{DepModEQ, 2},
		{DepModGE, 3},
		{DepModLE, 4},
		{DepModGT, 5},
		{DepModLT, 6},
	}

	for _, tt := range tests {
		if tt.mod != tt.expected {
			t.Errorf("DepMod constant value = %d, want %d", tt.mod, tt.expected)
		}
	}
}

func TestDepend_String_NameOnly(t *testing.T) {
	d := Depend{
		Name:    "glibc",
		Version: "",
		Mod:     DepModAny,
	}

	result := d.String()
	if result != "glibc" {
		t.Errorf("Depend.String() = %q, want %q", result, "glibc")
	}
}

func TestDepend_String_WithVersion(t *testing.T) {
	tests := []struct {
		name     string
		dep      Depend
		expected string
	}{
		{
			name: "equal",
			dep: Depend{
				Name:    "glibc",
				Version: "2.38",
				Mod:     DepModEQ,
			},
			expected: "glibc=2.38",
		},
		{
			name: "greater or equal",
			dep: Depend{
				Name:    "gcc",
				Version: "12.0",
				Mod:     DepModGE,
			},
			expected: "gcc>=12.0",
		},
		{
			name: "less or equal",
			dep: Depend{
				Name:    "python",
				Version: "3.12",
				Mod:     DepModLE,
			},
			expected: "python<=3.12",
		},
		{
			name: "greater than",
			dep: Depend{
				Name:    "openssl",
				Version: "1.1",
				Mod:     DepModGT,
			},
			expected: "openssl>1.1",
		},
		{
			name: "less than",
			dep: Depend{
				Name:    "nodejs",
				Version: "20.0",
				Mod:     DepModLT,
			},
			expected: "nodejs<20.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.dep.String()
			if result != tt.expected {
				t.Errorf("Depend.String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestDepend_String_AnyWithVersion(t *testing.T) {
	// If Mod is DepModAny but version is present, no operator is added
	d := Depend{
		Name:    "foo",
		Version: "1.0",
		Mod:     DepModAny,
	}

	result := d.String()
	// With DepModAny, no operator should be added
	if result != "foo1.0" {
		t.Errorf("Depend.String() = %q, want %q", result, "foo1.0")
	}
}

func TestDepend_String_EmptyName(t *testing.T) {
	d := Depend{
		Name:    "",
		Version: "1.0",
		Mod:     DepModEQ,
	}

	result := d.String()
	if result != "=1.0" {
		t.Errorf("Depend.String() = %q, want %q", result, "=1.0")
	}
}

func TestDepend_String_ComplexVersions(t *testing.T) {
	tests := []struct {
		name     string
		dep      Depend
		expected string
	}{
		{
			name: "epoch version",
			dep: Depend{
				Name:    "systemd",
				Version: "1:255-1",
				Mod:     DepModGE,
			},
			expected: "systemd>=1:255-1",
		},
		{
			name: "git version",
			dep: Depend{
				Name:    "mypackage",
				Version: "1.0.0.r10.gabcdef",
				Mod:     DepModEQ,
			},
			expected: "mypackage=1.0.0.r10.gabcdef",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := tt.dep.String()
			if result != tt.expected {
				t.Errorf("Depend.String() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestDepend_Fields(t *testing.T) {
	d := Depend{
		Name:        "testpkg",
		Version:     "1.2.3",
		Description: "Test package description",
		NameHash:    12345,
		Mod:         DepModGE,
	}

	if d.Name != "testpkg" {
		t.Errorf("Name = %q, want %q", d.Name, "testpkg")
	}
	if d.Version != "1.2.3" {
		t.Errorf("Version = %q, want %q", d.Version, "1.2.3")
	}
	if d.Description != "Test package description" {
		t.Errorf("Description = %q, want %q", d.Description, "Test package description")
	}
	if d.NameHash != 12345 {
		t.Errorf("NameHash = %d, want %d", d.NameHash, 12345)
	}
	if d.Mod != DepModGE {
		t.Errorf("Mod = %d, want %d", d.Mod, DepModGE)
	}
}

func TestDepend_ZeroValue(t *testing.T) {
	var d Depend

	result := d.String()
	if result != "" {
		t.Errorf("zero Depend.String() = %q, want empty string", result)
	}
}

// Test Dependency interface (newDependency requires libalpm)
// These tests verify the interface contract

func TestDependency_Interface(t *testing.T) {
	// Verify Depend can be used where dependency info is needed
	d := Depend{
		Name:    "test",
		Version: "1.0",
		Mod:     DepModEQ,
	}

	// Depend should provide all necessary info through String()
	str := d.String()
	if str != "test=1.0" {
		t.Errorf("got %q, want %q", str, "test=1.0")
	}
}

// Benchmark dependency string generation
func BenchmarkDepend_String_Simple(b *testing.B) {
	d := Depend{Name: "glibc", Version: "", Mod: DepModAny}
	for b.Loop() {
		_ = d.String()
	}
}

func BenchmarkDepend_String_WithVersion(b *testing.B) {
	d := Depend{Name: "glibc", Version: "2.38", Mod: DepModGE}
	for b.Loop() {
		_ = d.String()
	}
}
