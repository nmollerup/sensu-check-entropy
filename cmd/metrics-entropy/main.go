package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sensu/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-go/types"
)

// Config represents the metric plugin config.
type Config struct {
	sensu.PluginConfig
	Scheme string
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "metrics-entropy",
			Short:    "Sensu metric plugin to collect system entropy metrics",
			Keyspace: "sensu.io/plugins/metrics-entropy/config",
		},
	}

	options = []sensu.ConfigOption{
		&sensu.PluginConfigOption[string]{
			Path:      "scheme",
			Env:       "METRICS_ENTROPY_SCHEME",
			Argument:  "scheme",
			Shorthand: "s",
			Default:   getDefaultScheme(),
			Usage:     "Metric naming scheme, text to prepend to metric",
			Value:     &plugin.Scheme,
		},
	}
)

func main() {
	metric := sensu.NewCheck(&plugin.PluginConfig, options, checkArgs, executeMetric, false)
	metric.Execute()
}

func checkArgs(event *types.Event) (int, error) {
	return sensu.CheckStateOK, nil
}

func executeMetric(event *types.Event) (int, error) {
	entropy, err := readEntropy()
	if err != nil {
		return sensu.CheckStateUnknown, fmt.Errorf("failed to read entropy: %v", err)
	}

	timestamp := time.Now().Unix()
	fmt.Printf("%s %d %d\n", plugin.Scheme, entropy, timestamp)

	return sensu.CheckStateOK, nil
}

func readEntropy() (int, error) {
	data, err := os.ReadFile("/proc/sys/kernel/random/entropy_avail")
	if err != nil {
		return 0, err
	}

	entropy := strings.TrimSpace(string(data))
	value, err := strconv.Atoi(entropy)
	if err != nil {
		return 0, fmt.Errorf("failed to parse entropy value: %v", err)
	}

	return value, nil
}

func getDefaultScheme() string {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "unknown"
	}
	return hostname + ".entropy"
}
