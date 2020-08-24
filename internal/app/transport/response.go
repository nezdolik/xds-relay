package transport

// Response is the generic response interface
type Response interface {
	GetPayload() interface{}
}
