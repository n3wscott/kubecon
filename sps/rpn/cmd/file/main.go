package main

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go"
	"github.com/kelseyhightower/envconfig"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type envConfig struct {
	Broker string `envconfig:"BROKER" default:"http://localhost:8080/" required:"true"`
}

type fileReader struct {
	ce cloudevents.Client
}

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}

	c, err := NewCloudEventsClient(env.Broker)
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	f := &fileReader{
		ce: c,
	}

	log.Printf("will listen on :8080\n")
	log.Fatalf("failed to start receiver: %s", c.StartReceiver(context.Background(), f.ReadFile))
}

func NewCloudEventsClient(target string) (cloudevents.Client, error) {
	t, err := cloudevents.NewHTTPTransport(cloudevents.WithBinaryEncoding(), cloudevents.WithTarget(target))
	if err != nil {
		return nil, err
	}
	c, err := cloudevents.NewClient(t, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (f *fileReader) ReadFile(event cloudevents.Event) {

	data := &gcs{}
	_ = event.DataAs(data)

	if data.Attributes.EventType == "OBJECT_DELETE" {
		return
	}

	client, err := storage.NewClient(context.Background())
	if err != nil {
		fmt.Printf("failed to create client: %v", err)
		return
	}
	defer client.Close()

	bucket := client.Bucket(data.Attributes.BucketID)
	filename := data.Attributes.ObjectID

	fmt.Printf("reading: %v/%v/\n", bucket, filename)

	rc, err := bucket.Object(filename).NewReader(context.Background())
	if err != nil {
		fmt.Printf("unable to open file from bucket %v, file %q: %v", bucket, filename, err)
		return
	}
	defer rc.Close()
	slurp, err := ioutil.ReadAll(rc)
	if err != nil {
		fmt.Printf("unable to read data from bucket %v, file %q: %v", bucket, filename, err)
		return
	}

	file := string(slurp)

	for i, exp := range strings.Split(file, "\n") {
		fmt.Printf("exp: %q\n", exp)

		ne := cloudevents.NewEvent()
		ne.SetSubject(fmt.Sprintf("%s[%d]#%s", filename, i, event.ID()))
		ne.SetSource("github.com/n3wscott/kubecon/sps/rpn/cmd/file/")
		ne.SetType("sps.demo.rpn.expression")
		_ = ne.SetData(&Expression{
			Exp: exp,
		})

		go func(ne cloudevents.Event) {
			_, _ = f.ce.Send(context.Background(), ne) // TODO
		}(ne)
	}
}

type Expression struct {
	Exp string `json:"exp,omitempty"`
}

type gcs struct {
	Attributes attributes `json:"Attributes,omitempty"`
}

type attributes struct {
	BucketID  string `json:"bucketId,omitempty"`
	ObjectID  string `json:"objectId,omitempty"`
	EventType string `json:"eventType,omitempty"`
}

/*
   "Attributes": {
     "bucketId": "reverse-polish-notation",
     "ce-type": "google.gcs",
     "eventTime": "2019-05-13T22:34:44.952580Z",
     "eventType": "OBJECT_METADATA_UPDATE",
     "notificationConfig": "projects/_/buckets/reverse-polish-notation/notificationConfigs/1",
     "objectGeneration": "1557786739057215",
     "objectId": "simple.txt",
     "payloadFormat": "JSON_API_V1"
   },
*/
