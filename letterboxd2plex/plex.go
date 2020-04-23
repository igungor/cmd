package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const defaultUserAgent = "github.com/igungor/plex - Plex Go bindings"

const (
	moviesSection     = "1"
	collectionSection = "18"
)

var ErrUnauthorized = fmt.Errorf("invalid grant")

type Client struct {
	client    *http.Client
	baseURL   *url.URL
	token     string
	userAgent string
}

func New(plexAddr, token string, httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	if plexAddr == "" {
		return nil, fmt.Errorf("plex address must be provided")
	}

	if token == "" {
		return nil, fmt.Errorf("token must be provided")
	}

	baseURL, _ := url.Parse(plexAddr)
	c := &Client{
		client:    httpClient,
		baseURL:   baseURL,
		token:     token,
		userAgent: defaultUserAgent,
	}

	return c, nil
}

func (c *Client) Search(title string) (SearchResult, error) {
	title = url.QueryEscape(title)
	req := c.newRequest("GET", "/search?query="+title, nil)

	var result SearchResult
	_, err := c.do(req, &result)
	return result, err
}

func (c *Client) Libraries() (LibrariesResult, error) {
	req := c.newRequest("GET", "/library/sections", nil)

	var result LibrariesResult
	_, err := c.do(req, &result)
	return result, err
}

func (c *Client) LibraryContent(section string) (LibraryContentResult, error) {
	url := fmt.Sprintf("/library/sections/%v/all", section)
	req := c.newRequest("GET", url, nil)

	var result LibraryContentResult
	_, err := c.do(req, &result)
	return result, err
}

