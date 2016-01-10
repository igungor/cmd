package main

import (
	"strconv"
	"unicode/utf8"

	"github.com/igungor/quackle"
	"github.com/nsf/termbox-go"
)

const boardsize = 15

type board struct {
	qb        quackle.Board
	x, y      int
	w, h      int
	showScore bool
}

func (b *board) draw() {
	sw, sh := termbox.Size()
	x := (sw - b.w*2 + 2 + 1) / 2
	y := (sh - b.h + 1 + 1) / 2

	// columns on the top
	for dx := 0; dx < b.w; dx++ {
		termbox.SetCell(x+1+dx*2, y-2, rune('A'+dx), fgcolor, bgcolor)
		termbox.SetCell(x+1+dx*2+1, y-2, ' ', fgcolor, bgcolor)
	}

	// rows on the left
	for dy := 0; dy < b.h; dy++ {
		if dy+1 < 10 {
			tbprint(strconv.Itoa(dy+1), x-2, y+dy, fgcolor, bgcolor)
		} else {
			tbprint(strconv.Itoa(dy+1), x-3, y+dy, fgcolor, bgcolor)
		}
	}

	// draw borders
	drawRect(x, y, b.w, b.h)

	// multipliers and letters
	for row := 0; row < b.h; row++ {
		for col := 0; col < b.w; col++ {
			// mark letters
			bl := b.qb.Letter(row, col)
			if dm.AlphabetParameters().IsPlainLetter(bl) {
				letter := dm.AlphabetParameters().UserVisible(bl)
				score := dm.AlphabetParameters().Score(bl)
				r, _ := utf8.DecodeRuneInString(letter)
				termbox.SetCell(x+col*2, y+row, r, fgcolor, bgcolor)
				if b.showScore {
					termbox.SetCell(x+col*2+1, y+row, getScoreRune(score), fgcolor, bgcolor)
				}
				continue
			}

			// mark multipliers
			letterMult := dm.BoardParameters().LetterMultiplier(row, col)
			wordMult := dm.BoardParameters().WordMultiplier(row, col)
			switch {
			case letterMult == 2:
				tbprint("H²", x+col*2, y+row, termbox.ColorWhite, termbox.ColorBlue)
			case letterMult == 3:
				tbprint("H³", x+col*2, y+row, termbox.ColorWhite, termbox.ColorMagenta)
			case wordMult == 2:
				tbprint("K²", x+col*2, y+row, termbox.ColorWhite, termbox.ColorGreen)
			case wordMult == 3:
				tbprint("K³", x+col*2, y+row, termbox.ColorWhite, termbox.ColorBlack)
			default:
				termbox.SetCell(x+col*2, y+row, ' ', fgcolor, bgcolor)
			}
		}
	}
}

type legend struct {
	x, y int
	w, h int
}

func (l *legend) draw() {
	sw, sh := termbox.Size()
	l.x = sw/2 + boardsize/2 + 1
	l.y = ((sh - boardsize + 1 + 1 + 1) / 2) + boardsize

	tbprint("H²", l.x+0, l.y, termbox.ColorWhite, termbox.ColorBlue)
	tbprint("H³", l.x+2, l.y, termbox.ColorWhite, termbox.ColorMagenta)
	tbprint("K²", l.x+4, l.y, termbox.ColorWhite, termbox.ColorGreen)
	tbprint("K³", l.x+6, l.y, termbox.ColorWhite, termbox.ColorBlack)
}

var score2rune = []rune{' ', '₁', '₂', '₃', '₄', '₅', '₆', '₇', '₈', '₉', '⏨'}

func getScoreRune(score int) (r rune) {
	return score2rune[score]
}

// kelimelik board
var (
	boardLetterMult = [boardsize][boardsize]int{
		{1, 1, 1, 1, 1, 2, 1, 1, 1, 2, 1, 1, 1, 1, 1},
		{1, 3, 1, 1, 1, 1, 2, 1, 2, 1, 1, 1, 1, 3, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 3, 1, 1, 1, 1, 1, 3, 1, 1, 1, 1},
		{2, 1, 1, 1, 1, 2, 1, 1, 1, 2, 1, 1, 1, 1, 2},
		{1, 2, 1, 1, 1, 1, 2, 1, 2, 1, 1, 1, 1, 2, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 2, 1, 1, 1, 1, 2, 1, 2, 1, 1, 1, 1, 2, 1},
		{2, 1, 1, 1, 1, 2, 1, 1, 1, 2, 1, 1, 1, 1, 2},
		{1, 1, 1, 1, 3, 1, 1, 1, 1, 1, 3, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 3, 1, 1, 1, 1, 2, 1, 2, 1, 1, 1, 1, 3, 1},
		{1, 1, 1, 1, 1, 2, 1, 1, 1, 2, 1, 1, 1, 1, 1},
	}
	boardWordMult = [boardsize][boardsize]int{
		{1, 1, 3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{3, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 3},
		{1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 2, 1, 1, 1, 1, 2, 1, 1, 1, 1, 2, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1},
		{3, 1, 1, 1, 1, 1, 1, 2, 1, 1, 1, 1, 1, 1, 3},
		{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		{1, 1, 3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 3, 1, 1},
	}
)
