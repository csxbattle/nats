*nats* - библиотека, позволяющая использовать [NATS](https://nats.io/).

## Пример использования
```go
package main

import (
	"context"
	"log"
	"time"

	csxnats "github.com/csxbattle/nats"
	protonats "github.com/zoobr/csxproto/go-kit/nats"
)

func main() {
	// connect
	nc, err := csxnats.Connect(&csxnats.Config{
		ClientName:       "clientName",
		URL:              "nats://localhost:4222",
		ReconnectionTime: 5 * time.Second,
		User:             "user",
		Password:         "password",
	})
	if err != nil {
		panic(err)
	}

	// subscribe go-kit endpoint
	ctx := context.Background()
	goKitEndpoint := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.([]byte) // request if []byte
		log.Print(req)
		return nil, nil
	}

	go protonats.SubscribeEndpoint(nc, "queue-name", ctx, goKitEndpoint)

	// publish
	data := struct {
		Foo string
		Bar int
	}{"qwerty", "123"}

	err := nc.Publish("queue-name", data)
	if err != nil {
		panic(err)
	}
}
```