package registry

import (
	"fmt"
	"net"
	"net/url"
	"strconv"

	"gitlab.vailsys.com/vail-cloud-services/coffer/internal/logger"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/go-cleanhttp"
)

const (
	StatusPassing = "passing"
)

type ConsulRegistry struct {
	client *api.Client
	T      string
}

func newConsulRegistry(conf map[string]string) (Registry, error) {
	consulConf := api.DefaultConfig()

	if addr, ok := conf["address"]; ok {
		consulConf.Address = addr
	}

	transport := cleanhttp.DefaultPooledTransport()

	transport.MaxIdleConnsPerHost = 4
	consulConf.HttpClient.Transport = transport

	client, err := api.NewClient(consulConf)

	if err != nil {
		return nil, fmt.Errorf("consul client setup failed: %s", err)
	}

	c := &ConsulRegistry{
		client: client,
		T:      "consul",
	}

	return c, nil
}

func (c *ConsulRegistry) Register(reg Registration) error {
	p, _ := strconv.Atoi(reg.Port)
	as := &api.AgentServiceRegistration{
		Address: reg.Address,
		Port:    p,
		ID:      reg.Id,
		Name:    reg.Name,
		Tags:    reg.Tags,
		Check:   c.newAgentCheck(reg),
	}

	agent := c.client.Agent()

	err := agent.ServiceRegister(as)
	if err != nil {
		logger.Logger.Debugf("error registering service: %s", err)
		return err
	}

	return nil
}

func (c *ConsulRegistry) DeRegister(id string) error {
	agent := c.client.Agent()

	err := agent.ServiceDeregister(id)

	if err != nil {
		logger.Logger.Debugf("error deregistering service: %s", err)
		return err
	}

	return nil
}

func (c *ConsulRegistry) Type() string {
	return c.T
}

func (c *ConsulRegistry) newAgentCheck(reg Registration) *api.AgentServiceCheck {
	address := net.JoinHostPort(reg.Address, reg.Port)
	u := url.URL{
		Scheme: "http",
		Host:   address,
		Path:   "health",
	}

	logger.Logger.Debugf("registring http check endpoint: %s interval: %s timeout: %s ", u.String(), "15s", "1s")

	return &api.AgentServiceCheck{
		Status:   StatusPassing,
		HTTP:     u.String(),
		Interval: "15s",
		Timeout:  "1s",
	}
}
