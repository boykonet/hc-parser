package main

import (
	"math/rand"
	"time"

	"vbtor/app/conf"
	"vbtor/app/conf/funda"
	"vbtor/app/conf/vbt"
	"vbtor/app/run"
	"vbtor/modules/logger"

	"go.uber.org/zap"
)

const (
	vbtConfigurationPath = "configurations/vbt/conf.yaml"
	fundaConfigurationPath = "configurations/funda/conf.yaml"
	mainConfigurationPath = "configurations/main.conf.yaml"
)

func main() {
	log := logger.NewLogger()

	// Parse a main configuration file with max and min values
	mconf, err := conf.ParceConfiguration(mainConfigurationPath)
	if err != nil {
		log.Fatal("fatal error parse main configuration", zap.Error(err))
		return
	}
	
	// Parse configuration file with cookies
	vconf, err := vbt.ParceConfiguration(vbtConfigurationPath)
	if err != nil {
		log.Fatal("fatal error parse VB&T configuration", zap.Error(err))
		return
	}

	fconf, err := funda.ParceConfiguration(fundaConfigurationPath)
	if err != nil {
		log.Fatal("fatal error parse FUNDA configuration", zap.Error(err))
		return
	}

	runner := run.NewRunner(*vconf, *fconf, log)

	for {
		runner.Run()

		// Getting time for sleep
		timeSleep := rand.Intn(mconf.Max - mconf.Min) + mconf.Min
		rand.Seed(time.Now().UnixNano())
		
		// Sleep on this time
		time.Sleep(time.Duration(timeSleep) * time.Second)
	}
}
