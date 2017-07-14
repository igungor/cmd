package main

import (
	"bufio"
	"encoding/base64"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var (
	ErrUnrecognizedType = errors.New("unrecognized type")
)

func main() {
	// flags
	var (
		flagEncoding = flag.String("enc", "std", "base64 encoding, either 'std' or 'url'")
	)
	flag.Parse()

	var encoding *base64.Encoding
	switch *flagEncoding {
	case "url":
		encoding = base64.URLEncoding
	case "std":
		encoding = base64.StdEncoding
	default:
		flag.Usage()
		os.Exit(2)
	}

	if flag.NArg() == 2 {
		arg0 := flag.Arg(0)
		arg1 := flag.Arg(1)
		b, err := scan(encoding, arg0, arg1)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(b)
		return
	}

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		f := strings.Fields(s.Text())
		if len(f) != 2 {
			fmt.Println("Input should be in the form of 'dec/enc src'")
			continue
		}
		b, err := scan(encoding, f[0], f[1])
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		fmt.Println(b)
	}

	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
}

func scan(encoding *base64.Encoding, typ, src string) (string, error) {
	switch typ {
	case "d", "dec", "decode":
		b, err := encoding.DecodeString(src)
		if err != nil {
			return "", err
		}
		return string(b), nil
	case "e", "enc", "encode":
		b := encoding.EncodeToString([]byte(src))
		return b, nil
	}
	return "", ErrUnrecognizedType
}
