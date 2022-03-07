package main

import (
	"github.com/nats-io/nats.go"
	"io/ioutil"
	"log"
	"os"
	"time"
)

var (
	streamName    = "ORDERS"
	streamSubject = "ORDERS"
	subjectName   = "ORDERS.sent"
)

func main() {
	// Connect to a server
	conn, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	nc, err := nats.NewEncodedConn(conn, nats.JSON_ENCODER)
	if err != nil {
		log.Fatal(err)
	}

	// JET STREAM
	js, err := nc.Conn.JetStream(
		nats.PublishAsyncMaxPending(256),
	)
	if err != nil {
		log.Fatal(err)
	}

	createStream(js)

	publishMsg(js, "model3.json")

}

func createStream(js nats.JetStreamContext) {
	// Check if the ORDERS stream already exists; if not, create it.
	if s, err := js.StreamInfo(streamName); err != nil || s == nil {
		log.Println("Stream ", streamName, " not found: ", err)
		log.Println("Creating a stream> ", streamName)
		sAdded, err := js.AddStream(
			&nats.StreamConfig{
				Name: streamName,
				Subjects: []string{
					streamSubject + ".>",
				},
				MaxMsgs:      5,
				MaxConsumers: 5,
				Discard:      nats.DiscardOld,
				Retention:    nats.WorkQueuePolicy,
				MaxAge:       365 * 24 * time.Hour,
				Duplicates:   1 * time.Hour,
			},
		)
		log.Println("StreamInfo:", sAdded, "error", err)
	}

}

func publishMsg(js nats.JetStreamContext, fileName string) {
	log.Print("publishing an order")
	jsonFile, err := os.Open(fileName)
	defer jsonFile.Close()
	noerr(err)

	byteArray, err := ioutil.ReadAll(jsonFile)
	noerr(err)

	_, err = js.Publish(subjectName, byteArray)
	noerr(err)
}

func noerr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
