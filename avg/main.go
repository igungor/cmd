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

		// date parsing
		lastdate   time.Time
		datelayout string
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
		case date:
			if lastdate.IsZero() {
				lastdate, datelayout = layout(line)
				continue
			}

			t, err := time.Parse(datelayout, line)
			if err != nil {
				log.Fatalf("could not parse %q with layout %q", line, datelayout)
			}
			total += float64(t.Sub(lastdate).Nanoseconds())
			count++
			lastdate = t
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
	case date:
		fmt.Println((time.Duration(total) / time.Duration(count)).String())
	}
}

type typ int

const (
	unassigned typ = iota
	unknown
	duration
	integer
	float
	date
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
	_, layout := layout(s)
	if layout != "" {
		return date
	}
	return unknown
}

func layout(s string) (time.Time, string) {
	for _, layout := range layouts {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, layout
		}
	}
	return time.Time{}, ""
}

var layouts = []string{
	"2006-01-02 15:04:05.999",
	time.RFC3339,
	time.RFC3339Nano,
}
