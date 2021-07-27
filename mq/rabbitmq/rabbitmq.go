package rabbitmq

import (
	"fmt"
	"net/url"
	"time"

	"github.com/streadway/amqp"
	"github.com/sunmi-OS/gocore/v2/conf/viper"
)

var ch *amqp.Channel

func connRbbitmq() error {

	host := viper.GetEnvConfig("rabbitmq.host").String()
	port := viper.GetEnvConfig("rabbitmq.port").String()
	vhost := viper.GetEnvConfig("rabbitmq.vhost").String()
	user := url.QueryEscape(viper.GetEnvConfig("rabbitmq.user").String())
	password := url.QueryEscape(viper.GetEnvConfig("rabbitmq.password").String())

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

func UpdateRabbitmq() error {

	r1 := ch
	err := connRbbitmq()

	if err != nil {
		return err
	} else {
		err := r1.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

// 普通模式发布消息
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

// 普通模式消费消息
func Consume(queue string) (<-chan amqp.Delivery, error) {
	if ch == nil {
		err := connRbbitmq()
		if err != nil {
			return nil, err
		}
	}

	q, err := ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	msgs, _ := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto ack
		false,  // exclusive
		false,  // no local
		false,  // no wait
		nil,    // args
	)
	return msgs, nil
}

// 发布订阅模式发布消息
func Publish(exchange, msg string, durable bool) error {
	if ch == nil {
		err := connRbbitmq()
		if err != nil {
			return err
		}
	}

	err := ch.ExchangeDeclare(
		exchange, // name
		"fanout", // type
		durable,  // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return err
	}

	err = ch.Publish(
		exchange, // exchange
		"",       // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		})

	return err
}

// 发布订阅模式消费消息
func SubScribe(exchange string, durable bool) (<-chan amqp.Delivery, error) {
	if ch == nil {
		err := connRbbitmq()
		if err != nil {
			return nil, err
		}
	}

	err := ch.ExchangeDeclare(
		exchange, // name
		"fanout", // type
		durable,  // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"",      // name
		durable, // durable
		false,   // delete when unused
		true,    // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return nil, err
	}

	err = ch.QueueBind(
		q.Name,   // queue name
		"",       // routing key
		exchange, // exchange
		false,
		nil)
	if err != nil {
		return nil, err
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}
