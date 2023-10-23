package tfuncs

import "testing"

func Test_rndgen(t *testing.T) {
	got1 := rndgen(5)
	got2 := rndgen(5)

	if got1 == got2 {
		t.Errorf("rndgen() values must be different, got = %q, %q", got1, got2)
	}
}
