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
package main

import (
	"fmt"
	"github.com/MediaExchange/config"
	"github.com/MediaExchange/log"
	"github.com/MediaExchange/router"
	"github.com/MediaExchange/mex/api"
	"github.com/MediaExchange/mex/clients/tmdb"
	"github.com/MediaExchange/mex/clients/tvdb"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	var err error

	// Read the main configuration file.
	conf := new(MexConfig)
	err = config.FromFile("./mex_config.yaml", conf)
	if err != nil {
		log.Error("Error reading ./mex_config.yaml", log.Err(err))
		os.Exit(1)
	}

	// Authenticate with TMDB
	err = tmdb.Login(conf.Clients.TmdbApiKey)
	if err != nil {
		log.Error("TMDB authentication error", log.Err(err))
		os.Exit(1)
	}

	// Authenticate with TVDB
	err = tvdb.Login(conf.Clients.TvdbApiKey)
	if err != nil {
		log.Error("TVDB authentication error", log.Err(err))
		os.Exit(1)
	}

	port := conf.Server.Port
	addr := fmt.Sprintf(":%d", port)
	log.Info("Starting HTTP server", log.Int16("port", port))

	// Configure the router
	handler := router.NewRouter().
		AddRoute("GET", "/api/details", api.GetDetails).
		AddRoute("GET", "/api/proxy",   api.Proxy).
		AddRoute("GET", "/api/search",  api.Search).
		AddRoute("GET", "/.*",          FileServer)

	// Start the HTTP server
	err = http.ListenAndServe(addr, handler)
	if err != nil {
		log.Error("Fatal HTTP server error", log.Err(err))
		os.Exit(1)
	}

	os.Exit(0)
}

// FileServer serves static files from a directory.
func FileServer(writer http.ResponseWriter, request *http.Request) {
	// http.ServeFile doesn't support the concept of a base path to serve from.
	path := filepath.Join("ui/dist/mex", request.URL.Path)

	// Serve the file requested.
	http.ServeFile(writer, request, path)
}
