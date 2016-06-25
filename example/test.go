package main

import (
	"fmt"
	"github.com/mlctrez/udpchan"
	"time"
)

func receive(rec chan []byte) {
	for {
		if msg, ok := <-rec; !ok {
			return
		} else {
			fmt.Printf("got message %q\n", string(msg))
		}
	}
}

func main() {

	rec := make(chan []byte, 10)
	uc, err := udpchan.NewUdpChan(9000, &rec)
	if err != nil {
		panic(err)
	}

	go receive(rec)

	for i := 0; i < 10; i++ {
		uc.Send([]byte(fmt.Sprintf("sending at %v", time.Now().In(time.UTC).String())))
		time.Sleep(500 * time.Millisecond)
	}

	uc.Close()
	close(rec)
}
