package server

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"gopkg.in/tylerb/graceful.v1"

	"github.com/julienschmidt/httprouter"

	"gitlab.vailsys.com/jerny/coffer/cmd/coffer/options"
	"gitlab.vailsys.com/jerny/coffer/pkg/logger"
	"gitlab.vailsys.com/jerny/coffer/pkg/storage"
)

type CofferServer struct {
	recordingRepo   recording.Repo
	storageProvider *storage.GridFSProvider
	config          *options.CofferConfig
}

func NewCofferServer(opts *options.CofferConfig, provider *storage.GridFSProvider) (*CofferServer, error) {
	return &CofferServer{
		recordingRepo:   nil,
		storageProvider: provider,
	}, nil
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

	recordings, err := c.RecordingRepo.ListRecordings(accountId)
	if err != nil {
		c.writeError(w, err)
		return
	}

	writeResponseWithBody(w, http.StatusOK, recordings)
}

func (c *CofferServer) getRecording(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//accountId := ps.ByName("accountId")
}

func (c *CofferServer) downloadRecording(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//accountId := ps.ByName("accountId")
}

func (c *CofferServer) writeError(w http.ResponseWriter, err error) {
	logger.Logger.Errorf("error calling coffer api: %v: ", err)

	//if apiErr, ok := err.(api.RecordingError); ok {
	//writeAPIError(w, apiErr.Code, "")
	//return
	//}

	writeAPIError(w, http.StatusInternalServerError, fmt.Errorf("change me"))
}
