package main

import (
	"fmt"
	"log"

	"github.com/opinedajr/go-learning-api-gobank/internal"
)

func main() {
	fmt.Println("You are ready to Go!")

	repo, err := internal.NewPostGresRepository()
	if err != nil {
		log.Fatal(err)
	}

	server := internal.InitializeApiServer(":8888", repo)
	server.Run()
}
