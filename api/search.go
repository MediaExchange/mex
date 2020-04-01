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
package api

import (
	"encoding/json"
	"github.com/MediaExchange/log"
	"github.com/MediaExchange/router"
	"github.com/MediaExchange/mex/clients/tmdb"
	"github.com/MediaExchange/mex/clients/tvdb"
	"github.com/MediaExchange/mex/models"
	"net/http"
)

// Search finds media from all the search providers that matches the requested name.
func Search(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Access-Control-Allow-Origin", "*")

	// Get the name from the `?q=` query parameter.
	params := router.GetParams(request.Context())
	name := params["q"]
	if len(name) == 0 {
		log.Error("api.Search `q` query parameter is empty.")
		writer.WriteHeader(400)
		return
	}

	results := make([]models.SearchResult, 0)

	err := tmdb.Search(name, &results)
	if err != nil {
		// Error was already logged in tmdb.Search. Just report it back to the client.
		writer.Header().Set("Content-Type", "text/plain")
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}

	err = tvdb.Search(name, &results)
	if err != nil {
		// Error was already logged in tvdb.Search. Just report it back to the client.
		writer.Header().Set("Content-Type", "text/plain")
		writer.WriteHeader(http.StatusInternalServerError)
		_, _ = writer.Write([]byte(err.Error()))
		return
	}

	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(writer).Encode(results)
	return

	/*
	// Search The Movie Database.
	res, err := tmdb.Search(name)
	if err != nil {
		log.Error("api.Search unexpected error", log.Err(err))
		writer.WriteHeader(500)
		return
	}

	results := res

	// Search the TV Database
	res, err = tvdb.Search(name)
	if err != nil {
		log.Error("api.Search unexpected error", log.Err(err))
		writer.WriteHeader(500)
		return
	}

	results = append(results, res...)

	// Sort the results from newest to oldest.
	sort.Sort(sort.Reverse(clients.DateSorter(results)))

	// Respond with the combined results.
	writer.Header().Add("Content-Type", "application/json")
	writer.Header().Add("Access-Control-Allow-Origin", "*")

	if err := json.NewEncoder(writer).Encode(results); err != nil {
		log.Error("api.Search: serializer error", log.Err(err))
		writer.WriteHeader(500)
	}
	*/
}
