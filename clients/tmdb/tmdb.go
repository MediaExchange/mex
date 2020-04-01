/*
   Copyright 2019 Paul Howes

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/
package tmdb

import (
	"errors"
	"fmt"
	"github.com/MediaExchange/log"
	"github.com/MediaExchange/mex/clients/rest"
	"github.com/MediaExchange/mex/models"
	"strconv"
)

const (
	BaseUri = "https://api.themoviedb.org/3"
	//ImageUri = "https://image.tmdb.org/t/p/w154"
	//ImageUri = "https://image.tmdb.org/t/p/w185"
	ImageUri = "https://image.tmdb.org/t/p/w342"
)

var (
	ApiKey string
)

type searchResult struct {
	Id                  int     `json:"id"`
	Adult               bool    `json:"adult"`
	Title               string  `json:"title"`
	Video               bool    `json:"video"`
	ImdbId              string  `json:"imdb_id"`
	Status              string  `json:"status"`
	Runtime             int     `json:"runtime"`
	Homepage            string  `json:"homepage"`
	Overview            string  `json:"overview"`
	VoteCount           int     `json:"vote_count"`
	Popularity          float32 `json:"popularity"`
	PosterPath          string  `json:"poster_path"`
	ReleaseDate         string  `json:"release_date"`
	VoteAverage         float32 `json:"vote_average"`
	BackdropPath        string  `json:"backdrop_path"`
	OriginalTitle       string  `json:"original_title"`
	OriginalLanguage    string  `json:"original_language"`
}

type pagedSearchResult struct {
	Page            int            `json:"page"`
	TotalResults    int            `json:"total_results"`
	TotalPages      int            `json:"total_pages"`
	Results         []searchResult `json:"results"`
}

// detailResult contains some of the fields from the Get Movie Details endpoint:
// https://developers.themoviedb.org/3/movies/get-movie-details
type detailResult struct {
	Id                  int     `json:"id"`
	Adult               bool    `json:"adult"`
	Title               string  `json:"title"`
	Video               bool    `json:"video"`
	ImdbId              string  `json:"imdb_id"`
	Status              string  `json:"status"`
	Runtime             int     `json:"runtime"`
	Homepage            string  `json:"homepage"`
	Overview            string  `json:"overview"`
	VoteCount           int     `json:"vote_count"`
	Popularity          float32 `json:"popularity"`
	PosterPath          string  `json:"poster_path"`
	ReleaseDate         string  `json:"release_date"`
	VoteAverage         float64 `json:"vote_average"`
}

// Login retrieves the API key from the environment. There is no actual login
// like TVDB uses.
func Login(apikey string) error {
	log.Info("tmdb:Login")
	ApiKey = apikey
	return nil
}

func Search(name string, results *[]models.SearchResult) error {
	// Must log in before searching.
	if len(ApiKey) == 0 {
		s := "tmdb.Search: Must login with API key before searching"
		log.Error(s)
		return errors.New(s)
	}

	// Name must be provided.
	if len(name) == 0 {
		s := "tmdb.Search: name must be provided"
		log.Error(s)
		return errors.New(s)
	}

	// Result array must exist
	if results == nil {
		s := "tmdb.Search: results array must exist"
		log.Error(s)
		return errors.New(s)
	}

	log.Info("tmdb.Search: Starting movie search", log.String("name", name))
	return pagedSearch(name, 1, results)
}

func pagedSearch(name string, page int, results *[]models.SearchResult) error {
	reply := new(pagedSearchResult)
	_, err := newRequest().
		AddQuery("page", strconv.Itoa(page)).
		AddQuery("query", name).
		AddQuery("include_adult", "true").
		SetReplyBody(reply).
		Get(BaseUri + "/search/movie")
	if err != nil {
		log.Error("tmdb:pagedSearch: unexpected error", log.Err(err))
		return err
	}

	// Iterate through the results and convert each to a generic models.SearchResult object.
	for _, r := range reply.Results {
		// Only respond with results that have an image.
		if len(r.PosterPath) > 0 {
			sr := models.SearchResult {
				Id:          fmt.Sprintf("tmdb:%d", r.Id),
				Type:        models.Movie,
				Adult:       r.Adult,
				Title:       r.Title,
				Overview:    r.Overview,
				PosterUri:   ImageUri + r.PosterPath,
				ReleaseDate: r.ReleaseDate,
			}
			*results = append(*results, sr)
		}
	}

	// Recursively process the next page.
	if reply.Page < reply.TotalPages {
		return pagedSearch(name, reply.Page + 1, results)
	}

	return nil
}

func Details(id int) (*models.Details, error) {
	// TODO: log.Int() doesn't work here for some reason when id=299534
	log.Info("tmdb.Details", log.Int64("id", int64(id)))

	// Call the TMDB movie details service
	url := fmt.Sprintf("%s/movie/%d", BaseUri, id)
	reply := new(detailResult)
	_, err := newRequest().
		SetReplyBody(reply).
		Get(url)
	if err != nil {
		log.Error("tmdb.Details: unexpected error", log.String("url", url), log.Err(err))
		return nil, err
	}

	// Build the details model.
	d := &models.Details {
		Id:          fmt.Sprintf("tmdb:%d", reply.Id),
		Type:        models.Movie,
		Adult:       reply.Adult,
		Title:       reply.Title,
		Status:      reply.Status,
		Runtime:     reply.Runtime,
		Overview:    reply.Overview,
		ReleaseDate: reply.ReleaseDate,
	}

	if len(reply.PosterPath) > 0 {
		d.PosterUri = ImageUri + reply.PosterPath
	}

	if len(reply.Homepage) > 0 {
		d.Links = append(d.Links, models.Link {
			Name: "Homepage",
			Url: reply.Homepage,
		})
	}

	if len(reply.ImdbId) > 0 {
		d.Links = append(d.Links, models.Link {
			Name: "IMDB",
			Url:  "https://www.imdb.com/title/" + reply.ImdbId,
		})
	}

	d.Links = append(d.Links, models.Link {
		Name: "TheMovieDB",
		Url:  fmt.Sprintf("https://www.themoviedb.org/movie/%d", reply.Id),
	})

	return d, nil
}

// newRequest returns a new REST request with the API key set.
func newRequest() *rest.RestRequest {
	return rest.NewRequest().
		AddQuery("api_key", ApiKey)
}
