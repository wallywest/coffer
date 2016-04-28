package recording_test

import (
	"os"

	"gitlab.vailsys.com/vail-cloud-services/coffer/internal/logger"
	"gitlab.vailsys.com/vail-cloud-services/coffer/internal/testutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2"

	"testing"
)

var testSession *mgo.Session

func TestStorage(t *testing.T) {
	logger.TestLogger()
	os.Setenv("CHECK_SESSIONS", "0")
	RegisterFailHandler(Fail)
	RunSpecs(t, "Recording Suite")
	os.Setenv("CHECK_SESSIONS", "1")

}

var _ = BeforeSuite(func() {
	testutil.SetupMongo()
	err := testutil.SeedWithFile("testdata/recordings.json")
	Expect(err).ToNot(HaveOccurred())

	Seed()

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

func Seed() {
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

	var err error

	for _, f := range files {
		err = testutil.SeedAsset(f.FileName, f.FilePath)
		Expect(err).ToNot(HaveOccurred())
	}
}
