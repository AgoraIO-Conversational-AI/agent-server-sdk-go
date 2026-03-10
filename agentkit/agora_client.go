package agentkit

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	Agora "github.com/AgoraIO-Conversational-AI/agent-server-sdk-go"
	"github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/agents"
	"github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/client"
	"github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/core"
	"github.com/AgoraIO-Conversational-AI/agent-server-sdk-go/option"
)

// AuthMode represents the authentication mode for the Agora client.
type AuthMode string

const (
	// AuthModeAppCredentials uses AppID + AppCertificate to auto-generate a
	// short-lived ConvoAI token per API call — no Customer Secret required.
	AuthModeAppCredentials AuthMode = "app-credentials"
	// AuthModeBasic uses HTTP Basic Auth with a Customer ID and Customer Secret.
	AuthModeBasic AuthMode = "basic"
	// AuthModeToken uses a pre-built Agora token passed as the Authorization header.
	AuthModeToken AuthMode = "token"
)

// AgoraClientOptions contains configuration for creating an AgoraClient.
type AgoraClientOptions struct {
	// Area is the geographic region for the client (US, EU, AP, CN).
	Area option.Area
	// AppID is the Agora App ID, used as the path parameter and for token generation.
	AppID string
	// AppCertificate is used to sign tokens. Keep this value secret.
	AppCertificate string
	// CustomerID and CustomerSecret are used for Basic auth (from Agora Console).
	// Mutually exclusive with Token.
	CustomerID     string
	CustomerSecret string
	// Token is a pre-built raw Agora token for REST API authentication.
	// The SDK sets "Authorization: agora token=<Token>" automatically.
	// Mutually exclusive with CustomerID/CustomerSecret.
	Token string
}

// AgoraClient wraps the Fern-generated client with AppID, AppCertificate, and
// auth mode tracking—providing the same client API as the TypeScript and Python SDKs.
//
// Agents is the underlying agents REST client and is available for advanced usage.
type AgoraClient struct {
	Agents         *agents.Client
	AppID          string
	AppCertificate string
	AuthMode       AuthMode
}

// NewAgoraClient creates a new AgoraClient with the given options.
//
// App credentials mode (recommended) — no Customer Secret needed:
//
//	c := agentkit.NewAgoraClient(agentkit.AgoraClientOptions{
//	    Area:           option.AreaUS,
//	    AppID:          "your-app-id",
//	    AppCertificate: "your-app-certificate",
//	})
//
// Basic auth mode:
//
//	c := agentkit.NewAgoraClient(agentkit.AgoraClientOptions{
//	    Area:           option.AreaUS,
//	    AppID:          "your-app-id",
//	    AppCertificate: "your-app-certificate",
//	    CustomerID:     "your-customer-id",
//	    CustomerSecret: "your-customer-secret",
//	})
//
// Token auth mode (pre-built token):
//
//	c := agentkit.NewAgoraClient(agentkit.AgoraClientOptions{
//	    Area:           option.AreaUS,
//	    AppID:          "your-app-id",
//	    AppCertificate: "your-app-certificate",
//	    Token:          "your-raw-token",
//	})
func NewAgoraClient(opts AgoraClientOptions) *AgoraClient {
	reqOpts := []option.RequestOption{option.WithArea(opts.Area)}

	authMode := AuthModeAppCredentials
	if opts.CustomerID != "" {
		authMode = AuthModeBasic
		reqOpts = append(reqOpts, option.WithBasicAuth(opts.CustomerID, opts.CustomerSecret))
	} else if opts.Token != "" {
		authMode = AuthModeToken
		reqOpts = append(reqOpts, option.WithToken(opts.Token))
	}

	c := client.NewClient(reqOpts...)
	return &AgoraClient{
		Agents:         c.Agents,
		AppID:          opts.AppID,
		AppCertificate: opts.AppCertificate,
		AuthMode:       authMode,
	}
}

// StopAgent stops an agent by its ID without requiring an AgentSession reference.
//
// Use this when handling a stop request from your client app (e.g., an end-call
// button). Create a client with the same credentials used to start the agent and
// call this method.
//
// If the agent has already stopped (e.g., timed out or crashed), StopAgent
// returns nil rather than an error.
//
//	c := agentkit.NewAgoraClient(agentkit.AgoraClientOptions{
//	    Area:           option.AreaUS,
//	    AppID:          "your-app-id",
//	    AppCertificate: "your-app-certificate",
//	})
//	if err := c.StopAgent(ctx, agentID); err != nil {
//	    log.Printf("stop failed: %v", err)
//	}
func (c *AgoraClient) StopAgent(ctx context.Context, agentID string) error {
	var reqOpts []option.RequestOption
	if c.AuthMode == AuthModeAppCredentials && c.AppCertificate != "" {
		token, err := GenerateConvoAIToken(GenerateConvoAITokenOptions{
			AppID:          c.AppID,
			AppCertificate: c.AppCertificate,
			ChannelName:    "stop",
			Account:        agentID,
		})
		if err != nil {
			return fmt.Errorf("agentkit: failed to generate token for StopAgent: %w", err)
		}
		h := make(http.Header)
		h.Set("Authorization", "agora token="+token)
		reqOpts = []option.RequestOption{option.WithHTTPHeader(h)}
	}

	err := c.Agents.Stop(ctx, &Agora.StopAgentsRequest{
		Appid:   c.AppID,
		AgentID: agentID,
	}, reqOpts...)
	if err != nil {
		var apiErr *core.APIError
		if errors.As(err, &apiErr) && apiErr.StatusCode == http.StatusNotFound {
			return nil // Agent already stopped — treat as success
		}
		return err
	}
	return nil
}
