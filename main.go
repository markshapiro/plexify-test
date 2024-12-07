package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"plexify-test/handlers"
	"plexify-test/repos"
	"plexify-test/services"
	"plexify-test/utils"
)

func main() {

	jobRepo := repos.NewJobRepo()
	jobProcessor := utils.NewStringJobProcessor()

	jobService := services.NewJobService(jobRepo, jobProcessor)
	jobHandler := handlers.NewJobHandler(jobService)

	jobService.StartWorkers()

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	jobHandler.MountEndpoints(mux)

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

	jobService.StopWorkers()

	fmt.Println("Workers stopped")

}
