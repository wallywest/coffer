package options

import (
	"net"

	"github.com/spf13/pflag"
	mongostorage "gitlab.vailsys.com/jerny/coffer/pkg/storage/mongo"
)

type CofferConfig struct {
	AppName          string
	BindAddress      net.IP
	AdvertiseAddress net.IP
	Port             int
	SkipRegistration bool
	HeartBeatTTL     string
	DiscoveryMode    string
	DiscoveryTTL     string
	EnableProfiling  bool
	LogLevel         string
	MongoConfig      mongostorage.MongoConfig
}

func NewCofferConfig() *CofferConfig {
	c := &CofferConfig{
		AppName:          "coffer",
		BindAddress:      net.ParseIP("0.0.0.0"),
		AdvertiseAddress: net.ParseIP("127.0.0.1"),
		Port:             6000,
		HeartBeatTTL:     "15s",
		DiscoveryMode:    "static",
		DiscoveryTTL:     "15s",
		LogLevel:         "INFO",
		EnableProfiling:  false,
		SkipRegistration: false,
		MongoConfig: mongostorage.MongoConfig{
			DB:           "vcs",
			GridFSPrefix: "vcsfs",
			ServerList:   []string{"localhost:27017"},
		},
	}

	return c
}

func (c *CofferConfig) AddFlags(fs *pflag.FlagSet) {
	fs.IPVar(&c.BindAddress, "bind-address", c.BindAddress, "add me")
	fs.IPVar(&c.AdvertiseAddress, "advertise-address", c.AdvertiseAddress, "add me")
	fs.IntVar(&c.Port, "port", c.Port, "add me")
	fs.BoolVar(&c.SkipRegistration, "skip-registration", c.SkipRegistration, "add me")
	fs.BoolVar(&c.EnableProfiling, "enable-profiling", c.EnableProfiling, "add me")
	fs.StringVar(&c.LogLevel, "log-level", c.LogLevel, "add me")
	fs.StringSliceVar(&c.MongoConfig.ServerList, "mongo-servers", c.MongoConfig.ServerList, "add me")
}
