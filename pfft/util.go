package main

import "github.com/nsf/termbox-go"

// tbprint prints the msg at (x,y) position of the grid.
func tbprint(msg string, x, y int, fg, bg termbox.Attribute) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

// fill fills a rectanngle at position (x,y) with area of w*h.
// the grid.
func fill(x, y, w, h int, r rune) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, r, fgcolor, bgcolor)
		}
	}
}

// drawRect draws a rectangle with unicode borders at position (x,y) with area of
// w*h.
func drawRect(x, y, w, h int) {
	// top border
	termbox.SetCell(x-1, y-1, '┌', fgcolor, bgcolor)
	fill(x, y-1, w*2, 1, '─')
	termbox.SetCell(x+w*2, y-1, '┐', fgcolor, bgcolor)

	// body border
	fill(x-1, y, 1, h, '│')
	fill(x+w*2, y, 1, h, '│')

	// bottom border
	termbox.SetCell(x-1, y+h, '└', fgcolor, bgcolor)
	fill(x, y+h, w*2, 1, '─')
	termbox.SetCell(x+w*2, y+h, '┘', fgcolor, bgcolor)
}