func (c *Client) AddToCollection(collection, section, key string) error {
	params := make(url.Values)
	params.Add("type", moviesSection)
	params.Add("id", key)
	params.Add("collection[0].tag.tag", collection)
	params.Add("collection.locked", "1")

	url := fmt.Sprintf("/library/sections/%v/all?%v", section, params.Encode())
	req := c.newRequest("PUT", url, nil)

	resp, err := c.do(req, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) UpdateCollectionSummary(collectionKey, summary string) error {
	params := make(url.Values)
	params.Add("type", collectionSection)
	params.Add("id", collectionKey)
	params.Add("includeExternalMedia", "1")
	params.Add("summary.value", summary)

	url := fmt.Sprintf("/library/sections/%v/all?%v", moviesSection, params.Encode())
	req := c.newRequest("PUT", url, nil)

	resp, err := c.do(req, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (c *Client) Collections() (CollectionsResult, error) {
	params := make(url.Values)
	params.Add("includeCollections", "1")
	params.Add("includeAdvanced", "1")
	params.Add("includeMeta", "1")

	url := fmt.Sprintf("/library/sections/%v/collections?%v", moviesSection, params.Encode())
	req := c.newRequest("GET", url, nil)

	var result CollectionsResult
	_, err := c.do(req, &result)
	return result, err
}

func (c *Client) CollectionByTitle(title string) (CollectionByTitleResult, error) {
	all, err := c.Collections()
	if err != nil {
		return CollectionByTitleResult{}, err
	}

	for _, coll := range all.MediaContainer.Metadata {
		if coll.Title == title {
			return coll, nil
		}
	}

	return CollectionByTitleResult{}, fmt.Errorf("no collection found for given title")
}

func (c *Client) newRequest(method, relURL string, body io.Reader) *http.Request {
	rel, _ := url.Parse(relURL)
	u := c.baseURL.ResolveReference(rel)

	req, _ := http.NewRequest(method, u.String(), body)
	req.Header.Add("X-Plex-Token", c.token)

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", c.userAgent)

	return req
}

func (c *Client) do(r *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(r)
	if err != nil {
		return nil, err
	}

	err = checkResponse(resp)
	if err != nil {
		// close the body at all times if there is an http error
		resp.Body.Close()
		return resp, err
	}

	if v == nil {
		return resp, nil
	}

	// close the body for all cases from here
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(v)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

type ErrorResponse struct {
	Response *http.Response `json:"-"`

	Message string `json:"error_message"`
	Type    string `json:"error_type"`
}

func (e *ErrorResponse) Error() string {
	return fmt.Sprintf(
		"Type: %v Message: %q. Original error: %v %v: %v",
		e.Type,
		e.Message,
		e.Response.Request.Method,
		e.Response.Request.URL,
		e.Response.Status,
	)
}

func checkResponse(r *http.Response) error {
	statusCode := r.StatusCode
	if statusCode >= 200 && statusCode <= 299 {
		return nil
	}

	switch statusCode {
	case http.StatusUnauthorized:
		return ErrUnauthorized
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		err = json.Unmarshal(data, errorResponse)
		if err != nil {
			// unexpected error
			return fmt.Errorf("unexpected HTTP status: %v. details: %v", statusCode, string(data[:250]))
		}
	}
	return errorResponse
}

type LibraryContentResult struct {
	MediaContainer struct {
		Metadata []struct {
			Genre []struct {
				Tag string `json:"tag"`
			} `json:"Genre"`
			Role []struct {
				Tag string `json:"tag"`
			} `json:"Role"`
			AddedAt               int64   `json:"addedAt"`
			Art                   string  `json:"art"`
			Banner                string  `json:"banner"`
			ChildCount            int64   `json:"childCount"`
			ContentRating         string  `json:"contentRating"`
			Duration              int64   `json:"duration"`
			GUID                  string  `json:"guid"`
			Index                 int64   `json:"index"`
			Key                   string  `json:"key"`
			LastViewedAt          int64   `json:"lastViewedAt"`
			LeafCount             int64   `json:"leafCount"`
			OriginallyAvailableAt string  `json:"originallyAvailableAt"`
			Rating                float64 `json:"rating"`
			RatingKey             string  `json:"ratingKey"`
			Studio                string  `json:"studio"`
			Summary               string  `json:"summary"`
			Theme                 string  `json:"theme"`
			Thumb                 string  `json:"thumb"`
			Title                 string  `json:"title"`
			Type                  string  `json:"type"`
			UpdatedAt             int64   `json:"updatedAt"`
			ViewCount             int64   `json:"viewCount"`
			ViewedLeafCount       int64   `json:"viewedLeafCount"`
			Year                  int64   `json:"year"`
		} `json:"Metadata"`
		AllowSync           bool   `json:"allowSync"`
		Art                 string `json:"art"`
		Identifier          string `json:"identifier"`
		LibrarySectionID    int64  `json:"librarySectionID"`
		LibrarySectionTitle string `json:"librarySectionTitle"`
		LibrarySectionUUID  string `json:"librarySectionUUID"`
		MediaTagPrefix      string `json:"mediaTagPrefix"`
		MediaTagVersion     int64  `json:"mediaTagVersion"`
		Nocache             bool   `json:"nocache"`
		Size                int64  `json:"size"`
		Thumb               string `json:"thumb"`
		Title1              string `json:"title1"`
		Title2              string `json:"title2"`
		ViewGroup           string `json:"viewGroup"`
		ViewMode            int64  `json:"viewMode"`
	} `json:"MediaContainer"`
}

type LibrariesResult struct {
	MediaContainer struct {
		Directory []struct {
			Location []struct {
				ID   int64  `json:"id"`
				Path string `json:"path"`
			} `json:"Location"`
			Agent            string `json:"agent"`
			AllowSync        bool   `json:"allowSync"`
			Art              string `json:"art"`
			Composite        string `json:"composite"`
			Content          bool   `json:"content"`
			ContentChangedAt int64  `json:"contentChangedAt"`
			CreatedAt        int64  `json:"createdAt"`
			Directory        bool   `json:"directory"`
			Filters          bool   `json:"filters"`
			Key              string `json:"key"`
			Language         string `json:"language"`
			Refreshing       bool   `json:"refreshing"`
			ScannedAt        int64  `json:"scannedAt"`
			Scanner          string `json:"scanner"`
			Thumb            string `json:"thumb"`
			Title            string `json:"title"`
			Type             string `json:"type"`
			UpdatedAt        int64  `json:"updatedAt"`
			UUID             string `json:"uuid"`
		} `json:"Directory"`
		AllowSync       bool   `json:"allowSync"`
		Identifier      string `json:"identifier"`
		MediaTagPrefix  string `json:"mediaTagPrefix"`
		MediaTagVersion int64  `json:"mediaTagVersion"`
		Size            int64  `json:"size"`
		Title1          string `json:"title1"`
	} `json:"MediaContainer"`
}

type SearchResult struct {
	MediaContainer struct {
		Metadata []struct {
			Collection []struct {
				Tag string `json:"tag"`
			} `json:"Collection"`
			Country []struct {
				Tag string `json:"tag"`
			} `json:"Country"`
			Director []struct {
				Tag string `json:"tag"`
			} `json:"Director"`
			Genre []struct {
				Tag string `json:"tag"`
			} `json:"Genre"`
			Media []struct {
				Part []struct {
					AudioProfile string `json:"audioProfile"`
					Container    string `json:"container"`
					Duration     int64  `json:"duration"`
					File         string `json:"file"`
					ID           int64  `json:"id"`
					Key          string `json:"key"`
					Size         int64  `json:"size"`
					VideoProfile string `json:"videoProfile"`
				} `json:"Part"`
				AspectRatio     float64 `json:"aspectRatio"`
				AudioChannels   int64   `json:"audioChannels"`
				AudioCodec      string  `json:"audioCodec"`
				AudioProfile    string  `json:"audioProfile"`
				Bitrate         int64   `json:"bitrate"`
				Container       string  `json:"container"`
				Duration        int64   `json:"duration"`
				Height          int64   `json:"height"`
				ID              int64   `json:"id"`
				VideoCodec      string  `json:"videoCodec"`
				VideoFrameRate  string  `json:"videoFrameRate"`
				VideoProfile    string  `json:"videoProfile"`
				VideoResolution string  `json:"videoResolution"`
				Width           int64   `json:"width"`
			} `json:"Media"`
			Role []struct {
				Tag string `json:"tag"`
			} `json:"Role"`
			Writer []struct {
				Tag string `json:"tag"`
			} `json:"Writer"`
			AddedAt                int64   `json:"addedAt"`
			AllowSync              bool    `json:"allowSync"`
			Art                    string  `json:"art"`
			ChapterSource          string  `json:"chapterSource"`
			Duration               int64   `json:"duration"`
			GUID                   string  `json:"guid"`
			HasPremiumPrimaryExtra string  `json:"hasPremiumPrimaryExtra"`
			Key                    string  `json:"key"`
			LastViewedAt           int64   `json:"lastViewedAt"`
			LibrarySectionID       int64   `json:"librarySectionID"`
			LibrarySectionTitle    string  `json:"librarySectionTitle"`
			LibrarySectionUUID     string  `json:"librarySectionUUID"`
			OriginalTitle          string  `json:"originalTitle"`
			OriginallyAvailableAt  string  `json:"originallyAvailableAt"`
			Personal               bool    `json:"personal"`
			Rating                 float64 `json:"rating"`
			RatingImage            string  `json:"ratingImage"`
			RatingKey              string  `json:"ratingKey"`
			SourceTitle            string  `json:"sourceTitle"`
			Studio                 string  `json:"studio"`
			Summary                string  `json:"summary"`
			Tagline                string  `json:"tagline"`
			Thumb                  string  `json:"thumb"`
			Title                  string  `json:"title"`
			Type                   string  `json:"type"`
			UpdatedAt              int64   `json:"updatedAt"`
			UserRating             float64 `json:"userRating"`
			ViewCount              int64   `json:"viewCount"`
			Year                   int64   `json:"year"`
		} `json:"Metadata"`
		Provider []struct {
			Key   string `json:"key"`
			Title string `json:"title"`
			Type  string `json:"type"`
		} `json:"Provider"`
		Identifier      string `json:"identifier"`
		MediaTagPrefix  string `json:"mediaTagPrefix"`
		MediaTagVersion int64  `json:"mediaTagVersion"`
		Size            int64  `json:"size"`
	} `json:"MediaContainer"`
}

type CollectionsResult struct {
	MediaContainer struct {
		Metadata []struct {
			AddedAt       int64  `json:"addedAt"`
			ChildCount    string `json:"childCount"`
			ContentRating string `json:"contentRating"`
			GUID          string `json:"guid"`
			Index         int64  `json:"index"`
			Key           string `json:"key"`
			MaxYear       string `json:"maxYear"`
			MinYear       string `json:"minYear"`
			RatingKey     string `json:"ratingKey"`
			Subtype       string `json:"subtype"`
			Summary       string `json:"summary"`
			Thumb         string `json:"thumb"`
			Title         string `json:"title"`
			TitleSort     string `json:"titleSort"`
			Type          string `json:"type"`
			UpdatedAt     int64  `json:"updatedAt"`
		} `json:"Metadata"`
		AllowSync           bool   `json:"allowSync"`
		Art                 string `json:"art"`
		Identifier          string `json:"identifier"`
		LibrarySectionID    int64  `json:"librarySectionID"`
		LibrarySectionTitle string `json:"librarySectionTitle"`
		LibrarySectionUUID  string `json:"librarySectionUUID"`
		MediaTagPrefix      string `json:"mediaTagPrefix"`
		MediaTagVersion     int64  `json:"mediaTagVersion"`
		Size                int64  `json:"size"`
		Thumb               string `json:"thumb"`
		Title1              string `json:"title1"`
		Title2              string `json:"title2"`
		ViewGroup           string `json:"viewGroup"`
		ViewMode            int64  `json:"viewMode"`
	} `json:"MediaContainer"`
}

type CollectionByTitleResult struct {
	AddedAt       int64  `json:"addedAt"`
	ChildCount    string `json:"childCount"`
	ContentRating string `json:"contentRating"`
	GUID          string `json:"guid"`
	Index         int64  `json:"index"`
	Key           string `json:"key"`
	MaxYear       string `json:"maxYear"`
	MinYear       string `json:"minYear"`
	RatingKey     string `json:"ratingKey"`
	Subtype       string `json:"subtype"`
	Summary       string `json:"summary"`
	Thumb         string `json:"thumb"`
	Title         string `json:"title"`
	TitleSort     string `json:"titleSort"`
	Type          string `json:"type"`
	UpdatedAt     int64  `json:"updatedAt"`
}
