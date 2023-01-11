package main

import (
	"context"
	"net"
	"time"
)

type Cache interface {
	Set(ctx context.Context, key []byte, value []net.IPAddr, expiration time.Duration) error
	Get(ctx context.Context, key []byte) ([]net.IPAddr, error)
	Delete(ctx context.Context, key []byte) error
}

type Resolver interface {
	LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error)
}

type CachingResolver interface {
	Resolver
	ForgetIP(ctx context.Context, host string, ip net.IPAddr) error
	ForgetHost(ctx context.Context, host string) error
}

type cachingResolver struct {
	resolver        Resolver
	cache           Cache
	cacheExpiration time.Duration // num of seconds
}

func NewCachingResolver(cache Cache, resolver Resolver, cacheExpiration int) CachingResolver {
	return &cachingResolver{
		cache:           cache,
		resolver:        resolver,
		cacheExpiration: time.Duration(cacheExpiration) * time.Second,
	}
}

func (c *cachingResolver) LookupIPAddr(ctx context.Context, host string) ([]net.IPAddr, error) {
	cacheKey := []byte(host)

	ipAddrs, err := c.cache.Get(ctx, cacheKey)
	if err != nil {
		return nil, err // @TODO owrapovet error
	}

	if ipAddrs == nil || len(ipAddrs) == 0 {
		resolvedIPs, err := c.resolver.LookupIPAddr(ctx, host)
		if err != nil {
			return nil, err // @TODO owrapovet error
		}

		err = c.cache.Set(ctx, cacheKey, resolvedIPs, c.cacheExpiration)
		if err != nil {
			return nil, err // @TODO Owrapovet error
		}

		ipAddrs = resolvedIPs
	}

	return ipAddrs, nil
}

func (c *cachingResolver) ForgetHost(ctx context.Context, host string) error {
	return c.cache.Delete(ctx, []byte(host)) // @TODO Owrapovet error
}

func (c *cachingResolver) ForgetIP(ctx context.Context, host string, ip net.IPAddr) error {
	h := []byte(host)
	cachedIPs, err := c.cache.Get(ctx, h)
	if err != nil { // address is not stored in cache
		return nil
	}

	if len(cachedIPs) == 1 {
		return c.cache.Delete(ctx, h) // @TODO Owrapovet error
	}

	return c.cache.Set(ctx, h, findAndDelete(cachedIPs, ip), c.cacheExpiration)
}

func indexOf(ip net.IPAddr, data []net.IPAddr) {
	for
}

func remove(s []net.IPAddr, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func findAndDelete(s []net.IPAddr, deletedIP net.IPAddr) []net.IPAddr {
	index := 0
	for _, i := range s {
		if i.String() != deletedIP.String() {
			s[index] = i
			index++
		}
	}
	return s[:index]
}
