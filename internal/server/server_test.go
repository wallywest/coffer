package server_test

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"

	"gitlab.vailsys.com/vail-cloud-services/coffer/internal/logger"
	"gitlab.vailsys.com/vail-cloud-services/coffer/internal/options"
	"gitlab.vailsys.com/vail-cloud-services/coffer/internal/recording"
	"gitlab.vailsys.com/vail-cloud-services/coffer/internal/server"
	"gitlab.vailsys.com/vail-cloud-services/coffer/internal/storage/driver/mongo"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Server", func() {
	It("should have a health handler registered", func() {
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

		validUrl := ts.URL + "/health"

		req, _ := http.NewRequest("GET", validUrl, nil)

		res, err := http.DefaultClient.Do(req)
		Expect(err).ToNot(HaveOccurred())
		Expect(res.StatusCode).To(Equal(200))
	})

	Context("Service Registration", func() {
		It("should not try to register when registration disabled", func() {
		})
		It("should register with backoff when registration enabled", func() {
		})
	})

	Context("/Accounts/:accountId/Recordings", func() {
		XIt("should be able to list the recording records", func() {
			accountId := "AC56445f9d0b977d270d02b7026719484c2b6bf369"

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

			validUrl := ts.URL + "/Accounts/" + accountId + "/Recordings"

			req, _ := http.NewRequest("GET", validUrl, nil)

			res, err := http.DefaultClient.Do(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(http.StatusOK))
		})

		XIt("should be able to list recordings by callId", func() {
			accountId := "AC56445f9d0b977d270d02b7026719484c2b6bf369"
			callId := "CAc000ffe439109e79fc386bd4140b0c9e75585f55"

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

			validUrl := ts.URL + "/Accounts/" + accountId + "/Calls/" + callId + "/Recordings"

			req, _ := http.NewRequest("GET", validUrl, nil)

			res, err := http.DefaultClient.Do(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(http.StatusOK))
		})

		XIt("should paginate the responses when bigger max limit", func() {
		})

		XIt("should be able to list the recordings by a given cursor", func() {
		})
	})

	Context("/Accounts/:accountId/Recordings/:recordingId", func() {
		It("should be able to get the recording", func() {
			accountId := "AC56445f9d0b977d270d02b7026719484c2b6bf369"
			recordingId := "RE3a01435a34c2f288d2804d14f48e2731fbfb72bf"

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

			validUrl := ts.URL + "/Accounts/" + accountId + "/Recordings/" + recordingId

			req, _ := http.NewRequest("GET", validUrl, nil)

			res, err := http.DefaultClient.Do(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(http.StatusOK))
		})
	})

	Context("/Accounts/:accountId/Recordings/:recordingId/Download", func() {
		It("should return a not found for an invalid resource", func() {
			logger.SetLogLevel("DEBUG")
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

		It("should serve up the asset for streaming", func() {
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

			seed := RecordingSeeds[0]
			accountId := seed.AccountId
			recordingId := "RE" + seed.Id

			validUrl := ts.URL + "/Accounts/" + accountId + "/Recordings/" + recordingId + "/Stream"

			req, _ := http.NewRequest("GET", validUrl, nil)
			res, err := http.DefaultClient.Do(req)
			Expect(err).ToNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			Expect(res.ContentLength).To(Equal(int64(18924)))
			Expect(res.Header.Get("Content-Type")).To(Equal("audio/wav"))
		})
	})
})
