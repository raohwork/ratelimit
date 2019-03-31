package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/raohwork/ratelimit"
)

func mockserver() {
	content := strings.Repeat(".", 100*ratelimit.KB)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, content)
	})

	go func() {
		if err := http.ListenAndServe(":8000", nil); err != nil {
			log.Fatalf("Failed to start web server to test client rate limit: %s", err)
		}
	}()
}

func clientExample() error {
	resp, err := http.Get("http://127.0.0.1:8000/")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bucket := ratelimit.NewFromRate(10*ratelimit.KB, 10*ratelimit.KB, 0)
	wrappedReader := bucket.NewReader(resp.Body)
	if _, err := ioutil.ReadAll(wrappedReader); err != nil {
		return err
	}
	return nil
}
