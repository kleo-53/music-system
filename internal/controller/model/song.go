package model

// Song is a title and group with optional data
// @Description 	Represents a music song entity
// @property 		Group 		The group name
// @property 		Song 		The title of the song
// @property 		Text 		(Optional) 	The text of the song
// @property 		ReleaseDate	(Optional) 	The release date of the song
// @property 		Link 		(Optional) 	A link to video for the song
type Song struct {
	Group       string `json:"group"`
	Song        string `json:"song"`
	Text        string `json:"text,omitempty"`
	ReleaseDate string `json:"releaseDate,omitempty"`
	Link        string `json:"link,omitempty"`
}

// SongFilters defines the optional filtering criteria for songs
// @Description 	Used for filtering songs by fields below
// @property 		Group 		The group name
// @property 		Song 		The title of the song
// @property 		Text 		(Optional) 	The text of the song
// @property 		ReleaseDate	(Optional) 	The release date of the song
// @property 		Link 		(Optional) 	A link to video for the song
type SongFilters struct {
	Group       string `json:"group,omitempty"`
	Song        string `json:"song,omitempty"`
	Text        string `json:"text,omitempty"`
	ReleaseDate string `json:"release_date,omitempty"`
	Link        string `json:"link,omitempty"`
}

// SongDetail represents details about song
// @Description Details about song
// @property ReleaseDate The release date of the song
// @property Text The text of the song
// @property Link A link to video for the song
type SongDetail struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

// SongCommon represents the common data for a song
// @Description Minimal required data to represent a song
// @property Group The group name
// @property Song The title of the song
type SongCommon struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}
