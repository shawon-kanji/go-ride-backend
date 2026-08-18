package places

// Suggestion is one autocomplete result — a candidate place the caller
// hasn't resolved to coordinates yet (that's PlaceDetails, a second call,
// matching how the design's R01/R02 search-then-confirm flow works).
type Suggestion struct {
	PlaceID       string
	PrimaryText   string
	SecondaryText string
}

// ResolvedPlace is a place with real coordinates — the output of either
// resolving a Suggestion's PlaceID or reverse-geocoding a dropped pin.
type ResolvedPlace struct {
	PlaceID          string
	Name             string
	FormattedAddress string
	Lat              float64
	Lng              float64
}
