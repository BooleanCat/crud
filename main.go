package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/ping", func(_ http.ResponseWriter, _ *http.Request) {})
	_ = http.ListenAndServe(":9092", nil)
}
