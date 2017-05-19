package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver"
)

func main() {
	log.SetFlags(0)
	log.SetPrefix("sevil: ")

	var (
		flagHost = flag.String("host", "0.0.0.0", "host")
		flagPort = flag.String("port", "1987", "port")
	)
	flag.Parse()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			handleIndex(w, r)
			return
		case "POST":
			handlePost(w, r)
			return
		default:
			http.NotFound(w, r)
			return
		}
	})

	log.Fatal(http.ListenAndServe(net.JoinHostPort(*flagHost, *flagPort), nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(indexHtml))
}

func handlePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form: %v", err), http.StatusInternalServerError)
		return
	}

	yt := r.FormValue("urlInput")
	if yt == "" {
		http.Error(w, "no address given", http.StatusBadRequest)
		return
	}

	p, err := download(yt)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error downloading %q: %v", yt, err), http.StatusInternalServerError)
		return
	}
	defer os.Remove(p)

	_, fname := filepath.Split(p)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename='%v'", fname))
	http.ServeFile(w, r, p)
	return
}

func download(yt string) (string, error) {
	log.Printf("New download request for %q\n", yt)
	u, err := url.Parse(yt)
	if err != nil {
		return "", fmt.Errorf("Error parsing youtube url: %v", err)
	}

	playlist := u.Query().Get("list")
	var isPlaylist bool
	if playlist != "" {
		isPlaylist = true
	}

	tmpdir, err := ioutil.TempDir("/tmp/", "sevil-")
	if err != nil {
		return "", fmt.Errorf("Error creating temp directory: %v", err)
	}
	defer os.RemoveAll(tmpdir)

	var (
		bin  = "youtube-dl"
		args = []string{
			"--ignore-errors",
			"--extract-audio",
			"--audio-format",
			"mp3",
			"-o",
			tmpdir + "/" + "%(title)s.%(ext)s",
			u.String(),
		}
	)

	ctx := context.Background()
	cmd := exec.CommandContext(ctx, bin, args...)

	// cmd can fail since we passed the '--ignore-errors' parameter. even if
	// there is a download error, such as a private video in a playlist,
	// youtube-dl returns exit code 1. hence ignore the error.
	output, err := cmd.CombinedOutput()
	if err != nil {
		if strings.Contains(string(output), "is not a valid") {
			return "", fmt.Errorf("%q gecerli bir URL degil", yt)
		}
		log.Printf("youtube-dl failed for some reason: %v. Output: %v\n\nIgnoring the error...\n", err, string(output))
	}

	if isPlaylist {
		err = archiver.Zip.Make(tmpdir+".zip", []string{tmpdir})
		if err != nil {
			return "", fmt.Errorf("Error creating zip file: %v", err)
		}
		return tmpdir + ".zip", nil
	}

	files, err := ioutil.ReadDir(tmpdir)
	if err != nil {
		return "", fmt.Errorf("Error reading tmpdir: %v", err)
	}

	if len(files) == 0 {
		return "", fmt.Errorf("No file under %q", tmpdir)
	}

	fname := files[0].Name()
	fpath := filepath.Join(tmpdir, fname)

	os.Rename(fpath, "/tmp/"+fname)

	return filepath.Join("/tmp", fname), nil
}

const indexHtml = `
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="utf-8">
        <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
        <link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
        <style>
            .top-buffer {margin-top:60px;}
        </style>
    </head>

    <body>
        <form class="form-group" method="POST">
            <div class="row top-buffer">
                <div class="col-lg-offset-4 col-lg-4">
                    <div class="input-group">
                        <input id="urlInput" type="text" class="form-control" name="urlInput">
                        <span class="input-group-btn">
                            <button id="submitBtn" class="btn btn-primary" type="submit">indir</button>
                        </span>
                    </div>
                </div>
            </div>
        </form>
    </body>
</html>
`
