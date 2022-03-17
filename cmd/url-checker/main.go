package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/zrobisho/url-checker/internal/checkurl"
	"github.com/zrobisho/url-checker/internal/command"
	"github.com/zrobisho/url-checker/internal/metrics"
)

const (
	DefaultPrometheusPort  = 9090
	DefaultCheckInterval   = 1
	DefaultUpstreamTimeout = 5
)

func main() {
	var input command.CLIInput
	flag.StringVar(&input.UpstreamUrls, "upstream-urls", "", "comma separated list of urls to check status")
	flag.IntVar(&input.CheckInterval, "check-interval", DefaultCheckInterval, "interval between checks, in seconds, to query the upstream URL")
	flag.IntVar(&input.UpstreamTimeout, "timeout", DefaultUpstreamTimeout, "timeout to wait for upstream URL to respond, in seconds")
	flag.IntVar(&input.PrometheusPort, "prometheus-port", DefaultPrometheusPort, "Port to start prometheus server on")
	flag.Parse()

	log.SetFormatter(&log.JSONFormatter{})

	commandArgs, err := command.ParseInput(input)
	if err != nil {
		fmt.Println("Usage: ./url-checker")
		fmt.Println("Arguments:")
		flag.PrintDefaults()
		fmt.Println(err)
		os.Exit(1)
	}

	log.Infof("url-checker configuration : %+v", commandArgs)

	client := &http.Client{
		Timeout: commandArgs.UpstreamTimeout,
	}

	urlChecker := checkurl.URLChecker{
		Client: client,
	}

	for _, urlToCheck := range commandArgs.UpstreamUrls {
		log.WithFields(log.Fields{"url": urlToCheck}).Debugf("starting upstream check for url: %s", urlToCheck)
		go func(urlToCheck string) {
			ticker := time.NewTicker(commandArgs.CheckInterval)
			for range ticker.C {
				err := urlChecker.CheckURL(urlToCheck)
				if err != nil {
					metrics.ExternalURLDown(urlToCheck)
					log.WithFields(log.Fields{"url": urlToCheck}).Errorf("url is down %s", err)
				} else {
					metrics.ExternalURLUp(urlToCheck)
					log.WithFields(log.Fields{"url": urlToCheck}).Infof("url is up")
				}
			}
		}(urlToCheck)
	}

	log.Debugf("starting prometheus server on port %d", commandArgs.PrometheusPort)
	metrics.StartMetricsServer(commandArgs.PrometheusPort)
}
