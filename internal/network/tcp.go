package network

import (
	"fmt"
	"net"
	"time"
)

var _ Listener = (*TCPServer)(nil)

// TCPServer simple TCP listener adapter.
type TCPServer struct {
	host string
	port string
}

// NewTCPServer returns a TCPServer instance.
func NewTCPServer(host, port string) *TCPServer {
	return &TCPServer{host, port}
}

// Listen waits for incoming connection requests on the provided
// hostname and port.
func (ts *TCPServer) Listen(handlerFn HandlerFunc) error {
	l, err := net.Listen("tcp", ts.host+":"+ts.port)
	if err != nil {
		return fmt.Errorf("Error listening to %s on port %s: %w", ts.host, ts.port, err)
	}

	fmt.Printf("Listening on %s at port %s\n", ts.host, ts.port)

	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("Error accepting connection to %s on port %s: %w", ts.host, ts.port, err)
		}

		// TODO: Pass in a channel to pass errors on vs using an anon function.
		go func(conn net.Conn) {
			if err := ts.handle(conn, handlerFn); err != nil {
				// A logger call would go here...
				fmt.Printf("Error handling connection: %s\n", err)
			}
		}(conn)
	}
}

func (ts *TCPServer) handle(conn net.Conn, handlerFn HandlerFunc) error {
	defer conn.Close()

	conn.(*net.TCPConn).SetKeepAlivePeriod(300 * time.Second)
	conn.SetReadDeadline(time.Now().Add(5 * time.Second))

	// TODO: This isn't great and should use bytes.Buffer to read in
	// data in chunks vs one large allocation, but I didn't have time
	// to get this working with CRLF correctly.
	bytes := make([]byte, 4096)

	if _, err := conn.Read(bytes); err != nil {
		return fmt.Errorf("Error reading request: %w", err)
	}

	if handlerFn == nil {
		return fmt.Errorf("Error processing request: request handler not set")
	}

	out, err := handlerFn(bytes)
	if err != nil {
		return fmt.Errorf("Error processing request: %w", err)
	}

	_, err = conn.Write(out)
	if err != nil {
		return fmt.Errorf("Error sending response: %w", err)
	}
	return nil
}
