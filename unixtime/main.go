package main

import (
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
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}

	arg := flag.Arg(0)
	unixtime, err := strconv.ParseInt(arg, 10, 64)
	if err != nil {
		log.Fatal(err)
	}
	t := time.Unix(unixtime, 0)

	fmt.Println(t)
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: unixtime [unixtime]\n")
	flag.PrintDefaults()
}
