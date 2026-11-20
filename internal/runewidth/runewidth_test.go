package runewidth_test

import (
	"testing"

	"github.com/wader/ansisvg/internal/runewidth"
)

func TestRuneWidth(t *testing.T) {
	cases := []struct {
		r    rune
		want int
	}{
		// narrow
		{'a', 1},
		{'1', 1},
		{' ', 1},
		{'\n', 1},
		// block elements used in existing fixtures are narrow
		{'█', 1},
		{'▄', 1},
		// wide: CJK ideographs, kana, fullwidth, hangul, wide emoji
		{'汉', 2},
		{'あ', 2},
		{'한', 2},
		{'Ａ', 2},
		{'，', 2},
		{'😀', 2},
		{'🫶', 2},
		// outside any table
		{'\U00010000', 1},
	}
	for _, c := range cases {
		if got := runewidth.RuneWidth(c.r); got != c.want {
			t.Errorf("RuneWidth(%q) = %d, want %d", c.r, got, c.want)
		}
	}
}
