package main

import (
	"fmt"
	config "healthcheck/configs"
	"healthcheck/internal/colors"
	"healthcheck/internal/healthchecker"
)

func main() {
	config := config.Load()
	services := config.Services

	results := make(chan healthchecker.PingResult, len(services))
	for _, s := range services {
		go healthchecker.Ping(s.URL, results)
	}

	for range services {
		r := <-results
		color := colors.Red
		if r.Status == healthchecker.StatusUp {
			color = colors.Green
		}
		if r.Err != nil {
			color = colors.Red
			fmt.Printf("[%s] error: %v\n", colors.Wrap(r.URL, color), r.Err)
			continue
		}
		fmt.Printf("[%s] status: %s, last event: %s\n", colors.Wrap(r.URL, color), r.Status, r.Event)
	}

	fmt.Println(colors.Wrap("All services are checked", colors.Yellow))
}
