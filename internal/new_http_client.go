package internal

import (
	"context"
	"net"
	"net/http"
	"time"
)

func NewHTTPClient() *http.Client {
	return &http.Client{
		// Don't follow redirects
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			// Don't complain if we aren't able to connect at first, for at least 10 seconds.
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				attempts := 0

				for {
					var err error
					var conn net.Conn

					// Used DialTimeout instead of Dial to return an error if the connection is not established within 10 seconds
					// This is to prevent the program from hanging indefinitely if there are problems with the server
					conn, err = net.DialTimeout("tcp", addr, clientTimeout/time.Duration(attempts + 1))

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
