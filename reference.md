# Reference
## Agent Management
<details><summary><code>client.AgentManagement.StartAgent(Appid, request) -> *Agora.StartAgentResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Create and start a Conversational AI agent instance.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &Agora.StartAgentRequest{
        Name: "unique_name",
        Properties: &Agora.StartAgentRequestProperties{
            Channel: "channel_name",
            Token: "token",
            AgentRtcUID: "1001",
            RemoteRtcUIDs: []string{
                "1002",
            },
            IdleTimeout: Agora.Int(
                120,
            ),
            AdvancedFeatures: &Agora.StartAgentRequestPropertiesAdvancedFeatures{
                EnableAivad: Agora.Bool(
                    true,
                ),
            },
            Asr: &Agora.StartAgentRequestPropertiesAsr{
                Language: Agora.String(
                    "en-US",
                ),
            },
            Tts: &Agora.StartAgentRequestPropertiesTts{
                Vendor: Agora.StartAgentRequestPropertiesTtsVendorMicrosoft,
                Params: map[string]any{
                    "key": "<your_tts_api_key>",
                    "region": "eastus",
                    "voice_name": "en-US-AndrewMultilingualNeural",
                },
            },
            Llm: &Agora.StartAgentRequestPropertiesLlm{
                URL: "https://api.openai.com/v1/chat/completions",
                APIKey: Agora.String(
                    "<your_llm_key>",
                ),
                SystemMessages: []map[string]any{
                    map[string]any{
                        "role": "system",
                        "content": "You are a helpful chatbot.",
                    },
                },
                Params: map[string]any{
                    "model": "gpt-4o-mini",
                },
                MaxHistory: Agora.Int(
                    32,
                ),
                GreetingMessage: Agora.String(
                    "Hello, how can I assist you today?",
                ),
                FailureMessage: Agora.String(
                    "Please hold on a second.",
                ),
            },
        },
    }
