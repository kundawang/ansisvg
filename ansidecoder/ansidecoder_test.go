package ansidecoder_test

import (
	"io"
	"strings"
	"testing"

	"github.com/wader/ansisvg/ansidecoder"
)

type decoded struct {
	r    rune
	x, y int
	link string
}

func decodeAll(t *testing.T, input string) []decoded {
	t.Helper()
	d := ansidecoder.NewDecoder(strings.NewReader(input))
	var out []decoded
	for {
		r, _, err := d.ReadRune()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		// control bytes of escape sequences must never be emitted
		if r == 0x1b || r == 0x07 {
			t.Fatalf("control byte %#x was emitted", r)
		}
		out = append(out, decoded{r: r, x: d.X, y: d.Y, link: d.Link})
	}
	return out
}

func TestOSC8HyperlinkST(t *testing.T) {
	// OSC 8 ; ; URI ST, spanning an SGR color change
	out := decodeAll(t, "\x1b]8;;https://example.com/p\x1b\\\x1b[31mab\x1b]8;;\x1b\\c")
	want := []decoded{
		{r: 'a', x: 0, y: 0, link: "https://example.com/p"},
		{r: 'b', x: 1, y: 0, link: "https://example.com/p"},
		{r: 'c', x: 2, y: 0, link: ""},
	}
	assertDecoded(t, out, want)
}

func TestOSC8HyperlinkBEL(t *testing.T) {
	// BEL terminator variant, including params (id=...)
	out := decodeAll(t, "\x1b]8;id=42;https://e.test/\x07ab\x1b]8;;\x07c")
	want := []decoded{
		{r: 'a', x: 0, y: 0, link: "https://e.test/"},
		{r: 'b', x: 1, y: 0, link: "https://e.test/"},
		{r: 'c', x: 2, y: 0, link: ""},
	}
	assertDecoded(t, out, want)
}

func TestOSC8HyperlinkAcrossNewline(t *testing.T) {
	out := decodeAll(t, "\x1b]8;;https://nl.test/\x1b\\ab\ncd\x1b]8;;\x1b\\")
	want := []decoded{
		{r: 'a', x: 0, y: 0, link: "https://nl.test/"},
		{r: 'b', x: 1, y: 0, link: "https://nl.test/"},
		{r: '\n', x: 2, y: 0, link: "https://nl.test/"},
		{r: 'c', x: 0, y: 1, link: "https://nl.test/"},
		{r: 'd', x: 1, y: 1, link: "https://nl.test/"},
	}
	assertDecoded(t, out, want)
}

func TestOSC8UnclosedLink(t *testing.T) {
	// no closing OSC 8: link remains active until EOF
	out := decodeAll(t, "x\x1b]8;;https://open.test/\x1b\\yz")
	want := []decoded{
		{r: 'x', x: 0, y: 0, link: ""},
		{r: 'y', x: 1, y: 0, link: "https://open.test/"},
		{r: 'z', x: 2, y: 0, link: "https://open.test/"},
	}
	assertDecoded(t, out, want)
}

func TestOSC8NestedStart(t *testing.T) {
	// a second OSC 8 start replaces the active link, then close ends it
	out := decodeAll(t, "\x1b]8;;https://outer/\x1b\\O\x1b]8;;https://inner/\x1b\\I\x1b]8;;\x1b\\N")
	want := []decoded{
		{r: 'O', x: 0, y: 0, link: "https://outer/"},
		{r: 'I', x: 1, y: 0, link: "https://inner/"},
		{r: 'N', x: 2, y: 0, link: ""},
	}
	assertDecoded(t, out, want)
}

func TestOSC8EmptyAndSemicolonURI(t *testing.T) {
	// URI may itself contain ';'
	out := decodeAll(t, "\x1b]8;;https://e.test/a;b?c=d;e=f\x1b\\z")
	if len(out) != 1 || out[0].link != "https://e.test/a;b?c=d;e=f" {
		t.Fatalf("got %+v", out)
	}
}

func TestOSC8MalformedAndOtherOSCIgnored(t *testing.T) {
	// malformed OSC 8 and unrelated OSC commands must not change link state
	out := decodeAll(t,
		"\x1b]8Xbad\x1b\\"+ // not "8;"
			"\x1b]1337;File=abc\x1b\\"+ // different OSC command
			"\x1b]8;;https://ok.test/\x1b\\"+
			"\x1b]89nope\x1b\\"+ // malformed while link active: state untouched
			"z")
	if len(out) != 1 || out[0].r != 'z' || out[0].link != "https://ok.test/" {
		t.Fatalf("got %+v", out)
	}
}

func TestOSCBadEscapeReprocessed(t *testing.T) {
	// ESC inside OSC not followed by '\\' is a bad escape: the payload is
	// discarded and the following rune is processed as ordinary input
	out := decodeAll(t, "\x1b]8;;https://b.test/\x1bXabc")
	want := []decoded{
		{r: 'X', x: 0, y: 0, link: ""},
		{r: 'a', x: 1, y: 0, link: ""},
		{r: 'b', x: 2, y: 0, link: ""},
		{r: 'c', x: 3, y: 0, link: ""},
	}
	assertDecoded(t, out, want)
}

func TestWideRuneAdvancesTwoColumns(t *testing.T) {
	out := decodeAll(t, "a汉b")
	want := []decoded{
		{r: 'a', x: 0, y: 0, link: ""},
		{r: '汉', x: 1, y: 0, link: ""},
		{r: 'b', x: 3, y: 0, link: ""},
	}
	assertDecoded(t, out, want)
}

func assertDecoded(t *testing.T, got, want []decoded) {
	t.Helper()
	if len(got) != len(want) {
		t.Fatalf("length mismatch\ngot  %+v\nwant %+v", got, want)
	}
	for i := range got {
		if got[i] != want[i] {
			t.Errorf("index %d\ngot  %+v\nwant %+v", i, got[i], want[i])
		}
	}
}
