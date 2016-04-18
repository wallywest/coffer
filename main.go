package main

import (
	"math/rand"
	"runtime"
	"time"

	"gitlab.vailsys.com/jerny/coffer/internal/logger"
	"gitlab.vailsys.com/jerny/coffer/internal/options"
	"gitlab.vailsys.com/jerny/coffer/internal/recording"
	"gitlab.vailsys.com/jerny/coffer/internal/server"
	"gitlab.vailsys.com/jerny/coffer/internal/storage/driver/mongo"

	"github.com/spf13/pflag"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UTC().UnixNano())

	opt := options.NewCofferConfig()
	opt.AddFlags(pflag.CommandLine)

	pflag.Parse()

	logger.SetLogLevel(opt.LogLevel)

	provider, err := mongo.NewSessionProvider(opt.MongoConfig)
	if err != nil {
		logger.Logger.Fatalf(err.Error())
	}
	defer provider.Close()

	recordingRepo := recording.NewMongoRecordingRepo(opt.MongoConfig, provider)
	assetRepo := recording.NewGridFSRepo(opt.MongoConfig, provider)

	if err != nil {
		logger.Logger.Fatalf(err.Error())
	}

	s := server.NewCofferServer(opt, recordingRepo, assetRepo)

	if err := s.Run(); err != nil {
		logger.Logger.Fatalf(err.Error())
	}

	logger.Logger.Info("shutting down")
}