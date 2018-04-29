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
	timeoutflag := flag.Int("timeout", 4, "timeout in seconds")

	flag.Parse()
	log.SetOutput(os.Stderr)

	// Channels OMIT
	doneUrls := make(map[string]int) // HL1

	steps := make(chan CrawlerStep) // HL2
	errors := make(chan error)      // HL2
	// END OMIT

	// Control  OMIT
	alldone := make(chan bool)
	wg := sync.WaitGroup{}

	timeoutDuration := time.Second * time.Duration(*timeoutflag)
	timeout := time.After(timeoutDuration)
	// END OMIT
	log.Printf("Start at %s, depth of %d, timeout of %s", *initialUrl, *depthflag, timeoutDuration)

	// wgwait OMIT
	wg.Add(1)
	go FindUrlsIn(*initialUrl, *depthflag, steps, errors)

	go func() {
		wg.Wait()
		alldone <- true
	}()
	// END OMIT

	// Part1 OMIT
	for {
		select {
		case step := <-steps: // HL3
			if _, exist := doneUrls[step.childUrl]; step.depth > 0 && !exist { // HL4
				//log.Printf("Depth level %d\n", step.depth-1) OMIT
				doneUrls[step.childUrl] = 1 // HL5
				wg.Add(1)
				go FindUrlsIn(step.childUrl, step.depth-1, steps, errors) // HL6
			}

		case err := <-errors: // HL7
			if err != nil {
				log.Println(err)
			}
			wg.Done()
			// END OMIT
		case <-timeout:
			log.Println("Timeout")
			go func() {
				alldone <- false
			}()
			// END OMIT

		case done := <-alldone:
			log.Println("Finishing")
			data, _ := json.MarshalIndent(doneUrls, "", "   ")
			fmt.Print(string(data))
			log.Printf("Found: %d ,Completed %t\n", len(doneUrls), done)
			return
			// END OMIT
		}
	}

}
