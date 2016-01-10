package main

import (
	"strconv"
	"unicode/utf8"

	"github.com/igungor/quackle"
	termbox "github.com/nsf/termbox-go"
)

const boardsize = 15

type board struct {
	qb        quackle.Board
	w, h      int
	showScore bool
}

func (b *board) draw(x, y int) {
	// columns on the top
	for dx := 0; dx < b.w; dx++ {
		termbox.SetCell(x+dx*2, y-2, rune('A'+dx), fgcolor, bgcolor)
		termbox.SetCell(x+dx*2+1, y-2, ' ', fgcolor, bgcolor)
	}

	// rows on the left
	for dy := 0; dy < b.h; dy++ {
		if dy+1 < 10 {
			tbprint(strconv.Itoa(dy+1), x-2, y+dy, fgcolor, bgcolor)
		} else {
			tbprint(strconv.Itoa(dy+1), x-3, y+dy, fgcolor, bgcolor)
		}
	}

	// borders
	drawRect(x, y, b.w*2, b.h)

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
				if b.showScore {
					tbprint("H²", x+col*2, y+row, termbox.ColorWhite, termbox.ColorBlue)
				} else {
					tbprint("★", x+col*2, y+row, termbox.ColorBlue, bgcolor)
				}
			case letterMult == 3:
				if b.showScore {
					tbprint("H³", x+col*2, y+row, termbox.ColorWhite, termbox.ColorMagenta)
				} else {
					tbprint("★", x+col*2, y+row, termbox.ColorMagenta, bgcolor)
				}
			case wordMult == 2:
				if b.showScore {
					tbprint("K²", x+col*2, y+row, termbox.ColorWhite, termbox.ColorGreen)
				} else {
					tbprint("★", x+col*2, y+row, termbox.ColorGreen, bgcolor)
				}
			case wordMult == 3:
				if b.showScore {
					tbprint("K³", x+col*2, y+row, termbox.ColorWhite, termbox.ColorBlack)
				} else {
					tbprint("★", x+col*2, y+row, termbox.ColorBlack, bgcolor)
				}
			default:
				termbox.SetCell(x+col*2, y+row, ' ', fgcolor, bgcolor)
			}
		}
	}
}

type legend struct {
}

func (l *legend) draw(x, y int) {
	tbprint("H²", x+0, y, termbox.ColorWhite, termbox.ColorBlue)
	tbprint("H³", x+2, y, termbox.ColorWhite, termbox.ColorMagenta)
	tbprint("K²", x+4, y, termbox.ColorWhite, termbox.ColorGreen)
	tbprint("K³", x+6, y, termbox.ColorWhite, termbox.ColorBlack)
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
