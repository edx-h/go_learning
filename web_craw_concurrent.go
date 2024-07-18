package main

import (
	"fmt"
	"sync"
	"time"
)

type Fetcher interface {
	// Fetch returns the body of URL and
	// a slice of URLs found on that page.
	Fetch(url string) (body string, urls []string, err error)
}

var wg sync.WaitGroup

// Crawl uses fetcher to recursively crawl
// pages starting with url, to a maximum of depth.
func Crawl(url string, depth int, fetcher Fetcher, fetchedCache *SafeCacheMap) {
	// TODO: Fetch URLs in parallel.
	// TODO: Don't fetch the same URL twice.
	// This implementation doesn't do either:
	defer wg.Done()
	if depth <= 0 {
		return
	}

	var body string
	var urls []string
	var err error
	fetchedCache.mu.Lock()

	if res, ok := fetchedCache.cachemap[url]; !ok {
		fmt.Printf("Cache for %s does not exist!\n", url)
		body, urls, err = fetcher.Fetch(url)
		if err != nil {
			defer fetchedCache.mu.Unlock()
			fmt.Println(err)
			fmt.Printf("Newly found: %s with error %q\n", url, err)
			fetchedCache.cachemap[url] = &fakeResult{"", nil}
			fmt.Printf("Stored %s result with error into cache!\n", url)
			return
		} else {
			//defer fetchedCache.mu.Unlock()
			fmt.Printf("Newly found: %s %q\n", url, body)
			fetchedCache.cachemap[url] = &fakeResult{body, urls}
			fmt.Printf("Stored %s result into cache!\n", url)
		}
	} else {
		fmt.Printf("Cache found for %s. No need to fetch again!\n", url)
		if res.body != "" {
			fmt.Printf("Cache result for %s: url exists!\n", url)
			body = res.body
			urls = res.urls
			//defer fetchedCache.mu.Unlock()
		} else {
			fmt.Printf("Cache result for %s: url not exists!\n", url)
			defer fetchedCache.mu.Unlock()
			return
		}
	}

	// urls are not null
	fetchedCache.mu.Unlock()
	fmt.Println("Unlocked!")
	for _, u := range urls {
		wg.Add(1)
		go Crawl(u, depth-1, fetcher, fetchedCache)
	}
	return
}

func main() {
	start := time.Now()
	fmt.Println("Start time is", start.Format("Mon Jan 2 15:04:05.1234 MST 2006"))
	fetchedCache := SafeCacheMap{cachemap: make(map[string]*fakeResult)}
	wg.Add(1)
	Crawl("https://golang.org/", 20, fetcher, &fetchedCache)
	wg.Wait()
	end := time.Now()
	fmt.Println("End time is", end.Format("Mon Jan 2 15:04:05.1234 MST 2006"))
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}
type SafeCacheMap struct {
	cachemap map[string]*fakeResult
	mu       sync.Mutex
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	time.Sleep(1 * time.Second)
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
