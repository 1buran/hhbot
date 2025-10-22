// This package declares DTO for API responses,
// all structs are defined as close as possible to original data fields names to avoid
// chatty json tags.
//
// To avoid bloated of single file, all entities splitted into small related files,
// the main package file (this file) contains only high level structures: response types.
package dto

import "net/http"

type ResponseType interface {
	Vacancy | Dictionary | Resume
}

type ResponseItems[T ResponseType] struct {
	Found, Per_page, Pages, Page int
	Items                        []T
}

type Response struct {
	Code    int
	Status  string
	Headers http.Header
	Error   string
}
