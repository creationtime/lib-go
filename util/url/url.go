package url

import (
	"fmt"
	"net/url"
	"strings"
)

const (
	httpPrefix  = "http://"
	httpsPrefix = "https://"
	prefix      = "//"
	cut         = "/"
)

func Join(uri string, baseUrl string) string {
	if len(uri) == 0 {
		return uri
	}
	if strings.HasPrefix(uri, httpPrefix) || strings.HasPrefix(uri, httpsPrefix) || strings.HasPrefix(uri, prefix) {
		return uri
	}
	return fmt.Sprintf("%s/%s", strings.TrimSuffix(baseUrl, cut), strings.TrimPrefix(uri, cut))
}

func Trim(uri string) string {
	if len(uri) == 0 {
		return uri
	}
	if strings.HasPrefix(uri, httpPrefix) || strings.HasPrefix(uri, httpsPrefix) || strings.HasPrefix(uri, prefix) {
		u, err := url.Parse(uri)
		if err != nil {
			return ""
		}
		return u.RequestURI()
	}
	return uri
}
