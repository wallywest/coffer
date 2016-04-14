package recording_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.vailsys.com/jerny/coffer/pkg/logger"
	"gitlab.vailsys.com/jerny/coffer/pkg/testutil"
	"gopkg.in/mgo.v2"

	"testing"
)

var testSession *mgo.Session

func TestStorage(t *testing.T) {
	logger.TestLogger()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Recording Suite")
}

var _ = BeforeSuite(func() {
	testutil.SetupMongo()
	_, err := testutil.SeedWithFile("testdata/recordings.json")
	Expect(err).ToNot(HaveOccurred())

	files := []*FileSeed{
		{
			FileName: "a",
			FilePath: "testdata/a.wav",
		},
		{
			FileName: "b",
			FilePath: "testdata/b.wav",
		},
		{
			FileName: "c",
			FilePath: "testdata/c.wav",
		},
	}
	for _, f := range files {
		err = testutil.SeedAsset(f.FileName, f.FilePath)
		Expect(err).ToNot(HaveOccurred())
	}

	testSession = testutil.MongoSession()
})

var _ = AfterSuite(func() {
	if testSession != nil {
		testSession.Close()
	}
	testutil.WipeMongo()
	testutil.StopMongo()
})

type FileSeed struct {
	FileName string
	FilePath string
}
