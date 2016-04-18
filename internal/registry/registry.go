package registry

import "fmt"

type Config struct {
	Type             string
	Nodes            []string
	SkipRegistration bool
}

func DefaultConfig() Config {
	return Config{
		Type:             "inmem",
		Nodes:            []string{},
		SkipRegistration: false,
	}
}

type Registry interface {
	Register(Registration) error
	DeRegister(Id string) error
	Type() string
}

type Registration struct {
	Name    string
	Address string
	Port    string
	Id      string
	Tags    []string
}

type Factory func(map[string]string) (Registry, error)

var RegistryMap = map[string]Factory{
	"inmem": func(map[string]string) (Registry, error) {
		return NewInMemRegistry(), nil
	},
	"consul": newConsulRegistry,
}

func NewRegistry(t string, conf map[string]string) (Registry, error) {
	f, ok := RegistryMap[t]
	if !ok {
		return nil, fmt.Errorf("invalid registry backend")
	}
	return f(conf)
}
