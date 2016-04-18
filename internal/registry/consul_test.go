package registry_test

import (
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"gitlab.vailsys.com/jerny/coffer/internal/registry"
	"gitlab.vailsys.com/jerny/coffer/internal/testutil"

	"github.com/nats-io/nuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Consul Backend", func() {
	It("should throw an error when the consul cluster is down", func() {
		consul, err := registry.NewRegistry("consul", nil)
		Expect(err).ToNot(HaveOccurred())

		reg := registry.Registration{
			Name:    "test",
			Port:    "2000",
			Address: "0.0.0.0",
			Id:      nuid.Next(),
		}

		err = consul.Register(reg)
		Expect(err).To(HaveOccurred())
	})

	It("should register the service with the consul cluster", func() {
		t, _ := GinkgoT().(*testing.T)

		cluster := testutil.NewConsulCluster(t)
		defer cluster.Stop()

		conf := map[string]string{
			"address": cluster.Leader.HTTPAddr,
		}

		consul, err := registry.NewRegistry("consul", conf)
		Expect(err).ToNot(HaveOccurred())

		reg := registry.Registration{
			Name:    "test",
			Port:    "2000",
			Address: "0.0.0.0",
			Id:      nuid.Next(),
		}

		err = consul.Register(reg)
		Expect(err).ToNot(HaveOccurred())

		err = consul.DeRegister(reg.Id)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should hit the health endpoint to verify the service is running", func(done Done) {
		c := make(chan int, 0)

		mux := http.NewServeMux()
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
			c <- http.StatusOK
		})

		ts := httptest.NewServer(mux)
		defer ts.Close()

		t, _ := GinkgoT().(*testing.T)

		cluster := testutil.NewConsulCluster(t)
		defer cluster.Stop()

		conf := map[string]string{
			"address": cluster.Leader.HTTPAddr,
		}

		consul, err := registry.NewRegistry("consul", conf)
		Expect(err).ToNot(HaveOccurred())

		u, _ := url.Parse(ts.URL)
		addr, port, _ := net.SplitHostPort(u.Host)

		reg := registry.Registration{
			Name:    "test",
			Port:    port,
			Address: addr,
			Id:      nuid.Next(),
		}

		err = consul.Register(reg)
		Expect(err).ToNot(HaveOccurred())

		Expect(<-c).To(Equal(200))
		close(done)
	}, 20)
})
