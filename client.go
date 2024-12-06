package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"os/signal"
	"plexify-test/app"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	numClients        = 10
	RequestIntervalMS = 500
)

var (
	mtx sync.Mutex
	wg  sync.WaitGroup
	ids []int64

	jobCreateRequests503 int64 = 0

	sucessfulJobRequests   int64 = 0
	unsucessfulJobRequests int64 = 0

	sucessfulStatusRequests   int64 = 0
	unsucessfulStatusRequests int64 = 0
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(numClients)
	for w := 0; w < numClients; w++ {
		go worker(ctx)
	}

	wg.Add(1)
	go printStatsWorker(ctx)

	<-stop

	fmt.Println("Stopping Client Workers...")

	cancel()

	wg.Wait()

	fmt.Println("Workers stopped")
}

func worker(ctx context.Context) {

	defer wg.Done()

	for true {

		res, status, err := PostRequest(`http://localhost:8080/job/`, []byte(`{"payload":"Process this job!"}`))
		if err != nil {
			panic(err.Error())
		}
		if status == http.StatusAccepted {

			atomic.AddInt64(&sucessfulJobRequests, 1)

			var resp app.JobIDDto
			err = json.Unmarshal(res, &resp)
			if err != nil {
				panic(err.Error())
			}

			mtx.Lock()
			ids = append(ids, resp.JobID)
			mtx.Unlock()
		} else {

			atomic.AddInt64(&unsucessfulJobRequests, 1)

			if status == http.StatusServiceUnavailable {
				atomic.AddInt64(&jobCreateRequests503, 1)
			}
		}

		if rand.IntN(3) == 0 {

			chosenID := ids[rand.IntN(len(ids))]

			res, status, err := GetRequest(fmt.Sprintf(`http://localhost:8080/status/%d`, chosenID))
			if err != nil {
				panic(err.Error())
			}

			if status == http.StatusOK {

				atomic.AddInt64(&sucessfulStatusRequests, 1)

				var resp app.JobStatusDto
				err = json.Unmarshal(res, &resp)
				if err != nil {
					panic(err.Error())
				}

			} else {
				atomic.AddInt64(&unsucessfulStatusRequests, 1)
			}

		}

		select {
		case <-ctx.Done():
			return
		default:
		}

		time.Sleep(time.Millisecond * RequestIntervalMS)
	}

}

func printStatsWorker(ctx context.Context) {
	defer wg.Done()

	for true {

		fmt.Printf(`
-----------------------------------------
sucessful /job requests: %d
unsucessful /job Requests: %d
sucessful /status/:id requests: %d
unsucessful /status/:id Requests: %d
job create requests with status 503: %d
`,
			sucessfulJobRequests, unsucessfulJobRequests,
			sucessfulStatusRequests, unsucessfulStatusRequests,
			jobCreateRequests503)

		select {
		case <-ctx.Done():
			return
		default:
		}

		time.Sleep(time.Second)
	}
}

func GetRequest(url string) ([]byte, int, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	var body []byte

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, 0, err
		}
	}

	return body, resp.StatusCode, nil
}

func PostRequest(url string, data []byte) ([]byte, int, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	var body []byte

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, 0, err
		}
	}

	return body, resp.StatusCode, nil
}
