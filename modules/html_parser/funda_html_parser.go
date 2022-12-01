package html_parser

import (
	"bytes"
	"fmt"

	"github.com/antchfx/htmlquery"
)

type fundaParserHTML struct {
}

func NewFundaHTMLParser() IHTMLParser {
	return &fundaParserHTML{}
}

func (p *fundaParserHTML) ParseFlats(file []byte, cursor *string) ([]string, error) {
	var flats []string

	rootNode, err := htmlquery.Parse(bytes.NewBuffer(file))
	if err != nil {
		return nil, fmt.Errorf("parce error %v", err)
	}

	// looking new links
	flatURLs := htmlquery.Find(rootNode, `//div[@class="search-result-media"]/a[@data-object-url-tracking="resultlist"]`)
	if len(flatURLs) == 0 {
		return nil, fmt.Errorf("couldn't find new items")
	}
	for _, v := range flatURLs {
		flats = append(flats, htmlquery.SelectAttr(v, "href"))
	}


	// looking cursor for the next page
	currentCursor := htmlquery.Find(rootNode, `//a[@rel="next"]`)
	if len(currentCursor) == 0 {
		*cursor = ""
	} else {
		*cursor = htmlquery.SelectAttr(currentCursor[0], "href")
	}
	return flats, nil
}
