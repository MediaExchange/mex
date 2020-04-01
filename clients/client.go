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
package clients

import "sync"

// MediaType provides the type of media (movie, tv show, etc.)
type MediaType int

// Defines the different types of media handled
const (
	Movie   MediaType = iota
	TvShow
)

// SearchResult defines the common fields for the various clients to return.
type SearchResult struct {
	Id          string      `json:"id"`             // ID of the media.
	Type        MediaType   `json:"type"`           // Type of media.
	Adult       bool        `json:"adult"`          // True if the media is for adults.
	Title       string      `json:"title"`          // Name of the media found.
	ImdbId      string      `json:"imdbId"`         // ID of the media on IMDB.
	Status      string      `json:"status"`         // Current status of the media (released, in production, etc.)
	Runtime     int         `json:"runtime"`        // Runtime of the media in minutes.
	Episodes    []Episode   `json:"episodes"`       // Episodes if the media contains any.
	Homepage    string      `json:"homepage"`       // Homepage of the media.
	Overview    string      `json:"overview"`       // Overview description of the media.
	Zap2itId    string      `json:"zap2itId"`       // ID of the media on Zap2It
	PosterUri   string      `json:"posterUri"`      // URI of an image that can be displayed.
	ReleaseDate string      `json:"releaseDate"`    // When the media first aired on TV or was released in theaters.
}

// Episode contains information about a specific episode within a TV season.
type Episode struct {
	Name        string      `json:"name"`           // Episode name.
	Number      int         `json:"number"`         // Absolute (overall) episode number across all seasons.
	Season      int         `json:"season"`         // Season number.
	AirDate     string      `json:"airDate"`        // Date the episode first aired.
	Episode     int         `json:"episode"`        // Episode number.
	Overview    string      `json:"overview"`       // Overview description of the episode.
}

// SearchContext contains all of the channels used to make the operation asynchronous.
type SearchContext struct {
	ResultChan  chan SearchResult
	ErrorChan   chan error
	DoneChan    chan bool
	Waiter      sync.WaitGroup
}

// NewSearchContext returns a new context used to make the search asynchronous.
func NewSearchContext() *SearchContext {
	return &SearchContext {
		ResultChan: make(chan SearchResult),
		ErrorChan:  make(chan error),
		DoneChan:   make(chan bool),
	}
}
