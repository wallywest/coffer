package testutil_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.vailsys.com/jerny/coffer/pkg/testutil"

	"testing"
)

func TestTestutil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testutil Suite")
}

var _ = BeforeSuite(func() {
	testutil.SetupMongo()
})

var _ = AfterSuite(func() {
	testutil.WipeMongo()
	testutil.StopMongo()
})
