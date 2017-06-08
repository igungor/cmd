// torrent2magnet generates magnet link from the given torrent file.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/anacrolix/torrent/metainfo"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("torrent2magnet: ")
	flag.Parse()

	in := os.Stdin
	if flag.NArg() == 1 {
		f, err := os.Open(flag.Arg(0))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		in = f
	}

	mi, err := metainfo.Load(in)
	if err != nil {
		log.Fatal(err)
	}

	info, err := mi.UnmarshalInfo()
	if err != nil {
		log.Fatal(err)
	}

	magnet := mi.Magnet(info.Name, mi.HashInfoBytes())

	fmt.Println(magnet)
}
