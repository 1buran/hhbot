package apiendpoint

import "net/url"

type postnegotiations struct {
	vid, rid, msg string
}

func (n postnegotiations) SetPage(p int)            {}
func (n postnegotiations) Url(params ...any) string { return "/negotiations" }
func (n postnegotiations) Payload() url.Values {
	payload := url.Values{}
	payload.Set("resume_id", n.rid)
	payload.Set("vacancy_id", n.vid)
	payload.Set("message", "Интересно!")
	if n.msg != "" {
		payload.Set("message", n.msg)
	}

	return payload
}

func NewPostNegotinations(vacancyID, resumeID, message string) ApiEndpoint {
	return &postnegotiations{vid: vacancyID, rid: resumeID, msg: message}
}
