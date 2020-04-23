package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/rs/zerolog"
)

var logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).With().Timestamp().Logger()

func main() {
	if err := realMain(); err != nil {
		logger.Error().Msg(err.Error())
		os.Exit(1)
	}
}

func realMain() error {
	var (
		flagPlexAddr  = flag.String("plex.addr", "", "Plex server address (host:port)")
		flagPlexToken = flag.String("plex.token", "", "Plex API token")
		flagLogLevel  = flag.String("log", "info", "Log level (debug, info, error)")
	)
	flag.Parse()

	logLevel, _ := zerolog.ParseLevel(*flagLogLevel)
	logger = logger.Level(logLevel)

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
			logger.Debug().Str("movie", movie).Msg("Letterboxd movie not found in Plex library")
			// log.Printf("Letterboxd movie %q is not found in Plex library", movie)
			continue
		}

		logger.Debug().Str("movie", movie).Msg("found in Plex library")
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
