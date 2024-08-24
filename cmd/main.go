package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/ThailanTec/go-rabbitmq/entities"
	"github.com/ThailanTec/go-rabbitmq/queue"
	"github.com/streadway/amqp"
	"github.com/subosito/gotenv"
)

func init() {
	err := gotenv.Load()

	if err != nil {
		panic("Error to loading .env")
	}
}

func main() {
	in := make(chan []byte)
	ch := queue.Connect()

	queue.StartConsuming(in, ch)

	for m := range in {
		var order entities.Order
		err := json.Unmarshal(m, &order)
		if err != nil {
			fmt.Println(err.Error())
		}

		fmt.Println("Novo pedido feito: ", order.UUID)
		Start(order, ch)
	}
}

func Start(order entities.Order, ch *amqp.Channel) {
	go Worker(order, ch)

}

func Worker(order entities.Order, ch *amqp.Channel) {
	f, err := os.Open("destionations/" + order.Destionation + ".txt")
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		data := strings.Split(scanner.Text(), ",")
		json := destinationToJson(order, data[0], data[1])

		time.Sleep(2 * time.Second)

		queue.Notify(string(json), ch)
	}

	json := destinationToJson(order, "0", "0")
	queue.Notify(string(json), ch)
}

func destinationToJson(order entities.Order, lat, long string) []byte {
	dest := entities.Destionation{
		Order: order.UUID.String(),
		Lat:   lat,
		Long:  long,
	}
	json, _ := json.Marshal(dest)

	return json
}
