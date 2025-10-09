package apiendpoint

import (
	"net/url"
)

type ApiEndpoint interface {
	Url(params ...any) string
	Payload() url.Values
	SetPage(p int)
}
