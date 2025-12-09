package option

import (
	core "github.com/fern-demo/agoraio-go-sdk/core"
)

// Area type alias for global regions
type Area = core.Area

// Area constants
const (
	AreaUnknown = core.AreaUnknown
	AreaUS      = core.AreaUS
	AreaEU      = core.AreaEU
	AreaAP      = core.AreaAP
	AreaCN      = core.AreaCN
)

// WithArea creates a new AreaRequestOption with a pool for the specified area.
// The pool manages regional URL cycling and automatic domain selection.
func WithArea(area core.Area) *core.AreaRequestOption {
	return core.NewAreaRequestOption(area)
}

// WithPool creates a new AreaRequestOption with a pre-configured pool.
// Use this when you want to manage the pool lifecycle yourself.
func WithPool(pool *core.Pool) *core.AreaRequestOption {
	return &core.AreaRequestOption{Pool: pool}
}
