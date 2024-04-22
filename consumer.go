package main

import (
	"context"
	"log"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	registerStringParameter(consumerCmd.Flags(), TOPIC, "", "topic of pulsar consumer")
	registerStringParameter(consumerCmd.Flags(), SUBSCRIPTION, "", "subscription of pulsar topic")
	registerIntParameterP(consumerCmd.Flags(), MESSAGES, "n", 0, "number os messages to consume")
	registerStringParameter(consumerCmd.Flags(), CONSUME_DURATION, "1s", "duration of each message to consume")
}

var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "consume pulsar messages as a consumer",
	Run: func(cmd *cobra.Command, args []string) {
		runConsumer()
	},
}

func GetConsumer() pulsar.Consumer {
	topic := viper.GetString(TOPIC)
	if topic == "" {
		panic("topic is required")
	}
	subscription := viper.GetString(SUBSCRIPTION)
	if subscription == "" {
		panic("subscription is required")
	}
	client := GetPulsarClient()
	consumer, err := client.Subscribe(pulsar.ConsumerOptions{
		Topic:            topic,
		SubscriptionName: subscription,
		Type:             pulsar.Shared,
	})
	if err != nil {
		panic(err)
	}
	return consumer
}

func runConsumer() {
	interval := viper.GetDuration(CONSUME_DURATION)
	n := viper.GetInt(MESSAGES)
	consumer := GetConsumer()
	handleRecevie := func() error {
		msg, err := consumer.Receive(context.Background())
		if err != nil {
			return err
		}
		log.Printf("received msg (id:%s): %s\n", msg.ID(), string(msg.Payload()))
		time.Sleep(interval)
		err = consumer.AckID(msg.ID())
		if err != nil {
			return err
		}
		return nil
	}
	if n <= 0 {
		for {
			err := handleRecevie()
			if err != nil {
				log.Println("ERROR:", err)
				continue
			}
		}
	} else {
		for i := 0; i < n; i++ {
			err := handleRecevie()
			if err != nil {
				log.Println("ERROR:", err)
				continue
			}
		}
	}
}
