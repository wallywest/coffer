package mongo

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.vailsys.com/vail-cloud-services/coffer/internal/testutil"
)

var _ = Describe("Session", func() {
	It("should initialize the master session successfully", func() {
		testSession := testutil.MongoSession()
		defer testSession.Close()

		servers := testSession.LiveServers()

		opts := MongoConfig{
			ServerList: servers,
		}

		provider, err := NewSessionProvider(opts)

		Expect(err).ToNot(HaveOccurred())
		Expect(provider.masterSession).To(BeNil())
		defer provider.Close()

		session, err := provider.GetSession()
		Expect(err).ToNot(HaveOccurred())
		Expect(session).ToNot(BeNil())
		session.Close()
	})
})
