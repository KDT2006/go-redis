package client

import (
	"context"
	"fmt"
	"log"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient("localhost:5000")
	for i := 0; i < 10; i++ {
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
	}
}
