package transport

import (
	"context"
	"fmt"

	discovery "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	gcpv2 "github.com/envoyproxy/go-control-plane/pkg/cache/v2"
	"github.com/envoyproxy/xds-relay/internal/app/cache"
)

var _ gcpv2.Cache = &V2{}

// V2 cache wraps a shared cache
type V2 struct {
	cache *cache.Cache
}

// NewV2 creates a shared cache for v2 requests
func NewV2(cache *cache.Cache) *V2 {
	return &V2{
		cache: cache,
	}
}

// CreateWatch is the grpc backed xds handler
func (v *V2) CreateWatch(req gcpv2.Request) (chan gcpv2.Response, func()) {
	return nil, nil
}

// Fetch is the rest backed xds handler
func (v *V2) Fetch(ctx context.Context, req gcpv2.Request) (gcpv2.Response, error) {
	return nil, fmt.Errorf("Fetch is not supported")
}

var _ Request = &RequestV2{}

// RequestV2 is the v2.DiscoveryRequest impl of Request
type RequestV2 struct {
	r *discovery.DiscoveryRequest
}

// NewRequestV2 creates a Request objects which wraps v2.DiscoveryRequest
func NewRequestV2(r *discovery.DiscoveryRequest) *RequestV2 {
	return &RequestV2{
		r: r,
	}
}

var _ Response = &ResponseV2{}

// ResponseV2 is the v2.DiscoveryRequest impl of Request
type ResponseV2 struct {
	req                Request
	resources          []types.Resource
	marshaledResources []types.MarshaledResource
	version            string
}

// CreateResponse creates a versioned Response
func (r *RequestV2) CreateResponse(payloadVersion string, proto []types.Resource, marshaledResources []types.MarshaledResource) *Response {
	return &ResponseV2{
		req:                r,
		version:            version,
		resources:          protos,
		marshaledResources: marshaledResources,
	}
}
