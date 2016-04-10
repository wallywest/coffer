package main

import (
	"math/rand"
	"runtime"
	"time"

	"gitlab.vailsys.com/jerny/coffer/cmd/coffer/options"
	"gitlab.vailsys.com/jerny/coffer/pkg/logger"
	"gitlab.vailsys.com/jerny/coffer/pkg/storage"
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

	provider, err := storage.NewGridFSProvider(opt)
	if err != nil {
		logger.Logger.Fatalf(err.Error())
	}

	s := server.NewCofferServer(opt, provider)

	if err := s.Run(); err != nil {
		logger.Logger.Fatalf(err.Error())
	}
}
