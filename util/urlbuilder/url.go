package urlbuilder

import (
	"net/url"
	"strings"
)

const (
	// ModeJoin join parts into a baseURL
	ModeJoin = 1
	// ModeParse create by a baseURL
	ModeParse = 2
)

type UrlBuilder struct {
	mode int
	// url baseURL: host/path
	baseURL string

	schema   string
	host     string
	path     string
	fragment string

	parsed  *url.URL
	queries url.Values
}

type Query = url.Values

func New() *UrlBuilder {
	u := newBuilder()
	u.mode = ModeJoin

	return u
}

// Of create an url builder from a base url
func Of(base string) *UrlBuilder {
	u := newBuilder()
	u.mode = ModeParse
	u.baseURL = base

	return u
}

// Join create an url builder by multiple url parts
func Join(segments ...string) *UrlBuilder {
	u := newBuilder()
	u.mode = ModeParse

	baseURL := strings.Join(segments, "/")

	if !strings.HasPrefix(baseURL, "http") {
		baseURL = "http://" + baseURL
	}

	u.baseURL = baseURL

	return u
}

// Http create an url builder start of a http schema
func Http() *UrlBuilder {
	return newBuilder().Schema("http")
}

// Https create an url builder start of a https schema
func Https() *UrlBuilder {
	return newBuilder().Schema("https")
}

func newBuilder() *UrlBuilder {
	return &UrlBuilder{
		queries: url.Values{},
	}
}

// AddQuery app a parameter to query
func (ub *UrlBuilder) AddQuery(key, value string) *UrlBuilder {
	ub.queries.Add(key, value)

	return ub
}

// Schema if a wrong schema is inputted, https will default
func (ub *UrlBuilder) Schema(schema string) *UrlBuilder {
	schema = strings.Trim(schema, "/")

	if !strings.HasPrefix(schema, "http") {
		schema = "https"
	}

	ub.schema = schema

	return ub
}

func (ub *UrlBuilder) Host(host string) *UrlBuilder {
	ub.host = strings.Trim(host, "/")

	return ub
}

func (ub *UrlBuilder) Path(path string) *UrlBuilder {
	ub.path = strings.Trim(path, "/")

	return ub
}

func (ub *UrlBuilder) RawQuery(query string) *UrlBuilder {
	query = strings.TrimLeft("?", query)
	query = strings.TrimRight("&", query)

	ub.queries, _ = url.ParseQuery(query)

	return ub
}

func (ub *UrlBuilder) Fragment(fragment string) *UrlBuilder {
	ub.fragment = strings.TrimLeft("#", fragment)

	return ub
}

func (ub *UrlBuilder) Build() string {
	if ub.mode == ModeJoin {
		ub.buildBaseURL()
	}

	ub.parsed, _ = url.Parse(ub.baseURL)

	appendQuery := ub.queries.Encode()

	if appendQuery != "" {
		if ub.parsed.RawQuery == "" {
			ub.parsed.RawQuery = appendQuery
		} else {
			ub.parsed.RawQuery = ub.parsed.RawQuery + "&" + appendQuery
		}
	}

	if ub.fragment != "" {
		ub.parsed.Fragment = ub.fragment
	}

	return ub.parsed.String()
}

func (ub *UrlBuilder) buildBaseURL() {
	sb := strings.Builder{}

	if ub.schema == "" {
		ub.schema = "http"
	}

	sb.WriteString(ub.schema)
	sb.WriteString("://")
	sb.WriteString(ub.host)
	sb.WriteString("/")
	sb.WriteString(ub.path)

	ub.baseURL = sb.String()
}
