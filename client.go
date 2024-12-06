package main

import (
	"bytes"
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

func main() {

	var mtx sync.Mutex
	var ids []int64

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT)

	var wg sync.WaitGroup

	wg.Add(numClients)
	for w := 0; w < numClients; w++ {
		go func() {

			var done = false

			for !done {
				defer wg.Done()

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

				if rand.IntN(5) == 0 {
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
				case <-stop:
					fmt.Print("---")
					done = true
				default:
				}

				time.Sleep(time.Second)
			}
		}()
	}
	wg.Wait()
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
