package main

import (
	"math/rand"
	"runtime"
	"time"

	"gitlab.vailsys.com/jerny/coffer/cmd/coffer/options"
	"gitlab.vailsys.com/jerny/coffer/pkg/logger"
	"gitlab.vailsys.com/jerny/coffer/pkg/storage"
	"gitlab.vailsys.com/jerny/coffer/pkg/storage/mongo"
	"gitlab.vailsys.com/jerny/coffer/server"

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

	recordingRepo := storage.NewMongoRecordingRepo(opt.MongoConfig, provider)
	assetRepo := storage.NewGridFSRepo(opts.MongoConfig, provider)

	if err != nil {
		logger.Logger.Fatalf(err.Error())
	}

	s := server.NewCofferServer(opt, recordingRepo, assetRepo)

	if err := s.Run(); err != nil {
		logger.Logger.Fatalf(err.Error())
	}
}
