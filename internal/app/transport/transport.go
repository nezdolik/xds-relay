package transport

import (
	"context"

	api "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	endpointservice "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listenerservice "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	gcpv2 "github.com/envoyproxy/go-control-plane/pkg/server/v2"
	gcpv3 "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/envoyproxy/xds-relay/internal/app/cache"
	"google.golang.org/grpc"
)

// Service registers the xds endpoints with grpc server
type Service struct {
	gcpv2 gcpv2.Server
	gcpv3 gcpv3.Server
}

// NewService creates an instance of service, which can register the xds endpoints
func NewService(ctx context.Context, c *cache.Cache) *Service {
	return &Service{
		gcpv2: gcpv2.NewServer(ctx, NewV2(c), nil),
		gcpv3: gcpv3.NewServer(ctx, NewV3(c), nil),
	}
}

// Register the endpoints with grpc server
func (s *Service) Register(srv *grpc.Server) {
	// You must register your new xDS resource handler here for envoymanager
	// to comply with the server interface as well as accept requests
	api.RegisterRouteDiscoveryServiceServer(srv, s.gcpv2)
	api.RegisterClusterDiscoveryServiceServer(srv, s.gcpv2)
	api.RegisterEndpointDiscoveryServiceServer(srv, s.gcpv2)
	api.RegisterListenerDiscoveryServiceServer(srv, s.gcpv2)

	endpointservice.RegisterEndpointDiscoveryServiceServer(srv, s.gcpv3)
	clusterservice.RegisterClusterDiscoveryServiceServer(srv, s.gcpv3)
	routeservice.RegisterRouteDiscoveryServiceServer(srv, s.gcpv3)
	listenerservice.RegisterListenerDiscoveryServiceServer(srv, s.gcpv3)
}
