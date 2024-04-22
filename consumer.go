package main

import (
	"context"
	"log"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	registerStringParameter(consumerCmd.Flags(), TOPIC, "", "topic of pulsar consumer")
	registerStringParameter(consumerCmd.Flags(), SUBSCRIPTION, "", "subscription of pulsar topic")
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
	consumer := GetConsumer()
	for {
		msg, err := consumer.Receive(context.Background())
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}
		log.Printf("received msg (id:%s): %s\n", msg.ID(), string(msg.Payload()))
		err = consumer.AckID(msg.ID())
		if err != nil {
			log.Println("ERROR:", err)
			continue
		}
	}
}
