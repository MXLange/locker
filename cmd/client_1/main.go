package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/MXLange/locker"
)

var (
	serverAddr = flag.String("addr", "localhost:8080", "endereço do servidor")
	id         = flag.String("id", "teste", "id do lock")
	loops      = flag.Int("loops", 10, "número de loops")
	delay      = flag.Int("delay", 10, "delay entre loops em segundos")
)

func main() {
	flag.Parse()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	var wg sync.WaitGroup

	c := locker.NewClient("ws://" + *serverAddr + "/ws?id=")

	for i := 0; i < *loops; i++ {
		wg.Add(1)
		go func(loopNum int) {
			defer wg.Done()

			times := (loopNum % 3) + 1

			// Creates a unique ID for each loop
			uniqueID := fmt.Sprintf("%s-%d", *id, times)
			log.Printf("Loop %d: Connecting to server with id %s", times, uniqueID)

			// Connects to the server and locks
			conn, err := c.Lock(uniqueID)
			if err != nil {
				log.Printf("Loop %d: Error connecting to server: %v", times, err)
			}

			if conn != nil {
				defer conn.Close()
			}

			// Emulates some work
			time.Sleep(time.Duration(*delay*loopNum+1) * time.Second)

			// Unlocks the previous lock
			c.Unlock(conn)

			log.Printf("Loop %d: Lock released for id %s", times, uniqueID)
		}(i)

		// Pequeno delay entre inícios de loops
		time.Sleep(100 * time.Millisecond)
	}

	// Espera por todos os loops ou interrupção
	select {
	case <-interrupt:
	case <-func() chan struct{} {
		c := make(chan struct{})
		go func() {
			wg.Wait()
			close(c)
		}()
		return c
	}():
		log.Println("All jobs completed")
	}
}
