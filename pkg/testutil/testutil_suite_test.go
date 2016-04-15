package testutil_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlab.vailsys.com/jerny/coffer/pkg/testutil"

	"testing"
)

func TestTestutil(t *testing.T) {
	os.Setenv("CHECK_SESSIONS", "0")
	RegisterFailHandler(Fail)
	RunSpecs(t, "Testutil Suite")
	os.Setenv("CHECK_SESSIONS", "1")
}

var _ = BeforeSuite(func() {
	testutil.SetupMongo()
})

var _ = AfterSuite(func() {
	testutil.WipeMongo()
	testutil.StopMongo()
})
