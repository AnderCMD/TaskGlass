//go:build windows

package xaml

import "testing"

func TestEncode(t *testing.T) {
	cases := []struct {
		enabled bool
		pct     int
		want    int32
	}{
		{false, 50, 0}, {false, 0, 0}, {true, 0, 1}, {true, 60, 61}, {true, 100, 101}, {true, 250, 101}, {true, -9, 1},
	}
	for _, c := range cases {
		if got := Encode(c.enabled, c.pct); got != c.want {
			t.Errorf("Encode(%v, %d) = %d, want %d", c.enabled, c.pct, got, c.want)
		}
	}
}
