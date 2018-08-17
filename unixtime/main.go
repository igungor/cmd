package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	flag.Parse()
	flag.Usage = usage

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		line := s.Text()
		unixtime, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			log.Fatal(err)
		}
		t := time.Unix(unixtime, 0)
		fmt.Println(t)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: unixtime [unixtime]\n")
	flag.PrintDefaults()
}
