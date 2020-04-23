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

	if flag.NArg() > 0 {
		fmt.Println(unixtime(flag.Arg(0)))
		return
	}

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		fmt.Println(unixtime(s.Text()))
	}
}

func unixtime(s string) time.Time {
	t, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	return time.Unix(t, 0)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: unixtime [unixtime]\n")
	flag.PrintDefaults()
}
