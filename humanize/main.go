package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	humanizepkg "github.com/dustin/go-humanize"
)

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			line := s.Text()
			fmt.Println(humanize(line))
		}
		if err := s.Err(); err != nil {
			log.Fatal(err)
		}
		return
	}

	for _, arg := range flag.Args() {
		fmt.Println(humanize(arg))
	}

}

func humanize(s string) string {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		log.Fatalf("could not parse %v: %v", s, err)
	}
	return humanizepkg.Bytes(i)
}