client.AgentManagement.StartAgent(
        context.TODO(),
        "appid",
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**appid:** `string` — The App ID of the project.
    
</dd>
</dl>

<dl>
<dd>

**name:** `string` — The unique identifier of the agent. The same identifier cannot be used repeatedly.
    
</dd>
</dl>

<dl>
<dd>

**properties:** `*Agora.StartAgentRequestProperties` — Configuration details of the agent.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.AgentManagement.ListAgents(Appid) -> *Agora.ListAgentsResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Retrieve a list of agents that meet the specified conditions.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &Agora.ListAgentsRequest{
        Channel: Agora.String(
            "channel",
        ),
        FromTime: Agora.Float64(
            1.1,
        ),
        ToTime: Agora.Float64(
            1.1,
        ),
        State: Agora.ListAgentsRequestStateZero.Ptr(),
        Limit: Agora.Int(
            1,
        ),
        Cursor: Agora.String(
            "cursor",
        ),
    }
client.AgentManagement.ListAgents(
        context.TODO(),
        "appid",
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**appid:** `string` — The App ID of the project.
    
</dd>
</dl>

<dl>
<dd>

**channel:** `*string` — The channel to query for a list of agents.
    
</dd>
</dl>

<dl>
<dd>

**fromTime:** `*float64` — The start timestamp (in seconds) for the query. Default is 2 hours ago.
    
</dd>
</dl>

<dl>
<dd>

**toTime:** `*float64` — The end timestamp (in seconds) for the query. Default is current time.
    
</dd>
</dl>

<dl>
<dd>

**state:** `*Agora.ListAgentsRequestState` 

The agent state to filter by. Only one state can be specified per query:
- `IDLE` (0): Agent is idle.
- `STARTING` (1): The agent is being started.
- `RUNNING` (2): The agent is running.
- `STOPPING` (3): The agent is stopping.
- `STOPPED` (4): The agent has exited.
- `RECOVERING` (5): The agent is recovering.
- `FAILED` (6): The agent failed to execute.
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int` — The maximum number of entries returned per page.
    
</dd>
</dl>

<dl>
<dd>

**cursor:** `*string` — The paging cursor, indicating the starting position (`agent_id`) of the next page of results.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.AgentManagement.QueryAgentStatus(Appid, AgentID) -> *Agora.QueryAgentStatusResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Get the current state information of the specified agent instance.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
client.AgentManagement.QueryAgentStatus(
        context.TODO(),
        "appid",
        "agentId",
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**appid:** `string` — The App ID of the project.
    
</dd>
</dl>

<dl>
<dd>

**agentID:** `string` — The agent instance ID you obtained after successfully calling `join` to start a conversational AI agent.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.AgentManagement.GetAgentHistory(Appid, AgentID) -> *Agora.GetAgentHistoryResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Get the history of the conversation between the user and the agent.

Call this endpoint while the agent is running to retrieve the conversation history. You can set the maximum number of cached entries using the `llm.max_history` parameter when calling the start agent endpoint. The default value is `32`.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
client.AgentManagement.GetAgentHistory(
        context.TODO(),
        "appid",
        "agentId",
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**appid:** `string` — The App ID of the project.
    
</dd>
</dl>

<dl>
<dd>

**agentID:** `string` — The agent instance ID you obtained after successfully calling `join` to start a conversational AI agent.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.AgentManagement.StopAgent(Appid, AgentID) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Stop the specified conversational agent instance.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
client.AgentManagement.StopAgent(
        context.TODO(),
        "appid",
        "agentId",
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**appid:** `string` — The App ID of the project.
    
</dd>
</dl>

<dl>
<dd>

**agentID:** `string` — The agent instance ID you obtained after successfully calling `join` to start a conversational AI agent.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.AgentManagement.UpdateAgent(Appid, AgentID, request) -> *Agora.UpdateAgentResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Adjust Conversation AI Engine parameters at runtime.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &Agora.UpdateAgentRequest{
        Properties: &Agora.UpdateAgentRequestProperties{
            Token: Agora.String(
                "007eJxTYxxxxxxxxxxIaHMLAAAA0ex66",
            ),
            Llm: &Agora.UpdateAgentRequestPropertiesLlm{
                SystemMessages: []map[string]any{
                    map[string]any{
                        "role": "system",
                        "content": "You are a helpful assistant. xxx",
                    },
                    map[string]any{
                        "role": "system",
                        "content": "Previously, user has talked about their favorite hobbies with some key topics: xxx",
                    },
                },
                Params: map[string]any{
                    "model": "abab6.5s-chat",
                    "max_token": 1024,
                },
            },
        },
    }
client.AgentManagement.UpdateAgent(
        context.TODO(),
        "appid",
        "agentId",
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**appid:** `string` — The App ID of the project.
    
</dd>
</dl>

<dl>
<dd>

**agentID:** `string` — The agent instance ID you obtained after successfully calling `join` to start a conversational AI agent.
    
</dd>
</dl>

<dl>
<dd>

**properties:** `*Agora.UpdateAgentRequestProperties` — Configuration properties to update.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.AgentManagement.AgentSpeak(Appid, AgentID, request) -> *Agora.AgentSpeakResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Broadcast a custom message using the TTS module.

During a conversation with an agent, call this endpoint to immediately broadcast a custom message using the TTS module. Upon receiving the request, the system interrupts the agent's speech and thought process to deliver the message. This broadcast can be interrupted by human voice.

Note: The speak API is not supported when using `mllm` configuration.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &Agora.AgentSpeakRequest{
        Text: "Sorry, the conversation content is not compliant.",
        Priority: Agora.AgentSpeakRequestPriorityInterrupt.Ptr(),
        Interruptable: Agora.Bool(
            false,
        ),
    }
client.AgentManagement.AgentSpeak(
        context.TODO(),
        "appid",
        "agentId",
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**appid:** `string` — The App ID of the project.
    
</dd>
</dl>

<dl>
<dd>

**agentID:** `string` — The agent instance ID you obtained after successfully calling `join` to start a conversational AI agent.
    
</dd>
</dl>

<dl>
<dd>

**text:** `string` — The broadcast message text. The maximum length of the text content is 512 bytes.
    
</dd>
</dl>

<dl>
<dd>

**priority:** `*Agora.AgentSpeakRequestPriority` 

Sets the priority of the message broadcast:
- `INTERRUPT`: High priority. The agent immediately interrupts the current interaction to announce the message.
- `APPEND`: Medium priority. The agent announces the message after the current interaction ends.
- `IGNORE`: Low priority. If the agent is busy interacting, it ignores and discards the broadcast; the message is only announced if the agent is not interacting.
    
</dd>
</dl>

<dl>
<dd>

**interruptable:** `*bool` 

Whether to allow users to interrupt the agent's broadcast by speaking:
- `true`: Allow
- `false`: Don't allow
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.AgentManagement.AgentInterrupt(Appid, AgentID, request) -> *Agora.AgentInterruptResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Interrupt the specified agent while speaking or thinking.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &Agora.AgentInterruptRequest{}
client.AgentManagement.AgentInterrupt(
        context.TODO(),
        "appid",
        "agentId",
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**appid:** `string` — The App ID of the project.
    
</dd>
</dl>

<dl>
<dd>

**agentID:** `string` — The agent instance ID you obtained after successfully calling `join` to start a conversational AI agent.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Telephony
<details><summary><code>client.Telephony.RetrieveCallRecords(Appid) -> *Agora.RetrieveCallRecordsResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Query historical call records for a specified appid based on the filter criteria.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &Agora.RetrieveCallRecordsRequest{
        Number: Agora.String(
            "number",
        ),
        FromTime: Agora.Int(
            1,
        ),
        ToTime: Agora.Int(
            1,
        ),
        Type: Agora.RetrieveCallRecordsRequestTypeInbound.Ptr(),
        Limit: Agora.Int(
            1,
        ),
        Cursor: Agora.String(
            "cursor",
        ),
    }
client.Telephony.RetrieveCallRecords(
        context.TODO(),
        "appid",
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**appid:** `string` — The App ID of the project.
    
</dd>
</dl>

<dl>
<dd>

**number:** `*string` — Filter by phone number. Can be either the calling number or the called number.
    
</dd>
</dl>

<dl>
<dd>

**fromTime:** `*int` — Query list start timestamp (in seconds). Default is 60 days ago.
    
</dd>
</dl>

<dl>
<dd>

**toTime:** `*int` — Query list end timestamp (in seconds). Default is current time.
    
</dd>
</dl>

<dl>
<dd>

**type_:** `*Agora.RetrieveCallRecordsRequestType` 

Call type filter:
- `inbound`: Inbound call.
- `outbound`: Outbound call.

If not specified, all call types are returned.
    
</dd>
</dl>

<dl>
<dd>

**limit:** `*int` — Maximum number of items returned in a single page.
    
</dd>
</dl>

<dl>
<dd>

**cursor:** `*string` — Pagination cursor. Use the `agent_id` from the previous page as the cursor for the next page.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Telephony.InitiateOutboundCall(Appid, request) -> *Agora.InitiateOutboundCallResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Initiate an outbound call to a specified number and create an agent to join the specified RTC channel.

Use this endpoint to initiate an outbound call to the specified number and create an agent that joins the target RTC channel. The agent waits for the callee to answer.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &Agora.InitiateOutboundCallRequest{
        Name: "customer_service",
        Sip: &Agora.InitiateOutboundCallRequestSip{
            ToNumber: "+19876543210",
            FromNumber: "+11234567890",
            SipRtcUID: "100",
            SipRtcToken: "<agora_sip_rtc_token>",
        },
        PipelineID: Agora.String(
            "fzufjlweixxxxnlp",
        ),
        Properties: &Agora.InitiateOutboundCallRequestProperties{
            Channel: "<agora_channel>",
            Token: "<agora_channel_token>",
            AgentRtcUID: "111",
        },
    }
client.Telephony.InitiateOutboundCall(
        context.TODO(),
        "appid",
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**appid:** `string` — The App ID of the project.
    
</dd>
</dl>

<dl>
<dd>

**name:** `string` — The name identifier of the call session.
    
</dd>
</dl>

<dl>
<dd>

**sip:** `*Agora.InitiateOutboundCallRequestSip` — SIP (Session Initiation Protocol) call configuration object.
    
</dd>
</dl>

<dl>
<dd>

**pipelineID:** `*string` — The unique ID of a published project in AI Studio.
    
</dd>
</dl>

<dl>
<dd>

**properties:** `*Agora.InitiateOutboundCallRequestProperties` 

Call attribute configuration. The content of this field varies depending on the invocation method:
- **Using pipeline ID**: Simply pass in `channel`, `token`, and `agent_rtc_uid`.
- **Using complete configuration**: Pass in the complete parameters of the [Start a conversational AI agent](https://docs.agora.io/en/conversational-ai/rest-api/agent/join) `properties`, including all required fields such as `channel`, `token`, `agent_rtc_uid`, `remote_rtc_uids`, `tts`, and `llm`.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Telephony.RetrieveCallStatus(Appid, AgentID) -> *Agora.RetrieveCallStatusResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Retrieve the call status and related information of a specified agent.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
client.Telephony.RetrieveCallStatus(
        context.TODO(),
        "appid",
        "agent_id",
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**appid:** `string` — The App ID of the project.
    
</dd>
</dl>

<dl>
<dd>

**agentID:** `string` — The agent ID you obtained after successfully calling the API to initiate an outbound call.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.Telephony.HangupCall(Appid, AgentID, request) -> *Agora.HangupCallResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Instruct the agent to proactively hang up the ongoing call and leave the RTC channel.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &Agora.HangupCallRequest{}
client.Telephony.HangupCall(
        context.TODO(),
        "appid",
        "agent_id",
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**appid:** `string` — The App ID of the project.
    
</dd>
</dl>

<dl>
<dd>

**agentID:** `string` — The agent ID you obtained after successfully calling the API to initiate an outbound call.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

## Phone Number Management
<details><summary><code>client.PhoneNumberManagement.RetrieveNumberList() -> []*Agora.RetrieveNumberListResponseItem</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Retrieve a list of all imported phone numbers under the current account.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
client.PhoneNumberManagement.RetrieveNumberList(
        context.TODO(),
    )
}
```
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.PhoneNumberManagement.ImportNumber(request) -> *Agora.ImportNumberResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Import a pre-configured phone number that can be used for inbound or outbound calls.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &Agora.ImportNumberRequest{
        Provider: Agora.ImportNumberRequestProviderByo,
        PhoneNumber: "+19876543210",
        Label: "Sales Hotline",
        Inbound: Agora.Bool(
            true,
        ),
        Outbound: Agora.Bool(
            true,
        ),
        InboundConfig: &Agora.ImportNumberRequestInboundConfig{
            AllowedAddresses: []string{
                "112.126.15.64/27",
            },
        },
        OutboundConfig: &Agora.ImportNumberRequestOutboundConfig{
            Address: Agora.String(
                "xxx:xxx@sip.example.com",
            ),
            Transport: Agora.String(
                "tls",
            ),
        },
    }
client.PhoneNumberManagement.ImportNumber(
        context.TODO(),
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**provider:** `*Agora.ImportNumberRequestProvider` 

Number provider:
- `byo`: BYO (Bring Your Own)
- `twilio`: Twilio
    
</dd>
</dl>

<dl>
<dd>

**phoneNumber:** `string` — Telephone number in E.164 format.
    
</dd>
</dl>

<dl>
<dd>

**label:** `string` — A label used to identify the number.
    
</dd>
</dl>

<dl>
<dd>

**inbound:** `*bool` — Whether the number supports inbound calls.
    
</dd>
</dl>

<dl>
<dd>

**outbound:** `*bool` — Whether the number supports outbound calls.
    
</dd>
</dl>

<dl>
<dd>

**inboundConfig:** `*Agora.ImportNumberRequestInboundConfig` — SIP inbound call configuration.
    
</dd>
</dl>

<dl>
<dd>

**outboundConfig:** `*Agora.ImportNumberRequestOutboundConfig` — SIP outbound call configuration.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.PhoneNumberManagement.RetrieveNumberInformation(PhoneNumber) -> *Agora.RetrieveNumberInformationResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Retrieve detailed information for a specific phone number.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
client.PhoneNumberManagement.RetrieveNumberInformation(
        context.TODO(),
        "phone_number",
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**phoneNumber:** `string` — Telephone number in E.164 format. For example, +11234567890.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.PhoneNumberManagement.DeleteNumber(PhoneNumber) -> error</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Remove an imported phone number from the system.

Note: This operation only removes the number configuration from the Agora system; the number stored with the phone service provider is not deleted. After calling this endpoint, the number stops receiving calls routed through this system. To delete the number from the service provider, remove it in the service provider's console.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
client.PhoneNumberManagement.DeleteNumber(
        context.TODO(),
        "phone_number",
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**phoneNumber:** `string` — Telephone number in E.164 format. For example, +11234567890.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>

<details><summary><code>client.PhoneNumberManagement.UpdateNumberConfiguration(PhoneNumber, request) -> *Agora.UpdateNumberConfigurationResponse</code></summary>
<dl>
<dd>

#### 📝 Description

<dl>
<dd>

<dl>
<dd>

Update the configuration for a phone number.
</dd>
</dl>
</dd>
</dl>

#### 🔌 Usage

<dl>
<dd>

<dl>
<dd>

```go
request := &Agora.UpdateNumberConfigurationRequest{
        InboundConfig: &Agora.UpdateNumberConfigurationRequestInboundConfig{
            PipelineID: Agora.String(
                "xxxxx",
            ),
        },
        OutboundConfig: &Agora.UpdateNumberConfigurationRequestOutboundConfig{
            PipelineID: Agora.String(
                "xxxxx",
            ),
        },
    }
client.PhoneNumberManagement.UpdateNumberConfiguration(
        context.TODO(),
        "phone_number",
        request,
    )
}
```
</dd>
</dl>
</dd>
</dl>

#### ⚙️ Parameters

<dl>
<dd>

<dl>
<dd>

**phoneNumber:** `string` — Telephone number in E.164 format. For example, +11234567890.
    
</dd>
</dl>

<dl>
<dd>

**inboundConfig:** `*Agora.UpdateNumberConfigurationRequestInboundConfig` — Update inbound call configuration. Passing `null` will clear the configuration.
    
</dd>
</dl>

<dl>
<dd>

**outboundConfig:** `*Agora.UpdateNumberConfigurationRequestOutboundConfig` — Update outbound call configuration. Passing `null` will clear the configuration.
    
</dd>
</dl>
</dd>
</dl>


</dd>
</dl>
</details>
