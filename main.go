package main

import (
	"net/http"

	"github.com/BooleanCat/crud/router"
	"github.com/jessevdk/go-flags"
)

type opts struct {
	Store string `long:"store"`
}

func main() {
	var opts opts
	_, err := flags.Parse(&opts)
	if err != nil {
		panic(err)
	}

	router := router.Router(opts.Store)

	_ = http.ListenAndServe(":9092", router)

}
