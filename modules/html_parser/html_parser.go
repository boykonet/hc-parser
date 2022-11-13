package html_parser

import (
	"bytes"

	"golang.org/x/net/html"
)

type parserHTML struct {
	Tkn *html.Tokenizer
}

func NewHTMLParser(text []byte) IHTMLParser {
	return &parserHTML{
		Tkn: html.NewTokenizer(bytes.NewReader(text)),
	}
}

func getAttrs(t *html.Tokenizer) ([]string, []string) {
	var keys, vals []string
	for {
		key, val, moreAttrs := t.TagAttr()
		keys = append(keys, string(key))
		vals = append(vals, string(val))
		if !moreAttrs {
			break
		}
	}
	return keys, vals
}

func (p *parserHTML) ParseFlats(cursor *string) []string {
	var addrs []string
	var isCursor bool

	for {
		tt := p.Tkn.Next()
		switch tt {
		case html.ErrorToken:
			if !isCursor {
				*cursor = ""
			}
			return addrs
		case html.StartTagToken:
			t, ha := p.Tkn.TagName()
			if len(t) == 1 && t[0] == 'a' && ha {
				keys, vals := getAttrs(p.Tkn)
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