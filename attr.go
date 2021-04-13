package html

import (
	"strings"
)

func Id(id string) AttrNode {
	return Attr("id", id)
}

func Class(classes ...string) AttrNode {
	return Attr("class", strings.Join(classes, " "))
}

func Href(href string) AttrNode {
	return Attr("href", href)
}
