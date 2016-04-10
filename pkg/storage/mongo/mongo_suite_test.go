package mongo_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.vailsys.com/jerny/coffer/pkg/testutil"

	"testing"
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
