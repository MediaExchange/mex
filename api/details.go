package api

import (
	"encoding/json"
	"github.com/MediaExchange/log"
	"github.com/MediaExchange/router"
	"github.com/MediaExchange/mex/clients/tmdb"
	"github.com/MediaExchange/mex/clients/tvdb"
	"net/http"
	"strconv"
	"strings"
)

// GetDetails retrieves detailed information for media.
// The query parameter `id` contains the media provider and the provider's ID in the format `provider:id`.
func GetDetails(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Add("Access-Control-Allow-Origin", "*")

	// Get the provider and ID from the query string.
	params := router.GetParams(request.Context())
	param := params["id"]
	if len(param) == 0 {
		log.Error("api.GetDetails `id` query parameter is empty.")
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the provider and ID
	temp := strings.Split(param, ":")
	if len(temp) != 2 {
		log.Error("api.GetDetails `id` query parameter was not in the form provider:id", log.String("param", param))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	provider := temp[0]
	id, err := strconv.Atoi(temp[1])
	if err != nil {
		log.Error("api.getDetails `id` was not a number.", log.Err(err))
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Info("api.GetDetails", log.String("param", param))

	switch provider {
	case "tmdb":
		res, err := tmdb.Details(id)
		if err != nil {
			// Error was already logged in tmdb.Details. Just report it back to the client.
			writer.Header().Set("Content-Type", "text/plain")
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write([]byte(err.Error()))
			return
		}

		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(writer).Encode(res)

	case "tvdb":
		res, err := tvdb.Details(id)
		if err != nil {
			// Error was already logged in tmdb.Details. Just report it back to the client.
			writer.Header().Set("Content-Type", "text/plain")
			writer.WriteHeader(http.StatusInternalServerError)
			_, _ = writer.Write([]byte(err.Error()))
			return
		}

		writer.Header().Add("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(writer).Encode(res)

	default:
		log.Error("api.GetDetails: Unknown provider", log.String("provider", provider))
		writer.Header().Set("Content-Type", "text/plain")
		writer.WriteHeader(http.StatusBadRequest)
		_, _ = writer.Write([]byte("api.GetDetails: Unknown provider: " + provider))
	}
}
