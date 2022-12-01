package html_parser

type IHTMLParser interface {
	ParseFlats(file []byte, cursor *string) ([]string, error)
}