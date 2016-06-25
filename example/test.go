package main

import (
	"fmt"
	"github.com/mlctrez/udpchan"
	"time"
)

func main() {

	rec := make(chan []byte, 10)
	uc, err := udpchan.NewUdpChan(9000, &rec)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			msg := <-rec
			if msg != nil {
				fmt.Printf("got message %q\n", string(msg))
			}
		}
	}()

	uc.Send([]byte("hello"))

	time.Sleep(1 * time.Second)

	uc.Close()
	close(rec)

	time.Sleep(1 * time.Second)

}
