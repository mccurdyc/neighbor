package main

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../..")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Printf("error loading config file: %+v\n", err)
		os.Exit(1)
	}
}
