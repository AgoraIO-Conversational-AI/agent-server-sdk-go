package agents

import (
	"context"
	"net/http"

	Agora "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go"
	"github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/core"
	"github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/internal"
	"github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/option"
)

// StartWithMapBody starts an agent using an exact JSON body.
//
// AgentKit uses this for preset-backed sessions so provider-owned empty fields
// are not reintroduced by generated structs after preset resolution.
func StartWithMapBody(
	ctx context.Context,
	c *Client,
	appid string,
	body map[string]interface{},
	opts ...option.RequestOption,
) (*Agora.StartAgentsResponse, error) {
	response, err := c.WithRawResponse.StartWithMapBody(ctx, appid, body, opts...)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

// StartWithMapBody starts an agent using an exact JSON body and returns the raw response.
func (r *RawClient) StartWithMapBody(
	ctx context.Context,
	appid string,
	body map[string]interface{},
	opts ...option.RequestOption,
) (*core.Response[*Agora.StartAgentsResponse], error) {
	options := core.NewRequestOptions(opts...)
	baseURL := internal.ResolveBaseURL(
		options.BaseURL,
		r.baseURL,
		"https://api.agora.io/api/conversational-ai-agent",
	)
	endpointURL := internal.EncodeURL(
		baseURL+"/v2/projects/%v/join",
		appid,
	)
	headers := internal.MergeHeaders(
		r.options.ToHeader(),
		options.ToHeader(),
	)
	headers.Add("Content-Type", "application/json")
	var response *Agora.StartAgentsResponse
	raw, err := r.caller.Call(
		ctx,
		&internal.CallParams{
			URL:             endpointURL,
			Method:          http.MethodPost,
			Headers:         headers,
			MaxAttempts:     options.MaxAttempts,
			BodyProperties:  options.BodyProperties,
			QueryParameters: options.QueryParameters,
			Client:          options.HTTPClient,
			Request:         body,
			Response:        &response,
		},
	)
	if err != nil {
		return nil, err
	}
	return &core.Response[*Agora.StartAgentsResponse]{
		StatusCode: raw.StatusCode,
		Header:     raw.Header,
		Body:       response,
	}, nil
}
