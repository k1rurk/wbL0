package natJetStream

import (
	"encoding/json"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
	"log"
	"os"
	"os/signal"
	"time"
	"wb_l0/cache"
	"wb_l0/database"
	//"wb_l0/helperJson"
)

const (
	streamName    = "ORDERS"
	consumerName  = "ConOrder"
	streamSubject = "ORDERS"

	subjectName = "ORDERS.sent"
)

func Sub(cache *cache.Cache, db *gorm.DB) {
	// Connect to a server
	// opt := nats.Options{}
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	// JET STREAM
	natsJs, err := nc.JetStream(
		nats.PublishAsyncMaxPending(256),
	)
	if err != nil {
		log.Fatal(err)
	}
	if s, err := natsJs.StreamInfo(streamName); err != nil || s == nil {
		log.Fatal("Stream ", streamName, " not found: ", err)
	}

	if c, err := natsJs.ConsumerInfo(streamName, consumerName); err != nil || c == nil {
		log.Println("Consumer ", streamName, " not found: ", err)
		log.Println("Creating a consumer> ", consumerName)

		ci, err := natsJs.AddConsumer(streamName, &nats.ConsumerConfig{
			Durable:        consumerName,
			DeliverPolicy:  nats.DeliverAllPolicy,
			AckPolicy:      nats.AckExplicitPolicy,
			AckWait:        5 * time.Second,
			MaxAckPending:  1000,
			MaxDeliver:     2,
			ReplayPolicy:   nats.ReplayInstantPolicy,
			DeliverSubject: streamSubject + ".processed",
			FilterSubject:  streamSubject + ".>",
		})
		log.Println("ConsumerInfo:", ci, "error", err)
	}
	log.Println("attempting to receive orders")
	//done := make(chan bool, 1)
	asyncQSub, err := natsJs.QueueSubscribe(subjectName, consumerName, func(msg *nats.Msg) {
		var order []byte
		order = msg.Data
		var o database.Order
		err = json.Unmarshal(order, &o)
		if err != nil {
			log.Printf("error while unmarshaling bytes from jet stream: %v\n", err)
			err := msg.Ack()
			noerr(err)
			return
		}

		if err := cache.Check(o.OrderUid); err != nil {
			noerr(err)
			err := msg.Ack()
			noerr(err)
			return
		}

		err := database.Create(&o, db)
		if err != nil {
			log.Printf("error while create db row %v\n", err)
			err := msg.Ack()
			noerr(err)
			return
		}

		cache.Set(o.OrderUid, &o)

		err = msg.Ack()

		if err != nil {
			log.Printf("ACK message error: %v\n", msg.Header)
			return
		}

		log.Println("Message is accepted and acknowledged")
	},
		nats.Durable(consumerName),
		nats.DeliverAll(),
		nats.ManualAck(),
		nats.AckExplicit(),
		nats.AckWait(5*time.Second),
		nats.MaxAckPending(256),
		nats.MaxDeliver(2),
	)
	if err != nil {
		log.Fatal("async_queue_sub err: " + err.Error())
	}
	log.Printf("Listening on [%s], queue group [%s]\n", subjectName, consumerName)
	if !asyncQSub.IsValid() {
		log.Fatal("async_queue_sub closed")
	}

	// wait forever
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	go func() {
		<-ch
		log.Println("Draining...")
		if err := nc.Drain(); err != nil {
			log.Println("Drain err: ", err.Error())
		}
	}()

}

func noerr(err error) {
	if err != nil {
		log.Println(err)
	}
}
