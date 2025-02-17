package main

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/redis/go-redis/v9"
)

func TestOfficialRedisClient(t *testing.T) {
	listenAddr := ":5001"
	server := NewServer(Config{ListenAddress: listenAddr})
	go func() {
		log.Fatal(server.Start())
	}()

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost%s", listenAddr),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	defer rdb.Close()
	fmt.Println(rdb)

	testcases := map[string]string{
		"foo":    "bar",
		"go":     "pher",
		"game":   "play",
		"number": "10000",
	}

	for key, val := range testcases {
		if err := rdb.Set(context.Background(), key, val, 0).Err(); err != nil {
			log.Fatal(err)
		}

		newVal, err := rdb.Get(context.Background(), key).Result()
		if err != nil {
			t.Fatal(err)
		}

		if newVal != val {
			t.Fatalf("expected %s but got %s", val, newVal)
		}
	}
}

func TestMapToRESP(t *testing.T) {
	in := map[string]string{
		"first":  "1",
		"second": "2",
	}

	out := writeRespMap(in)
	fmt.Println(string(out))
}
