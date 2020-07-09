package main

import (
	"flag"
	"log"
	"wheely/test/internal/delivery"
)

func main() {
	port := flag.String("port", "8080", "port to listen")
	flag.Parse()

	deliveryServer := delivery.NewServer(*port)

	log.Printf("listening on port %s\n", *port)
	if err := deliveryServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
