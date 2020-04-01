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
	"github.com/MediaExchange/log"
	"github.com/MediaExchange/router"
	"io/ioutil"
	"net/http"
	"strings"
)

// Proxy proxies a GET from the URL specified in as `http://localhost:9000/proxy?url=...`
// Some of the search providers do not want to be used as an image server.
// They prevent this by looking at the http-referrer header, which is
// automatically set by all browsers. By proxying the GET request through
// the API, the raw image data is retrieved with no extraneous headers and
// then returned to the browser for display.
func Proxy(writer http.ResponseWriter, request *http.Request) {
	params := router.GetParams(request.Context())
	urlString := params["url"]
	if len(urlString) == 0 {
		log.Error("api.Proxy url is empty.")
		writer.WriteHeader(400)
		return
	}

	log.Info("api.Proxy", log.String("url", urlString))

	res, err := http.Get(urlString)
	if err != nil {
		log.Error("api.Proxy unexpected error", log.Err(err))
		writer.WriteHeader(500)
		return
	}

	// Copy all of the headers from the remote server.
	for k, v := range res.Header {
		writer.Header().Set(k, strings.Join(v, ","))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Error("api.Proxy error reading response", log.Err(err))
		writer.WriteHeader(500)
		return
	}

	_, err = writer.Write(body)
	if err != nil {
		log.Error("api.Proxy error writing response", log.Err(err))
		writer.WriteHeader(500)
		return
	}
}
