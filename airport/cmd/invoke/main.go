package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"github.com/kelseyhightower/envconfig"
)

type envConfig struct {
	Scheme string `envconfig:"WS_SCHEME" default:"ws"`
	Host   string `envconfig:"WS_HOST" required:"true"`
	Path   string `envconfig:"WS_PATH" default:"/"`
}

var (
	participantCount = 3
	offerCount       = 3
)

func main() {
	var env envConfig
	if err := envconfig.Process("", &env); err != nil {
		log.Printf("[ERROR] Failed to process env var: %s", err)
		os.Exit(1)
	}

	queue := 0
	requested := make(map[string]int, 0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: env.Scheme, Host: env.Host, Path: env.Path}

	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	log.Printf("Connected!")

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}

			msg := string(message)

			switch msg {
			case "o":
				req := fmt.Sprintf("o%d", rand.Intn(offerCount))
				if err := c.WriteMessage(websocket.TextMessage, []byte(req)); err != nil {
					log.Println("err:", err)
				}
				log.Printf("ordered: %s", req)

			case "c", "s", "f":
				log.Printf("back with %s", msg)
				queue--
			default:
				log.Printf("recv: %s", message)
			}
		}
	}()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	updateCounts := func() {

		data := make(map[string][]interface{})

		u := fmt.Sprintf("http://%s/airport/data", env.Host)

		resp, err := http.Get(u)
		if err != nil {
			log.Println("Retailers error", err)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)

		if err := json.Unmarshal(body, &data); err != nil {
			log.Println("Unmarshal error", err)
			return
		}

		participantCount = len(data["retailers"])
		log.Println("Retailers now at", participantCount)
	}

	for {
		select {
		case <-done:
			return
		case <-ticker.C:

			if queue > 0 {
				continue
			}

			updateCounts()

			req := fmt.Sprintf("r%d", rand.Intn(participantCount))
			//req := "r2"
			err = c.WriteMessage(websocket.TextMessage, []byte(req)) // todo: we can choose more options.
			if err != nil {
				log.Println("write:", err)
				return
			}
			requested[req]++
			queue++

			log.Println("---", queue)
			for k, v := range requested {
				log.Println(k, "-->", v)
			}

		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
