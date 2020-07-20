package main

import (
	"log"

	"wheely/test/pkg/transfer"
	"wheely/test/pkg/transfer/utils"
)

func main() {
	cfg, err := utils.NewConfig()
	if err != nil {
		log.Fatal("Error when initializing config: %s\n", err.Error())
	}

	transferServer := transfer.NewServer(cfg)
	if err := transferServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
