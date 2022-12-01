package html_parser

import (
	"testing"

)

type data struct {
	Data struct {
		Href string `xpath: a[@data-object-url-tracking="resultlist"]`
	} `xpath: //div[@class="search-result-media"]`
}

func TestFunda(t *testing.T) {
	// file, err := ioutil.ReadFile("funda")
	// if err != nil {
	// 	t.Fatalf("open %v", err)
	// }

	// rootNode, err := htmlquery.Parse(bytes.NewBuffer(file))
	// if err != nil {
	// 	t.Fatalf("parce error %v", err)
	// }

	// s := htmlquery.Find(rootNode, `//div[@class="search-result-media"]/a[@data-object-url-tracking="resultlist"]`)
	// if len(s) == 0 {
	// 	t.Fatalf("couldn't find")
	// }
	// for _, v := range s {
	// 	t.Logf("result: %v", htmlquery.SelectAttr(v, "href"))
	// }
}