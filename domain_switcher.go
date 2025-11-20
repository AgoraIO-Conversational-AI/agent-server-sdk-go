package Agora

import (
	"context"
	"errors"
	"fmt"
	"net"
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

type Region int

const (
	RegionUnknown Region = iota
	RegionUS
	RegionEU
	RegionAP
	RegionCN
)

type DomainConfig struct {
	RegionDomainPrefixes []string
	MajorDomainSuffixes  []string
}

var regionDomainMap = map[Region]DomainConfig{
	RegionUS: {
		RegionDomainPrefixes: []string{
			USWestRegionDomainPrefix,
			USEastRegionDomainPrefix,
		},
		MajorDomainSuffixes: []string{
			OverseaMajorDomain,
			ChineseMainlandMajorDomain,
		},
	},
	RegionEU: {
		RegionDomainPrefixes: []string{
			EUWestRegionDomainPrefix,
			EUCentralRegionDomainPrefix,
		},
		MajorDomainSuffixes: []string{
			OverseaMajorDomain,
			ChineseMainlandMajorDomain,
		},
	},
	RegionAP: {
		RegionDomainPrefixes: []string{
			APSoutheastRegionDomainPrefix,
			APNortheastRegionDomainPrefix,
		},
		MajorDomainSuffixes: []string{
			OverseaMajorDomain,
			ChineseMainlandMajorDomain,
		},
	},
	RegionCN: {
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

type DomainSwitcher struct {
	region                Region
	domainSuffixes        []string
	currentDomain         string
	regionPrefixes        []string
	currentRegionPrefixes []string
	locker                *sync.Mutex
	lastUpdate            time.Time
}

func NewDomainSwitcher(region Region) (*DomainSwitcher, error) {
	config, ok := regionDomainMap[region]
	if !ok {
		return nil, errors.New("invalid region")
	}

	ds := &DomainSwitcher{
		region:         region,
		domainSuffixes: config.MajorDomainSuffixes,
		locker:         &sync.Mutex{},
	}

	ds.regionPrefixes = append(ds.regionPrefixes, config.RegionDomainPrefixes...)
	ds.currentRegionPrefixes = ds.regionPrefixes
	ds.currentDomain = ds.domainSuffixes[0]

	return ds, nil
}

const updateDuration = 30 * time.Second

func (ds *DomainSwitcher) needsUpdate() bool {
	return time.Since(ds.lastUpdate) > updateDuration
}

func (ds *DomainSwitcher) resolveBestDomain(ctx context.Context) (string, error) {
	var wg sync.WaitGroup

	done := make(chan struct{}, 1)
	res := make(chan string, len(ds.domainSuffixes))

	for _, domain := range ds.domainSuffixes {
		wg.Add(1)

		go func(domain string, regionPrefix string) {
			defer wg.Done()
			url := regionPrefix + "." + domain
			_, err := net.LookupHost(url)
			if err == nil {
				res <- domain
			}
		}(domain, ds.currentRegionPrefixes[0])
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case domain := <-res:
		return domain, nil
	case <-done:
	}
	return "", errors.New("all DNS lookups failed")
}

func (ds *DomainSwitcher) SelectBestDomain(ctx context.Context) error {
	if !ds.needsUpdate() {
		return nil
	}

	ds.locker.Lock()
	defer ds.locker.Unlock()

	if ds.needsUpdate() {
		domain, err := ds.resolveBestDomain(ctx)
		if err != nil {
			return err
		}
		ds.selectDomain(domain)
	}
	return nil
}

func (ds *DomainSwitcher) NextRegion() {
	ds.locker.Lock()
	defer ds.locker.Unlock()

	ds.currentRegionPrefixes = ds.currentRegionPrefixes[1:]
	if len(ds.currentRegionPrefixes) == 0 {
		ds.currentRegionPrefixes = ds.regionPrefixes
	}
}

func (ds *DomainSwitcher) selectDomain(domain string) {
	for _, suffix := range ds.domainSuffixes {
		if suffix == domain {
			ds.currentDomain = domain
			ds.lastUpdate = time.Now()
			return
		}
	}
}

func (ds *DomainSwitcher) GetBaseURL() string {
	ds.locker.Lock()
	defer ds.locker.Unlock()

	currentRegion := ds.currentRegionPrefixes[0]
	currentDomain := ds.currentDomain

	return fmt.Sprintf("https://%s.%s", currentRegion, currentDomain)
}

func GetBaseURLForRegion(region Region) (string, error) {
	ds, err := NewDomainSwitcher(region)
	if err != nil {
		return "", err
	}
	return ds.GetBaseURL(), nil
}
