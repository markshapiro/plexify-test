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
	"syscall"
	"time"
)

const (
	numClients = 1
)

var (
	mtx sync.Mutex
	wg  sync.WaitGroup
	ids []int64
)

func main() {

	ctx, cancel := context.WithCancel(context.Background())

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)
	//signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	wg.Add(numClients)
	for w := 0; w < numClients; w++ {
		go worker(ctx)
	}

	<-stop

	cancel()

	wg.Wait()
}

func worker(ctx context.Context) {

	defer wg.Done()

	for true {

		res, status, err := PostRequest(`http://localhost:8080/job/`, []byte(`{"payload":"Process this job!"}`))
		if err != nil {
			panic(err.Error())
		}
		if status == http.StatusAccepted {
			var resp app.JobIDDto
			err = json.Unmarshal(res, &resp)
			if err != nil {
				panic(err.Error())
			}

			mtx.Lock()
			ids = append(ids, resp.JobID)
			mtx.Unlock()
		}

		if rand.IntN(3) == 0 {

			chosenID := ids[rand.IntN(len(ids))]

			res, status, err := GetRequest(fmt.Sprintf(`http://localhost:8080/status/%d`, chosenID))
			if err != nil {
				panic(err.Error())
			}

			if status == http.StatusAccepted {
				var resp app.JobStatusDto
				err = json.Unmarshal(res, &resp)
				if err != nil {
					panic(err.Error())
				}

				fmt.Println(resp)
			}

		}

		select {
		case <-ctx.Done():
			return
		default:
		}

		time.Sleep(time.Second * 2)
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
