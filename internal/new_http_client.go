package internal

import (
	"context"
	"net"
	"net/http"
	"time"
)

func NewHTTPClient() *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			// Don't complain if we aren't able to connect at first, for at least 10 seconds.
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				attempts := 0

				for {
					var err error
					var conn net.Conn

					conn, err = net.Dial("tcp", addr)

					if err == nil {
						return conn, nil
					}

					// Already a timeout
					if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
						return nil, err
					}

					// 100 * 100ms = 10s
					if attempts > 100 {
						return nil, err
					}

					attempts += 1
					time.Sleep(100 * time.Millisecond)
				}
			},
		},
	}
}
