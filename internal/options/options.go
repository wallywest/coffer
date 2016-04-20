package options

import (
	"net"

	"github.com/spf13/pflag"
	"gitlab.vailsys.com/vail-cloud-services/coffer/internal/registry"
	mongostorage "gitlab.vailsys.com/vail-cloud-services/coffer/internal/storage/driver/mongo"
)

type CofferConfig struct {
	AppName          string
	BindAddress      net.IP
	AdvertiseAddress net.IP
	Port             int
	HeartBeatTTL     string
	DiscoveryMode    string
	DiscoveryTTL     string
	EnableProfiling  bool
	LogLevel         string
	MongoConfig      mongostorage.MongoConfig
	RegistryConfig   registry.Config
}

func NewCofferConfig() *CofferConfig {
	c := &CofferConfig{
		AppName:         "coffer",
		Port:            6000,
		HeartBeatTTL:    "15s",
		BindAddress:     net.ParseIP("0.0.0.0"),
		DiscoveryMode:   "static",
		DiscoveryTTL:    "15s",
		LogLevel:        "INFO",
		EnableProfiling: false,
		MongoConfig: mongostorage.MongoConfig{
			DB:           "vcsdb",
			GridFSPrefix: "vcsfs",
			ServerList:   []string{"localhost:27017"},
		},
		RegistryConfig: registry.DefaultConfig(),
	}

	return c
}

func (c *CofferConfig) AddFlags(fs *pflag.FlagSet) {
	fs.IPVar(&c.BindAddress, "bind-address", c.BindAddress, "add me")
	fs.IPVar(&c.AdvertiseAddress, "advertise-address", c.AdvertiseAddress, "add me")
	fs.IntVar(&c.Port, "port", c.Port, "add me")
	fs.BoolVar(&c.EnableProfiling, "enable-profiling", c.EnableProfiling, "add me")
	fs.StringVar(&c.LogLevel, "log-level", c.LogLevel, "add me")
	fs.StringSliceVar(&c.MongoConfig.ServerList, "mongo-servers", c.MongoConfig.ServerList, "add me")
	fs.StringVar(&c.MongoConfig.DB, "mongo-db", c.MongoConfig.DB, "add me")
	fs.StringVar(&c.MongoConfig.GridFSPrefix, "mongo-prefix", c.MongoConfig.GridFSPrefix, "add me")
	fs.StringVar(&c.RegistryConfig.Type, "registry-type", c.RegistryConfig.Type, "add me")
	fs.StringSliceVar(&c.RegistryConfig.Nodes, "registry-nodes", c.RegistryConfig.Nodes, "add me")
	fs.BoolVar(&c.RegistryConfig.SkipRegistration, "skip-registration", c.RegistryConfig.SkipRegistration, "add me")
}
