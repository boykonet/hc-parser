package html_parser

import "golang.org/x/net/html"

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
