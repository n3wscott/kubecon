package person

import (
	"fmt"
	"github.com/cloudevents/sdk-go"
)

type Hello struct {
	Name string `json:"name"`
}

func Receive(event cloudevents.Event) {
	data := &Hello{}
	if err := event.DataAs(data); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Hello, %s!\n", data.Name)

	fmt.Printf("\n---☁️  Event---\n%s\n\n", event)
}

/*

curl -X POST -H "Content-Type: application/json" \
  -H "Content-Type: application/json" \
  -H "ce-specversion: 0.2" \
  -H "ce-source: curl-command" \
  -H "ce-type: curl.demo" \
  -H "ce-id: 123-abc" \
  -d '{"name":"Earl"}' \
  http://localhost:8080/


*/
