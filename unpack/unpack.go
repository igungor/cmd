package main

// TODO(ig): handle the spammer archives as well (archives that have multiple
// files without a container directory)

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	flagOutput = flag.String("o", "", "output directory")
)

func main() {
	log.SetFlags(0)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: unpack [options] files\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		os.Exit(2)
	}

	flag.Parse()

	if flag.NArg() == 0 {
		flag.Usage()
	}

	f := flag.Args()[0]
	filename, err := filepath.Abs(f)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	_, err = os.Stat(filename)
	if err != nil && os.IsNotExist(err) {
		log.Printf("Given file '%s' does not exist.\n", filename)
		os.Exit(1)
	}

	for _, unpacker := range unpackers {
		if unpacker.matches(f) {
			if err := unpacker.unpack(filename, *flagOutput); err != nil {
				log.Fatal(err)
			}
			return
		}
	}

	log.Fatalf("file '%s' is not recognized by this tool.\n", f)
}

var unpackers = []unpacker{
	// 'unrar vb test.rar': list file contents without all the crap
	&spec{
		id:         "rar",
		patterns:   []string{"*.rar"},
		executable: "unrar",
		args:       []string{"x"},
	},
	&spec{
		id:         "zip",
		patterns:   []string{"*.zip", "*.jar"},
		executable: "unzip",
		args:       []string{},
	},
	&spec{
		id:         "tar",
		patterns:   []string{"*.tar"},
		executable: "tar",
		args:       []string{"xvf"},
	},
	&spec{
		id:         "targz",
		patterns:   []string{"*.tar.gz", "*.tgz"},
		executable: "tar",
		args:       []string{"xvf"},
	},
	&spec{
		id:         "tarbz2",
		patterns:   []string{"*.tar.bz2"},
		executable: "tar",
		args:       []string{"xvf"},
	},
	&spec{
		id:         "gzip",
		patterns:   []string{"*.gz"},
		executable: "gunzip",
		args:       []string{"-c"},
	},
}

type unpacker interface {
	matches(string) bool
	unpack(string, string) error
}

// spec implements unpacker
type spec struct {
	id         string
	patterns   []string
	executable string
	args       []string
}

func (s spec) matches(filename string) bool {
	for _, pattern := range s.patterns {
		matched, err := filepath.Match(pattern, filename)
		if err != nil {
			log.Printf("Error while pattern matching. Err: %v\n", err)
			return false
		}
		if matched {
			return true
		}
	}
	return false
}

func (s spec) unpack(filename, output string) error {
	// os.Rename won't move from /tmp to /home if they are in separate
	// partitions. Create the tmpdir in current directory
	tmpdir, err := ioutil.TempDir(".", "unpack-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpdir)

	args := append(s.args, filename)
	cmd := exec.Command(s.executable, args...)
	cmd.Stdout = os.Stdout
	cmd.Dir = tmpdir

	if err := cmd.Run(); err != nil {
		return err
	}

	if output == "" {
		output = strings.TrimSuffix(filename, filepath.Ext(filename))
	}
	if err := os.MkdirAll(output, 0755); err != nil {
		return err
	}

	if err := os.Rename(tmpdir, output); err != nil {
		defer os.RemoveAll(tmpdir)
		return err
	}

	return nil
}

func isScattered(dir string) bool {
	files, _ := ioutil.ReadDir(dir)
	return len(files) > 1
}
