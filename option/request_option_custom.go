package option

import (
	core "github.com/fern-demo/agoraio-go-sdk/v505/core"
)

// Region type alias
type Region = core.Region

// Region constants
const (
	RegionUSWest    = core.RegionUSWest
	RegionUSEast    = core.RegionUSEast
	RegionUSCentral = core.RegionUSCentral
	RegionCanada    = core.RegionCanada
	RegionMexico    = core.RegionMexico
)

func WithRegion(region core.Region) *core.RegionRequestOption {
	return &core.RegionRequestOption{
		Region: region,
	}
}
