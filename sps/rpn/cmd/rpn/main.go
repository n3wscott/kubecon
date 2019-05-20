package main

import (
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/n3wscott/kubecon/sps/rpn/pkg/rpn"
	"log"
)

// 15 7 1 1 + − ÷ 3 × 2 1 1 + + − =
// 5

/*
 curl -H "ce-specversion: 0.3" \
-H "ce-source: curl" \
-H "ce-type: sps.demo.rpn.expression" \
-H "content-type: application/json" \
-X POST -d '{"exp":"4 5 + ="}' \
http://localhost:8080/
*/

type Expression struct {
	Exp string `json:"exp,omitempty"`
}

type Result struct {
	Value string `json:"value,omitempty"`
}

func main() {
	c, err := cloudevents.NewDefaultClient()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	log.Printf("will listen on :8080\n")
	log.Fatalf("failed to start receiver: %s",
		c.StartReceiver(context.Background(), calculate))
}

const (
	expressionType = "sps.demo.rpn.expression"
	resultType     = "sps.demo.rpn.result"
)

func calculate(event cloudevents.Event, resp *cloudevents.EventResponse) {
	data := &Expression{}
	if err := event.DataAs(data); err != nil {
		fmt.Printf("failed to get data as Example: %s\n", err.Error())
		return
	}

	// calculate one step using reverse polish notation
	exp, done := rpn.Calculate(data.Exp)

	r := cloudevents.NewEvent()
	r.SetSource("github.com/n3wscott/kubecon/sps/rpn/cmd/rpn/")
	r.SetSubject(event.Subject())
	if done {
		r.SetType(resultType)
		_ = r.SetData(&Result{
			Value: exp,
		})
	} else {
		r.SetType(expressionType)
		_ = r.SetData(&Expression{
			Exp: exp,
		})
	}

	resp.RespondWith(200, &r)
}
