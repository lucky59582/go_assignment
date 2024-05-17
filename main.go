package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
)

// Send an API request and send the response body length to the results channel
func fetchAPI(url string, results chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(url)
	if err != nil {
		results <- fmt.Sprintf("Error fetching %s: %v", url, err)
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		results <- fmt.Sprintf("Error reading response from %s: %v", url, err)
		return
	}
	results <- fmt.Sprintf("Fetched %s with response length %d", url, len(body))
}

func main() {
	urls := []string{
		"https://jsonplaceholder.typicode.com/posts/1",
		"https://jsonplaceholder.typicode.com/posts/2",
		"https://jsonplaceholder.typicode.com/posts/3",
		"https://jsonplaceholder.typicode.com/posts/30",
		"https://placeholder.typicode.com/posts/test",
	}

	results := make(chan string, len(urls))
	var wg sync.WaitGroup

	// Fan-out: start a goroutine for each URL
	for _, url := range urls {
		wg.Add(1)
		go fetchAPI(url, results, &wg)
	}

	// Fan-in: collect results
	go func() {
		wg.Wait()
		close(results)
	}()

	for result := range results {
		fmt.Println(result)
	}
}
