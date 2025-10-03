package etsi119612

import (
	"context"

	go_cache "github.com/eko/gocache/lib/v4/cache"
)

//patrial encapsulation, more encapsulation here 
func NewCachedTSLFetcher(cache go_cache.CacheInterface[[]byte]) *CachedTSLFetcher {
	return &CachedTSLFetcher{cache: cache}
}

type CachedTSLFetcher struct {
	cache go_cache.CacheInterface[[]byte]
}

func (f *CachedTSLFetcher) Fetcher(ctx context.Context, url string) (*TSL, []byte, error) {
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
