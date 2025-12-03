package Agora

import (
	"context"
	"errors"

	client "github.com/fern-demo/agoraio-go-sdk/v505/client"
	option "github.com/fern-demo/agoraio-go-sdk/v505/option"
)

// WrapperClient wraps the generated Agora client with automatic region selection.
// It embeds the base client, so all sub-clients (Agents, Telephony, PhoneNumbers)
// are accessible directly on the WrapperClient.
//
// Example usage:
//
//	ctx := context.Background()
//	client, err := Agora.NewWrapperClient(ctx, "your-app-id",
//	    option.WithBasicAuth("username", "password"),
//	)
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Access sub-clients directly
//	resp, err := client.Agents.Start(ctx, &Agora.StartAgentsRequest{...})
type WrapperClient struct {
	*client.Client

	// Region stores the detected or configured region for this client.
	Region Region
}

// WrapperClientConfig holds configuration options for creating a WrapperClient.
type WrapperClientConfig struct {
	// AppID is the Agora application ID used for region detection.
	AppID string

	// Region allows explicitly setting a region, bypassing auto-detection.
	// If set, the client will use this region instead of detecting one.
	Region *Region

	// RegionResolver is an optional custom function for determining the region.
	// If provided, it will be used instead of the default detection logic.
	// This is useful for implementing custom region selection strategies.
	RegionResolver func(ctx context.Context, appID string) (Region, error)
}

// NewWrapperClient creates a new WrapperClient with automatic region selection.
// It determines the appropriate region based on the provided app ID and sets
// the base URL accordingly.
//
// Parameters:
//   - ctx: Context for the region detection request (if needed)
//   - appID: The Agora application ID
//   - opts: Optional request options (authentication, custom HTTP client, etc.)
//
// Returns an error if region detection fails.
func NewWrapperClient(ctx context.Context, appID string, opts ...option.RequestOption) (*WrapperClient, error) {
	config := WrapperClientConfig{
		AppID: appID,
	}
	return NewWrapperClientWithConfig(ctx, config, opts...)
}

// NewWrapperClientWithConfig creates a new WrapperClient with the given configuration.
// This allows for more advanced configuration options like explicit region setting
// or custom region resolvers.
//
// Example with explicit region:
//
//	config := Agora.WrapperClientConfig{
//	    AppID:  "your-app-id",
//	    Region: &Agora.RegionEU,
//	}
//	client, err := Agora.NewWrapperClientWithConfig(ctx, config,
//	    option.WithBasicAuth("username", "password"),
//	)
//
// Example with custom resolver:
//
//	config := Agora.WrapperClientConfig{
//	    AppID: "your-app-id",
//	    RegionResolver: func(ctx context.Context, appID string) (Agora.Region, error) {
//	        // Custom logic to determine region
//	        return Agora.RegionUS, nil
//	    },
//	}
//	client, err := Agora.NewWrapperClientWithConfig(ctx, config, opts...)
func NewWrapperClientWithConfig(ctx context.Context, config WrapperClientConfig, opts ...option.RequestOption) (*WrapperClient, error) {
	var region Region
	var err error

	if config.Region != nil {
		region = *config.Region
	} else if config.RegionResolver != nil {
		region, err = config.RegionResolver(ctx, config.AppID)
		if err != nil {
			return nil, err
		}
	} else {
		region, err = determineRegion(ctx, config.AppID)
		if err != nil {
			return nil, err
		}
	}

	if !IsValidRegion(region) {
		return nil, errors.New("invalid region: " + string(region))
	}

	baseURL := BaseURLForRegion(region)
	opts = append([]option.RequestOption{option.WithBaseURL(baseURL)}, opts...)

	return &WrapperClient{
		Client: client.NewClient(opts...),
		Region: region,
	}, nil
}

// NewWrapperClientWithRegion creates a new WrapperClient with an explicitly specified region.
// This is a convenience function for when you know the region ahead of time
// and don't need auto-detection.
//
// Example:
//
//	client, err := Agora.NewWrapperClientWithRegion(Agora.RegionEU,
//	    option.WithBasicAuth("username", "password"),
//	)
func NewWrapperClientWithRegion(region Region, opts ...option.RequestOption) (*WrapperClient, error) {
	if !IsValidRegion(region) {
		return nil, errors.New("invalid region: " + string(region))
	}

	baseURL := BaseURLForRegion(region)
	opts = append([]option.RequestOption{option.WithBaseURL(baseURL)}, opts...)

	return &WrapperClient{
		Client: client.NewClient(opts...),
		Region: region,
	}, nil
}

// determineRegion determines the appropriate region for the given app ID.
// TODO: Implement actual region detection logic based on Agora's requirements.
// Possible implementations:
//   - Call a region discovery endpoint
//   - Parse app ID prefix for region hints
//   - Use account metadata
//   - Geographic detection based on caller location
func determineRegion(ctx context.Context, appID string) (Region, error) {
	if appID == "" {
		return "", errors.New("appID is required for region detection")
	}

	// TODO: Replace this stub with actual region detection logic.
	// For now, default to US region.
	//
	// Example implementation patterns:
	//
	// 1. Region discovery endpoint:
	//    resp, err := http.Get("https://api.agora.io/region?app_id=" + appID)
	//    // parse response and return region
	//
	// 2. App ID prefix parsing:
	//    if strings.HasPrefix(appID, "eu-") {
	//        return RegionEU, nil
	//    }
	//
	// 3. Account metadata lookup:
	//    account, err := fetchAccountInfo(ctx, appID)
	//    return account.Region, nil

	return RegionUS, nil
}
