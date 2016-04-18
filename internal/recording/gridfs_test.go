package recording_test

import (
	"gitlab.vailsys.com/jerny/coffer/internal/recording"
	"gitlab.vailsys.com/jerny/coffer/internal/storage/driver/mongo"
	"gopkg.in/mgo.v2/bson"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("AssetRepo interface", func() {
	var objectId bson.ObjectId
	var name string

	It("should LIST asset from GridFS", func() {
		servers := testSession.LiveServers()

		opts := mongo.MongoConfig{
			DB:           "test",
			GridFSPrefix: "testfs",
			ServerList:   servers,
		}

		provider, err := mongo.NewSessionProvider(opts)
		defer provider.Close()

		Expect(err).ToNot(HaveOccurred())

		repo := recording.NewGridFSRepo(opts, provider)

		list, err := repo.ListFiles()
		Expect(err).ToNot(HaveOccurred())
		Expect(len(list)).To(Equal(3))

		for _, r := range list {
			Expect(r.Id).ToNot(BeNil())
		}

		objectId = list[0].Id
		name = list[0].Name
	})

	It("should GET a Recording asset from GridFS by ObjectId", func() {
		servers := testSession.LiveServers()

		opts := mongo.MongoConfig{
			DB:           "test",
			GridFSPrefix: "testfs",
			ServerList:   servers,
		}

		provider, err := mongo.NewSessionProvider(opts)
		defer provider.Close()

		Expect(err).ToNot(HaveOccurred())

		repo := recording.NewGridFSRepo(opts, provider)

		_, err = repo.GetFileById("")
		Expect(err).To(HaveOccurred())

		f, err := repo.GetFileById(objectId)
		Expect(err).ToNot(HaveOccurred())
		Expect(f.Id).To(Equal(objectId))
		Expect(f.Name).To(Equal(name))
	})

	It("should OPEN a Recording asset from GridFS by name", func() {
		servers := testSession.LiveServers()

		opts := mongo.MongoConfig{
			DB:           "test",
			GridFSPrefix: "testfs",
			ServerList:   servers,
		}

		provider, err := mongo.NewSessionProvider(opts)
		defer provider.Close()

		Expect(err).ToNot(HaveOccurred())

		repo := recording.NewGridFSRepo(opts, provider)
		file, err := repo.OpenByName(name)
		Expect(err).ToNot(HaveOccurred())
		defer file.Close()
		Expect(file.Name()).To(Equal(name))
	})

	It("should OPEN a Recording asset from GridFS by id", func() {
		servers := testSession.LiveServers()

		opts := mongo.MongoConfig{
			DB:           "test",
			GridFSPrefix: "testfs",
			ServerList:   servers,
		}

		provider, err := mongo.NewSessionProvider(opts)
		defer provider.Close()

		Expect(err).ToNot(HaveOccurred())

		repo := recording.NewGridFSRepo(opts, provider)
		_, err = repo.OpenById(objectId)
		Expect(err).ToNot(HaveOccurred())
	})

	It("should return an error when OPENING a recordign asset from GridFS by id", func() {
		servers := testSession.LiveServers()

		opts := mongo.MongoConfig{
			DB:           "test",
			GridFSPrefix: "testfs",
			ServerList:   servers,
		}

		provider, err := mongo.NewSessionProvider(opts)
		defer provider.Close()

		Expect(err).ToNot(HaveOccurred())

		repo := recording.NewGridFSRepo(opts, provider)
		_, err = repo.OpenById("blahblah")
		Expect(err).To(HaveOccurred())

		_, ok := err.(recording.RepoError)
		Expect(ok).To(BeTrue())

		_, err = repo.OpenById("blahblah11")
		Expect(err).To(HaveOccurred())

		_, err = repo.OpenById("blahblah1111")
		Expect(err).To(HaveOccurred())
	})

})
