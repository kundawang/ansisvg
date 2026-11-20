package svgscreen

import "testing"

func TestSafeHref(t *testing.T) {
	safe := []string{
		"https://example.com/page",
		"http://example.com",
		"HTTPS://upper.example.com",
		"mailto:user@example.com",
		"ftp://host/file",
		"ssh://git@host/repo",
		"tel:+123456",
		"relative/path.svg",
		"/absolute/path",
		"#fragment",
		"//cdn.example.com/asset",
		"https://example.com/a;b?c=d&e=f",
	}
	for _, u := range safe {
		if !safeHref(u) {
			t.Errorf("expected safe: %q", u)
		}
	}

	unsafe := []string{
		"",                        // empty URI
		"javascript:alert(1)",     // script scheme
		"JaVaScRiPt:alert(1)",     // mixed case
		"java\x00script:alert(1)", // embedded control char
		"data:text/html,<script>", // data scheme
		"vbscript:msgbox(1)",      // vbscript
		"file:///etc/passwd",      // local file
		"blob:https://e.test/abc", // unknown scheme
		"custom-scheme:x",         // unknown scheme
		"https\x7f://e.test",      // DEL character
		"https://e.test/\x1b[31m", // escape character
		":bad-scheme",
		"1https://e.test", // digit-first scheme is not valid
		"ht tps://e.test", // whitespace inside URL rejected via control/space rules
	}
	for _, u := range unsafe {
		if safeHref(u) {
			t.Errorf("expected unsafe: %q", u)
		}
	}
}

func TestSetupLinkRectsWideAndGaps(t *testing.T) {
	s := &Screen{LineHeight: 1.0, Lines: []Line{
		{Y: 0, Chars: []Char{
			{Char: "汉", X: 0, Link: "https://w.test/"},
			{Char: "a", X: 2, Link: "https://w.test/"},
		}},
		// gap inside one link produces two rects
		{Y: 1, Chars: []Char{
			{Char: "a", X: 0, Link: "https://g.test/"},
			{Char: "b", X: 3, Link: "https://g.test/"},
		}},
		// unsafe link produces no rect at all
		{Y: 2, Chars: []Char{
			{Char: "x", X: 0, Link: "javascript:alert(1)"},
		}},
	}}

	s.setupLinkRects()

	if len(s.Dom.LinkRects) != 3 {
		t.Fatalf("got %d rects: %+v", len(s.Dom.LinkRects), s.Dom.LinkRects)
	}
	r0 := s.Dom.LinkRects[0]
	if r0.X != "0ch" || r0.Width != "3ch" || r0.Y != "0em" || string(r0.Href) != "https://w.test/" {
		t.Errorf("wide rect wrong: %+v", r0)
	}
	r1 := s.Dom.LinkRects[1]
	if r1.X != "0ch" || r1.Width != "1ch" {
		t.Errorf("gap rect 1 wrong: %+v", r1)
	}
	r2 := s.Dom.LinkRects[2]
	if r2.X != "3ch" || r2.Width != "1ch" || r2.Y != "1em" {
		t.Errorf("gap rect 2 wrong: %+v", r2)
	}
}
