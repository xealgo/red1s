package network

// HandlerFunc generic request handler function implemented
// by the client code.
type HandlerFunc func([]byte) ([]byte, error)

// Listener simple listener server interface.
type Listener interface {
	Listen(handlerFn HandlerFunc) error
}
