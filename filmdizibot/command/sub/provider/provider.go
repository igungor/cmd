package provider

import (
	"context"
	"fmt"
	"io"
)

type QueryItem int

const (
	QueryByTitle QueryItem = iota
	QueryByHash
)

type Provider interface {
	// Provider name
	Name() string

	// Provider accepts requests with title/hash or what.
	QueryType() QueryItem

	// Query subtitles for the episode
	Query(context.Context, string) ([]*Subtitle, error)

	// Download subtitle
	Download(context.Context, *Subtitle) (io.ReadCloser, error)
}

type Subtitle struct {
	Title           string
	Season          int
	Episode         int
	Language        string
	Release         string
	Status          string
	HearingImpaired bool
	DownloadURL     string
	PageURL         string
}

func (s *Subtitle) String() string {
	var hi string
	if s.HearingImpaired {
		hi = "(HI)"
	}
	return fmt.Sprintf("S%02dE%02d | %v | %v | %v %v\n", s.Season, s.Episode, s.Title, s.Language, s.Release, hi)
}
