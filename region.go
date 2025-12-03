package Agora

// Region represents an Agora API region.
type Region string

const (
	// RegionUS represents the United States region.
	RegionUS Region = "us"
	// RegionEU represents the European Union region.
	RegionEU Region = "eu"
	// RegionAPAC represents the Asia-Pacific region.
	RegionAPAC Region = "apac"
)

// regionBaseURLs maps each region to its corresponding API base URL.
// TODO: Replace placeholder URLs with actual Agora region-specific endpoints.
var regionBaseURLs = map[Region]string{
	RegionUS:   "https://api.agora.io/api/conversational-ai-agent",        // US endpoint (placeholder)
	RegionEU:   "https://api-eu.agora.io/api/conversational-ai-agent",     // EU endpoint (placeholder)
	RegionAPAC: "https://api-apac.agora.io/api/conversational-ai-agent",   // APAC endpoint (placeholder)
}

// BaseURLForRegion returns the API base URL for the given region.
// Returns the US region URL as default if the region is not recognized.
func BaseURLForRegion(region Region) string {
	if url, ok := regionBaseURLs[region]; ok {
		return url
	}
	return regionBaseURLs[RegionUS]
}

// AllRegions returns a slice of all available regions.
func AllRegions() []Region {
	return []Region{RegionUS, RegionEU, RegionAPAC}
}

// IsValidRegion checks if the given region is a valid Agora region.
func IsValidRegion(region Region) bool {
	_, ok := regionBaseURLs[region]
	return ok
}
