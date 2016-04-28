package recording

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.vailsys.com/vail-cloud-services/coffer/internal/storage/driver/mongo"
	"gitlab.vailsys.com/vail-cloud-services/coffer/internal/testutil"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ = Describe("Mongo Cursor", func() {
	var testSession *mgo.Session
	var provider *mongo.SessionProvider
	var opts mongo.MongoConfig
	var query *cursorQuery

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

		query = &cursorQuery{
			Session:    testSession,
			DB:         opts.DB,
			Collection: "recordings",
		}
	})

	AfterEach(func() {
		testSession.Close()
		provider.Close()
	})

	It("should return a cursor for the result set", func() {

		query.Selector = bson.M{}
		cursor, err := newRecordingCursor(query)

		Expect(err).ToNot(HaveOccurred())
		Expect(cursor.Id).ToNot(BeNil())
		Expect(cursor.Total()).To(Equal(67))

		recordings, info, err := cursor.NextPage()

		Expect(err).ToNot(HaveOccurred())
		Expect(len(recordings)).To(Equal(20))
		Expect(info.Numpages).To(Equal(4))
		Expect(info.Start).To(Equal(0))
		Expect(info.End).To(Equal(19))
		Expect(info.PageIndex).To(Equal(0))

		Expect(cursor.HasNext()).To(Equal(true))

		recordings, info, err = cursor.NextPage()
		Expect(err).ToNot(HaveOccurred())
		Expect(len(recordings)).To(Equal(20))
		Expect(info.Numpages).To(Equal(4))
		Expect(info.Start).To(Equal(20))
		Expect(info.End).To(Equal(39))
		Expect(info.PageIndex).To(Equal(1))
		Expect(cursor.HasNext()).To(Equal(true))

		recordings, info, err = cursor.NextPage()
		Expect(err).ToNot(HaveOccurred())
		Expect(len(recordings)).To(Equal(20))
		Expect(info.Numpages).To(Equal(4))
		Expect(info.Start).To(Equal(40))
		Expect(info.End).To(Equal(59))
		Expect(info.PageIndex).To(Equal(2))
		Expect(cursor.HasNext()).To(Equal(true))

		recordings, info, err = cursor.NextPage()
		Expect(err).ToNot(HaveOccurred())
		Expect(len(recordings)).To(Equal(7))
		Expect(info.Numpages).To(Equal(4))
		Expect(info.Start).To(Equal(60))
		Expect(info.End).To(Equal(66))
		Expect(info.PageIndex).To(Equal(3))
		Expect(cursor.HasNext()).To(Equal(false))

		err = cursor.Close()
		Expect(err).ToNot(HaveOccurred())
	})

	It("should return a cursor for a result set with 0 recordings", func() {
		query.Selector = bson.M{"accountId": "1"}
		cursor, err := newRecordingCursor(query)

		Expect(err).ToNot(HaveOccurred())
		Expect(cursor.Id).ToNot(BeNil())
		Expect(cursor.Total()).To(Equal(0))

		recordings, info, err := cursor.NextPage()

		Expect(err).ToNot(HaveOccurred())
		Expect(len(recordings)).To(Equal(0))
		Expect(info.Numpages).To(Equal(0))
		Expect(info.Start).To(Equal(0))
		Expect(info.End).To(Equal(0))

		Expect(cursor.HasNext()).To(Equal(false))
	})

})
