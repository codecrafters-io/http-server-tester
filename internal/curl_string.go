package internal

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
)

func httpRequestToCurlString(req *http.Request) string {
	curlCommand := fmt.Sprintf("curl -v -X %s %s%s%s",
		req.Method, req.URL.String(), formatHeaders(req.Header), formatBody(req))

	return curlCommand
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
