package main

import (
	"fmt"
	"time"

	"github.com/igungor/quackle"
)

const (
	lexicon  = "turkish"
	alphabet = "turkish"
	datadir  = "data"
)

var game quackle.Game
var dm quackle.DataManager

// hacky stuff.
// var datadir = fmt.Sprintf("%v/src/github.com/igungor/cmd/pfft/data", os.Getenv("GOPATH"))

func initGame() {
	// set up data manager
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
	board := quackle.NewBoardParameters()
	for y := 0; y < boardsize; y++ {
		for x := 0; x < boardsize; x++ {
			board.SetLetterMultiplier(x, y, quackle.QuackleBoardParametersLetterMultiplier(boardLetterMult[x][y]))
			board.SetWordMultiplier(x, y, quackle.QuackleBoardParametersWordMultiplier(boardWordMult[x][y]))
		}
	}
	dm.SetBoardParameters(board)

	// find lexicon
	dawg := quackle.LexiconParametersFindDictionaryFile(lexicon + ".dawg")
	gaddag := quackle.LexiconParametersFindDictionaryFile(lexicon + ".gaddag")
	dm.LexiconParameters().LoadDawg(dawg)
	dm.LexiconParameters().LoadGaddag(gaddag)
	dm.StrategyParameters().Initialize(lexicon)

	dm.SeedRandomNumbers(uint(time.Now().UnixNano()))

	// set up players and game
	game = quackle.NewGame()
	player1, _ := newComputerPlayer(dm, "comp1", 0)
	player2 := newHumanPlayer(dm, "iby", 1)
	players := quackle.NewPlayerList()
	players.Add(player1)
	players.Add(player2)
	game.SetPlayers(players)
	game.AssociateKnownComputerPlayers()
	game.AddPosition()
}

func getCurPos(game quackle.Game) quackle.GamePosition {
	return game.CurrentPosition().(quackle.GamePosition)
}

func newComputerPlayer(dm quackle.DataManager, name string, id int) (quackle.Player, error) {
	var player quackle.Player
	found := make([]bool, 1)
	player = dm.ComputerPlayers().PlayerForName("Speedy Player", found)
	if !found[0] {
		return player, fmt.Errorf("player could not be found")
	}
	comp := player.ComputerPlayer()

	player = quackle.NewPlayer(name)
	player.SetComputerPlayer(comp)
	player.SetType(int(quackle.PlayerComputerPlayerType))
	player.SetId(id)
	return player, nil
}

func newHumanPlayer(dm quackle.DataManager, name string, id int) quackle.Player {
	player := quackle.NewPlayer(name)
	player.SetType(int(quackle.PlayerHumanPlayerType))
	player.SetId(id)
	return player
}
