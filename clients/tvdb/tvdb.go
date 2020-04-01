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
package tvdb

import (
	"errors"
	"fmt"
	"github.com/MediaExchange/log"
	"github.com/MediaExchange/mex/clients/rest"
	"github.com/MediaExchange/mex/models"
	"net/http"
	"strconv"
)

const (
	BaseUri = "https://api.thetvdb.com"
	ImageUri = "https://www.thetvdb.com/banners/"
)

var (
	// Stores the token received by the Login function.
	Token string
)

// Request sent to retrieve a token.
type tokenRequest struct {
	ApiKey string `json:"apikey"`
}

// Login response.
type tokenReply struct {
	Token string `json:"token"`
}

// Search response.
type searchResult struct {
	Data []struct {
		Id              int         `json:"id"`
		Slug            string      `json:"slug"`
		Banner          string      `json:"banner"`
		Status          string      `json:"status"`
		Aliases         []string    `json:"aliases"`
		Network         string      `json:"network"`
		Overview        string      `json:"overview"`
		FirstAired      string      `json:"firstAired"`
		SeriesName      string      `json:"seriesName"`
	}                               `json:"data"`
}

// Image response.
type imageResult struct {
	Data []struct {
		FileName string `json:"fileName"`
		Id int `json:"id"`
		KeyType string `json:"keyType"`
		LanguageId int `json:"languageId"`
		RatingsInfo struct {
			Average float64 `json:"average"`
			Count int `json:"count"`
		} `json:"ratingsInfo"`
		Resolution string `json:"resolution"`
		SubKey string `json:"subKey"`
		Thumbnail string `json:"thumbnail"`
	} `json:"data"`
	Errors struct {
		InvalidFilters []string `json:"invalidFilters"`
		InvalidLanguage string `json:"invalidLanguage"`
		InvalidQueryParams []string `json:"invalidQueryParams"`
	} `json:"errors"`
}

// pagedEpisodeResult contains a page of episodeResult.
type pagedEpisodeResult struct {
	Links struct {
		First           int     `json:"first"`
		Last            int     `json:"last"`
		Next            int     `json:"next"`
		Prev            int     `json:"prev"`
	} `json:"links"`
	Data []episodeResult`json:"data"`
}

// episodeResult represents a single episode in a show.
type episodeResult struct {
	AbsoluteNumber      int     `json:"absoluteNumber"`
	AiredEpisodeNumber  int     `json:"airedEpisodeNumber"`
	AiredSeason         int     `json:"airedSeason"`
	AirsAfterSeason     int     `json:"airsAfterSeason"`
	AirsBeforeEpisode   int     `json:"airsBeforeEpisode"`
	AirsBeforeSeason    int     `json:"airsBeforeSeason"`
	EpisodeName         string  `json:"episodeName"`
	FirstAired          string  `json:"firstAired"`
	Id                  int     `json:"id"`
	ImdbId              string  `json:"imdbId"`
	Overview            string  `json:"overview"`
	ShowUrl             string  `json:"showUrl"`
}

// detailResult represents additional information about the show.
type detailResult struct {
	Data struct {
		Id              int     `json:"id"`             // ID of the show
		Slug            string  `json:"slug"`           // URL slug added to https://www.thetvdb.com/series/%s
		ImdbId          string  `json:"imdbId"`         // ID of the show in IMDB
		Status          string  `json:"status"`         // Current status of the show ("Cancelled", "In Production", etc.)
		Runtime         string  `json:"runtime"`        // Time in minutes each episode runs.
		Overview        string  `json:"overview"`       // Overview description of the show.
		Zap2itId        string  `json:"zap2itId"`       // ID of the show in Zap2It
		FirstAired      string  `json:"firstAired"`     // Date when the show first aired.
		SeriesName      string  `json:"seriesName"`     // Name of the series.
	}                           `json:"data"`
}

// Login returns an authentication token used in future API calls.
func Login(apikey string) error {
	log.Info("tvdb.Login")
	reply := new(tokenReply)
	_, err := rest.NewRequest().
		SetBody(tokenRequest {
			ApiKey: apikey,
		}).
		SetReplyBody(reply).
		Post(BaseUri + "/login")

	if err != nil {
		log.Error("TVDB Unexpected error", log.Err(err))
		return err
	}

	Token = reply.Token
	return nil
}

// Search for a show by name.
func Search(name string, results *[]models.SearchResult) error {
	log.Info("tvdb.Search", log.String("name", name))

	// Refresh the auth token if necessary.
	if err := refresh(); err != nil {
		return err
	}

	reply := new(searchResult)
	_, err := rest.NewRequest().
		SetBearerAuth(Token).
		AddQuery("name", name).
		SetReplyBody(reply).
		Get(BaseUri + "/search/series")
	if err != nil {
		log.Error("tvdb.Search: Unexpected error", log.Err(err))
		return err
	}

	// Add all of the shows to the search results.
	for _, r := range reply.Data {
		// Retrieves the poster image URL path for the show.
		posterUrl, err := posterImageUrl(r.Id)
		if err != nil {
			return err
		}

		if len(posterUrl) > 0 {
			sr := models.SearchResult{
				Id:          fmt.Sprintf("tvdb:%d", r.Id),
				Type:        models.TvShow,
				Adult:       false,
				Title:       r.SeriesName,
				Overview:    r.Overview,
				PosterUri:   posterUrl,
				ReleaseDate: r.FirstAired,
			}
			*results = append(*results, sr)
		}
	}

	return nil
}

