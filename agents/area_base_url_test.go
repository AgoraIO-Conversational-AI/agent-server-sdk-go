package agents_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	Agora "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go"
	"github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/client"
	"github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type captureHTTPClient struct {
	lastRequest *http.Request
}

func (c *captureHTTPClient) Do(req *http.Request) (*http.Response, error) {
	c.lastRequest = req.Clone(req.Context())
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(`{}`)),
	}, nil
}

func TestWithAreaStartRequestIncludesConvoAIBasePath(t *testing.T) {
	httpClient := &captureHTTPClient{}
	c := client.NewClient(
		option.WithArea(option.AreaUS),
		option.WithHTTPClient(httpClient),
	)

	_, err := c.Agents.Start(context.Background(), &Agora.StartAgentsRequest{
		Appid: "appid",
		Name:  "agent",
	})
	require.NoError(t, err)
	require.NotNil(t, httpClient.lastRequest)
	assert.True(t, strings.HasPrefix(httpClient.lastRequest.URL.Path, "/api/conversational-ai-agent/v2/projects/appid"))
}

func TestWithBaseURLDoesNotAppendConvoAIPathAutomatically(t *testing.T) {
	httpClient := &captureHTTPClient{}
	c := client.NewClient(
		option.WithBaseURL("https://example.test/custom"),
		option.WithHTTPClient(httpClient),
	)

	_, err := c.Agents.Start(context.Background(), &Agora.StartAgentsRequest{
		Appid: "appid",
		Name:  "agent",
	})
	require.NoError(t, err)
	require.NotNil(t, httpClient.lastRequest)
	assert.Equal(t, "/custom/v2/projects/appid/join", httpClient.lastRequest.URL.Path)
}
