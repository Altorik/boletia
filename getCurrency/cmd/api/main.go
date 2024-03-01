package main

import (
	"boletia/cmd/api/bootstrap"
	"log"
	_ "net/http/pprof"
)

func main() {
	if err := bootstrap.Run(); err != nil {
		log.Fatal(err)
	}
}
