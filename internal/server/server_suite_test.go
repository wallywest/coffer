package server_test

import (
	"os"

	"gitlab.vailsys.com/vail-cloud-services/coffer/internal/logger"
	"gitlab.vailsys.com/vail-cloud-services/coffer/internal/testutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"testing"
)

var testSession *mgo.Session

func TestServer(t *testing.T) {
	logger.TestLogger()
	os.Setenv("CHECK_SESSIONS", "0")
	RegisterFailHandler(Fail)
	RunSpecs(t, "Server Suite")
	os.Setenv("CHECK_SESSIONS", "1")
}

var _ = BeforeSuite(func() {
	testutil.SetupMongo()
	err := testutil.SeedWithFile("../../testdata/recordings.json")
	Expect(err).ToNot(HaveOccurred())

	for _, seed := range RecordingSeeds {
		err := seed.Verify()
		//seed, err := NewRecordingSeed(f)
		Expect(err).ToNot(HaveOccurred())

		err = testutil.SeedAssetWithMeta(seed.Id, seed.File, bson.M{"accountId": seed.AccountId, "fileId": seed.Id})
		Expect(err).ToNot(HaveOccurred())
		err = seed.Close()
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

var RecordingSeeds = []*RecordingSeed{
	{
		Id:        "7be3dd50daea70113910b786ce2b6dc9b6a9cf02",
		AccountId: "AC56445f9d0b977d270d02b7026719484c2b6bf369",
	},
	{
		Id:        "19cd21551fcdc1e19c506a510c6cf3cd7a422dc0",
		AccountId: "AC56445f9d0b977d270d02b7026719484c2b6bf369",
	}, {
		Id:        "8f7cb8f2d30bef3e748102b082d0702f38733cf7",
		AccountId: "AC56445f9d0b977d270d02b7026719484c2b6bf369",
	},
}

type RecordingSeed struct {
	Id        string
	AccountId string
	File      *os.File
	Size      int64
}

func (r *RecordingSeed) Verify() error {
	filePath := "../../testdata/files/" + r.Id
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	r.File = f
	stat, err := f.Stat()
	if err != nil {
		return err
	}
	r.Size = stat.Size()
	return nil
}

func (r *RecordingSeed) FilePath() string {
	return "../../testdata/files/" + r.Id
}

func (r *RecordingSeed) Close() error {
	return r.File.Close()
}
