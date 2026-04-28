package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPoolGetCurrentURLIncludesConvoAIPathForAllAreas(t *testing.T) {
	tests := []struct {
		name     string
		area     Area
		expected string
	}{
		{
			name:     "us",
			area:     AreaUS,
			expected: "https://api-us-west-1.agora.io/api/conversational-ai-agent",
		},
		{
			name:     "eu",
			area:     AreaEU,
			expected: "https://api-eu-west-1.agora.io/api/conversational-ai-agent",
		},
		{
			name:     "ap",
			area:     AreaAP,
			expected: "https://api-ap-southeast-1.agora.io/api/conversational-ai-agent",
		},
		{
			name:     "cn",
			area:     AreaCN,
			expected: "https://api-cn-east-1.sd-rtn.com/api/conversational-ai-agent",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pool, err := NewPool(tc.area)
			require.NoError(t, err)
			assert.Equal(t, tc.expected, pool.GetCurrentURL())
		})
	}
}

