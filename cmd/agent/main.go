package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var header = `
	
 ____  _     _____ ____  ____  ____  ____  _            ____  _____ _____ _      _____ 
/  _ \/ \   /  __//  _ \/  __\/  __\/  _ \/ \__/|      /  _ \/  __//  __// \  /|/__ __\
| / \|| |   | |  _| / \||  \/||  \/|| / \|| |\/||_____ | / \|| |  _|  \  | |\ ||  / \  
| |-||| |_/\| |_//| \_/||  __/|    /| \_/|| |  ||\____\| |-||| |_//|  /_ | | \||  | |  
\_/ \|\____/\____\\____/\_/   \_/\_\\____/\_/  \|      \_/ \|\____\\____\\_/  \|  \_/  
                                                                                       
                                                                                       

`

func main() {
	logger := log.Default()
	var configF = flag.String("c", "algoprom.json", "config file to use")
	flag.Parse()

	config := Config{}
	confData, err := os.ReadFile(*configF)
	if err != nil {
		logger.Fatalf("Error Reading Config File: %v", err)
	}
	if err = json.Unmarshal(confData, &config); err != nil {
		logger.Fatalf("Unable to Unmarshal Config: %v", err)
	}
	fmt.Print(header)
	run(&config, logger)
}

var contexts = map[string]*context.CancelFunc{}

func run(conf *Config, logger *log.Logger) {
	tickers := make([]*time.Ticker, 0, 1)
	done := make(chan bool)
	shutdown := make(chan error, 1)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	addr := conf.MetricsListenAddr
	if addr == "" {
		addr = "127.0.0.1:9967"
	}

	go func(addr string) {
		logger.Printf("Starting Metrics Server on: %s", addr)
		http.Handle("/metrics", promhttp.Handler())
		http.ListenAndServe(addr, nil)
	}(addr)

	for _, m := range conf.Measurements {
		ticker := time.NewTicker(m.Interval.Duration)
		tickers = append(tickers, ticker)
		logger.Printf("Starting Measurement: %s with interval:%s", m.Name, m.Interval.Duration)
		go func() {
			if m.Immediate {
				runMeasure(m, conf, logger)
			}
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					runMeasure(m, conf, logger)
				}
			}
		}()
	}

	select {
	case signalKill := <-interrupt:
		logger.Println("Received Interrupt:", signalKill)
		for name, cancel := range contexts {
			logger.Println("Cancelling:", name)
			(*cancel)()
		}
		for _, ticker := range tickers {
			ticker.Stop()
		}
	case err := <-shutdown:
		logger.Println("Error:", err)
	}

}

func runMeasure(m Measurement, conf *Config, logger *log.Logger) error {
	defer countProcessed.Inc()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	contexts[m.Name] = &cancel
	countSuccess.Inc()
	return nil
}
