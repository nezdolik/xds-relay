package transport

import (
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
)

// Request is the generic interface to abstract v2 and v3 DiscoveryRequest types
type Request interface {
	CreateResponse(string, []types.Resource, []types.MarshaledResource) Response
}
