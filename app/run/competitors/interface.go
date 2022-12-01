package competitors

import "vbtor/modules/logger"

type ICompetitorRunner interface {
	Run(f CompetitorFunc)
}

type CompetitorFunc func(name string, flats []string, countAllFlats int, log logger.ILogger)