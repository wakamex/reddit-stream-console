package app

import "testing"

func TestVersionGreater(t *testing.T) {
	cases := []struct {
		a, b string
		want bool
	}{
		{"1.10.0", "1.9.0", true},
		{"1.9.0", "1.10.0", false},
		{"2.0.0", "1.99.99", true},
		{"1.0.1", "1.0.0", true},
		{"1.0.0", "1.0.1", false},
		{"1.0.0", "1.0.0", false},
		{"0.2.0", "0.1.0", true},
	}
	for _, c := range cases {
		got := versionGreater(c.a, c.b)
		if got != c.want {
			t.Errorf("versionGreater(%q, %q) = %v, want %v", c.a, c.b, got, c.want)
		}
	}
}
