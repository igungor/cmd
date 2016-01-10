package main

// BUG(ig): drawBoard can't display some letters for some reason. 'NEYCE' appears as 'N YCE'

import (
	"log"

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

	game := newGame()
	game.draw()
	game.loop()
	return nil
}
