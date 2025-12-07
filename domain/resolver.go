package domain

import (
	"context"
	"errors"
	"net"
	"sync"
)

// Resolver defines the interface for DNS resolution.
type Resolver interface {
	Resolve(ctx context.Context, domains []string, regionPrefix string) (string, error)
}

// ResolverFunc is a function type that implements the Resolver interface.
type ResolverFunc func(ctx context.Context, domains []string, regionPrefix string) (string, error)

// Resolve implements the Resolver interface for ResolverFunc.
func (d ResolverFunc) Resolve(ctx context.Context, domains []string, regionPrefix string) (string, error) {
	return d(ctx, domains, regionPrefix)
}

type resolverImpl struct{}

func newResolverImpl() *resolverImpl {
	return &resolverImpl{}
}

// Resolve performs DNS resolution on all domains concurrently and returns the first successful one.
func (r *resolverImpl) Resolve(ctx context.Context, domains []string, regionPrefix string) (string, error) {
	var wg sync.WaitGroup

	done := make(chan struct{}, 1)
	res := make(chan string, len(domains))
	for _, d := range domains {
		wg.Add(1)

		go func(domain string, regionPrefix string) {
			defer wg.Done()
			url := regionPrefix + "." + domain
			_, err := net.LookupHost(url)
			if err == nil {
				res <- domain
			}
		}(d, regionPrefix)
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
	return "", errors.New("query all dns is failed")
}
