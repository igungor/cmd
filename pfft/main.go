package main

import (
	"log"
	"os"

	termbox "github.com/nsf/termbox-go"
)

var f *os.File

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

	var err error
	f, err = os.OpenFile("foo", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	game := newGame()
	game.draw()
	game.loop()
	return nil
}
