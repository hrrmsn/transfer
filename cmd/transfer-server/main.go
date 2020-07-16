package main

import (
	"log"

	"wheely/test/internal/transfer"
)

func main() {
	transferServer := transfer.NewServer()
	if err := transferServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
