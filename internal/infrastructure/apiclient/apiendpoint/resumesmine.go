package apiendpoint

import "net/url"

type getresumesmine struct{}

func (g getresumesmine) SetPage(p int)            {}
func (g getresumesmine) Url(params ...any) string { return "/resumes/mine" }
func (g getresumesmine) Payload() url.Values      { return nil }

func NewGetResumesMine() ApiEndpoint { return &getresumesmine{} }
