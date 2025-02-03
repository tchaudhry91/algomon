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
	"github.com/tchaudhry91/algoprom/algochecks"
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

	for _, c := range conf.Checks {
		ticker := time.NewTicker(c.Interval.Duration)
		tickers = append(tickers, ticker)
		logger.Printf("Starting Check: %s with interval:%s", c.Name, c.Interval.Duration)
		go func(c *algochecks.Check) {
			if c.Immediate {
				err := runCheck(c, conf, logger)
				if err != nil {
					logger.Printf("Error in Check %s : %v", c.Name, err)
				}
			}
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					err := runCheck(c, conf, logger)
					logger.Printf("Error in Check %s : %v", c.Name, err)
				}
			}
		}(&c)
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

func runCheck(c *algochecks.Check, conf *Config, logger *log.Logger) error {
	var algorithmer algochecks.Algorithmer
	for _, aa := range conf.Algorithmers {
		if aa.Type == c.AlgorithmerType {
			algorithmer = algochecks.Build(aa.Type, aa.Params, logger)
		}
	}
	if algorithmer == nil {
		return fmt.Errorf("Algorithmer Not Found:%s", c.AlgorithmerType)
	}

	processed := countProcessed.WithLabelValues(c.Name)
	succeeded := countSuccess.WithLabelValues(c.Name)
	failed := countFail.WithLabelValues(c.Name)
	defer processed.Inc()

	tempWorkDir, err := os.MkdirTemp(conf.BaseWorkingDir, c.Name+"-")
	if err != nil {
		logger.Printf("Unable to create Temp Dir for check %s: %v", c.Name, err)
		failed.Inc()
		return err
	}

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	contexts[c.Name] = &cancel

	// Fetch inputs

	_, err = algorithmer.ApplyAlgorithm(ctx, c.Algorithm, c.AlgorithmParams, c.Inputs, tempWorkDir)

	if err != nil {
		failed.Inc()
		logger.Printf("%s check failed: %v", c.Name, err)
		return err

	}
	logger.Printf("%s check exited", c.Name)
	succeeded.Inc()
	return nil
}
