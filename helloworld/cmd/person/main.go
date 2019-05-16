package main

import (
	"context"
	"fmt"
	"github.com/cloudevents/sdk-go"
	"github.com/n3wscott/kubecon/helloworld/pkg/person"
)

func main() {
	client, _ := cloudevents.NewDefaultClient()

	if err := client.StartReceiver(context.Background(), Receive); err != nil {
		panic(err)
	}
}

func Receive(event cloudevents.Event) {
	data := &person.Hello{}
	if err := event.DataAs(data); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Hello, %s!\n", data.Name)

	fmt.Printf("\n---☁️  Event---\n%s\n\n", event)
}
