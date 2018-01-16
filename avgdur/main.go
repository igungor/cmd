package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "avgdur calculates the average duration of the given input\n")
		fmt.Fprintf(os.Stderr, "  - input is accepted via stdin\n")
		fmt.Fprintf(os.Stderr, "  - duration format is the Go duration format\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, `Usage: echo '15m10s\n3m5s' | avgdur`)
		fmt.Fprintf(os.Stderr, "\n")
	}
	flag.Parse()
	var (
		total time.Duration
		count time.Duration
	)
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		dur, err := time.ParseDuration(line)
		if err != nil {
			log.Fatalf("could not parse duration %v: %v", line, err)
		}
		total += dur
		count++
	}
	if err := s.Err(); err != nil {
		log.Fatalf("scanner failed: %v", err)
	}

	fmt.Println((total / count).String())
}
