package agentkit

import (
	"context"

	Agora "github.com/AgoraIO/agora-agents-go"
	"github.com/AgoraIO/agora-agents-go/agents"
	"github.com/AgoraIO/agora-agents-go/option"
)

func startAgentsWithMapBody(
	ctx context.Context,
	client *agents.Client,
	appID string,
	name string,
	preset string,
	pipelineID string,
	properties map[string]interface{},
	opts ...option.RequestOption,
) (*Agora.StartAgentsResponse, error) {
	body := map[string]interface{}{
		"name":       name,
		"properties": properties,
	}
	if preset != "" {
		body["preset"] = preset
	}
	if pipelineID != "" {
		body["pipeline_id"] = pipelineID
	}
	return agents.StartWithMapBody(ctx, client, appID, body, opts...)
}
