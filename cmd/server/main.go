package main

import (
	"log"

	"github.com/MXLange/locker"
)

func main() {

	l, err := locker.NewLockServer(":8080")
	if err != nil {
		log.Fatal(err)
	}

	l.Start()
}
