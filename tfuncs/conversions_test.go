package tfuncs

import "testing"

func Test_raw(t *testing.T) {
	got := raw("hello")

	if got != "hello" {
		t.Errorf("raw() = %v, want %v", got, "hello")
	}
}

func Test_linebyline(t *testing.T) {
	lines := linebyline("hello\nworld")

	if len(lines) != 2 {
		t.Errorf("linebyline() = %v, want %v", len(lines), 2)
	}

	if lines[0] != "hello" {
		t.Errorf("linebyline() = %v, want %v", lines[0], "hello")
	}

	if lines[1] != "world" {
		t.Errorf("linebyline() = %v, want %v", lines[1], "world")
	}
}
