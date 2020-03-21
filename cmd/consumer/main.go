package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/go-redis/redis"
	"github.com/streadway/amqp"
)

var (
	client *redis.Client
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func newBD(bdName string, body string) {
	// TODO: add a row to redis db
	err := client.Set("bd-"+bdName, "success", 0).Err()
	if err != nil {
		panic(err)
	}

	fmt.Printf("bd created: %s, from the settings: %s\n", bdName, body)
}

func getBD(bdName string) {
	val, err := client.Get("bd-" + bdName).Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("bd-"+bdName, val)
}

func main() {
	client = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"BD",  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			req, ok := d.Headers["req"]
			if !ok {
				failOnError(errors.New("no request found"), "Failed to find request header")
			}

			bd, ok := d.Headers["entry"]
			if !ok {
				failOnError(errors.New("no entry found"), "Failed to find entry header")
			}

			switch req {
			case "post":
				newBD(bd.(string), string(d.Body))
			case "get":
				getBD(bd.(string))
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
