package main

import (
	"unicode/utf8"

	termbox "github.com/nsf/termbox-go"
)

const tabstopLength = 8
const preferredHorizontalThreshold = 5

type editbox struct {
	text             []byte
	lineVisOffset    int
	curByteOffset    int // cursor offset in bytes
	curVisOffset     int // visual cursor offset in termbox cells
	curUnicodeOffset int // cursor offset in unicode code points
	w, h             int
}

// Draws the editbox in the given location, 'h' is not used at the moment
func (eb *editbox) draw(x, y int) {
	drawRect(x, y, eb.w, eb.h)
	eb.AdjustVOffset(eb.w)

	fill(x, y, eb.w, eb.h, ' ')

	t := eb.text
	lx := 0
	tabstop := 0
	for {
		rx := lx - eb.lineVisOffset
		if len(t) == 0 {
			break
		}

		if lx == tabstop {
			tabstop += tabstopLength
		}

		if rx >= eb.w {
			termbox.SetCell(x+eb.w-1, y, '→',
				fgcolor, bgcolor)
			break
		}

		r, size := utf8.DecodeRune(t)
		if r == '\t' {
			for ; lx < tabstop; lx++ {
				rx = lx - eb.lineVisOffset
				if rx >= eb.w {
					goto next
				}

				if rx >= 0 {
					termbox.SetCell(x+rx, y, ' ', fgcolor, bgcolor)
				}
			}
		} else {
			if rx >= 0 {
				termbox.SetCell(x+rx, y, r, fgcolor, bgcolor)
			}
			lx += 1
		}
	next:
		t = t[size:]
	}

	if eb.lineVisOffset != 0 {
		termbox.SetCell(x, y, '←', fgcolor, bgcolor)
	}
	termbox.SetCursor(x+eb.CursorX(), y)
}

// Adjusts line visual offset to a proper value depending on width
func (eb *editbox) AdjustVOffset(width int) {
	ht := preferredHorizontalThreshold
	max_h_threshold := (width - 1) / 2
	if ht > max_h_threshold {
		ht = max_h_threshold
	}

	threshold := width - 1
	if eb.lineVisOffset != 0 {
		threshold = width - ht
	}
	if eb.curVisOffset-eb.lineVisOffset >= threshold {
		eb.lineVisOffset = eb.curVisOffset + (ht - width + 1)
	}

	if eb.lineVisOffset != 0 && eb.curVisOffset-eb.lineVisOffset < ht {
		eb.lineVisOffset = eb.curVisOffset - ht
		if eb.lineVisOffset < 0 {
			eb.lineVisOffset = 0
		}
	}
}

func (eb *editbox) MoveCursorTo(boffset int) {
	eb.curByteOffset = boffset
	eb.curVisOffset, eb.curUnicodeOffset = visOffset2codeOffset(eb.text, boffset)
}

func (eb *editbox) RuneUnderCursor() (rune, int) {
	return utf8.DecodeRune(eb.text[eb.curByteOffset:])
}

func (eb *editbox) RuneBeforeCursor() (rune, int) {
	return utf8.DecodeLastRune(eb.text[:eb.curByteOffset])
}

func (eb *editbox) MoveCursorOneRuneBackward() {
	if eb.curByteOffset == 0 {
		return
	}
	_, size := eb.RuneBeforeCursor()
	eb.MoveCursorTo(eb.curByteOffset - size)
}

func (eb *editbox) MoveCursorOneRuneForward() {
	if eb.curByteOffset == len(eb.text) {
		return
	}
	_, size := eb.RuneUnderCursor()
	eb.MoveCursorTo(eb.curByteOffset + size)
}

func (eb *editbox) MoveCursorToBeginningOfTheLine() {
	eb.MoveCursorTo(0)
}

func (eb *editbox) MoveCursorToEndOfTheLine() {
	eb.MoveCursorTo(len(eb.text))
}

func (eb *editbox) DeleteRuneBackward() {
	if eb.curByteOffset == 0 {
		return
	}

	eb.MoveCursorOneRuneBackward()
	_, size := eb.RuneUnderCursor()
	eb.text = byteSliceRemove(eb.text, eb.curByteOffset, eb.curByteOffset+size)
}

func (eb *editbox) DeleteRuneForward() {
	if eb.curByteOffset == len(eb.text) {
		return
	}
	_, size := eb.RuneUnderCursor()
	eb.text = byteSliceRemove(eb.text, eb.curByteOffset, eb.curByteOffset+size)
}

func (eb *editbox) DeleteTheRestOfTheLine() {
	eb.text = eb.text[:eb.curByteOffset]
}

func (eb *editbox) InsertRune(r rune) {
	var buf [utf8.UTFMax]byte
	n := utf8.EncodeRune(buf[:], r)
	eb.text = byteSliceInsert(eb.text, eb.curByteOffset, buf[:n])
	eb.MoveCursorOneRuneForward()
}

// Please, keep in mind that cursor depends on the value of lineVisOffset, which
// is being set on Draw() call, so.. call this method after Draw() one.
func (eb *editbox) CursorX() int {
	return eb.curVisOffset - eb.lineVisOffset
}
func runeAdvanceLen(r rune, pos int) int {
	if r == '\t' {
		return tabstopLength - pos%tabstopLength
	}
	return 1
}

func visOffset2codeOffset(text []byte, boffset int) (voffset, coffset int) {
	text = text[:boffset]
	for len(text) > 0 {
		r, size := utf8.DecodeRune(text)
		text = text[size:]
		coffset += 1
		voffset += runeAdvanceLen(r, voffset)
	}
	return
}

func byteSliceGrow(s []byte, desired_cap int) []byte {
	if cap(s) < desired_cap {
		ns := make([]byte, len(s), desired_cap)
		copy(ns, s)
		return ns
	}
	return s
}

func byteSliceRemove(text []byte, from, to int) []byte {
	size := to - from
	copy(text[from:], text[to:])
	text = text[:len(text)-size]
	return text
}

func byteSliceInsert(text []byte, offset int, what []byte) []byte {
	n := len(text) + len(what)
	text = byteSliceGrow(text, n)
	text = text[:n]
	copy(text[offset+len(what):], text[offset:])
	copy(text[offset:], what)
	return text
}
