package domain

import (
	"context"
	"testing"
	"time"
)

func TestNewDomainSwitcher(t *testing.T) {
	tests := []struct {
		name    string
		region  Region
		wantErr bool
	}{
		{
			name:    "US region",
			region:  RegionUS,
			wantErr: false,
		},
		{
			name:    "EU region",
			region:  RegionEU,
			wantErr: false,
		},
		{
			name:    "AP region",
			region:  RegionAP,
			wantErr: false,
		},
		{
			name:    "CN region",
			region:  RegionCN,
			wantErr: false,
		},
		{
			name:    "Unknown region",
			region:  RegionUnknown,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds, err := NewDomainSwitcher(tt.region)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewDomainSwitcher() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && ds == nil {
				t.Errorf("NewDomainSwitcher() returned nil for valid region")
			}
		})
	}
}

func TestGetBaseURL(t *testing.T) {
	tests := []struct {
		name   string
		region Region
		want   string
	}{
		{
			name:   "US region",
			region: RegionUS,
			want:   "https://api-us-west-1.agora.io/api/conversational-ai-agent",
		},
		{
			name:   "EU region",
			region: RegionEU,
			want:   "https://api-eu-west-1.agora.io/api/conversational-ai-agent",
		},
		{
			name:   "AP region",
			region: RegionAP,
			want:   "https://api-ap-southeast-1.agora.io/api/conversational-ai-agent",
		},
		{
			name:   "CN region",
			region: RegionCN,
			want:   "https://api-cn-east-1.sd-rtn.com/api/conversational-ai-agent",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ds, err := NewDomainSwitcher(tt.region)
			if err != nil {
				t.Fatalf("NewDomainSwitcher() error = %v", err)
			}
			got := ds.GetBaseURL()
			if got != tt.want {
				t.Errorf("GetBaseURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetBaseURLForRegion(t *testing.T) {
	tests := []struct {
		name    string
		region  Region
		want    string
		wantErr bool
	}{
		{
			name:    "US region",
			region:  RegionUS,
			want:    "https://api-us-west-1.agora.io/api/conversational-ai-agent",
			wantErr: false,
		},
		{
			name:    "Unknown region",
			region:  RegionUnknown,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBaseURLForRegion(tt.region)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBaseURLForRegion() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetBaseURLForRegion() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNextRegion(t *testing.T) {
	ds, err := NewDomainSwitcher(RegionUS)
	if err != nil {
		t.Fatalf("NewDomainSwitcher() error = %v", err)
	}

	firstURL := ds.GetBaseURL()
	if firstURL != "https://api-us-west-1.agora.io/api/conversational-ai-agent" {
		t.Errorf("Initial URL = %v, want https://api-us-west-1.agora.io/api/conversational-ai-agent", firstURL)
	}

	ds.NextRegion()
	secondURL := ds.GetBaseURL()
	if secondURL != "https://api-us-east-1.agora.io/api/conversational-ai-agent" {
		t.Errorf("After NextRegion() URL = %v, want https://api-us-east-1.agora.io/api/conversational-ai-agent", secondURL)
	}

	ds.NextRegion()
	thirdURL := ds.GetBaseURL()
	if thirdURL != "https://api-us-west-1.agora.io/api/conversational-ai-agent" {
		t.Errorf("After second NextRegion() URL = %v, want https://api-us-west-1.agora.io/api/conversational-ai-agent (should wrap around)", thirdURL)
	}
}

func TestSelectBestDomain(t *testing.T) {
	ds, err := NewDomainSwitcher(RegionUS)
	if err != nil {
		t.Fatalf("NewDomainSwitcher() error = %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ds.SelectBestDomain(ctx)
	if err != nil {
		t.Logf("SelectBestDomain() error = %v (this may be expected in test environments)", err)
	}

	baseURL := ds.GetBaseURL()
	if baseURL == "" {
		t.Errorf("GetBaseURL() returned empty string")
	}
}

func TestNeedsUpdate(t *testing.T) {
	ds, err := NewDomainSwitcher(RegionUS)
	if err != nil {
		t.Fatalf("NewDomainSwitcher() error = %v", err)
	}

	if !ds.needsUpdate() {
		t.Errorf("needsUpdate() = false, want true for new DomainSwitcher")
	}

	ds.lastUpdate = time.Now()
	if ds.needsUpdate() {
		t.Errorf("needsUpdate() = true, want false after recent update")
	}

	ds.lastUpdate = time.Now().Add(-31 * time.Second)
	if !ds.needsUpdate() {
		t.Errorf("needsUpdate() = false, want true after 31 seconds")
	}
}
