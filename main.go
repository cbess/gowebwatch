package main

import (
	"fmt"
	"github.com/cbess/webwatch/web"
	"github.com/spf13/viper"
)

func main() {
	loadConfig()

	fmt.Println("Watcher -", viper.GetString("name"))

	watchChan := make(chan web.StatusItem)
	urls := viper.GetStringSlice("urls")

	// watch each url
	for _, url := range urls {
		watcher := web.New(url, true)
		watcher.SetChan(watchChan)
		watcher.Interval = 10 // secs

		go func() {
			watcher.Watch()

			fmt.Println("Watching web...", watcher.Url)
		}()
	}

	// get the chan info
	for sItem := range watchChan {
		msg := "FAIL"
		if sItem.IsOk() {
			msg = "OK"
		}

		fmt.Println(sItem.Url, " ", sItem.StatusCode, msg)
	}
}

func loadConfig() {
	viper.SetConfigName("test") // config filename
	viper.AddConfigPath("$HOME/tmp")
	viper.SetConfigType("yaml")

	// defaults
	viper.SetDefault("name", "great")

	viper.ReadInConfig()
}
