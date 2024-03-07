package main

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/http"
	"os"
	"time"
)

const webPort = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {
	//try to connect to rabbitmq

	conn, err := connect()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	app := Config{
		Rabbit: conn,
	}

	fmt.Printf("Staring broker service on port %s \n", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	var counts int64
	var backoff = time.Second * 1
	var conn *amqp.Connection

	//wait until rabbit is ready
	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
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
