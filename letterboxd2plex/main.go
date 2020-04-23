package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	if err := realMain(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func realMain() error {
	var (
		flagPlexAddr  = flag.String("plex.addr", "", "Plex server address (host:port)")
		flagPlexToken = flag.String("plex.token", "", "Plex API token")
	)
	flag.Parse()

	if *flagPlexAddr == "" {
		return fmt.Errorf("Plex address must be provided")
	}

	if *flagPlexToken == "" {
		return fmt.Errorf("Plex API token must be provided")
	}

	if flag.NArg() != 1 {
		return fmt.Errorf("Letterboxd list URL must be provided")
	}

	plexc, err := New(*flagPlexAddr, *flagPlexToken, nil)
	if err != nil {
		return err
	}

	library, err := plexc.LibraryContent(moviesSection)
	if err != nil {
		return err
	}

	plexMovies := make(map[string]string)
	for _, metadata := range library.MediaContainer.Metadata {
		plexMovies[metadata.Title] = metadata.RatingKey
	}

	letterboxdList, err := FetchMovieList(flag.Arg(0))
	if err != nil {
		return err
	}

	var keys []string
	for _, movie := range letterboxdList.Movies {
		ratingKey, ok := plexMovies[movie]
		if !ok {
			log.Printf("Letterboxd movie %q is not found in Plex library", movie)
			continue
		}

		log.Printf("Found %q in Plex library", movie)
		keys = append(keys, ratingKey)
	}

	for _, key := range keys {
		if err := plexc.AddToCollection(letterboxdList.Title, moviesSection, key); err != nil {
			return err
		}
	}

	collection, err := plexc.CollectionByTitle(letterboxdList.Title)
	if err != nil {
		return err
	}

	return plexc.UpdateCollectionSummary(collection.RatingKey, letterboxdList.Summary)
}
