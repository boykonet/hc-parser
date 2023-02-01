package main

import (
	"math/rand"
	"os"
	"os/signal"
	"time"

	"vbtor/app/conf"
	confFunda "vbtor/app/conf/funda"
	confRijnland "vbtor/app/conf/rijnland"
	confVBT "vbtor/app/conf/vbt"
	"vbtor/app/run"
	"vbtor/app/run/competitors/funda"
	"vbtor/app/run/competitors/rijnland"
	"vbtor/app/run/competitors/vbt"
	"vbtor/modules/logger"

	"go.uber.org/zap"
)

const (
	vbtConfigurationPath		= "configurations/vbt/conf.yaml"
	fundaConfigurationPath		= "configurations/funda/conf.yaml"
	rijnlandConfigurationPath	= "configurations/rijnland/conf.yaml"
	mainConfigurationPath		= "configurations/main.conf.yaml"
)

func main() {
	s := make(chan os.Signal)
	signal.Notify(s, os.Interrupt)

	log := logger.NewLogger()

	// Parse a main configuration file with max and min values
	mconf, err := conf.ParseConfiguration(mainConfigurationPath)
	if err != nil {
		log.Fatal("fatal error parse main configuration", zap.Error(err))
		return
	}

	// VB&T
	vconf, err := confVBT.ParseConfiguration(vbtConfigurationPath)
	if err != nil {
		log.Fatal("fatal error parse VB&T configuration", zap.Error(err))
		return
	}
	vrunner := vbt.NewRunner(*vconf, log.With(zap.String("competitor", "vb&t")))
	err = vrunner.SetFlatsFromFile()
	if err != nil {
		log.Error("error set flats from file", zap.Error(err))
	}

	// FUNDA
	fconf, err := confFunda.ParseConfiguration(fundaConfigurationPath)
	if err != nil {
		log.Fatal("fatal error parse FUNDA configuration", zap.Error(err))
		return
	}
	frunner := funda.NewRunner(*fconf, log.With(zap.String("competitor", "funda")))
	err = frunner.SetFlatsFromFile()
	if err != nil {
		log.Error("error set flats from file", zap.Error(err))
	}

	// RIJNLAND
	rconf, err := confRijnland.ParseConfiguration(rijnlandConfigurationPath)
	if err != nil {
		log.Fatal("fatal error parse RIJNLAND configuration", zap.Error(err))
		return
	}
	rrunner := rijnland.NewRunner(*rconf, log.With(zap.String("competitor", "rijnland")))
	err = rrunner.SetFlatsFromFile()
	if err != nil {
		log.Error("error set flats from file", zap.Error(err))
	}

	go func(){
		for sig := range s {
			err1 := vrunner.SaveToFile()
			err2 := frunner.SaveToFile()
			// err3 := rrunner.SaveToFile()
			log.Error("signal", zap.Any("signal", sig), zap.Errors("error", []error{err1, err2}))
			os.Exit(0)
		}
	}()

	//
	runner := run.NewRunner(vrunner, frunner, rrunner, log)

	for {
		runner.Run()

		// Getting time for sleep
		timeSleep := rand.Intn(mconf.Max-mconf.Min) + mconf.Min
		rand.Seed(time.Now().UnixNano())

		// Sleep on this time
		time.Sleep(time.Duration(timeSleep) * time.Second)
	}
}
