package models

// SearchResult defines the common fields for the various clients to return.
type SearchResult struct {
	Id          string      `json:"id"`             // ID of the media.
	Type        MediaType   `json:"type"`           // Type of media (Movie, TV Show, etc.)
	Adult       bool        `json:"adult"`          // True if the media is for adults (Rated X, TV-MA, etc.)
	Title       string      `json:"title"`          // Name of the media found.
	Overview    string      `json:"overview"`       // Overview description of the media.
	PosterUri   string      `json:"posterUri"`      // URI of an image that can be displayed.
	ReleaseDate string      `json:"releaseDate"`    // When the media first aired on TV or was released in theaters.
}
