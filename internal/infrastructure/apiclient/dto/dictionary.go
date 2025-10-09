package dto

import (
	"strings"
)

type DictionaryItem struct{ Id, Name string }

type DictionaryItems []DictionaryItem

func (d DictionaryItems) String() string {
	var names []string
	for _, v := range d {
		names = append(names, v.Name)
	}
	return strings.Join(names, ", ")
}

func (d DictionaryItems) Get(k string) string {
	for _, v := range d {
		if v.Id == k {
			return v.Name
		}
	}
	return ""
}

type Dictionary map[string]DictionaryItems

func (d Dictionary) Get(k string) (items DictionaryItems, ok bool) {
	items, ok = d[k]
	return
}
