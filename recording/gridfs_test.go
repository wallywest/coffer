package recording_test

import (
	"gopkg.in/mgo.v2/bson"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.vailsys.com/jerny/coffer/recording"
	"gitlab.vailsys.com/jerny/coffer/storage/driver/mongo"
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

})
