package main

import (
	context "context"
	"log"

	Agora "github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk"
	client "github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/client"
	option "github.com/AgoraIO-Conversational-AI/agora-agent-go-sdk/option"
)

func main() {
	// NOTE: copied from telephony/telephony_test/telephony_test.go

	client := client.NewClient(
		option.WithBasicAuth(
			"<omitted>",
			"<omitted>",
		),
		option.WithArea(option.AreaUS),
	)
	request := &Agora.CallTelephonyRequest{
		Appid: "appid",
		Name:  "customer_service",
		Sip: &Agora.CallTelephonyRequestSip{
			ToNumber:   "+19876543210",
			FromNumber: "+11234567890",
			RtcUID:     "100",
			RtcToken:   "<agora_sip_rtc_token>",
		},
		PipelineID: Agora.String(
			"fzufjlweixxxxnlp",
		),
		Properties: &Agora.CallTelephonyRequestProperties{
			Channel:     "<agora_channel>",
			Token:       "<agora_channel_token>",
			AgentRtcUID: "111",
		},
	}

	response, invocationErr := client.Telephony.Call(
		context.TODO(),
		request,
	)

	if invocationErr != nil {
		log.Fatalf("Error calling telephony: %v", invocationErr)
	}

	_ = response
}
