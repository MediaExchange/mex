package models

// MediaType provides the type of media (movie, tv show, etc.)
type MediaType int

// Defines the different types of media handled
const (
	Movie   MediaType = iota
	TvShow
)
