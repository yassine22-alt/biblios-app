package model

type SearchCriteria struct {
	Title           string  `json:"title,omitempty"`
	Author          string  `json:"author_name,omitempty"`
	MinPrice        float64 `json:"min_price,omitempty"`
	MaxPrice        float64 `json:"max_price,omitempty"`
	Genre           string  `json:"genre,omitempty"`
	PublishedAfter  string  `json:"published_after,omitempty"`
	PublishedBefore string  `json:"published_before,omitempty"`
}
