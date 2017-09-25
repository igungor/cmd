// akilyu sends a SIGINT to given process after a given -t timeout.
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
	var (
		flagTimeout = flag.Duration("t", 0, "timeout")
	)
	flag.Usage = usage
	flag.Parse()
	log.SetFlags(0)

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	arg := flag.Arg(0)
	pid, err := strconv.Atoi(arg)
	if err != nil {
		log.Fatalf("Invalid PID: %v", err)
	}

	p, err := os.FindProcess(pid)
	if err != nil {
		log.Fatalf("Could not find the process (%v): %v", pid, err)
	}

	<-time.After(*flagTimeout)

	if err := p.Signal(os.Interrupt); err != nil {
		log.Fatalf("Could not send SIGINT to process %v: %v", p.Pid, err)
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: akilyu PID\n")
	flag.PrintDefaults()
}
