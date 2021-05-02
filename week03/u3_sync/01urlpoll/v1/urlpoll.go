package main

import (
	"log"
	"net/http"
	"sync"
	"time"
)

var urls = []string{
	"http://www.google.com/",
	"http://golang.org/",
	"http://blog.golang.org/",
}

type Resource struct {
	url        string
	polling    bool
	lastPolled int64
}

type Resources struct {
	data []*Resource
	lock *sync.Mutex
}

func Poller(res *Resources) {
	for {
		// get the least recently-polled Resource
		// and mark it as being polled
		res.lock.Lock()
		var r *Resource
		for _, v := range res.data {
			if v.polling {
				continue
			}
			if r == nil || v.lastPolled < r.lastPolled {
				r = v
			}
		}
		if r != nil {
			r.polling = true
		}
		res.lock.Unlock()
		if r == nil {
			continue
		}

		// poll the URL
		resp, err := http.Head(r.url)
		if err != nil {
			log.Println("Error", r.url, err)
		}
		log.Printf("poll url %v %v\n",r.url, resp.Status)
		time.Sleep(1000*time.Millisecond)

		// update the Resource's polling and lastPolled
		res.lock.Lock()
		r.polling = false
		r.lastPolled = time.Duration(time.Now().Nanosecond()).Nanoseconds()
		res.lock.Unlock()
	}
}


// Share Memory By Communicating Andrew Gerrand 13 July 2010
// - https://blog.golang.org/codelab-share
func main() {
	res := &Resources{data: make([]*Resource,0), lock:&sync.Mutex{}}
	for _, url := range urls {
		res.data = append(res.data, &Resource{url: url, polling: false, lastPolled: time.Duration(0).Nanoseconds()})
		go Poller(res)
	}
	select {}
}
