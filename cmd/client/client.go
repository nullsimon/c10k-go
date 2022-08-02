package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"
)

var (
	httpClient *http.Client
	counter    int
)

const (
	MaxIdleConnections int = 100
	RequestTimeout     int = 30
	//concurrent  goroutin
	WORKERNUM int = 4
	//max cpu num
	CPUNUM int = 4
	//all jobs num
	JOBS int = 1000000 * 4
)

func createHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: MaxIdleConnections,
		},
		Timeout: time.Duration(RequestTimeout) * time.Second,
	}

	return client
}

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		testHttpGet()
		results <- j
	}
}

func testHttpGet() {
	url := "http://127.0.0.1:1323"
	counter++
	fmt.Println(counter)
	fetch(url)
}

func init() {
	runtime.GOMAXPROCS(CPUNUM)
	httpClient = createHTTPClient()
}

func main() {

	jobs := make(chan int, JOBS)
	results := make(chan int, JOBS)

	for w := 1; w <= WORKERNUM; w++ {
		go worker(w, jobs, results)
	}

	for j := 1; j <= JOBS; j++ {
		jobs <- j
	}
	close(jobs)

	for a := 1; a <= JOBS; a++ {
		<-results
	}
}

func fetch(url string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	res, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}
