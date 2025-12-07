package domain

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

const (
	ChineseMainlandMajorDomain = "sd-rtn.com"
	OverseaMajorDomain         = "agora.io"
)

const GlobalDomainPrefix = "api"

const (
	USWestRegionDomainPrefix = "api-us-west-1"
	USEastRegionDomainPrefix = "api-us-east-1"
)

const (
	APSoutheastRegionDomainPrefix = "api-ap-southeast-1"
	APNortheastRegionDomainPrefix = "api-ap-northeast-1"
)

const (
	EUWestRegionDomainPrefix    = "api-eu-west-1"
	EUCentralRegionDomainPrefix = "api-eu-central-1"
)

const (
	CNEastRegionDomainPrefix  = "api-cn-east-1"
	CNNorthRegionDomainPrefix = "api-cn-north-1"
)

// Domain contains the region-specific domain prefixes and major domain suffixes.
type Domain struct {
	RegionDomainPrefixes []string
	MajorDomainSuffixes  []string
}

// RegionDomain maps each Area to its corresponding Domain configuration.
var RegionDomain = map[Area]Domain{
	US: {
		RegionDomainPrefixes: []string{
			USWestRegionDomainPrefix,
			USEastRegionDomainPrefix,
		},
		MajorDomainSuffixes: []string{
			OverseaMajorDomain,
			ChineseMainlandMajorDomain,
		},
	},
	EU: {
		RegionDomainPrefixes: []string{
			EUWestRegionDomainPrefix,
			EUCentralRegionDomainPrefix,
		},
		MajorDomainSuffixes: []string{
			OverseaMajorDomain,
			ChineseMainlandMajorDomain,
		},
	},
	AP: {
		RegionDomainPrefixes: []string{
			APSoutheastRegionDomainPrefix,
			APNortheastRegionDomainPrefix,
		},
		MajorDomainSuffixes: []string{
			OverseaMajorDomain,
			ChineseMainlandMajorDomain,
		},
	},
	CN: {
		RegionDomainPrefixes: []string{
			CNEastRegionDomainPrefix,
			CNNorthRegionDomainPrefix,
		},
		MajorDomainSuffixes: []string{
			ChineseMainlandMajorDomain,
			OverseaMajorDomain,
		},
	},
}

// Pool manages domain selection and failover for a specific region.
type Pool struct {
	domainArea            Area
	domainSuffixes        []string
	currentDomain         string
	regionPrefixes        []string
	currentRegionPrefixes []string
	locker                *sync.Mutex

	resolver   Resolver
	lastUpdate time.Time
}

// NewPool creates a new domain pool for the specified area.
func NewPool(domainArea Area) (*Pool, error) {
	if _, ok := RegionDomain[domainArea]; !ok {
		return nil, errors.New("invalid domain area")
	}
	d := &Pool{
		domainArea:     domainArea,
		domainSuffixes: RegionDomain[domainArea].MajorDomainSuffixes,
		resolver:       newResolverImpl(),
		locker:         &sync.Mutex{},
	}

	d.regionPrefixes = append(d.regionPrefixes, RegionDomain[domainArea].RegionDomainPrefixes...)

	d.currentRegionPrefixes = d.regionPrefixes
	d.currentDomain = d.domainSuffixes[0]

	return d, nil
}

const updateDuration = 30 * time.Second

func (d *Pool) domainNeedUpdate() bool {
	return time.Since(d.lastUpdate) > updateDuration
}

// SelectBestDomain performs DNS resolution to select the best available domain.
func (d *Pool) SelectBestDomain(ctx context.Context) error {
	if !d.domainNeedUpdate() {
		return nil
	}

	d.locker.Lock()
	defer d.locker.Unlock()

	if d.domainNeedUpdate() {
		domain, err := d.resolver.Resolve(ctx, d.domainSuffixes, d.currentRegionPrefixes[0])
		if err != nil {
			return err
		}
		d.selectDomain(domain)
	}
	return nil
}

// NextRegion rotates to the next region prefix for failover.
func (d *Pool) NextRegion() {
	d.locker.Lock()
	defer d.locker.Unlock()

	d.currentRegionPrefixes = d.currentRegionPrefixes[1:]
	if len(d.currentRegionPrefixes) == 0 {
		d.currentRegionPrefixes = d.regionPrefixes
	}
}

func (d *Pool) selectDomain(domain string) {
	if contains(d.domainSuffixes, domain) {
		d.currentDomain = domain
		d.lastUpdate = time.Now()
	}
}

// GetCurrentUrl returns the current base URL for API requests.
func (d *Pool) GetCurrentUrl() string {
	d.locker.Lock()
	defer d.locker.Unlock()

	currentRegion := d.currentRegionPrefixes[0]
	currentDomain := d.currentDomain

	return fmt.Sprintf("https://%s.%s", currentRegion, currentDomain)
}

// contains checks if a string slice contains a specific string.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
