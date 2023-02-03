package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

func finish(finished chan int) {
	finished <- 0
}

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, visited map[string]bool, mu sync.Mutex, finished chan int, fetcher Fetcher) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	defer finish(finished)
	if depth <= 0 {
		return
	}
	mu.Lock()
	_, ok := visited[url]
	mu.Unlock()
	if !ok {
		mu.Lock()
		visited[url] = true
		mu.Unlock()
		body, urls, err := fetcher.Fetch(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("found: %s %q\n", url, body)
		threads := 0
		finished_this := make(chan int)
		for _, u := range urls {
			mu.Lock()
			_, ok := visited[u]
			mu.Unlock()
			if !ok {
				threads++
				go Crawl(u, depth-1, visited, mu, finished_this, fetcher)
			}
		}
		for threads > 0 {
			<-finished_this
			threads--
		}
		return
	}
}

func main() {
	visited := make(map[string]bool)
	finished := make(chan int)
	var mu sync.Mutex
	go Crawl("https://golang.org/", 4, visited, mu, finished, fetcher)
	<-finished
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
