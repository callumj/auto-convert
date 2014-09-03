package shared

import (
	"fmt"
	"net/http"
	"net/url"
)

func GetFullUrl(req *http.Request) *url.URL {
	if req.URL.IsAbs() {
		return req.URL
	} else {
		var scheme string
		if req.TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
		baseURL, err := url.Parse(fmt.Sprintf("%s://%s", scheme, req.Host))
		if err != nil {
			return nil
		}
		return baseURL.ResolveReference(req.URL)
	}
}
