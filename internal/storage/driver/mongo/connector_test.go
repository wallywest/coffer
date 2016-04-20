package mongo

import (
	"time"

	"gitlab.vailsys.com/vail-cloud-services/coffer/internal/testutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Connector", func() {
	var connector *DefaultConnector

	It("should set the addrs and dial timeout without errors", func() {
		expectedTimeout := time.Duration(DefaultDialTimeoutSeconds * time.Second)
		connector = &DefaultConnector{}

		opts := MongoConfig{
			ServerList: []string{"host1", "host2"},
		}

		err := connector.Configure(opts)
		Expect(err).ToNot(HaveOccurred())
		Expect(connector.dialInfo.Timeout).To(Equal(expectedTimeout))
	})

	It("should call GetSession() with a running mongod and connect", func() {
		session := testutil.MongoSession()
		defer session.Close()

		servers := session.LiveServers()
		connector = &DefaultConnector{}

		opts := MongoConfig{
			ServerList: servers,
		}

		err := connector.Configure(opts)
		Expect(err).ToNot(HaveOccurred())

		connector, err := connector.NewSession()
		Expect(err).ToNot(HaveOccurred())
		Expect(connector).ToNot(BeNil())
		connector.Close()
	})

	XIt("should support authentication mechanisms", func() {})
})
