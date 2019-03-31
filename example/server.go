package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/raohwork/ratelimit"
)

func serverExample() (err error) {
	bucket := ratelimit.NewFromRate(10*ratelimit.KB, 10*ratelimit.KB, 0)
	content := strings.Repeat(".", 100*ratelimit.KB)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// all server threads share 10k bandwidth by using same bucket
		wrappedWriter := bucket.NewWriter(w)
		fmt.Fprint(wrappedWriter, content)
	})
	http.HandleFunc("/10k", func(w http.ResponseWriter, r *http.Request) {
		// each thread has 10k bandwidth
		bucket := ratelimit.NewFromRate(10*ratelimit.KB, 10*ratelimit.KB, 0)
		wrappedWriter := bucket.NewWriter(w)
		fmt.Fprint(wrappedWriter, content)
	})

	return http.ListenAndServe(":8000", nil)
}
