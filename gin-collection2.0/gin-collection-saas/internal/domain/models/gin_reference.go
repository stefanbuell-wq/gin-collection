package models

// GinReference represents a gin from the reference catalog
type GinReference struct {
	ID                 int64   `json:"id"`
	Name               string  `json:"name"`
	Brand              *string `json:"brand,omitempty"`
	Country            *string `json:"country,omitempty"`
	Region             *string `json:"region,omitempty"`
	GinType            *string `json:"gin_type,omitempty"`
	ABV                *float64 `json:"abv,omitempty"`
	BottleSize         *int    `json:"bottle_size,omitempty"`
	Description        *string `json:"description,omitempty"`
	NoseNotes          *string `json:"nose_notes,omitempty"`
	PalateNotes        *string `json:"palate_notes,omitempty"`
	FinishNotes        *string `json:"finish_notes,omitempty"`
	RecommendedTonic   *string `json:"recommended_tonic,omitempty"`
	RecommendedGarnish *string `json:"recommended_garnish,omitempty"`
	ImageURL           *string `json:"image_url,omitempty"`
	Barcode            *string `json:"barcode,omitempty"`
}

// GinReferenceSearchParams holds search parameters
type GinReferenceSearchParams struct {
	Query   string
	Country string
	GinType string
	Limit   int
	Offset  int
}
