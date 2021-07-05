package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"sync"
	"time"
)

func main() {
	var (
		flagDumpRequest  = flag.Bool("d", false, "dump full request to stdout")
		flagLogRequest   = flag.Bool("l", false, "dump request path to stdout")
		flagCountRequest = flag.Bool("c", false, "count requests")
		flagAddr         = flag.String("p", ":8080", "address to listen to")
		flagDelay        = flag.Duration("t", 0, "delay the response")
		flagResponse     = flag.String("w", "", "response")
	)
	flag.Parse()

	logger := log.New(os.Stdout, "", log.LstdFlags)

	ln, err := net.Listen("tcp", *flagAddr)
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

		w.Write([]byte(*flagResponse))
	})
	logger.Fatal(http.Serve(ln, nil))
}
