package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"plexify-test/app"
	"regexp"
	"strconv"
	"syscall"
	"time"
)

var (
	getStatusRequestArgs = regexp.MustCompile("^/status/([0-9]+)$")
)

func statusHandler(w http.ResponseWriter, r *http.Request) {

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

		resp, err := app.GetJobStatus(int64(id))
		if err != nil {

			if err == app.ErrNotFound {
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

func jobHandler(w http.ResponseWriter, r *http.Request) {

	switch {
	case r.Method == http.MethodPost:

		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
			return
		}

		var newJob app.JobCreateDto

		err = json.Unmarshal(body, &newJob)
		if err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		if len(newJob.Payload) == 0 {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		resp, err := app.JobCreate(newJob)
		if err != nil {

			if err == app.ErrQueueFull {
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

func main() {

	app.Start()

	mux := http.NewServeMux()

	mux.HandleFunc("/job/", jobHandler)

	mux.HandleFunc("/status/", statusHandler)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	go func() {
		fmt.Println("Server is running on port 8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Println("Server error:", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	fmt.Println("Stopping Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Server forced to shutdown:", err)
	} else {
		fmt.Println("Server gracefully stopped")
	}

	fmt.Println("Stopping Workers...")

	app.Stop()

	fmt.Println("Workers stopped")
}
