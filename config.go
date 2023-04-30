package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/tasksolo/gosolo"
)

var configPath = flag.String("config", fmt.Sprintf("%s/.solo.conf", os.Getenv("HOME")), "config file path")

func readConfig() (map[string]*gosolo.Config, error) {
	fh, err := os.Open(*configPath)
	if err != nil {
		return map[string]*gosolo.Config{}, nil //nolint:nilerr
	}

	defer fh.Close()

	cfg := map[string]*gosolo.Config{}

	dec := json.NewDecoder(fh)

	err = dec.Decode(&cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func writeConfig(cfg map[string]*gosolo.Config) error {
	fh, err := os.OpenFile(*configPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return err
	}

	defer fh.Close()

	enc := json.NewEncoder(fh)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")

	err = enc.Encode(cfg)
	if err != nil {
		return err
	}

	return nil
}
