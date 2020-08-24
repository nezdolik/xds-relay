package transport

import (
	"context"
	"fmt"

	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	gcpv3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/xds-relay/internal/app/cache"
)

var _ gcpv3.Cache = &V3{}

// V3 cache wraps a shared cache
type V3 struct {
	cache *cache.Cache
}

// NewV3 creates a shared cache for v2 requests
func NewV3(cache *cache.Cache) *V3 {
	return &V3{
		cache: cache,
	}
}

// CreateWatch is the grpc backed xds handler
func (v *V3) CreateWatch(req gcpv3.Request) (chan gcpv3.Response, func()) {
	return nil, nil
}

// Fetch is the rest backed xds handler
func (v *V3) Fetch(ctx context.Context, req gcpv3.Request) (gcpv3.Response, error) {
	return nil, fmt.Errorf("Fetch is not supported")
}

var _ Request = &RequestV2{}

// RequestV3 is the v2.DiscoveryRequest impl of Request
type RequestV3 struct {
	r *discovery.DiscoveryRequest
}

// NewRequestV3 creates a Request objects which wraps v2.DiscoveryRequest
func NewRequestV3(r *discovery.DiscoveryRequest) *RequestV3 {
	return &RequestV3{
		r: r,
	}
}
