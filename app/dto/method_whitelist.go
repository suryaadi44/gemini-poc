package dto

import "net/http"

type MethodWhitelist struct {
	GET     bool
	HEAD    bool
	POST    bool
	PUT     bool
	DELETE  bool
	CONNECT bool
	OPTIONS bool
	TRACE   bool
	PATCH   bool
}

func NewMethodWhitelist(methods []string) MethodWhitelist {
	var whitelist MethodWhitelist
	for _, method := range methods {
		switch method {
		case http.MethodGet:
			whitelist.GET = true
		case http.MethodHead:
			whitelist.HEAD = true
		case http.MethodPost:
			whitelist.POST = true
		case http.MethodPut:
			whitelist.PUT = true
		case http.MethodDelete:
			whitelist.DELETE = true
		case http.MethodConnect:
			whitelist.CONNECT = true
		case http.MethodOptions:
			whitelist.OPTIONS = true
		case http.MethodTrace:
			whitelist.TRACE = true
		case http.MethodPatch:
			whitelist.PATCH = true
		}
	}
	return whitelist
}

func (m *MethodWhitelist) Append(methods []string) {
	for _, method := range methods {
		switch method {
		case http.MethodGet:
			m.GET = true
		case http.MethodHead:
			m.HEAD = true
		case http.MethodPost:
			m.POST = true
		case http.MethodPut:
			m.PUT = true
		case http.MethodDelete:
			m.DELETE = true
		case http.MethodConnect:
			m.CONNECT = true
		case http.MethodOptions:
			m.OPTIONS = true
		case http.MethodTrace:
			m.TRACE = true
		case http.MethodPatch:
			m.PATCH = true
		}
	}
}
