package printer

type IPrinter interface {
	Preview(competitor string)
	Flats(month string, countAllFlats int, newFlats []string)
}