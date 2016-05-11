package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	log.SetPrefix("ffp: ")
	log.SetFlags(0)

	flag.Usage = usage
	flag.Parse()
	if flag.NArg() != 1 {
		log.Fatal("missing argument")
	}

	cmd := "ffprobe"
	args := []string{"-print_format", "json", "-show_format", "-show_streams"}
	args = append(args, flag.Arg(0))

	out, err := exec.Command(cmd, args...).Output()
	if err != nil {
		log.Fatal(err)
	}

	var ff FFProbeOutput
	err = json.NewDecoder(bytes.NewReader(out)).Decode(&ff)
	if err != nil {
		log.Fatalf("decoding json failed: %v", err)
	}
	fmt.Println(ff)
}

type FFProbeOutput struct {
	Format struct {
		BitRate    string `json:"bit_rate"`
		Duration   string `json:"duration"`
		Filename   string `json:"filename"`
		FormatName string `json:"format_name"`
	} `json:"format"`
	Streams []struct {
		CodecName          string `json:"codec_name"`
		CodecType          string `json:"codec_type"`
		DisplayAspectRatio string `json:"display_aspect_ratio"`
		Height             int    `json:"height"`
		Width              int    `json:"width"`
		Tags               struct {
			Language string `json:"language"`
			Title    string `json:"title"`
		} `json:"tags"`
	} `json:"streams"`
}

func (f FFProbeOutput) String() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("format: %q\n", f.Format.FormatName))
	buf.WriteString(fmt.Sprintf("duration: %q\n", f.Format.Duration))

	for i, s := range f.Streams {
		buf.WriteString(fmt.Sprintf("Stream '%v'\n", i))
		buf.WriteString(fmt.Sprintf("\tcodec-type: %q\n", s.CodecType))
		buf.WriteString(fmt.Sprintf("\tcodec-name: %q\n", s.CodecName))
		if s.CodecType == "audio" || s.CodecType == "subtitle" {
			if s.Tags.Title != "" && s.Tags.Language != "" {
				buf.WriteString(fmt.Sprintf("\tlanguage: %q | %q\n", s.Tags.Language, s.Tags.Title))
			}
		}
	}

	return buf.String()
}

func usage() {
	fmt.Fprintf(os.Stderr, "ffp [options]\n")
	fmt.Fprintf(os.Stderr, "options:\n")
	flag.PrintDefaults()
}
