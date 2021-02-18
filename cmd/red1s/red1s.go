package main

import (
	"flag"
	"fmt"

	"github.com/xealgo/red1s/internal"
	"github.com/xealgo/red1s/internal/network"
	"github.com/xealgo/red1s/internal/resp"
	"github.com/xealgo/red1s/internal/store"
)

const (
	defaultHost = "localhost"
	defaultPort = "6379"
)

func main() {
	fmt.Println("Red1s Server")

	var connHost, connPort string

	useHost := flag.String("host", defaultHost, "Network hostname")
	usePort := flag.String("port", defaultPort, "Network port")
	flag.Parse()

	if useHost != nil {
		connHost = *useHost
	}

	if usePort != nil {
		connPort = *usePort
	}

	// Setup the in-memory store.
	ds := store.NewMemory()

	// Setup the TCP server.
	tcp := network.NewTCPServer(connHost, connPort)

	// TODO: Listen for sigint and implement graceful shutdown using a
	// channel to determine if we have any requests in-flight.
	if err := tcp.Listen(handler(ds)); err != nil {
		panic(err.Error())
	}
}

// handler returns a network.HandlerFunc that will be used by the tcp server
// to handle requests.
//
// TODO: This would obviously grow quickly as more commands were added. We'd
// instead want to implement a router down the road.
func handler(ds store.DataStore) network.HandlerFunc {
	return func(b []byte) ([]byte, error) {
		tokens, err := resp.Decode(b)
		if err != nil {
			return nil, fmt.Errorf("Error decoding request: %w", err)
		}

		p := resp.New(tokens)

		cmd, err := p.Parse()
		if err != nil {
			return nil, fmt.Errorf("Error parsing data: %w", err)
		}

		switch cmd.Name {
		case internal.CmdGet:
			return internal.Get(ds, cmd)
		case internal.CmdSet:
			return internal.Set(ds, cmd)
		case internal.CmdDel:
			return internal.Del(ds, cmd)
		default:
			return internal.Unknown(cmd)
		}
	}
}