func Details(id int) (*models.Details, error) {
	// TODO: log.Int() doesn't work here for some reason when id=264030
	log.Info("tvdb.Details", log.Int64("id", int64(id)))

	// Refresh the auth token if necessary.
	if err := refresh(); err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/series/%d", BaseUri, id)
	reply := new(detailResult)
	_, err := newRequest().
		SetReplyBody(reply).
		Get(url)
	if err != nil {
		log.Error("tvdb.Details: unexpected error", log.String("url", url), log.Err(err))
		return nil, err
	}

	// Build the details model.
	d := &models.Details {
		Id:          fmt.Sprintf("tvdb:%d", reply.Data.Id),
		Type:        models.TvShow,
		Title:       reply.Data.SeriesName,
		Status:      reply.Data.Status,
		Overview:    reply.Data.Overview,
		ReleaseDate: reply.Data.FirstAired,
	}

	// The runtime appears to be minutes stored in a string. Fingers crossed that this works.
	runtime, err := strconv.Atoi(reply.Data.Runtime)
	if err != nil {
		d.Runtime = runtime
	}

	// IMDB ID converted to a link.
	if len(reply.Data.ImdbId) > 0 {
		d.Links = append(d.Links, models.Link {
			Name: "IMDB",
			Url:  "https://www.imdb.com/title/" + reply.Data.ImdbId,
		})
	}

	// Zap2It ID converted to a link.
	if len(reply.Data.Zap2itId) > 0 {
		d.Links = append(d.Links, models.Link {
			Name: "Zap2It",
			Url:  "https://tvlistings.zap2it.com/overview.html?programSeriesId=" + reply.Data.Zap2itId,
		})
	}

	// TVDB slug converted to a link.
	if len(reply.Data.Slug) > 0 {
		d.Links = append(d.Links, models.Link {
			Name: "TheTVDB",
			Url:  "https://www.thetvdb.com/series/" + reply.Data.Slug,
		})
	}

	// TVDB makes the poster image URL a separate API call.
	posterUrl, err := posterImageUrl(id)
	if err != nil {
		return nil, err
	}
	d.PosterUri = posterUrl

	// Get all of the episodes and add them to the details
	d.Episodes = make([]models.Episode, 0)
	if err := getEpisodes(id, 1, &d.Episodes); err != nil {
		return nil, err
	}


	return d, nil
}

// getEpisodes returns a page of episode information for a series.
func getEpisodes(id int, page int, episodes *[]models.Episode) error {
	// Retrieve a page of getEpisodes.
	path := fmt.Sprintf("/series/%d/episodes", id)
	reply := new(pagedEpisodeResult)
	_, err := rest.NewRequest().
		SetBearerAuth(Token).
		AddQuery("page", strconv.Itoa(page)).
		SetReplyBody(reply).
		Get(BaseUri + path)
	if err != nil {
		// TODO: log.Int crashed with a value of id=264030
		// TODO: log.Int crashed with a value of page=1
		log.Error("tvdb.getEpisodes: unexpected error", log.Int64("id", int64(id)), log.Int64("page", int64(page)), log.Err(err))
		return err
	}

	for _, r := range reply.Data {
		e := models.Episode {
			Name:     r.EpisodeName,
			Number:   r.AbsoluteNumber,
			Season:   r.AiredSeason,
			AirDate:  r.FirstAired,
			Episode:  r.AiredEpisodeNumber,
			Overview: r.Overview,
		}
		*episodes = append(*episodes, e)
	}

	if reply.Links.Next > 0 {
		if err := getEpisodes(id, page + 1, episodes); err != nil {
			// The error was already logged. Just pass it back up the call stack.
			return err
		}
	}

	return nil
}

// Refresh updates the token expiration without performing a full authentication.
func refresh() error {
	if len(Token) == 0 {
		s := "tvdb.refresh: login before using API"
		log.Error(s)
		return errors.New(s)
	}

	reply := new(tokenReply)
	_, err := rest.NewRequest().
		SetBearerAuth(Token).
		SetReplyBody(reply).
		Get(BaseUri + "/refresh_token")
	if err != nil {
		log.Error("tvdb.refresh: unexpected error", log.Err(err))
		return err
	}

	Token = reply.Token
	return nil
}

// posterImageUrl returns the URL of the cover art image for the series.
func posterImageUrl(id int) (string, error) {
	// Query the web service.
	reply := new(imageResult)
	path := fmt.Sprintf("/series/%d/images/query", id)
	res, err := rest.NewRequest().
		SetBearerAuth(Token).
		AddQuery("keyType", "poster").
		SetReplyBody(reply).
		Get(BaseUri + path)
	if err != nil {
		// An image not being available is acceptable.
		if res.StatusCode == http.StatusNotFound {
			return "", nil
		}

		// Other errors can't be handled.
		log.Error("tvdb.posterImageUrl: Unexpected error", log.Err(err))
		return "", err
	}

	// Sort the images by rating and only return the highest rated image.
	var average float64 = 0.0
	var imagePath string
	for _, data := range reply.Data {
		if data.RatingsInfo.Average > average {
			imagePath = data.FileName
		}
	}

	// TVDB doesn't like to host images, so we have to proxy the URL.
	if len(imagePath) > 0 {
		imagePath = "http://localhost:9000/api/proxy?url=" + ImageUri + imagePath
	}

	return imagePath, nil
}

// newRequest returns a new RestRequest object with the bearer authentication token already added.
func newRequest() *rest.RestRequest {
	return rest.NewRequest().SetBearerAuth(Token)
}