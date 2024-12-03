package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/humans-group/graphqlgen/plugin/resolver/decorator"
	"github.com/humans-group/graphqlgen/plugin/resolver/timeout"
)

func main() {
	flag.Parse()

	cfg, err := loadConfig()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "load config: %v", err)
		os.Exit(2)
	}

	err = api.Generate(cfg,
		plugins()...)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "generate: %v", err)
		os.Exit(3)
	}
}

func loadConfig() (*config.Config, error) {
	location := *configFile
	if location == "" {
		return config.LoadConfigFromDefaultLocations()
	}

	return config.LoadConfig(location)
}

func plugins() []api.Option {
	oo := make([]api.Option, 0)
	oo = append(oo,
		api.AddPlugin(&decorator.Plugin{}),
	)

	if *timeoutsPluginEnabled {
		outFile := *timeoutsPluginOutput
		oo = append(oo,
			api.AddPlugin(timeout.NewPlugin(outFile)))
	}
	return oo
}

var (
	configFile            = flag.String("config", "", "file with the graphqlgen configuration")
	timeoutsPluginEnabled = flag.Bool("timeouts_plugin", false, "enables generation of resolver timeouts wrapper")
	timeoutsPluginOutput  = flag.String("timeouts_plugin_out", "resolver_timeouts.go", "")
)
