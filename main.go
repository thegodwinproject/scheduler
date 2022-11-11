package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gobuffalo/envy"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
)

func main() {

	runSchedule := envy.Get("TASK_SCHEDULE", "@daily")
	taskURL := envy.Get("TASK_SCHEDULER_URL", "http://localhost:3000/task_run")
	key := envy.Get("TASK_KEY", "1hw553tdyye65ymv77d")

	fmt.Printf("URL: %s\n", taskURL)
	fmt.Printf("Schedule: %s\n", runSchedule)
	fmt.Printf("key: %s\n", key)
	// URL with key
	URL := fmt.Sprintf("%s/%s", taskURL, key)

	// set up cron schedule
	c := cron.New()
	_, err := c.AddFunc(runSchedule, func() {

		fmt.Printf("Running task\n")
		// load page
		res, err := http.Get(URL)
		if err != nil {
			log.WithFields(log.Fields{"service": "task_scheduler"}).Warnf("error reaching %s: %v", URL, err)
		} else {
			if res.Status != "200 OK" {
				log.WithFields(log.Fields{"service": "task_scheduler"}).Warnf("error reaching %s: %v", URL, res.Status)
			}
			fmt.Printf("Task run successfully: received %s for %s\n", res.Status, URL)
			defer res.Body.Close()
		}
	})
	// log cron error
	if err != nil {
		log.WithFields(log.Fields{"service": "task_scheduler"}).Fatalf("error creating cron job: %v", err)
	}
	// create waitgroup
	var wg sync.WaitGroup
	wg.Add(1)
	fmt.Printf("Starting cron scheduler...\n")
	// start cron
	c.Start()
	// wait forever
	wg.Wait()
}
