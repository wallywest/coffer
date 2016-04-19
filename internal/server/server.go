package server

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gitlab.vailsys.com/jerny/coffer/internal/logger"
	"gitlab.vailsys.com/jerny/coffer/internal/options"
	"gitlab.vailsys.com/jerny/coffer/internal/recording"
	"gitlab.vailsys.com/jerny/coffer/internal/registry"
	"gitlab.vailsys.com/jerny/coffer/version"

	"gopkg.in/tylerb/graceful.v1"

	"github.com/cenkalti/backoff"
	"github.com/codegangsta/negroni"
	"github.com/julienschmidt/httprouter"
	"github.com/nats-io/nuid"
)

const DEFAULT_MAX_REGISTRATION_TIME = 10 * time.Second
const assetCacheControlMaxAge = 365 * 24 * time.Hour

type CofferServer struct {
	Config *options.CofferConfig

	recordingRepo recording.RecordingRepo
	assetRepo     recording.AssetRepo

	registryDriver   registry.Registry
	registration     registry.Registration
	skipRegistration bool
}

func NewCofferServer(opts *options.CofferConfig, recordingRepo recording.RecordingRepo, assetRepo recording.AssetRepo) *CofferServer {
	return &CofferServer{
		recordingRepo:    recordingRepo,
		assetRepo:        assetRepo,
		Config:           opts,
		skipRegistration: opts.RegistryConfig.SkipRegistration,
	}
}

func (c *CofferServer) HTTPHandler() http.Handler {
	r := httprouter.New()

	r.PanicHandler = panicHandler()

	r.GET("/health", healthHandler)
	r.GET("/version", versionHandler)
	r.GET("/Accounts/:accountId/Recordings", c.listRecordings)
	r.GET("/Accounts/:accountId/Recordings/:recordingId", c.getRecording)
	r.GET("/Accounts/:accountId/Recordings/:recordingId/Download", c.downloadRecording)

	n := negroni.New(loggerMiddleware())
	n.UseHandler(r)

	return n
}

func (c *CofferServer) Run() error {
	logger.Logger.WithField("service", c.Config.AppName).WithField("port", c.Config.Port).Info("starting")

	location := net.JoinHostPort(c.Config.BindAddress.String(), strconv.Itoa(c.Config.Port))

	srv := &graceful.Server{
		Timeout: 5 * time.Second,
		Server: &http.Server{
			Addr:           location,
			Handler:        c.HTTPHandler(),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}

	go c.registerService()

	return srv.ListenAndServe()
}

func (c *CofferServer) listRecordings(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	accountId := ps.ByName("accountId")

	recordings, _, err := c.recordingRepo.List(accountId)
	if err != nil {
		c.writeError(w, err)
		return
	}

	writeResponseWithBody(w, http.StatusOK, recordings)
}

func (c *CofferServer) getRecording(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	accountId := ps.ByName("accountId")
	recordingId := ps.ByName("recordingId")

	recording, err := c.recordingRepo.Get(accountId, recordingId)
	if err != nil {
		c.writeError(w, err)
		return
	}

	writeResponseWithBody(w, http.StatusOK, recording)
}

func (c *CofferServer) downloadRecording(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	accountId := ps.ByName("accountId")
	recordingId := stripRecordingPrefix(ps.ByName("recordingId"))

	logger.Logger.Debugf("fetching recording file: %v", recordingId)

	gfsmeta, err := c.assetRepo.GetFile(accountId, recordingId)

	if err != nil {
		c.writeError(w, err)
		return
	}

	gfsfile, err := c.assetRepo.OpenById(gfsmeta.Id)

	if err != nil {
		logger.Logger.Debugf("error opening file: %v", err)
		c.writeError(w, err)
		return
	}

	w.Header().Set("ETag", fmt.Sprintf(`"%s"`, gfsfile.Md5)) // If-None-Match handled by ServeContent
	w.Header().Set("Cache-Control", fmt.Sprintf("max-age=%.f", assetCacheControlMaxAge.Seconds()))

	if w.Header().Get("Content-Type") == "" {
		// Set the content type if not already set.
		w.Header().Set("Content-Type", gfsfile.ContentType)
	}

	if w.Header().Get("Content-Length") == "" {
		// Set the content length if not already set.
		w.Header().Set("Content-Length", fmt.Sprint(gfsmeta.Length))
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+gfsmeta.Name+".wav")

	http.ServeContent(w, r, gfsmeta.Name, time.Time{}, gfsfile.FileReader)
}

func (c *CofferServer) writeError(w http.ResponseWriter, err error) {
	logger.Logger.Errorf("error calling coffer api: %v", err)

	if rErr, ok := err.(recording.RepoError); ok {
		writeAPIError(w, rErr.Status, rErr)
		return
	}

	writeAPIError(w, http.StatusInternalServerError, fmt.Errorf("change me"))
}

func (c *CofferServer) registerService() {
	if c.skipRegistration {
		return
	}

	rConfig := c.Config.RegistryConfig
	conf := map[string]string{
		"address": rConfig.Nodes[0],
	}

	driver, err := registry.NewRegistry(rConfig.Type, conf)
	if err != nil {
		logger.Logger.Errorf("error configuring registry %s", c.Config.AppName)
	}

	p := strconv.Itoa(c.Config.Port)

	reg := registry.Registration{
		Name:    c.Config.AppName,
		Port:    p,
		Address: c.Config.AdvertiseAddress.String(),
		Id:      c.Config.AppName + ":" + nuid.Next(),
	}

	c.registration = reg
	c.registryDriver = driver

	var count = 0
	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.MaxElapsedTime = DEFAULT_MAX_REGISTRATION_TIME

	operation := func() error {
		count = count + 1
		return driver.Register(reg)
	}

	notifier := func(e error, t time.Duration) {
		logger.Logger.Errorf("error registering service: %s elapsed: %s attempt: %v", e, t, count)
	}

	err = backoff.RetryNotify(operation, expBackoff, notifier)
	if err != nil {
		logger.Logger.Errorf("error registering service: %s elapsed: %s attempt: %v", err, expBackoff.GetElapsedTime(), count)
	}
}

func (c *CofferServer) ShutDown() error {
	logger.Logger.Debug("shutting down")

	if c.skipRegistration {
		return nil
	}

	logger.Logger.Infof("deregisering service: %s", c.registration.Id)
	err := c.registryDriver.DeRegister(c.registration.Id)
	if err != nil {
		logger.Logger.Errorf("error deregistering service: %s err: %s", c.Config.AppName, err)
		return err
	}
	return nil
}

func panicHandler() func(http.ResponseWriter, *http.Request, interface{}) {
	return func(w http.ResponseWriter, r *http.Request, err interface{}) {
		logger.Logger.Error(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	writeResponseWithBody(w, http.StatusOK, nil)
}

func stripRecordingPrefix(recordingId string) string {
	return strings.TrimPrefix(recordingId, "RE")
}

func versionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	writeResponseWithBody(w, http.StatusOK, version.Map)
}
