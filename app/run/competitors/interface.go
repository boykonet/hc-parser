package competitors

import "vbtor/modules/logger"

type ICompetitorRunner interface {
	Run(f CompetitorFunc)
	SetFlatsFromFile() error
	SaveToFile() error
}

type CompetitorFunc func(name string, flats []string, countAllFlats int, log logger.ILogger)

type Properties struct {
	Properties []string `yaml: properties`
}