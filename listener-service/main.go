package main

import (
	amqp "github.com/rabbitmq/amqp091-go"
	"listener/event"
	"log"
	"os"
	"time"
)

func main() {
	//try to connect to rabbitmq

	conn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	//start listening for messages
	log.Println("listening and consuming rabbitmq messages")

	//create consumer
	consumer, err := event.NewConsumer(conn)
	if err != nil {
		panic(err)
	}

	//watch the queue and consume events
	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backoff = time.Second * 1
	var conn *amqp.Connection

	//wait until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@localhost")
		if err != nil {
			log.Println("rabbitmq is not ready yet")
			counts++
		} else {
			log.Println("connected to rabbitmq")
			conn = c
			break
		}
		log.Println("backing off for ", backoff, " seconds")
		time.Sleep(backoff)
		backoff *= 2

		if counts > 5 {
			log.Println(err)
			return nil, err
		}
	}
	return conn, nil
}
