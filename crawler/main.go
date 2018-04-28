package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	initialUrl := flag.String("url", "http://www.google.es", "Url to start the crawler")
	depthflag := flag.Int("depth", 2, "Max depth to iterate")
	timeoutflag := flag.Int("timeout", 10, "timeout in seconds")

	flag.Parse()
	log.SetOutput(os.Stderr)
	log.Println("Start")

	doneUrls := make(map[string]int)
	wg := sync.WaitGroup{}

	steps := make(chan CrawlerStep)
	errors := make(chan error)

	timeoutDuration := time.Second * time.Duration(*timeoutflag)
	log.Println(timeoutDuration)
	timeout := time.After(timeoutDuration)
	alldone := make(chan bool)

	wg.Add(1)
	go FindUrlsIn(*initialUrl, *depthflag, steps, errors)

	go func() {
		wg.Wait()
		alldone <- true
	}()

	for {
		select {
		case step := <-steps:
			if _, exist := doneUrls[step.childUrl]; step.depth > 0 && !exist {
				//log.Printf("Depth level %d\n", step.depth-1)
				doneUrls[step.childUrl] = 1
				wg.Add(1)
				go FindUrlsIn(step.childUrl, step.depth-1, steps, errors)
			}

		case err := <-errors:
			if err != nil {
				log.Println(err)
			}
			wg.Done()
		case <-timeout:
			log.Println("Timeout")
			go func() {
				alldone <- false
			}()

		case done := <-alldone:
			log.Println("Finishing")
			data, _ := json.MarshalIndent(doneUrls, "", "   ")
			fmt.Print(string(data))
			log.Printf("Found: %d ,Completed %t\n", len(doneUrls), done)
			return
		}
	}

}
