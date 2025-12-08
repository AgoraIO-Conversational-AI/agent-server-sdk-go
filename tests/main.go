package main

import (
	context "context"
	"fmt"
	"log"

	Agora "github.com/fern-demo/agoraio-go-sdk/v505"
	client "github.com/fern-demo/agoraio-go-sdk/v505/client"
	option "github.com/fern-demo/agoraio-go-sdk/v505/option"
)

func main() {
	// NOTE: copied from telephony/telephony_test/telephony_test.go
	client := client.NewClient(
		option.WithBasicAuth(
			"<omitted>",
			"<omitted>",
		),
		option.WithRegion(option.RegionUSWest),
	)
	request := &Agora.CallTelephonyRequest{
		Appid: "appid",
		Name:  "customer_service",
		Sip: &Agora.CallTelephonyRequestSip{
			ToNumber:    "+19876543210",
			FromNumber:  "+11234567890",
			SipRtcUID:   "100",
			SipRtcToken: "<agora_sip_rtc_token>",
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

	fmt.Printf("Call response: %v", response)
}
