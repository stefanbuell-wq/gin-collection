package models

import "time"

// Gin represents a gin bottle in a collection
type Gin struct {
	ID                 int64      `json:"id"`
	TenantID           int64      `json:"tenant_id"`
	UUID               string     `json:"uuid"`
	Name               string     `json:"name" binding:"required,min=1,max=255"`
	Brand              *string    `json:"brand,omitempty"`
	Country            *string    `json:"country,omitempty"`
	Region             *string    `json:"region,omitempty"`
	GinType            *string    `json:"gin_type,omitempty"` // London Dry, Old Tom, etc.
	ABV                *float64   `json:"abv,omitempty"`
	BottleSize         *int       `json:"bottle_size,omitempty"` // ml
	FillLevel          *int       `json:"fill_level,omitempty"`  // 0-100%
	Price              *float64   `json:"price,omitempty"`
	CurrentMarketValue *float64   `json:"current_market_value,omitempty"`
	PurchaseDate       *time.Time `json:"purchase_date,omitempty"`
	PurchaseLocation   *string    `json:"purchase_location,omitempty"`
	Barcode            *string    `json:"barcode,omitempty"`
	Rating             *int       `json:"rating,omitempty"` // 1-5
	NoseNotes          *string    `json:"nose_notes,omitempty"`
	PalateNotes        *string    `json:"palate_notes,omitempty"`
	FinishNotes        *string    `json:"finish_notes,omitempty"`
	GeneralNotes       *string    `json:"general_notes,omitempty"`
	Description        *string    `json:"description,omitempty"`
	PhotoURL           *string    `json:"photo_url,omitempty"`
	IsFinished         bool       `json:"is_finished"`
	RecommendedTonic   *string    `json:"recommended_tonic,omitempty"`
	RecommendedGarnish *string    `json:"recommended_garnish,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`

	// Primary photo URL (from photos table, set by repository)
	PrimaryPhotoURL *string `json:"primary_photo_url,omitempty"`

	// Related data (loaded separately)
	Botanicals []*GinBotanical `json:"botanicals,omitempty"`
	Photos     []*GinPhoto     `json:"photos,omitempty"`
	Cocktails  []*Cocktail     `json:"cocktails,omitempty"`
}

// GinFilter represents filtering options for gin queries
type GinFilter struct {
	TenantID   int64
	IsFinished *bool
	GinType    *string
	Country    *string
	MinRating  *int
	MaxRating  *int
	SortBy     string // name, rating, price, date, country, fill_level
	SortOrder  string // asc, desc
	Limit      int
	Offset     int
}

// Botanical represents a botanical ingredient
type Botanical struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Category    *string `json:"category,omitempty"`
	Description *string `json:"description,omitempty"`
}

// GinBotanical represents the many-to-many relationship between gins and botanicals
type GinBotanical struct {
	ID          int64      `json:"id"`
	TenantID    int64      `json:"tenant_id"`
	GinID       int64      `json:"gin_id"`
	BotanicalID int64      `json:"botanical_id"`
	Prominence  Prominence `json:"prominence"`

	// Loaded botanical info
	Botanical *Botanical `json:"botanical,omitempty"`
}

// Prominence indicates how prominent a botanical is in a gin
type Prominence string

const (
	ProminenceDominant Prominence = "dominant"
	ProminenceNotable  Prominence = "notable"
	ProminenceSubtle   Prominence = "subtle"
)

// GinPhoto represents a photo of a gin bottle
type GinPhoto struct {
	ID         int64      `json:"id"`
	TenantID   int64      `json:"tenant_id"`
	GinID      int64      `json:"gin_id"`
	PhotoURL   string     `json:"photo_url"`
	PhotoType  PhotoType  `json:"photo_type"`
	Caption    *string    `json:"caption,omitempty"`
	IsPrimary  bool       `json:"is_primary"`
	StorageKey *string    `json:"storage_key,omitempty"`
	FileSizeKB *int       `json:"file_size_kb,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
}

// PhotoType represents the type of photo
type PhotoType string

const (
	PhotoTypeBottle  PhotoType = "bottle"
	PhotoTypeLabel   PhotoType = "label"
	PhotoTypeMoment  PhotoType = "moment"
	PhotoTypeTasting PhotoType = "tasting"
)

// Cocktail represents a cocktail recipe
type Cocktail struct {
	ID           int64            `json:"id"`
	Name         string           `json:"name"`
	Description  *string          `json:"description,omitempty"`
	Instructions *string          `json:"instructions,omitempty"`
	GlassType    *string          `json:"glass_type,omitempty"`
	IceType      *string          `json:"ice_type,omitempty"`
	Difficulty   CocktailDifficulty `json:"difficulty"`
	PrepTime     *int             `json:"prep_time,omitempty"` // minutes
	CreatedAt    time.Time        `json:"created_at"`

	// Related data
	Ingredients []*CocktailIngredient `json:"ingredients,omitempty"`
}

// CocktailDifficulty represents cocktail preparation difficulty
type CocktailDifficulty string

const (
	DifficultyEasy   CocktailDifficulty = "easy"
	DifficultyMedium CocktailDifficulty = "medium"
	DifficultyHard   CocktailDifficulty = "hard"
)

// CocktailIngredient represents an ingredient in a cocktail
type CocktailIngredient struct {
	ID         int64   `json:"id"`
	CocktailID int64   `json:"cocktail_id"`
	Ingredient string  `json:"ingredient"`
	Amount     *string `json:"amount,omitempty"`
	Unit       *string `json:"unit,omitempty"`
	IsGin      bool    `json:"is_gin"`
}

// TastingSession represents a tasting event
type TastingSession struct {
	ID        int64      `json:"id"`
	TenantID  int64      `json:"tenant_id"`
	GinID     int64      `json:"gin_id"`
	UserID    *int64     `json:"user_id,omitempty"`
	Date      time.Time  `json:"date"`
	Notes     *string    `json:"notes,omitempty"`
	Rating    *int       `json:"rating,omitempty"`
	Tonic     *string    `json:"tonic,omitempty"`
	Botanicals *string   `json:"botanicals,omitempty"` // Comma-separated list
	CreatedAt time.Time  `json:"created_at"`

	// Loaded user info
	UserName *string `json:"user_name,omitempty"`
}

// TastingSessionWithGin includes gin information for list views
type TastingSessionWithGin struct {
	TastingSession
	GinName  string  `json:"gin_name"`
	GinBrand *string `json:"gin_brand,omitempty"`
}

// GinStats represents statistics about a gin collection
type GinStats struct {
	TotalGins           int                    `json:"total_gins"`
	AvailableGins       int                    `json:"available_gins"`
	FinishedGins        int                    `json:"finished_gins"`
	AverageRating       float64                `json:"average_rating"`
	TotalValue          float64                `json:"total_value"`
	TotalMarketValue    float64                `json:"total_market_value"`
	GinsByType          map[string]int         `json:"gins_by_type"`
	GinsByCountry       map[string]int         `json:"gins_by_country"`
	TopRatedGins        []*Gin                 `json:"top_rated_gins"`
	TopBotanicals       []*BotanicalCount      `json:"top_botanicals"`
	FillLevelDistribution map[string]int       `json:"fill_level_distribution"`
}

// BotanicalCount represents a botanical with its usage count
type BotanicalCount struct {
	Botanical *Botanical `json:"botanical"`
	Count     int        `json:"count"`
}
