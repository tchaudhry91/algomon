package main

import (
	"context"
	"embed"
	"_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/charmbracelet/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tchaudhry91/algomon/actions"
	"github.com/tchaudhry91/algomon/algochecks"
	"github.com/tchaudhry91/algomon/measure"
	"github.com/tchaudhry91/algomon/store"
)

var header = `
	
 ____  _     _____ ____  ____  ____  ____  _            ____  _____ _____ _      _____ 
/  _ \/ \   /  __//  _ \/  __\\/  __\\/  _ \/ \__/|      /  _ \/  __//  __// \  /|/__ __\
| / \|| |   | |  _| / \||  \/||  \/|| / \|| |\/||_____ | / \|| |  _|  \  | |\ ||  / \  
| |-||| |_/\| |_//| \_/||  __/|    /| \_/|| |  |\|\____\| |-||| |_//|  /_ | | \||  | |  
\_/ \|\____/\____\\____/\_/   \_/\_\\____/\_/  \|      \_/ \|\____\\____\\_/  \|  \_/  
                                                                                       
                                                                                       

`

//go:embed static/*
var uiFiles embed.FS

func main() {
	logger := log.Default()
	logger.SetPrefix("algomon")
	logger.SetReportCaller(true)
	var configF = flag.String("c", "algomon.json", "config file to use")
	var debugMode = flag.Bool("d", false, "debug logs on")
	flag.Parse()
	if *debugMode {
		logger.SetLevel(log.DebugLevel)
	}

	config := Config{}
	confData, err := os.ReadFile(*configF)
	if err != nil {
		logger.Fatal("Error Reading Config File", "err", err)
	}
	if err = json.Unmarshal(confData, &config); err != nil {
		logger.Fatal("Unable to Unmarshal Config", "err", err)
	}
	fmt.Print(header)
	run(&config, logger)
}

var contexts = map[string]*context.CancelFunc{}

func validateConfig(conf *Config) error {
	datasources := make(map[string]struct{})
	for _, d := range conf.Datasources {
		if d.Name == "" {
			return fmt.Errorf("datasource name cannot be empty")
		}
		datasources[d.Name] = struct{}{}
	}

	algorithmers := make(map[string]struct{})
	for _, a := range conf.Algorithmers {
		if a.Type == "" {
			return fmt.Errorf("algorithmer type cannot be empty")
		}
		algorithmers[a.Type] = struct{}{}
	}

	actioners := make(map[string]struct{})
	for _, a := range conf.Actioners {
		if a.Type == "" {
			return fmt.Errorf("actioner type cannot be empty")
		}
		actioners[a.Type] = struct{}{}
	}

	for _, c := range conf.Checks {
		if c.Name == "" {
			return fmt.Errorf("check name cannot be empty")
		}
		if _, ok := algorithmers[c.AlgorithmerType]; !ok {
			return fmt.Errorf("check %q uses undefined algorithmer type %q", c.Name, c.AlgorithmerType)
		}
		for _, i := range c.Inputs {
			if _, ok := datasources[i.Datasource]; !ok {
				return fmt.Errorf("check %q uses undefined datasource %q", c.Name, i.Datasource)
			}
		}
		for _, a := range c.Actions {
			if _, ok := actioners[a.Actioner]; !ok {
				return fmt.Errorf("check %q uses undefined actioner type %q", c.Name, a.Actioner)
			}
		}
	}
	return nil
}

func run(conf *Config, logger *log.Logger) {
	if err := validateConfig(conf); err != nil {
		logger.Fatal("Invalid configuration", "err", err)
	}

	tickers := make([]*time.Ticker, 0, 1)
	done := make(chan bool)
	shutdown := make(chan error, 1)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	addr := conf.APIListenAddr
	if addr == "" {
		addr = "127.0.0.1:9967"
	}

	s, err := store.NewBoltStore(conf.DatabaseFile, logger)
	if err != nil {
		logger.Fatal("Could not open database", "err", err)
	}

	slogHandler := slog.New(logger.WithPrefix("APIServer"))
	apiServer := NewAPIServer(s, conf, slogHandler)

	apiMux := http.NewServeMux()
	apiMux.Handle("/metrics", promhttp.Handler())
	apiMux.Handle("/api/", apiServer.Mux())
	uiSub, err := fs.Sub(uiFiles, "static")
	if err == nil {
		logger.Info("Registered UI..")
		apiMux.Handle("/", http.FileServer(http.FS(uiSub)))
	}

	server := &http.Server{
		Addr:    addr,
		Handler: apiMux,
	}

	go func() {
		logger.Info("Starting API Server", "addr", addr)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logger.Fatal(err)
		}
	}()

	algorithmers := make(map[string]algochecks.Algorithmer)
	for _, aa := range conf.Algorithmers {
		algorithmers[aa.Type] = algochecks.Build(aa, logger)
	}

	actioners := make(map[string]actions.Actioner)
	for _, aa := range conf.Actioners {
		actioners[aa.Type] = actions.Build(aa, logger)
	}

	for _, c := range conf.Checks {
		ticker := time.NewTicker(c.Interval.Duration)
		tickers = append(tickers, ticker)
		logger.Info("Starting Check", "name", c.Name, "interval", c.Interval.Duration)
		go func(c *algochecks.Check, logger *log.Logger) {
			if c.Immediate {
				err := runCheck(c, conf, logger, s, algorithmers, actioners)
				if err != nil {
					logger.Error("err", err)
				}
			}
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					err := runCheck(c, conf, logger, s, algorithmers, actioners)
					if err != nil {
						logger.Error("err", err)
					}
				}
			}
		}(&c, logger.WithPrefix(c.Name))
	}

	select {
	case signalKill := <-interrupt:
		logger.Info("Received Interrupt", "signal", signalKill)
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer shutdownCancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("HTTP server shutdown error", "err", err)
		}
		for name, cancel := range contexts {
			logger.Info("Cancelling", "name", name)
			(*cancel)()
		}
		for _, ticker := range tickers {
			ticker.Stop()
		}
	case err := <-shutdown:
		logger.Error("err", err)
	}

}

