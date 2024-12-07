package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strconv"

	"plexify-test/models"
	"plexify-test/services"
)

var (
	getStatusRequestArgs = regexp.MustCompile("^/status/([0-9]+)$")
)

type JobHandler struct {
	service services.JobService
}

func NewJobHandler(service services.JobService) JobHandler {
	return JobHandler{service}
}

func (t JobHandler) MountEndpoints(mux *http.ServeMux) {
	mux.HandleFunc("/job/", t.jobHandler)
	mux.HandleFunc("/status/", t.statusHandler)
}

func (t JobHandler) statusHandler(w http.ResponseWriter, r *http.Request) {

	switch {
	case r.Method == http.MethodGet:

		args := getStatusRequestArgs.FindStringSubmatch(r.URL.Path)

		if len(args) < 2 {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(args[1])
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		resp, err := t.service.GetJobStatus(int64(id))
		if err != nil {

			if err == services.ErrNotFound {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}

			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		writeResponse(w, resp, http.StatusOK)

	default:
		notFoundHandler(w)
	}
}

func (t JobHandler) jobHandler(w http.ResponseWriter, r *http.Request) {

	switch {
	case r.Method == http.MethodPost:

		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		var newJob models.JobCreateDto

		err = json.Unmarshal(body, &newJob)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if len(newJob.Payload) == 0 {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		resp, err := t.service.JobCreate(newJob)
		if err != nil {

			if err == services.ErrQueueFull {
				http.Error(w, "Unavailable", http.StatusServiceUnavailable)
				return
			}

			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		writeResponse(w, resp, http.StatusAccepted)

	default:
		notFoundHandler(w)
	}
}

func writeResponse(w http.ResponseWriter, resp any, status int) {
	b, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Write(b)
}

func notFoundHandler(w http.ResponseWriter) {
	http.Error(w, "Not Found", http.StatusNotFound)
}
