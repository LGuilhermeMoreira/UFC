package main

import (
	"fmt"
	"log"
	"sockets/entity"
	"sync"
	"time"
)

func main() {
	wg := sync.WaitGroup{}
	wg.Add(100)
	start := time.Now()
	for i := 0; i < 100; i++ {
		go func() {
			defer wg.Done()
			user := entity.GetUser()
			if err := user.SendRequest(); err != nil {
				log.Println(err)
			}
			fmt.Println(user.Result, i)
		}()
	}
	wg.Wait()
	duration := time.Since(start)
	fmt.Printf("Tempo total: %d ms (%.2f segundos)\n", duration.Milliseconds(), duration.Seconds())
}
