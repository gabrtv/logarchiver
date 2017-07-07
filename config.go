package main

import (
	"github.com/kelseyhightower/envconfig"
)

const (
	appName = "logarchiver"
)

type config struct {
	QueueURL string `envconfig:"QUEUE_URL"`
}

func parseConfig() (*config, error) {
	ret := new(config)
	if err := envconfig.Process(appName, ret); err != nil {
		return nil, err
	}
	return ret, nil
}
