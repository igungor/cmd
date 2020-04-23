package main

import (
	"flag"
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
	flagConfig := flag.String("c", "config.yaml", "Configuration file path")
	flag.Parse()

	cfg, err := decodeConfig(*flagConfig)
	if err != nil {
		return err
	}

	plexc, err := New(cfg.Plex.Addr, cfg.Plex.Token, nil)
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

	logLevel, _ := zerolog.ParseLevel(cfg.LogLevel)
	logger = logger.Level(logLevel)

	for _, link := range cfg.Lists {
		letterboxdList, err := FetchMovieList(link)
		if err != nil {
			return err
		}

		var keys []string
		for _, movie := range letterboxdList.Movies {
			ratingKey, ok := plexMovies[movie]
			if !ok {
				logger.Debug().Str("movie", movie).Msg("Letterboxd movie not found in Plex library")
				continue
			}

			logger.Debug().Str("movie", movie).Msg("found in Plex library")
			keys = append(keys, ratingKey)
		}

		listTitle := letterboxdList.Title
		for _, key := range keys {
			if err := plexc.AddToCollection(listTitle, moviesSection, key); err != nil {
				return err
			}
		}

		collection, err := plexc.CollectionByTitle(listTitle)
		if err != nil {
			return err
		}

		return plexc.UpdateCollectionSummary(collection.RatingKey, letterboxdList.Summary)
	}

	return nil
}
