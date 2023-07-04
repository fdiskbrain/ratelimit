package main

import (
	"github.com/envoyproxy/ratelimit/src/service_cmd/runner"
	"github.com/envoyproxy/ratelimit/src/settings"
)

func main() {
	cfg := settings.NewSettings()
	runner := runner.NewRunner(cfg)
	runner.Run()
}
