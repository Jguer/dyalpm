package dyalpm

import "testing"

// io.Reader allows empty reads; they must not index into the buffer.
func TestChangelogReader_EmptyRead(t *testing.T) {
	r := &changelogReader{pkg: &package_{}, fp: 1}
	if n, err := r.Read(nil); n != 0 || err != nil {
		t.Fatalf("Read(nil) = %d, %v; want 0, nil", n, err)
	}
}
