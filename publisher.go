package publisher

import (
  "log"
  "os"

  "github.com/streadway/amqp"
  "github.com/VioletGrey/error-handler"
)

func Emit(exchange_name string, routing_key string, message []byte) {
  conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
  vgError.FailOnError(err, "Failed to connect to RabbitMQ", "")
  defer conn.Close()

  ch, err := conn.Channel()
  vgError.FailOnError(err, "Failed to open a channel", "")
  defer ch.Close()

  err = ch.ExchangeDeclare(
    exchange_name, // name
    "direct",      // type
    true,          // durable
    false,         // auto-deleted
    false,         // internal
    false,         // no-wait
    nil,           // arguments
  )
  vgError.FailOnError(err, "Failed to declare an exchange", "")

  body := message
  err = ch.Publish(
    exchange_name,         // exchange
    routing_key, // routing key
    false, // mandatory
    false, // immediate
    amqp.Publishing{
      ContentType: "text/plain",
      Body:        []byte(body),
    })
  vgError.FailOnError(err, "Failed to publish a message", string(body))

  log.Printf(" [x] Sent %s to %s", body, routing_key)
}