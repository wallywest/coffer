package testutil_test

import (
	"gitlab.vailsys.com/jerny/coffer/pkg/testutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Testutil", func() {
	Context("SeedWithFile", func() {
		It("should throw an error with a file that doesn't exist", func() {
			_, err := testutil.SeedWithFile("test.json")
			Expect(err).To(HaveOccurred())
		})

		It("should import a file from a mongoexport", func() {
			//add options for db/collection
			num, err := testutil.SeedWithFile("testdata/seed.json")
			Expect(err).ToNot(HaveOccurred())
			Expect(num).To(Equal(uint64(5)))

			session := testutil.MongoSession()
			defer session.Close()

			//remove this hard coding for collection/db
			q := session.DB("test").C("recordings").Find(nil)
			total, err := q.Count()
			Expect(err).ToNot(HaveOccurred())
			Expect(total).To(Equal(5))
		})
	})

	Context("SeedAsset", func() {
		It("should throw an error", func() {
			err := testutil.SeedAsset("blah", "blah.wav")
			Expect(err).To(HaveOccurred())
		})

		It("should insert an asset into gridfs", func() {
			err := testutil.SeedAsset("a", "testdata/a.wav")
			Expect(err).ToNot(HaveOccurred())
		})
	})

})
