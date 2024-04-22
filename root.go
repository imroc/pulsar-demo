package main

import (
	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	fs := rootCmd.PersistentFlags()
	registerStringParameter(fs, URL, "", "pulsar url")
	registerStringParameter(fs, TOKEN, "", "pulsar authentication jwt token")
	rootCmd.AddCommand(producerCmd)
}

var rootCmd = &cobra.Command{
	Use:   "pulsar-demo",
	Short: "pulsar demo app",
}

func GetPulsarClient() pulsar.Client {
	url := viper.GetString(URL)
	if url == "" {
		panic("url is requred")
	}
	opts := pulsar.ClientOptions{
		URL: url,
	}
	token := viper.GetString(TOKEN)
	if token != "" {
		opts.Authentication = pulsar.NewAuthenticationToken(token)
	}
	client, err := pulsar.NewClient(opts)
	if err != nil {
		panic(err)
	}
	return client
}
