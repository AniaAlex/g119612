package etsi119612

import (
	"context"

	go_cache "github.com/eko/gocache/lib/v4/cache"
)

// we do not need code that we do not use
// add here some other parts
type ByteCache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, b []byte, opts ...Option) error
	Delete(ctx context.Context, key string) error
}

type Option struct{}

// reduce the number of functions for go cache using interface segregation principle
type GoCacheAdapter struct {
	cache go_cache.CacheInterface[[]byte]
}

func (g *GoCacheAdapter) Get(ctx context.Context, key string) ([]byte, error) {
	return g.cache.Get(ctx, key)
}

func (g *GoCacheAdapter) Set(ctx context.Context, key string, b []byte, opts ...Option) error {
	// process options if needed
	return g.cache.Set(ctx, key, b)
}

func (g *GoCacheAdapter) Delete(ctx context.Context, key string) error {
	return g.cache.Delete(ctx, key)
}

// patrial encapsulation, more encapsulation is needed here
func NewCachedTSLFetcher(cache ByteCache) *cachedTSLFetcher {
	return &cachedTSLFetcher{cache: cache}
}

// adatper for go cache
func NewCachedTSLFetcherWithGoCache(cache go_cache.CacheInterface[[]byte]) *cachedTSLFetcher {
	return &cachedTSLFetcher{cache: &GoCacheAdapter{cache: cache}}
}

type cachedTSLFetcher struct {
	cache ByteCache
}

func (f *cachedTSLFetcher) Fetcher(ctx context.Context, url string) (*TSL, []byte, error) {
	if cachedValue, err := f.cache.Get(ctx, url); err == nil && len(cachedValue) > 0 {
		tsl, err := UnmarshalCleanCerts(cachedValue, url)
		if err != nil {
			return nil, nil, err
		}
		return tsl, cachedValue, nil

	}
	bodyBytes, err := FetchTSLBytes(url)
	if err != nil {
		return nil, nil, err
	}
	if len(bodyBytes) == 0 {
		return nil, nil, nil
	}
	setErr := f.cache.Set(ctx, url, bodyBytes)
	if setErr != nil {
		return nil, nil, setErr
	}
	tsl, err := UnmarshalCleanCerts(bodyBytes, url)
	if err != nil {
		return nil, nil, err
	}
	return tsl, bodyBytes, nil
}
