package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func main() {
	var (
		flagTimeout = flag.Duration("t", 0, "timeout")
		flagSignal  = flag.String("s", "sigint", "the signal to be sent")
	)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "akilyu sends the specified signal to the given process after a certain amount of time\n\n")
		fmt.Fprintf(os.Stderr, "Usage: akilyu PID\n")
		flag.PrintDefaults()
	}
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

	if err := p.Signal(signal(*flagSignal)); err != nil {
		log.Fatalf("Could not send SIGINT to process %v: %v", p.Pid, err)
	}
}

func signal(sig string) os.Signal {
	sig = strings.ToLower(sig)
	if strings.HasPrefix(sig, "sig") {
		sig = strings.TrimPrefix(sig, "sig")
	}

	switch sig {
	case "int":
		return os.Interrupt
	case "hup":
		return syscall.SIGHUP
	case "kill":
		return syscall.SIGKILL
	case "usr1":
		return syscall.SIGUSR1
	case "usr2":
		return syscall.SIGUSR2
	}
	return os.Interrupt
}
