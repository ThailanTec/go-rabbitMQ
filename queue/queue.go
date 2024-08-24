package queue

import (
	"fmt"
	"log"
	"os"

	"github.com/streadway/amqp"
)

func Connect() *amqp.Channel {
	amqpServer := os.Getenv("RABBITMQ_DEFAULT_HOST")
	amqpPort := os.Getenv("RABBITMQ_DEFAULT_PORT")
	amqpUser := os.Getenv("RABBITMQ_DEFAULT_USER")
	amqpPassword := os.Getenv("RABBITMQ_DEFAULT_PASSWORD")

	// Cria a URL de conexão
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%s/", amqpUser, amqpPassword, amqpServer, amqpPort)

	conn, err := amqp.Dial(dsn)
	failOnError(err, "Erro ao conectar com o RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Falha ao abrir o canal")

	return ch
}

func StartConsuming(in chan []byte, ch *amqp.Channel) {
	q, err := ch.QueueDeclare(
		os.Getenv("RABBITMQ_CONSUMER_QUEUE"),
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Filed to declare queue")

	msg, err := ch.Consume(
		q.Name,
		"go-worker-simultor",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register consumer")

	go func() {
		for m := range msg {
			in <- []byte(m.Body)
		}
		close(in)
	}()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}

}

func Notify(payload string, ch *amqp.Channel) {

	err := ch.Publish(
		os.Getenv("RABBITMQ_CONSUMER_DESTINATION"),
		os.Getenv("RABBITMQ_CONSUMER_ROUTING_KEY"),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(payload),
		},
	)
	failOnError(err, "Erro to publish message")
	fmt.Print("Messagem enviada: ", payload)
}
