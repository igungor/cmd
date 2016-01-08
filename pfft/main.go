package main

// TODO(ig): draw board on the center. positions are hardcoded currently.

import (
	"log"
	"strconv"

	"github.com/nsf/termbox-go"
)

func main() {
	if err := realMain(); err != nil {
		log.Fatal(err)
	}
}

func realMain() error {
	if err := termbox.Init(); err != nil {
		return err
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	drawBoard(110, 25, 15, 15)
	termbox.Flush()
mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				break mainloop
			}
		case termbox.EventError:
			return ev.Err
		}
		drawBoard(110, 25, 15, 15)
		termbox.Flush()
	}
	return nil
}

// drawBoard draws the board onto (x,y) position of the grid.
func drawBoard(x, y, w, h int) {
	// columns on the top
	for dx := 0; dx < w; dx++ {
		termbox.SetCell(x+1+dx*2, y-2, rune('A'+dx), fgcolor, bgcolor)
		termbox.SetCell(x+1+dx*2+1, y-2, ' ', fgcolor, bgcolor)
	}

	// rows on the left
	for dy := 0; dy < h; dy++ {
		if dy < 10 {
			tbprint(strconv.Itoa(dy+1), x-2, y+dy, fgcolor, bgcolor)
		} else {
			tbprint(strconv.Itoa(dy+1), x-3, y+dy, fgcolor, bgcolor)
		}
	}

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

	// TODO(ig): mark multipliers
	// TODO(ig): draw letters
}

// tbprint prints the msg onto (x,y) position of the grid.
func tbprint(msg string, x, y int, fg, bg termbox.Attribute) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

// fill fills a box of w*h area with r runes starting from (x,y) position of
// the grid.
func fill(x, y, w, h int, r rune) {
	for ly := 0; ly < h; ly++ {
		for lx := 0; lx < w; lx++ {
			termbox.SetCell(x+lx, y+ly, r, fgcolor, bgcolor)
		}
	}
}

const (
	fgcolor = termbox.ColorDefault
	bgcolor = termbox.ColorDefault
)
