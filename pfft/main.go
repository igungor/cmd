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

	drawBoard(getCurPos(game).Board())
	drawLegend()
	termbox.Flush()

mainloop:
	for {
		if getCurPos(game).GameOver() {
			// TODO(ig): handle gameover
			time.Sleep(3 * time.Second)
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
			drawBoard(getCurPos(game).Board())
			drawLegend()
			termbox.Flush()
		case termbox.EventMouse:
			// TODO(ig): handle mouse clicks
		case termbox.EventError:
			return ev.Err
		}
		drawBoard(getCurPos(game).Board())
		drawLegend()
		termbox.Flush()
	}
	return nil
}

// drawBoard draws the board at the center of the grid.
func drawBoard(board quackle.Board) {
	sw, sh := termbox.Size()
	x := (sw - boardsize*2 + 2 + 1 + 1) / 2
	y := (sh - boardsize + 1 + 1 + 1) / 2
	w, h := boardsize, boardsize

	termbox.Clear(fgcolor, bgcolor)
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
				// letter scores
				if showScore {
					termbox.SetCell(x+col*2+1, y+row, getLetterScoreSubscript(score), fgcolor, bgcolor)
				}
				continue
			}

			// mark multipliers
			letterMult := dm.BoardParameters().LetterMultiplier(row, col)
			wordMult := dm.BoardParameters().WordMultiplier(row, col)
			switch {
			case letterMult == 2:
				// termbox.SetCell(x+col*2, y+row, '\'', termbox.ColorWhite, termbox.ColorBlue)
				tbprint("H²", x+col*2, y+row, termbox.ColorWhite, termbox.ColorBlue)
			case letterMult == 3:
				// termbox.SetCell(x+col*2, y+row, '"', termbox.ColorWhite, termbox.ColorMagenta)
				tbprint("H³", x+col*2, y+row, termbox.ColorWhite, termbox.ColorMagenta)
			case letterMult == 4: // unused for kelimelik
				// r = '^'
			case wordMult == 2:
				// termbox.SetCell(x+col*2, y+row, '-', termbox.ColorWhite, termbox.ColorGreen)
				tbprint("K²", x+col*2, y+row, termbox.ColorWhite, termbox.ColorGreen)
			case wordMult == 3:
				// termbox.SetCell(x+col*2, y+row, '=', termbox.ColorWhite, termbox.ColorBlack)
				tbprint("K³", x+col*2, y+row, termbox.ColorWhite, termbox.ColorBlack)
			case wordMult == 4: // unused for kelimelik
				// r = '~'
			default:
				termbox.SetCell(x+col*2, y+row, ' ', fgcolor, bgcolor)
			}
		}
	}
}

func getLetterScoreSubscript(score int) (r rune) {
	switch score {
	case 0:
		r = '₀'
	case 1:
		r = '₁'
	case 2:
		r = '₂'
	case 3:
		r = '₃'
	case 4:
		r = '₄'
	case 5:
		r = '₅'
	case 6:
		r = '₆'
	case 7:
		r = '₇'
	case 8:
		r = '₈'
	case 9:
		r = '₉'
	case 10:
		r = '⏨'
	}
	return r
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

// func move() {
// 	// play
// 	for {
// 		if getCurPos(game).GameOver() {
// 			fmt.Println("GAME OVER!")
// 			break
// 		}
// 		player := getCurPos(game).CurrentPlayer().(quackle.Player)

// 		// go on cold heartless machine!
// 		if quackle.QuacklePlayerPlayerType(player.Xtype()) == quackle.PlayerComputerPlayerType {
// 			move := game.HaveComputerPlay()
// 			fmt.Println("Rack: ", player.Rack().ToString())
// 			fmt.Println("Move: ", move.ToString())
// 			fmt.Printf("Board:\n %v\n", getCurPos(game).Board().ToString())
// 			continue
// 		}

// 		game.AdvanceToNoncomputerPlayer()
// 		fmt.Println("Rack: ", player.Rack().ToString())

// 		// read input
// 		var move quackle.Move
// 	MOVELOOP:
// 		for {
// 			r := bufio.NewReader(os.Stdin)
// 			input, _ := r.ReadString('\n')
// 			input = strings.TrimSuffix(input, "\n")
// 			fields := strings.Fields(input)
// 			switch len(fields) {
// 			case 1:
// 				// pass
// 				if fields[0] == "-" {
// 					move = quackle.MoveCreatePassMove()
// 					game.CommitMove(move)
// 					break MOVELOOP
// 				}
// 				fmt.Println("NEIN! gecerli biseyler yaz")
// 				continue MOVELOOP
// 			case 2:
// 				place, word := fields[0], fields[1]
// 				move = quackle.MoveCreatePlaceMove(place, dm.AlphabetParameters().Encode(word))
// 				if getCurPos(game).ValidateMove(move) == int(quackle.GamePositionValidMove) {
// 					game.CommitMove(move)
// 					break MOVELOOP
// 				}
// 				fmt.Println("NEIN! gecerli bir hamle degil")
// 				continue MOVELOOP
// 			default:
// 				fmt.Println("NEIN! gecerli biseyler yaz")
// 				continue MOVELOOP
// 			}
// 		}

// 		fmt.Println("Rack: ", player.Rack().ToString())
// 		fmt.Println("Move: ", move.ToString())
// 		fmt.Printf("Board:\n %v\n", getCurPos(game).Board().ToString())
// 	}
// }
