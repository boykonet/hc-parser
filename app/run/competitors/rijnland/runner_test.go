package rijnland

import (
	"testing"
	confRijnland "vbtor/app/conf/rijnland"
	"vbtor/modules/logger"

	"go.uber.org/zap"
)

func TestRijnland(t *testing.T) {
	conf, err := confRijnland.ParseConfiguration("../../../../configurations/rijnland/conf.yaml")
	if err != nil {
		t.Fatalf("parse config: %v", err)
	}

	// fmt.Println(conf)
	runner := NewRunner(*conf, logger.NewLogger().With(zap.String("competitor", "rijnland")))
	runner.Run(func(name string, flats []string, countAllFlats int, log logger.ILogger) {
	})
}