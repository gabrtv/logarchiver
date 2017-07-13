package main

import (
	"io"
	"log"
	"os"

	"sync"

	"github.com/nats-io/nats"
)

const (
	appName = "logarchiver"
)

func main() {
	log.Println("starting log archiver...")

	natsConfig, err := parseNATSConfig()
	if err != nil {
		log.Fatalf("failed to parse NATS configuration: %v\n", err)
	}

	minioConfig, err := parseMinioConfig()
	if err != nil {
		log.Fatalf("failed to parse Minio configuration: %v\n", err)
	}

	r, w := io.Pipe()
	defer w.Close()

	var wg sync.WaitGroup

	wg.Add(2)
	go recv(natsConfig, w, &wg)
	go send(minioConfig, r, &wg)
	wg.Wait()

}

func recv(cfg *natsConfig, w io.Writer, wg *sync.WaitGroup) {
	defer wg.Done()

	c, err := nats.Connect(cfg.URL)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v\n", err)
	}
	defer c.Close()

	c.Subscribe(cfg.Topic, func(msg *nats.Msg) {
		w.Write(msg.Data)
	})
	for c.IsConnected() && !c.IsClosed() {
	}
}

func send(cfg *minioConfig, r io.Reader, wg *sync.WaitGroup) {
	defer wg.Done()

	mc, err := newMinioClient(cfg)
	if err != nil {
		log.Fatalf("failed to connect to Minio: %v\n", err)
	}

	err = createBucket(cfg, mc)
	if err != nil {
		log.Fatalf("failed to create bucket: %v\n", err)
	}

	for {
		_, err := io.Copy(os.Stdout, r)
		if err == io.EOF {
			log.Println("returning after EOF")
			return
		}
		if err != nil {
			log.Fatalf("failed to io.copy: %v\n", err)
			os.Exit(1)
		}
	}

}
