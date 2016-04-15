package recording_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.vailsys.com/jerny/coffer/recording"
	"gitlab.vailsys.com/jerny/coffer/storage/driver/mongo"
)

var _ = Describe("MongoRecordingRepo interface", func() {
	Context("LIST", func() {
		It("should return an empty LIST for a missing accountId", func() {
			servers := testSession.LiveServers()

			opts := mongo.MongoConfig{
				DB:         "test",
				ServerList: servers,
			}

			provider, err := mongo.NewSessionProvider(opts)
			defer provider.Close()

			Expect(err).ToNot(HaveOccurred())

			repo := recording.NewMongoRecordingRepo(opts, provider)
			list, _, err := repo.List("")
			Expect(err).ToNot(HaveOccurred())
			Expect(list).To(BeEmpty())
		})

		It("should be able to LIST all recordings from a repo", func() {
			servers := testSession.LiveServers()

			opts := mongo.MongoConfig{
				DB:         "test",
				ServerList: servers,
			}

			provider, err := mongo.NewSessionProvider(opts)
			defer provider.Close()

			Expect(err).ToNot(HaveOccurred())

			repo := recording.NewMongoRecordingRepo(opts, provider)
			list, _, err := repo.List("AC56445f9d0b977d270d02b7026719484c2b6bf369")
			Expect(err).ToNot(HaveOccurred())
			Expect(list).ToNot(BeEmpty())
			Expect(len(list)).To(Equal(65))
			list2, _, err := repo.List("ACa57d943eba574316d2769ae146f8b34e2810f3db")
			Expect(err).ToNot(HaveOccurred())
			Expect(list2).ToNot(BeEmpty())
			Expect(len(list2)).To(Equal(2))
		})
	})

	Context("GET", func() {
		It("should be able to GET a recording from the repo", func() {
			servers := testSession.LiveServers()

			opts := mongo.MongoConfig{
				DB:         "test",
				ServerList: servers,
			}

			provider, err := mongo.NewSessionProvider(opts)
			defer provider.Close()

			Expect(err).ToNot(HaveOccurred())

			repo := recording.NewMongoRecordingRepo(opts, provider)

			recording, err := repo.Get("AC56445f9d0b977d270d02b7026719484c2b6bf369", "RE3a01435a34c2f288d2804d14f48e2731fbfb72bf")

			Expect(err).ToNot(HaveOccurred())
			Expect(recording).ToNot(BeNil())
			Expect(recording.RecordingId).To(Equal("RE3a01435a34c2f288d2804d14f48e2731fbfb72bf"))
			Expect(recording.AccountId).To(Equal("AC56445f9d0b977d270d02b7026719484c2b6bf369"))
			Expect(recording.DownloadUrl).To(Equal("/Accounts/AC56445f9d0b977d270d02b7026719484c2b6bf369/recordings/3a01435a34c2f288d2804d14f48e2731fbfb72bf"))
			Expect(recording.FileId).To(Equal("3a01435a34c2f288d2804d14f48e2731fbfb72bf"))
			Expect(recording.FileName).To(Equal("callrec_0_O143_172.20.152.36_5237_1.0.14_1454604276.wav"))
			Expect(recording.MimeType).To(Equal("audio/wav"))
			Expect(recording.FileSize).To(Equal(int64(38124)))
		})
	})

	It("should return an error with a missed lookup", func() {
		servers := testSession.LiveServers()

		opts := mongo.MongoConfig{
			DB:         "test",
			ServerList: servers,
		}

		provider, err := mongo.NewSessionProvider(opts)
		defer provider.Close()

		Expect(err).ToNot(HaveOccurred())

		repo := recording.NewMongoRecordingRepo(opts, provider)

		_, err = repo.Get("AC56445f9d0b977d270d02b7026719484c2b6bf369", "RE5")
		Expect(err).To(HaveOccurred())
		_, ok := err.(recording.Error)
		Expect(ok).To(BeTrue())

		_, err = repo.Get("AC4", "RE3a01435a34c2f288d2804d14f48e2731fbfb72bf")
		Expect(err).To(HaveOccurred())
		_, ok = err.(recording.Error)
		Expect(ok).To(BeTrue())
	})

	Context("DELETE", func() {
		XIt("should be able to DELETE a recording from a repo", func() {
		})
	})
})
