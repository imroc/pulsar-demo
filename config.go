package main

import (
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
}

func registerStringParameter(fs *pflag.FlagSet, name, value, usage string) {
	fs.String(name, value, usage)
	viper.BindPFlag(name, fs.Lookup(name))
}

func registerIntParameterP(fs *pflag.FlagSet, name, shorthand string, value int, usage string) {
	fs.IntP(name, shorthand, value, usage)
	viper.BindPFlag(name, fs.Lookup(name))
}
