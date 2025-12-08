package core

type Region string

const (
	RegionUSWest    Region = "us-west"
	RegionUSEast    Region = "us-east"
	RegionUSCentral Region = "us-central"
	RegionCanada    Region = "canada"
	RegionMexico    Region = "mexico"
)

type RegionRequestOption struct {
	Region Region
}

var regionalBaseURLs = map[Region]string{
	RegionUSWest:    "https://api-us-west.agora.io",
	RegionUSEast:    "https://api-us-east.agora.io",
	RegionUSCentral: "https://api-us-central.agora.io",
	RegionCanada:    "https://api-canada.agora.io",
	RegionMexico:    "https://api-mexico.agora.io",
}

func (o *RegionRequestOption) applyRequestOptions(opts *RequestOptions) {
	if baseURL, ok := regionalBaseURLs[o.Region]; ok {
		opts.BaseURL = baseURL
	}
}
