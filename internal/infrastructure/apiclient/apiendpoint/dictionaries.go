package apiendpoint

import "net/url"

type getdictionaries struct{}

func (d getdictionaries) SetPage(p int)            {}
func (d getdictionaries) Url(params ...any) string { return "/dictionaries" }
func (d getdictionaries) Payload() url.Values      { return nil }

func NewGetDictionary() ApiEndpoint { return &getdictionaries{} }
