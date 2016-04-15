package server_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"gitlab.vailsys.com/jerny/coffer/options"
	"gitlab.vailsys.com/jerny/coffer/pkg/logger"
	"gitlab.vailsys.com/jerny/coffer/recording"
	"gitlab.vailsys.com/jerny/coffer/server"
	"gitlab.vailsys.com/jerny/coffer/storage/driver/mongo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	Context("/Accounts/:accountId/Recordings", func() {
		XIt("should be able to list the recording records", func() {
			accountId := "AC56445f9d0b977d270d02b7026719484c2b6bf369"

			opts := options.NewCofferConfig()
			opts.MongoConfig.DB = "test"
			opts.MongoConfig.GridFSPrefix = "testfs"
			opts.MongoConfig.ServerList = testSession.LiveServers()

			logger.SetLogLevel("DEBUG")

			provider, err := mongo.NewSessionProvider(opts.MongoConfig)
			Expect(err).ToNot(HaveOccurred())
			defer provider.Close()

			rrepo := recording.NewMongoRecordingRepo(opts.MongoConfig, provider)
			arepo := recording.NewGridFSRepo(opts.MongoConfig, provider)
			s := server.NewCofferServer(opts, rrepo, arepo)

			ts := httptest.NewServer(s.HTTPHandler())
			defer ts.Close()

			validUrl := ts.URL + "/Accounts/" + accountId + "/Recordings"

			req, _ := http.NewRequest("GET", validUrl, nil)

			res, err := http.DefaultClient.Do(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(http.StatusOK))
		})

	})
	Context("/Accounts/:accountId/Recordings/:recordingId", func() {})
	Context("/Accounts/:accountId/Recordings/:recordingId/Download", func() {
		It("should return a not found for an invalid resource", func() {
			opts := options.NewCofferConfig()
			provider, err := mongo.NewSessionProvider(opts.MongoConfig)
			Expect(err).ToNot(HaveOccurred())
			defer provider.Close()

			rrepo := recording.NewMongoRecordingRepo(opts.MongoConfig, provider)
			arepo := recording.NewGridFSRepo(opts.MongoConfig, provider)
			s := server.NewCofferServer(opts, rrepo, arepo)

			ts := httptest.NewServer(s.HTTPHandler())
			defer ts.Close()

			invalidUrl := ts.URL + "/Accounts/ACa57d943eba574316d2769ae146f8b34e2810f3db/Recordings/RE001/Download"

			req, _ := http.NewRequest("GET", invalidUrl, nil)

			res, err := http.DefaultClient.Do(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(http.StatusNotFound))
		})

		It("should serve up an asset for download", func() {
			opts := options.NewCofferConfig()
			opts.MongoConfig.DB = "test"
			opts.MongoConfig.GridFSPrefix = "testfs"
			opts.MongoConfig.ServerList = testSession.LiveServers()

			provider, err := mongo.NewSessionProvider(opts.MongoConfig)
			Expect(err).ToNot(HaveOccurred())
			defer provider.Close()

			rrepo := recording.NewMongoRecordingRepo(opts.MongoConfig, provider)
			arepo := recording.NewGridFSRepo(opts.MongoConfig, provider)
			s := server.NewCofferServer(opts, rrepo, arepo)

			ts := httptest.NewServer(s.HTTPHandler())
			defer ts.Close()

			for _, seed := range RecordingSeeds {
				accountId := seed.AccountId
				recordingId := "RE" + seed.Id

				validUrl := ts.URL + "/Accounts/" + accountId + "/Recordings/" + recordingId + "/Download"

				req, _ := http.NewRequest("GET", validUrl, nil)

				res, err := http.DefaultClient.Do(req)
				Expect(err).ToNot(HaveOccurred())
				Expect(res.StatusCode).To(Equal(http.StatusOK))

				out, err := ioutil.ReadAll(res.Body)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(out)).To(Equal(int(seed.Size)))
			}
		})
	})
})
