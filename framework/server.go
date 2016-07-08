package framework

import (
	"encoding/json"
	"net/http"
	"strings"

	"errors"
	"github.com/yanzay/log"
)

type Server interface {
	Start() error
}

type HTTPServer struct {
	address   string
	scheduler Scheduler
}

func NewHttpServer(address string, scheduler Scheduler) *HTTPServer {
	if strings.HasPrefix(address, "http://") {
		address = address[len("http://"):]
	}
	return &HTTPServer{
		address:   address,
		scheduler: scheduler,
	}
}

func (s *HTTPServer) Start() error {
	log.Infof("Starting HTTP server at %s", s.address)
	http.HandleFunc("/api/group/add", s.groupAdd)
	http.HandleFunc("/api/group/list", s.groupList)
	return http.ListenAndServe(s.address, nil)
}

func (s *HTTPServer) groupAdd(w http.ResponseWriter, r *http.Request) {
	cluster := s.scheduler.Cluster()
	queryParams := r.URL.Query()

	groupID := queryParams.Get(ParamGroupID)
	if groupID == "" {
		respond(w, http.StatusBadRequest, ErrGroupIDRequired)
	}

	if cluster.ExistsGroup(groupID) {
		respond(w, http.StatusBadRequest, ErrGroupExists)
	}

	subscription := queryParams.Get(ParamSubscription)
	bootstrapBrokers := queryParams.Get(ParamBootstrapBrokers)

	group := &Group{
		ID:               groupID,
		Subscriptions:    strings.Split(subscription, ","),
		BootstrapBrokers: strings.Split(bootstrapBrokers, ","),
	}

	cluster.AddGroup(group)
	respond(w, http.StatusOK, nil)
}

func (s *HTTPServer) groupList(w http.ResponseWriter, r *http.Request) {
	cluster := s.scheduler.Cluster()

	respond(w, http.StatusOK, cluster.GetGroups())
}

func respond(w http.ResponseWriter, statusCode int, body interface{}) {
	errBody, ok := body.(error)
	if ok {
		body = NewErrorResponse(errBody.Error())
	}

	bytes, err := json.Marshal(body)
	if err != nil {
		panic(err) //this shouldn't happen
	}

	w.WriteHeader(statusCode)
	_, err = w.Write(bytes)
	if err != nil {
		log.Errorf("Http server failed to respond: %s", err)
	}
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewErrorResponse(msg string) *ErrorResponse {
	return &ErrorResponse{
		Error: msg,
	}
}

var (
	ErrGroupIDRequired = errors.New("Missing required parameter " + ParamGroupID)
	ErrGroupExists     = errors.New("Group already exists")
	ErrInternal        = errors.New("An error occurred")
)
