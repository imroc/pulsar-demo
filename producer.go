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
	registerStringParameter(producerCmd.Flags(), TOPIC, "", "topic of pulsar producer")
	registerIntParameterP(producerCmd.Flags(), MESSAGES, "n", 0, "number of messages to send")
	registerStringParameter(producerCmd.Flags(), PRODUCE_DURATION, "1s", "duration of each message to produce")
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
	interval := viper.GetDuration(PRODUCE_DURATION)
	n := viper.GetInt(MESSAGES)
	producer := GetProuducer()

	payload := []byte(`{"msg": "test"}`)

	log.Println("msg send start")
	for i := 1; ; i++ {
		log.Println("send msg", i)
		_, err := producer.Send(context.Background(), &pulsar.ProducerMessage{
			Payload: payload,
		})
		if err != nil {
			log.Println("Send error:", err)
			continue
		}
		time.Sleep(interval)
		if n <= 0 {
			continue
		} else {
			if i <= n {
				continue
			} else {
				break
			}
		}
	}
	log.Println("msg send end")
	return nil
}
