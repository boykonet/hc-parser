package html_parser

import (
	"bytes"

	"golang.org/x/net/html"
)

type vbtParserHTML struct {
}

func NewVbtHTMLParser() IHTMLParser {
	return &vbtParserHTML{}
}

func (p *vbtParserHTML) ParseFlats(file []byte, cursor *string) ([]string, error) {
	var addrs []string
	var isCursor bool

	tkn := html.NewTokenizer(bytes.NewReader(file))

	for {
		tt := tkn.Next()
		switch tt {
		case html.ErrorToken:
			if !isCursor {
				*cursor = ""
			}
			return addrs, nil
		case html.StartTagToken:
			t, ha := tkn.TagName()
			if len(t) == 1 && t[0] == 'a' && ha {
				keys, vals := getAttrs(tkn)
				if len(keys) == 2 && keys[1] == "class" && vals[1] == "property" {
					addrs = append(addrs, vals[0])
				} else if len(keys) == 3 &&
					keys[0] == "class" && vals[0] == "shiftpage" &&
					keys[2] == "rel" && vals[2] == "prefetch" {
					*cursor = vals[1]
					isCursor = true
				}
			}
		}
	}
}
