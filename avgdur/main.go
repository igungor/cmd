package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "avg calculates the average of the given input\n")
		fmt.Fprintf(os.Stderr, "  - input is accepted via stdin\n")
		fmt.Fprintf(os.Stderr, "  - type of the input is automatically detected (accepted inputs are Go time duration, integer and float)\n")
		fmt.Fprintf(os.Stderr, "\n")
		fmt.Fprintf(os.Stderr, `Usage: echo '15m10s\n3m5s' | avgdur`)
		fmt.Fprintf(os.Stderr, "\n")
	}
	flag.Parse()
	var (
		total float64
		count int64

		firsttype typ
	)
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}
		if firsttype == unassigned {
			firsttype = typeof(line)
		}
		switch firsttype {
		case duration:
			dur, _ := time.ParseDuration(line)
			total += float64(dur.Nanoseconds())
			count++
		case integer:
			i, _ := strconv.ParseInt(line, 10, 64)
			total += float64(i)
			count++
		case float:
			f, _ := strconv.ParseFloat(line, 64)
			total += f
			count++
		case unknown:
			log.Fatalf("unrecognized input %q", line)
		}
	}

	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	switch firsttype {
	case duration:
		fmt.Println((time.Duration(total) / time.Duration(count)).String())
	case integer:
		fmt.Println(int64(total) / count)
	case float:
		fmt.Println(total / float64(count))
	}
}

type typ int

const (
	unassigned typ = iota
	unknown
	duration
	integer
	float
)

func typeof(s string) typ {
	_, err := time.ParseDuration(s)
	if err == nil {
		return duration
	}
	_, err = strconv.ParseInt(s, 10, 64)
	if err == nil {
		return integer
	}
	_, err = strconv.ParseFloat(s, 64)
	if err == nil {
		return float
	}
	return unknown
}
