package main

// BUG(ig): drawBoard can't display some letters for some reason. 'NEYCE' appears as 'N YCE'

import (
	"log"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/igungor/quackle"
	"github.com/nsf/termbox-go"
)

const (
	fgcolor = termbox.ColorDefault
	bgcolor = termbox.ColorDefault
)

var showScore bool

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
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	initGame()
	drawGame()

mainloop:
	for {
		if getCurPos(game).GameOver() {
			drawGameover()
			break mainloop
		}
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				break mainloop
			case termbox.KeyEnter:
				game.HaveComputerPlay()
			case termbox.KeyCtrlS:
				showScore = !showScore
			default:
			}
		case termbox.EventResize:
			drawGame()
		case termbox.EventMouse:
			// TODO(ig): handle mouse clicks
		case termbox.EventError:
			return ev.Err
		}
		drawGame()
	}
	return nil
}

func drawGame() {
	termbox.Clear(fgcolor, bgcolor)
	drawBoard(getCurPos(game).Board())
	drawLegend()
	drawRack1()
	drawRack2()
	termbox.Flush()
}

// drawBoard draws the board at the center of the grid.
func drawBoard(board quackle.Board) {
	sw, sh := termbox.Size()
	x := (sw - boardsize*2 + 2 + 1) / 2
	y := (sh - boardsize + 1 + 1) / 2
	w, h := boardsize, boardsize

	// columns on the top
	for dx := 0; dx < w; dx++ {
		termbox.SetCell(x+1+dx*2, y-2, rune('A'+dx), fgcolor, bgcolor)
		termbox.SetCell(x+1+dx*2+1, y-2, ' ', fgcolor, bgcolor)
	}

	// rows on the left
	for dy := 0; dy < h; dy++ {
		if dy+1 < 10 {
			tbprint(strconv.Itoa(dy+1), x-2, y+dy, fgcolor, bgcolor)
		} else {
			tbprint(strconv.Itoa(dy+1), x-3, y+dy, fgcolor, bgcolor)
		}
	}

	// draw borders
	drawRect(x, y, w, h)

	// multipliers and letters
	for row := 0; row < h; row++ {
		for col := 0; col < w; col++ {
			// mark letters
			bl := board.Letter(row, col)
			if dm.AlphabetParameters().IsPlainLetter(bl) {
				letter := dm.AlphabetParameters().UserVisible(bl)
				score := dm.AlphabetParameters().Score(bl)
				r, _ := utf8.DecodeRuneInString(letter)
				termbox.SetCell(x+col*2, y+row, r, fgcolor, bgcolor)
				if showScore {
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

func drawLegend() {
	sw, sh := termbox.Size()
	x := sw/2 + boardsize/2 + 2
	y := ((sh - boardsize + 1 + 1 + 1) / 2) + 15 + 1

	tbprint("H²", x+0, y, termbox.ColorWhite, termbox.ColorBlue)
	tbprint("H³", x+2, y, termbox.ColorWhite, termbox.ColorMagenta)
	tbprint("K²", x+4, y, termbox.ColorWhite, termbox.ColorGreen)
	tbprint("K³", x+6, y, termbox.ColorWhite, termbox.ColorBlack)
}

func drawRack1() {
}

func drawRack2() {
}

func drawGameover() {
	termbox.Clear(fgcolor, bgcolor)
	sw, sh := termbox.Size()
	tbprint("GAME OVER", sw/2-4, sh/2, fgcolor, bgcolor)
	termbox.Flush()
	time.Sleep(1 * time.Second)
}

func getScoreRune(score int) (r rune) {
	return score2rune[score]
}

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
