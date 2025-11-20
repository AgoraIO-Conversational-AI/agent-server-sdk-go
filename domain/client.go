package domain

import (
	"github.com/fern-demo/agoraio-go-sdk/v505/client"
	"github.com/fern-demo/agoraio-go-sdk/v505/core"
	"github.com/fern-demo/agoraio-go-sdk/v505/option"
)

func NewClientWithRegion(region Region, opts ...option.RequestOption) (*client.Client, error) {
	baseURL, err := GetBaseURLForRegion(region)
	if err != nil {
		return nil, err
	}

	hasBaseURL := false
	for _, opt := range opts {
		if _, ok := opt.(*core.BaseURLOption); ok {
			hasBaseURL = true
			break
		}
	}

	if !hasBaseURL {
		opts = append([]option.RequestOption{option.WithBaseURL(baseURL)}, opts...)
	}

	return client.NewClient(opts...), nil
}
