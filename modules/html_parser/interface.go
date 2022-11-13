package html_parser

type IHTMLParser interface {
	ParseFlats(cursor *string) []string
}