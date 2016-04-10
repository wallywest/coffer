package storage

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.vailsys.com/jerny/coffer/pkg/storage/mongo"
	"gitlab.vailsys.com/jerny/coffer/pkg/testutil"
)

var _ = Describe("MongoRecordingRepo interface", func() {
	It("should be able to GET a recording from the repo", func() {
		testSession := testutil.MongoSession()
		defer testSession.Close()

		servers := testSession.LiveServers()

		opts := mongo.MongoConfig{
			ServerList: servers,
		}

		provider, err := mongo.NewSessionProvider(opts)

		Expect(err).ToNot(HaveOccurred())

		repo := NewMongoRecordingRepo(provider)

		recording, err := repo.Get("12345", "12346")

		Expect(err).ToNot(HaveOccurred())
		fmt.Println(recording)
	})

	It("should be able to LIST all recordings from a repo", func() {
	})

	It("should be able to DELETE a recording from a repo", func() {
	})
})
