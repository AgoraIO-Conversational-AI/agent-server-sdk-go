package agents

import (
	"context"
	"net/http"

	Agora "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go"
	"github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/core"
	"github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/internal"
	"github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/option"
)

// StartWithPropertiesMap starts an agent with a dynamic properties payload.
// AgentKit uses this for managed presets so provider-owned fields removed from
// the map do not reappear through generated zero-value structs.
func (c *Client) StartWithPropertiesMap(
	ctx context.Context,
	request *Agora.StartAgentsRequest,
	properties map[string]interface{},
	opts ...option.RequestOption,
) (*Agora.StartAgentsResponse, error) {
	response, err := c.WithRawResponse.StartWithPropertiesMap(ctx, request, properties, opts...)
	if err != nil {
		return nil, err
	}
	return response.Body, nil
}

func (r *RawClient) StartWithPropertiesMap(
	ctx context.Context,
	request *Agora.StartAgentsRequest,
	properties map[string]interface{},
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
		request.Appid,
	)
	headers := internal.MergeHeaders(
		r.options.ToHeader(),
		options.ToHeader(),
	)
	headers.Add("Content-Type", "application/json")

	body := map[string]interface{}{
		"name":       request.Name,
		"properties": properties,
	}
	if request.Preset != nil {
		body["preset"] = *request.Preset
	}
	if request.PipelineID != nil {
		body["pipeline_id"] = *request.PipelineID
	}

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
