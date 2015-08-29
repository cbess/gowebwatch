package main

import (
	"fmt"
	"github.com/cbess/gowebwatch/web"
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

			fmt.Println("Watching web...", watcher.URL)
		}()
	}

	if len(urls) == 0 {
		fmt.Println("No URLs! Check config.")
		return
	}

	// get the chan info
	for sItem := range watchChan {
		msg := "FAIL"
		if sItem.IsOk() {
			msg = "OK"
		}

		fmt.Println(sItem.URL, " ", sItem.StatusCode, msg)
	}
}

func loadConfig() {
	viper.SetConfigName("test") // config filename
	viper.AddConfigPath("$HOME/tmp")
	viper.SetConfigType("yaml")

	// defaults
	viper.SetDefault("name", "No Name")

	viper.ReadInConfig()
}
