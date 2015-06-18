// guesslang is a command line tool that tries to detect the language of a
// given characters.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/igungor/chardet"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("guesslang: ")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s FILE1 FILE2 ...\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()

	c := &detector{chardet.NewTextDetector()}

	if flag.NArg() > 0 {
		for _, f := range flag.Args() {
			content, err := ioutil.ReadFile(flag.Arg(0))
			if err != nil {
				log.Fatal(err)
			}

			res := c.Detect(content)
			fmt.Printf("%s: %s [%s]\n", f, res.Charset, res.Language)
		}
		return
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		res := c.Detect(scanner.Bytes())
		fmt.Printf("%s [%s]\n", res.Charset, res.Language)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

type detector struct {
	*chardet.Detector
}

func (c *detector) Detect(content []byte) *chardet.Result {
	res, err := c.DetectBest(content)
	if err != nil {
		res.Charset = "not-detected"
		res.Language = "not-detected"
	}
	return res
}
