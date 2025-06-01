package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/exec"

	"github.com/BurntSushi/toml"

	ga "saml.dev/gome-assistant"
)

const (
	defaultServerUrl = "http://homeassistant.local:8123"
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Fprintf(os.Stderr, "config file is the one expected argument\n")
		os.Exit(1)
	}

	cfg := loadCfg()

	entityListners := []ga.EntityListener{}
	for k, v := range cfg.EntityRunners {
		l := ga.NewEntityListener().EntityIds(k).Call(
			buildEnityRunner(
				v,
			),
		).Build()
		entityListners = append(entityListners, l)
	}

	app, err := ga.NewApp(ga.NewAppRequest{
		URL:              cfg.ServerURL,
		HAAuthToken:      os.Getenv("HA_AUTH_TOKEN"),
		HomeZoneEntityId: "zone.home",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not connect: %v\n", err)
		os.Exit(1)
	}

	app.RegisterEntityListeners(entityListners...)

	app.Start()
}

func loadCfg() config {
	var cfg config
	_, err := toml.DecodeFile(flag.Arg(0), &cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "config failed:\n%v\n", err)
		os.Exit(1)
	}

	if cfg.ServerURL == "" {
		slog.Warn("defaulting server URL", "default", defaultServerUrl)
		cfg.ServerURL = defaultServerUrl
	}

	return cfg
}

func buildEnityRunner(c entityRunnerConfig) ga.EntityListenerCallback {
	return func(service *ga.Service, state ga.State, sensor ga.EntityData) {
		err := exec.Command(c.Command, c.Args...).Run()
		if err != nil {
			slog.Error("running command:", "error", err)
		}
	}
}

type config struct {
	ServerURL     string                        `toml:"server_url"`
	EntityRunners map[string]entityRunnerConfig `toml:"entity_runners"`
}

type entityRunnerConfig struct {
	Command string
	Args    []string
}
