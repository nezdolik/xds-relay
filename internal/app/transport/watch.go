package transport

import (
	"fmt"

	gcp "github.com/envoyproxy/go-control-plane/pkg/cache/v2"
	gcpv3 "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
)

// Watch interface abstracts v2 and v3 watches
type Watch interface {
	Close()
	GetCh() interface{}
	Send(s Response) (bool, error)
	GetRequest() Request
}

var _ Watch = &WatchV2{}

// WatchV2 wraps the v2 DiscoveryRequest/Response
type WatchV2 struct {
	Req Request
	out chan gcp.Response
}

// NewWatchV2 creates a new instance of v2 watch
func NewWatchV2(req Request) *WatchV2 {
	return &WatchV2{
		Req: req,
		out: make(chan gcp.Response, 1),
	}
}

func (w *WatchV2) Close() {
	close(w.out)
}

func (w *WatchV2) GetCh() interface{} {
	return w.out
}

func (w *WatchV2) Send(s Response) (bool, error) {
	resp, ok := s.GetPayload().(*gcp.Response)
	if !ok {
		return false, fmt.Errorf("payload %s could not be casted to gcp.Response", s)
	}

	select {
	case w.out <- *resp:
		return true, nil
	default:
		return false, nil
	}
}

func (w *WatchV2) GetRequest() Request {
	return w.Req
}

var _ Watch = &WatchV3{}

// WatchV3 wraps the v3 DiscoveryRequest/Response
type WatchV3 struct {
	Req Request
	out chan gcpv3.Response
}

// NewWatchV3 creates a new instace of V3 watch
func NewWatchV3(req Request) *WatchV3 {
	return &WatchV3{
		Req: req,
		out: make(chan gcpv3.Response, 1),
	}
}

func (w *WatchV3) Close() {
	close(w.out)
}

func (w *WatchV3) GetCh() interface{} {
	return w.out
}

func (w *WatchV3) Send(s Response) (bool, error) {
	resp, ok := s.GetPayload().(*gcpv3.Response)
	if !ok {
		return false, fmt.Errorf("payload %s could not be casted to gcp.Response", s)
	}

	select {
	case w.out <- *resp:
		return true, nil
	default:
		return false, nil
	}
}

func (w *WatchV3) GetRequest() Request {
	return w.Req
}
