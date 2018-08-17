package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"sync"
	"time"
)

func main() {
	var (
		flagDumpRequest  = flag.Bool("d", false, "dump full request to stdout")
		flagLogRequest   = flag.Bool("l", false, "dump request path to stdout")
		flagCountRequest = flag.Bool("c", false, "count requests")
		flagPort         = flag.Uint("p", 0, "port to listen to")
		flagDelay        = flag.Duration("t", 0, "delay the response")
	)
	flag.Parse()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	ln, err := net.Listen("tcp", net.JoinHostPort("0.0.0.0", strconv.Itoa(int(*flagPort))))
	if err != nil {
		logger.Fatalf("Could not listen to the port: %v", err)
	}
	logger.Printf("Running on %v\n", ln.Addr())

	var (
		mu      sync.Mutex
		counter int64
	)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if flagDelay.Nanoseconds() > 0 {
			time.Sleep(*flagDelay)
		}
		if *flagLogRequest {
			logger.Printf("%v\n", r.URL)
		}
		if *flagDumpRequest {
			b, _ := httputil.DumpRequest(r, true)
			logger.Println(string(b))
		}

		if *flagCountRequest {
			mu.Lock()
			counter++
			mu.Unlock()
			logger.Printf("counter: %d\n", counter)
		}
	})
	logger.Fatal(http.Serve(ln, nil))
}
