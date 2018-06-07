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
	var (
		flagComma = flag.Bool("c", false, "")
	)
	flag.Parse()

	if flag.NArg() == 0 {
		s := bufio.NewScanner(os.Stdin)
		for s.Scan() {
			line := s.Text()
			fmt.Println(humanize(line, *flagComma))
		}
		if err := s.Err(); err != nil {
			log.Fatal(err)
		}
		return
	}

	for _, arg := range flag.Args() {
		fmt.Println(humanize(arg, *flagComma))
	}

}

func humanize(s string, comma bool) string {
	if comma {
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			log.Fatalf("could not parse %v: %v", s, err)
		}
		return humanizepkg.Comma(int64(i))
	}

	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		log.Fatalf("could not parse %v: %v", s, err)
	}
	return humanizepkg.Bytes(i)
}
