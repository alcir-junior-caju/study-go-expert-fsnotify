package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
)

type DBConfig struct {
	DB       string `json:"db"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
}

var config DBConfig

func main() {
	pwd, _ := os.Getwd()
	file := pwd + "/config.json"
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		panic(err)
	}
	defer watcher.Close()
	MarshalConfig(file)
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fmt.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					MarshalConfig(file)
					fmt.Println("modified file:", event.Name)
					fmt.Println(config)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			}
		}
	}()
	err = watcher.Add(file)
	if err != nil {
		panic(err)
	}
	<-done
}

func MarshalConfig(file string) {
	data, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &config)
	if err != nil {
		panic(err)
	}
}
