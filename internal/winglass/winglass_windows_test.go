//go:build windows

package winglass

import "testing"

func TestAlpha(t *testing.T) {
	cases := map[int]uint8{100: 255, 150: 255, 85: 216, 50: 127, 10: 25, 5: 25, 0: 25, -4: 25}
	for pct, want := range cases {
		if got := Alpha(pct); got != want {
			t.Errorf("Alpha(%d) = %d, want %d", pct, got, want)
		}
	}
}
