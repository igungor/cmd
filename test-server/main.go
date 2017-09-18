package main

import (
	"flag"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"strconv"
)

func main() {
	var (
		flagDumpRequest = flag.Bool("d", false, "dump full request to stdout")
		flagLogRequest  = flag.Bool("l", true, "dump request path to stdout")
		flagPort        = flag.Uint("p", 0, "port to listen to")
	)
	flag.Parse()
	ln, err := net.Listen("tcp", net.JoinHostPort("0.0.0.0", strconv.Itoa(int(*flagPort))))
	if err != nil {
		log.Fatalf("Could not listen to the port: %v", err)
	}
	log.Printf("Running on %v\n", ln.Addr())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if *flagLogRequest {
			log.Printf("%v\n", r.URL)
		}
		if *flagDumpRequest {
			b, _ := httputil.DumpRequest(r, true)
			log.Println(string(b))
		}
	})
	log.Fatal(http.Serve(ln, nil))
}
