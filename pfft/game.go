package main

import (
	"time"

	"github.com/igungor/quackle"
	"github.com/nsf/termbox-go"
)

const (
	lexicon  = "turkish"
	alphabet = "turkish"
	datadir  = "data"
)

// hacky stuff.
// var datadir = fmt.Sprintf("%v/src/github.com/igungor/cmd/pfft/data", os.Getenv("GOPATH"))

const (
	fgcolor = termbox.ColorDefault
	bgcolor = termbox.ColorDefault
)

var dm quackle.DataManager

type game struct {
	qg     quackle.Game
	board  board
	rack1  rack
	rack2  rack
	legend legend

	isOver     bool
	showLegend bool
}

func (g *game) draw() {
	// renew quackle board
	g.board.qb = g.pos().Board()

	termbox.Clear(fgcolor, bgcolor)
	g.board.draw()
	g.legend.draw()
	g.rack1.draw()
	g.rack2.draw()
	termbox.Flush()
}

func (g *game) loop() {
mainloop:
	for {
		if g.pos().GameOver() {
			g.isOver = true
			g.over()
			break mainloop
		}
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc, termbox.KeyCtrlC:
				break mainloop
			case termbox.KeyEnter:
				g.qg.HaveComputerPlay()
			case termbox.KeyCtrlS:
				g.board.showScore = !g.board.showScore
			default:
			}
		case termbox.EventResize:
			g.draw()
		case termbox.EventMouse:
			// TODO(ig): handle mouse clicks
		case termbox.EventError:
			panic(ev.Err)
		}
		g.draw()
	}
}

func (g *game) over() {
	termbox.Clear(fgcolor, bgcolor)
	sw, sh := termbox.Size()
	tbprint("GAME OVER", sw/2-4, sh/2, fgcolor, bgcolor)
	termbox.Flush()
	time.Sleep(1 * time.Second)
}

func (g *game) pos() quackle.GamePosition {
	return g.qg.CurrentPosition().(quackle.GamePosition)
}

func newGame() *game {
	dm = quackle.NewDataManager()
	dm.SetComputerPlayers(quackle.ComputerPlayerCollectionFullCollection().SwigGetPlayerList())
	dm.SetBackupLexicon(lexicon)
	dm.SetAppDataDirectory(datadir)

	// set up alphabet
	abc := quackle.AlphabetParametersFindAlphabetFile(alphabet)
	qabc := quackle.UtilStdStringToQString(abc)
	flexAbc := quackle.NewFlexibleAlphabetParameters()
	flexAbc.Load(qabc)
	dm.SetAlphabetParameters(flexAbc)

	// set up board
	bp := quackle.NewBoardParameters()
	for y := 0; y < boardsize; y++ {
		for x := 0; x < boardsize; x++ {
			bp.SetLetterMultiplier(x, y, quackle.QuackleBoardParametersLetterMultiplier(boardLetterMult[x][y]))
			bp.SetWordMultiplier(x, y, quackle.QuackleBoardParametersWordMultiplier(boardWordMult[x][y]))
		}
	}
	dm.SetBoardParameters(bp)

	// find lexicon
	dawg := quackle.LexiconParametersFindDictionaryFile(lexicon + ".dawg")
	gaddag := quackle.LexiconParametersFindDictionaryFile(lexicon + ".gaddag")
	dm.LexiconParameters().LoadDawg(dawg)
	dm.LexiconParameters().LoadGaddag(gaddag)
	dm.StrategyParameters().Initialize(lexicon)

	dm.SeedRandomNumbers(uint(time.Now().UnixNano()))

	newCompPlayer := func(name string, id int) quackle.Player {
		found := make([]bool, 1)
		player := dm.ComputerPlayers().PlayerForName("Speedy Player", found)
		if !found[0] {
			panic("player could not be found")
		}
		comp := player.ComputerPlayer()

		player = quackle.NewPlayer(name)
		player.SetComputerPlayer(comp)
		player.SetType(int(quackle.PlayerComputerPlayerType))
		player.SetId(id)
		return player
	}

	newHumanPlayer := func(name string, id int) quackle.Player {
		player := quackle.NewPlayer(name)
		player.SetType(int(quackle.PlayerHumanPlayerType))
		player.SetId(id)
		return player
	}

	// set up players and game
	g := quackle.NewGame()
	player1 := newCompPlayer("comp1", 0)
	player2 := newHumanPlayer("iby", 1)
	players := quackle.NewPlayerList()
	players.Add(player1)
	players.Add(player2)
	g.SetPlayers(players)
	g.AssociateKnownComputerPlayers()
	g.AddPosition()

	b := board{
		qb: g.CurrentPosition().(quackle.GamePosition).Board(),
		w:  boardsize,
		h:  boardsize,
	}

	return &game{
		qg:    g,
		board: b,
	}
}
