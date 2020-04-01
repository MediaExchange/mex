package models

// Details defines the common fields for various search providers to return.
type Details struct {
	Id          string      `json:"id"`             // ID of the media.
	Type        MediaType   `json:"type"`           // Type of media.
	Adult       bool        `json:"adult"`          // True if the media is for adults.
	Links       []Link      `json:"links"`          // Links to external information about the media.
	Title       string      `json:"title"`          // Name of the media found.
	Status      string      `json:"status"`         // Current status of the media (released, in production, etc.)
	Runtime     int         `json:"runtime"`        // Runtime of the media in minutes.
	Episodes    []Episode   `json:"episodes"`       // Episodes if the media contains any.
	Overview    string      `json:"overview"`       // Overview description of the media.
	PosterUri   string      `json:"posterUri"`      // URI of the poster image to display.
	ReleaseDate string      `json:"releaseDate"`    // When the media first aired on TV or was released in theaters.
}

// Link contains a reference to external information about the media.
type Link struct {
	Name        string      `json:"name"`           // Name of linked service.
	Url         string      `json:"url"`            // Service URL for the media.
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
