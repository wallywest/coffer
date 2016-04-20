package mongo_test

import (
	"testing"

	"gitlab.vailsys.com/vail-cloud-services/coffer/internal/testutil"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMongo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mongo Suite")
}

var _ = BeforeSuite(func() {
	testutil.SetupMongo()
})

var _ = AfterSuite(func() {
	testutil.WipeMongo()
	testutil.StopMongo()
})
