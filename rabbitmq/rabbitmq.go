package rabbitmq

import (
	"fmt"
	"time"
	"net/url"
	"github.com/streadway/amqp"
	"github.com/sunmi-OS/gocore/viper"
)

var ch *amqp.Channel

func connRbbitmq() error {

	host := viper.C.GetString("rabbitmq.host")
	port := viper.C.GetString("rabbitmq.port")
	vhost := viper.C.GetString("rabbitmq.vhost")
	user := url.QueryEscape(viper.C.GetString("rabbitmq.user"))
	password := url.QueryEscape(viper.C.GetString("rabbitmq.password"))

	amqpcoinf := amqp.Config{
		Vhost:     vhost,
		Heartbeat: 10 * time.Second,
		Locale:    "en_US",
	}

	conn, err := amqp.DialConfig(fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port), amqpcoinf)
	if err != nil {
		return err
	}

	ch, err = conn.Channel()
	return err
}

func Push(msg string, msgName string) error {

	if ch == nil {
		err := connRbbitmq()
		if err != nil {
			return err
		}
	}

	q, err := ch.QueueDeclare(
		msgName, // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})

	return err
}
