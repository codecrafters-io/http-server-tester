package http_connection

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
)

func httpRequestToCurlString(req *http.Request) string {
	method := req.Method
	var curlCommand string

	if method == "GET" {
		curlCommand = fmt.Sprintf("curl -v %s%s%s",
			req.URL.String(), formatHeaders(req.Header), formatBody(req))
	} else {
		curlCommand = fmt.Sprintf("curl -v -X %s %s%s%s",
			method, req.URL.String(), formatHeaders(req.Header), formatBody(req))
	}

	return curlCommand
}

func HttpKeepAliveRequestToCurlString(requests []*http.Request) string {
	// While sending multiple requests on the same connection,
	// We need to log all requests at once to properly emulate with curl
	var requestsBuilder strings.Builder
	for _, req := range requests {
		requestsBuilder.WriteString(fmt.Sprintf("%s%s%s ",
			req.URL.String(), formatHeaders(req.Header), formatBody(req)))
	}
	return fmt.Sprintf("curl --http1.1 -v %s", requestsBuilder.String())
}

func formatHeaders(headers http.Header) string {
	var formattedHeaders string

	// The sorting stuff is to make the output reproducible as hash iteration
	// is not guaranteed generate the same result every time
	var headerKeys = make([]string, 0, len(headers))
	for key := range headers {
		headerKeys = append(headerKeys, key)
	}
	sort.Strings(headerKeys)

	for _, key := range headerKeys {
		values := headers[key]
		sort.Strings(values)

		for _, value := range values {
			formattedHeaders += fmt.Sprintf(" -H \"%s: %s\"", key, value)
		}
	}
	return formattedHeaders
}

func formatBody(req *http.Request) string {
	if req.Body == nil || bodyToString(req) == "" {
		return ""

	}
	return fmt.Sprintf(" -d '%s'", escapeSingleQuotes(bodyToString(req)))
}

func bodyToString(req *http.Request) string {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return ""
	}

	req.Body = io.NopCloser(bytes.NewBuffer(body))

	return string(body)
}

func escapeSingleQuotes(s string) string {
	return strings.ReplaceAll(s, "'", `\'`)
}
