package http_assertions

import (
	http_parser "github.com/codecrafters-io/http-server-tester/internal/http/parser"
)

type HTTPAssertion interface {
	Run(value http_parser.HTTPResponse) error
}
