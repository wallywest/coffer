package server

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"gopkg.in/tylerb/graceful.v1"

	"github.com/julienschmidt/httprouter"

	"gitlab.vailsys.com/jerny/coffer/cmd/coffer/options"
	"gitlab.vailsys.com/jerny/coffer/pkg/logger"
	"gitlab.vailsys.com/jerny/coffer/recording"
)

const assetCacheControlMaxAge = 365 * 24 * time.Hour

type CofferServer struct {
	recordingRepo recording.RecordingRepo
	assetRepo     recording.AssetRepo
	Config        *options.CofferConfig
}

func NewCofferServer(opts *options.CofferConfig, recordingRepo recording.RecordingRepo, assetRepo recording.AssetRepo) *CofferServer {
	return &CofferServer{
		recordingRepo: recordingRepo,
		assetRepo:     assetRepo,
		Config:        opts,
	}
}

func (c *CofferServer) HTTPHandler() http.Handler {
	r := httprouter.New()

	r.PanicHandler = panicHandler()

	r.GET("/Accounts/:accountId/Recordings", c.listRecordings)
	r.GET("/Accounts/:accountId/Recordings/:recordingId", c.getRecording)
	r.GET("/Accounts/:accountId/Recordings/:recordingId/Download", c.downloadRecording)

	return r
}

func (c *CofferServer) Run() error {
	logger.Logger.Info("running service %s", c.Config.AppName)

	location := net.JoinHostPort(c.Config.BindAddress.String(), strconv.Itoa(c.Config.Port))

	srv := &graceful.Server{
		Timeout: 10 * time.Second,
		Server: &http.Server{
			Addr:           location,
			Handler:        c.HTTPHandler(),
			MaxHeaderBytes: 1 << 20,
		},
	}

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

	gfsfile, err := c.assetRepo.GetFile(accountId, recordingId)

	if err != nil {
		logger.Logger.Debugf("error finding file: %v", err)
		c.writeError(w, err)
		return
	}

	reader, err := c.assetRepo.OpenById(gfsfile.Id)

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
		w.Header().Set("Content-Length", fmt.Sprint(gfsfile.Length))
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+gfsfile.Name+".wav")

	http.ServeContent(w, r, gfsfile.Name, time.Time{}, reader)
}

func (c *CofferServer) writeError(w http.ResponseWriter, err error) {
	logger.Logger.Errorf("error calling coffer api: %v", err)

	//if apiErr, ok := err.(api.RecordingError); ok {
	//writeAPIError(w, apiErr.Code, "")
	//return
	//}

	writeAPIError(w, http.StatusInternalServerError, fmt.Errorf("change me"))
}

func panicHandler() func(http.ResponseWriter, *http.Request, interface{}) {
	return func(w http.ResponseWriter, r *http.Request, err interface{}) {
		panic(err)
		logger.Logger.Error(err)
	}
}

func stripRecordingPrefix(recordingId string) string {
	return strings.TrimPrefix(recordingId, "RE")
}
