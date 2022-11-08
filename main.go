package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	log.Println("Running")

	rdb := newRDB()

	s := newServer(rdb)
	addr := fmt.Sprintf(":%s", os.Getenv("PORT"))
	log.Printf("Listening on %s\n", addr)
	err := s.Listen(addr)
	if err != nil {
		log.Fatal(err)
	}

}
