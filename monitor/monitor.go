package monitor

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"log"
)

func Watch(loc string) *fsnotify.Watcher {
	if loc == "" {
		loc = "./tmp"
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Panic(err)
	}

	//done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				fmt.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("modified file:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			}
		}
	}()
	err = watcher.Add(loc)
	if err != nil {
		log.Panic(err)
	}
	//<-done
	return watcher
}