func runCheck(c *algochecks.Check, conf *Config, logger *log.Logger, s *store.BoltStore, algorithmers map[string]algochecks.Algorithmer, actioners map[string]actions.Actioner) error {
	algorithmer := algorithmers[c.AlgorithmerType]
	if algorithmer == nil {
		return fmt.Errorf("AlgorithmerType:%s not found", c.AlgorithmerType)
	}

	processed := countProcessed.WithLabelValues(c.Name)
	succeeded := countSuccess.WithLabelValues(c.Name)
	failed := countFail.WithLabelValues(c.Name)
	defer processed.Inc()

	tempWorkDir, err := os.MkdirTemp(conf.BaseWorkingDir, c.Name+"-")
	if err != nil {
		failed.Inc()
		return fmt.Errorf("Unable to create Temp Dir: %v", err)
	}
	defer os.RemoveAll(tempWorkDir)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	contexts[c.Name] = &cancel

	// Fetch inputs
	inputs := map[string]measure.Result{}
	for _, i := range c.Inputs {
		d := fetchDatasourceByName(conf, i.Datasource)
		if d == nil {
			failed.Inc()
			return fmt.Errorf("Datasource Not Found: %s", i.Datasource)
		}
		api, err := measure.GetPromAPIClient(d.URL)
		if err != nil {
			failed.Inc()
			return fmt.Errorf("Failed to create Prom API Client: %v", err)
		}
		res, err := i.MeasureProm(ctx, api)
		if err != nil {
			failed.Inc()
			return fmt.Errorf("Failed to measure prometheus query: %v", err)
		}
		inputs[i.Name] = res
	}
	output, err := algorithmer.ApplyAlgorithm(ctx, c.Algorithm, c.AlgorithmParams, inputs, tempWorkDir)
	output.Name = c.Name
	if c.Debug {
		defer logger.Debugf("Output: %s", output.CombinedOut)
	}
	if err != nil || output.RC != 0 {
		failed.Inc()
		logger.Error("Check failed", "name", c.Name, "err", err, "rc", output.RC)

		for _, a := range c.Actions {
			actioner := actioners[a.Actioner]
			if actioner == nil {
				logger.Error("Actioner not found", "type", a.Actioner)
				continue
			}
			logger.Info("Dispatching Action", "action", a.Name)
			out, err := actioner.Action(ctx, a.Action, output.CombinedOut, a.Params, tempWorkDir)
			if c.Debug {
				logger.Debugf("Action Output: %s", out.CombinedOut)
			}
			if err != nil {
				logger.Error("Action Failed with error", "name", a.Name, "err", err)
			}
			// Store Values to Database
			actionKey, err := s.PutAction(ctx, c.Name, &a, &out)
			if err != nil {
				logger.Error("Action Storage Failed with error", "name", a.Name, "err", err)
			}
			output.ActionKeys = append(output.ActionKeys, actionKey)
			outputKey, err := s.PutCheck(ctx, c, &output)
			if err != nil {
				logger.Error("Check Storage Failed", "err", err)

				logger.Error("Exited with failure", "storage_key", outputKey)
			}
			logger.Info("Exited with Error. Output Stored to Key", "storage_key", outputKey)
			return err
		}
	}

	outputKey, err := s.PutCheck(ctx, c, &output)
	if err != nil {
		logger.Error("Check Storage Failed", "err", err)
	}
	logger.Info("Exited successfully. Output Stored to Key", "storage_key", outputKey)
	succeeded.Inc()
	return nil
}
