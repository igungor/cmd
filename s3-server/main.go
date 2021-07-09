package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"net/http/httptest"
	"os"
	"os/signal"
	"syscall"

	"github.com/johannesboyne/gofakes3"
	"github.com/johannesboyne/gofakes3/backend/s3bolt"
)

func main() {
	var (
		flagBoltFile = flag.String("db", "", "Path to boltdb file")
	)
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, os.Kill)
	defer cancel()

	logger := log.New(os.Stdout, "", 0)

	var dbpath string
	if *flagBoltFile != "" {
		dbpath = *flagBoltFile
	} else {
		f, err := ioutil.TempFile("", "")
		if err != nil {
			logger.Fatalf("tmpfile: %v", err)
		}
		f.Close()

		dbpath = f.Name()
		logger.Printf("DB file: %v", f.Name())
	}

	backend, err := s3bolt.NewFile(dbpath)
	if err != nil {
		logger.Fatalf("bolt: %v", err)
	}
	faker := gofakes3.New(backend)

	ts := httptest.NewServer(faker.Server())
	defer ts.Close()

	logger.Println(ts.URL)

	<-ctx.Done()
}
