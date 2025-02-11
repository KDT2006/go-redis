package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/KDT2006/go-redis/client"
)

func TestServerWithMultiplePeers(t *testing.T) {
	server := NewServer(Config{
		ListenAddress: ":5001",
	})
	go func() {
		log.Fatal(server.Start())
	}()
	time.Sleep(time.Millisecond)

	nClients := 10
	wg := sync.WaitGroup{}
	wg.Add(nClients)
	for i := 0; i < nClients; i++ {
		go func() {
			client, err := client.NewClient("localhost:5001")
			if err != nil {
				t.Errorf("NewClient() error: %+v", err)
			}
			defer client.Close()

			setValue := fmt.Sprintf("bar_%d", i)
			log.Println("SET => ", setValue)
			if err := client.Set(context.Background(), fmt.Sprintf("foo_%d", i), setValue); err != nil {
				log.Fatal(err)
			}

			value, err := client.Get(context.Background(), fmt.Sprintf("foo_%d", i))
			if err != nil {
				log.Fatal(err)
			}
			log.Println("GET => ", value)

			if setValue != value {
				t.Errorf("SET:%s and GET:%s values don't match", setValue, value)
			}

			wg.Done()
		}()
	}

	wg.Wait()

	time.Sleep(time.Millisecond)
	if len(server.peers) != 0 {
		t.Fatalf("expected 0 peers, but got %d", len(server.peers))
	}
}
