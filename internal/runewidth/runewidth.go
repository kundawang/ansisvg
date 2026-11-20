package runewidth

import "sort"

// runeInterval is a half-open Unicode code point interval [Lo, Hi].
type runeInterval struct {
	Lo rune
	Hi rune
}

// wideRuneIntervals contains code points that occupy two terminal columns
// (East Asian Wide/Fullwidth and wide-presentation emoji), following the
// Unicode East Asian Width property (W and F) together with the emoji
// presentation ranges used by common terminals. The table is sorted and
// searched with a binary search in RuneWidth.
var wideRuneIntervals = []runeInterval{
	{0x1100, 0x115F}, // Hangul Jamo init. consonants
	{0x231A, 0x231B},
	{0x2329, 0x232A},
	{0x23E9, 0x23EC},
	{0x23F0, 0x23F0},
	{0x23F3, 0x23F3},
	{0x25FD, 0x25FE},
	{0x2614, 0x2615},
	{0x2648, 0x2653},
	{0x267F, 0x267F},
	{0x2693, 0x2693},
	{0x26A1, 0x26A1},
	{0x26AA, 0x26AB},
	{0x26BD, 0x26BE},
	{0x26C4, 0x26C5},
	{0x26CE, 0x26CE},
	{0x26D4, 0x26D4},
	{0x26EA, 0x26EA},
	{0x26F2, 0x26F3},
	{0x26F5, 0x26F5},
	{0x26FA, 0x26FA},
	{0x26FD, 0x26FD},
	{0x2705, 0x2705},
	{0x270A, 0x270B},
	{0x2728, 0x2728},
	{0x274C, 0x274C},
	{0x274E, 0x274E},
	{0x2753, 0x2755},
	{0x2757, 0x2757},
	{0x2795, 0x2797},
	{0x27B0, 0x27B0},
	{0x27BF, 0x27BF},
	{0x2B1B, 0x2B1C},
	{0x2B50, 0x2B50},
	{0x2B55, 0x2B55},
	{0x2E80, 0x2FDF}, // CJK Radicals Supplement, Kangxi Radicals
	{0x2FF0, 0x303E}, // Ideographic Description Characters, CJK symbols
	{0x3041, 0x33FF}, // Hiragana, Katakana, Bopomofo, Hangul Compatibility Jamo, enclosed CJK
	{0x3400, 0x4DBF}, // CJK Unified Ideographs Extension A
	{0x4DC0, 0x4DFF}, // Yijing Hexagram Symbols
	{0x4E00, 0x9FFF}, // CJK Unified Ideographs
	{0xA000, 0xA4CF}, // Yi Syllables and Radicals
	{0xA960, 0xA97F}, // Hangul Jamo Extended-A
	{0xAC00, 0xD7A3}, // Hangul Syllables
	{0xF900, 0xFAFF}, // CJK Compatibility Ideographs
	{0xFE10, 0xFE19}, // Vertical Forms
	{0xFE30, 0xFE6F}, // CJK Compatibility Forms and small form variants
	{0xFF01, 0xFF60}, // Fullwidth ASCII variants
	{0xFFE0, 0xFFE6}, // Fullwidth signs
	{0x16FE0, 0x16FE0},
	{0x17000, 0x18AFF}, // Tangut, Tangut Components, Khitan Small Script
	{0x1B000, 0x1B0FF},
	{0x1B100, 0x1B122},
	{0x1B132, 0x1B132},
	{0x1B150, 0x1B152},
	{0x1B164, 0x1B167},
	{0x1B170, 0x1B2FF}, // Nushu
	{0x1AFF0, 0x1AFF3},
	{0x1AFF5, 0x1AFFB},
	{0x1AFFD, 0x1AFFE},
	{0x1F004, 0x1F004},
	{0x1F0CF, 0x1F0CF},
	{0x1F18E, 0x1F18E},
	{0x1F191, 0x1F19A},
	{0x1F200, 0x1F202},
	{0x1F210, 0x1F23B},
	{0x1F240, 0x1F248},
	{0x1F250, 0x1F251},
	{0x1F260, 0x1F265},
	{0x1F300, 0x1F64F}, // Misc Symbols and Pictographs, Emoticons
	{0x1F680, 0x1F6FF}, // Transport and Map
	{0x1F780, 0x1F7FF},
	{0x1F900, 0x1F9FF}, // Supplemental Symbols and Pictographs
	{0x1FA00, 0x1FA6F},
	{0x1FA70, 0x1FAFF},
	{0x1FB00, 0x1FBFF},
	{0x20000, 0x2FFFD}, // CJK Unified Ideographs Extensions B-F and compat
	{0x30000, 0x3FFFD}, // CJK Compatibility Ideographs Supplement
	{0x31350, 0x3135F},
}

// RuneWidth returns the number of fixed-pitch columns a rune occupies.
// Wide glyphs (CJK, fullwidth and wide emoji) cover two columns, every
// other rune covers one column.
func RuneWidth(r rune) int {
	i := sort.Search(len(wideRuneIntervals), func(i int) bool {
		return wideRuneIntervals[i].Hi >= r
	})
	if i < len(wideRuneIntervals) && r >= wideRuneIntervals[i].Lo {
		return 2
	}
	return 1
}
