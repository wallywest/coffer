package recording_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.vailsys.com/vail-cloud-services/coffer/internal/recording"
	"gitlab.vailsys.com/vail-cloud-services/coffer/internal/storage/driver/mongo"
	"gitlab.vailsys.com/vail-cloud-services/coffer/internal/testutil"
	"gopkg.in/mgo.v2"
)

var _ = Describe("MongoRecordingRepo interface", func() {
	var testSession *mgo.Session
	var provider *mongo.SessionProvider
	var opts mongo.MongoConfig

	BeforeEach(func() {
		testSession = testutil.MongoSession()
		servers := testSession.LiveServers()

		opts = mongo.MongoConfig{
			DB:         "test",
			ServerList: servers,
		}

		var err error
		provider, err = mongo.NewSessionProvider(opts)
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		testSession.Close()
		provider.Close()
	})

	Context("LIST", func() {
		It("should return an empty LIST for a missing accountId", func() {
			repo := recording.NewMongoRecordingRepo(opts, provider)
			list, _, err := repo.List("")
			Expect(err).ToNot(HaveOccurred())
			Expect(list).To(BeEmpty())
		})

		It("should be able to LIST all recordings from a repo", func() {
			repo := recording.NewMongoRecordingRepo(opts, provider)
			list, _, err := repo.List("AC56445f9d0b977d270d02b7026719484c2b6bf369")
			Expect(err).ToNot(HaveOccurred())
			Expect(list).ToNot(BeEmpty())
			Expect(len(list)).To(Equal(20))
			list2, info, err := repo.List("ACa57d943eba574316d2769ae146f8b34e2810f3db")
			Expect(err).ToNot(HaveOccurred())
			Expect(list2).ToNot(BeEmpty())
			Expect(len(list2)).To(Equal(info.End + 1))
		})

		It("should be able to list by cursorId", func() {
			repo := recording.NewMongoRecordingRepo(opts, provider)
			list, info, err := repo.List("AC56445f9d0b977d270d02b7026719484c2b6bf369")
			Expect(err).ToNot(HaveOccurred())
			Expect(list).ToNot(BeEmpty())
			Expect(len(list)).To(Equal(20))
			Expect(info.CursorId).ToNot(BeNil())

			list, _, err = repo.ListByCursor(info.CursorId)
			Expect(err).ToNot(HaveOccurred())
			Expect(list).ToNot(BeEmpty())
			Expect(len(list)).To(Equal(20))
		})

		It("should be able to LIST all recordings from a repo by callId", func() {
			repo := recording.NewMongoRecordingRepo(opts, provider)
			list, _, err := repo.ListByCall("AC56445f9d0b977d270d02b7026719484c2b6bf369", "CAc000ffe439109e79fc386bd4140b0c9e75585f55")
			Expect(err).ToNot(HaveOccurred())
			Expect(list).ToNot(BeEmpty())
			Expect(len(list)).To(Equal(1))
		})
	})

	Context("GET", func() {
		It("should be able to GET a recording from the repo", func() {
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
		repo := recording.NewMongoRecordingRepo(opts, provider)

		_, err := repo.Get("AC56445f9d0b977d270d02b7026719484c2b6bf369", "RE5")
		Expect(err).To(HaveOccurred())
		_, ok := err.(recording.RepoError)
		Expect(ok).To(BeTrue())

		_, err = repo.Get("AC4", "RE3a01435a34c2f288d2804d14f48e2731fbfb72bf")
		Expect(err).To(HaveOccurred())
		_, ok = err.(recording.RepoError)
		Expect(ok).To(BeTrue())
	})

	XContext("DELETE", func() {
		It("should be able to DELETE a recording from a repo", func() {
			repo := recording.NewMongoRecordingRepo(opts, provider)

			err := repo.Delete("AC56445f9d0b977d270d02b7026719484c2b6bf369", "RE3a01435a34c2f288d2804d14f48e2731fbfb72bf")

			Expect(err).ToNot(HaveOccurred())
		})

		It("should throw an error with a missing recordingId", func() {
			repo := recording.NewMongoRecordingRepo(opts, provider)

			err := repo.Delete("AC56445f9d0b977d270d02b7026719484c2b6bf369", "RE3")

			Expect(err).To(HaveOccurred())
		})
	})
})
