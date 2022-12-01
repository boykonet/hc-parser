package run

import (
	"fmt"
	"sync"

	fundaConf "vbtor/app/conf/funda"
	vbtConf "vbtor/app/conf/vbt"
	"vbtor/app/run/competitors"
	"vbtor/app/run/competitors/funda"
	"vbtor/app/run/competitors/vbt"
	"vbtor/app/utils"

	"vbtor/modules/logger"
	"vbtor/modules/printer"

	"go.uber.org/zap"
)

const (
	vbtCompetitor = "VB&T"
	fundaCompetitor = "FUNDA"
)

type runner struct {
	competitors map[string]competitors.ICompetitorRunner

	print printer.IPrinter

	wg *sync.WaitGroup
	mutex *sync.Mutex
}

func NewRunner(
	vconf vbtConf.Configuration,
	fconf fundaConf.Configuration,
	logger logger.ILogger,
) IRunner {
	return &runner{
		competitors: map[string]competitors.ICompetitorRunner{
			vbtCompetitor: vbt.NewRunner(vconf, logger.With(zap.String("competitor", vbtCompetitor))),
			fundaCompetitor: funda.NewRunner(fconf, logger.With(zap.String("competitor", fundaCompetitor))),
		},

		print: printer.NewPrinter(),
		
		wg: &sync.WaitGroup{},
		mutex: &sync.Mutex{},
	}
}

func (r *runner) Run() error {
	var flatsFlag bool

	for _, customRunner := range r.competitors {
		r.wg.Add(1)

		go customRunner.Run(func(name string, flats []string, countAllFlats int, log logger.ILogger) {
			r.mutex.Lock()

			r.print.Preview(name)
			r.print.Flats("ALL MONTHS", countAllFlats, flats)

			if len(flats) > 0 {
				flatsFlag = true
			}

			r.mutex.Unlock()

			r.wg.Done()
		})
	}

	r.wg.Wait()

	fmt.Printf("-----------------------\n\n")

	if flatsFlag {
		if err := utils.RunMusic("./music/vbt_notification.mp3"); err != nil {
			return fmt.Errorf("falat error while playing music: %v", err)
		}
	}
	return nil
}