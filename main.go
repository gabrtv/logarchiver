package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"github.com/nats-io/go-nats"
)

// BufferMax size of in-memory buffer before flushing to storage
const BufferMax = 8

func main() {
	log.Println("starting log archiver...")

	cfg, err := parseConfig()
	if err != nil {
		log.Fatalf("failed to parse configuration: %v\n", err)
	}

	c, err := nats.Connect(cfg.QueueURL)
	ckerr(err, fatal)
	defer c.Close()

	r, w := io.Pipe()
	defer w.Close()

	var wg sync.WaitGroup

	wg.Add(2)
	go recv(c, w, &wg)
	go send(r, &wg)
	wg.Wait()

}

func recv(c *nats.Conn, w io.Writer, wg *sync.WaitGroup) {
	defer wg.Done()

	c.Subscribe("topic", func(msg *nats.Msg) {
		w.Write(msg.Data)
	})
	for c.IsConnected() && !c.IsClosed() {
	}
}

func send(r io.Reader, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		_, err := io.CopyN(os.Stdout, r, BufferMax)
		if err == io.EOF {
			return
		}
		ckerr(err, debug)
	}
}

func ckerr(err error, fn func(error)) {
	if err != nil {
		fn(err)
	}
}

func debug(err error) {
	fmt.Fprintf(os.Stderr, "error: %v\n", err)
}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "fatal: %v\n", err)
	os.Exit(1)
}
