package recording_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.vailsys.com/jerny/coffer/pkg/testutil"
	"gopkg.in/mgo.v2"

	"testing"
)

var testSession *mgo.Session

func TestStorage(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Storage Suite")
}

var _ = BeforeSuite(func() {
	testutil.SetupMongo()
	_, err := testutil.SeedWithFile("testdata/recordings.json")
	Expect(err).ToNot(HaveOccurred())

	err = testutil.SeedAsset("testdata/a.wav")
	Expect(err).ToNot(HaveOccurred())

	testSession = testutil.MongoSession()
})

var _ = AfterSuite(func() {
	if testSession != nil {
		testSession.Close()
	}
	testutil.WipeMongo()
	testutil.StopMongo()
})
