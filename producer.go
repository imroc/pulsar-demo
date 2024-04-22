package main

import (
	"context"
	"errors"
	"log"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	registerStringParameter(producerCmd.Flags(), TOPIC, "", "topic of pulsar producer")
	registerIntParameterP(producerCmd.Flags(), MESSAGES, "n", 0, "number of messages to send")
}

var producerCmd = &cobra.Command{
	Use:   "producer",
	Short: "produce pulsar messages as a producer",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runProducer()
	},
}

func GetProuducer() pulsar.Producer {
	topic := viper.GetString(TOPIC)
	if topic == "" {
		panic("topic is required")
	}
	client := GetPulsarClient()
	producer, err := client.CreateProducer(pulsar.ProducerOptions{
		Topic: topic,
	})
	if err != nil {
		panic(err)
	}
	return producer
}

func runProducer() error {
	n := viper.GetInt(MESSAGES)
	if n <= 0 {
		return errors.New("no messsages need to send")
	}
	producer := GetProuducer()
	payload := []byte(`{"msg": "test"}`)
	log.Println("msg send start")
	for i := 0; i < n; i++ {
		log.Println("send msg", i)
		producer.Send(context.Background(), &pulsar.ProducerMessage{
			Payload: payload,
		})
	}
	log.Println("msg send end")
	return nil
}
