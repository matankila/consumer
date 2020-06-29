package main

import (
	"errors"
	"fmt"
	"github.com/matankilla/consumer/environment"
	r "github.com/matankilla/consumer/redis"
	"github.com/matankilla/consumer/utils"
	"log"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/streadway/amqp"
)

var (
	client  *redis.Client
	conn    *amqp.Connection
	ch      *amqp.Channel
	replies <-chan amqp.Delivery
)

func init() {
	initAmqp()
	initRedis()
}

func initAmqp() {
	var (
		err     error
		connUrl = "amqp://" + environment.RabbitUser +
			":" + environment.RabbitPassword + "@localhost:" + environment.RabbitPort
	)

	conn, err = amqp.Dial(connUrl)
	utils.FailOnError(err, "Failed to connect to RabbitMQ")

	ch, err = conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		"BD",  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	utils.FailOnError(err, "Failed to declare a queue")
	log.Printf("declared Queue (%q %d messages, %d consumers), binding to Exchange (key %q)",
		q.Name, q.Messages, q.Consumers, "go-test-key")

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	utils.FailOnError(err, "Failed to set QoS")

	replies, err = ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	utils.FailOnError(err, "Error consuming the Queue")
}

func initRedis() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:" + environment.RedisPort,
		Password: environment.RedisPassword, // no password set
		DB:       0,                         // use default DB
	})
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func server() {
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":"+environment.ServerPort, nil))
}

func main() {
	go server()
	log.Println("Start consuming the Queue...")
	for d := range replies {
		req, ok := d.Headers["request"]
		if !ok {
			utils.FailOnError(errors.New("no request found"), "Failed to find request header")
		}

		bd, ok := d.Headers["entry"]
		if !ok {
			utils.FailOnError(errors.New("no entry found"), "Failed to find entry header")
		}

		switch req {
		case http.MethodPost:
			r.NewBD(client, bd.(string), string(d.Body))
			d.Ack(false)
		default:
			log.Println("none valid request from consumer")
			d.Nack(false, true)
		}
	}
}
