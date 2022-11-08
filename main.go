package main

import (
	"log"
)

func main() {
	log.Println("Running")

	rdb := newRDB()

	s := newServer(rdb)
	addr := ":16379"
	log.Printf("Listening on %s\n", addr)
	err := s.Listen(addr)
	if err != nil {
		log.Fatal(err)
	}

}
