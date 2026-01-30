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

// Config represents the check plugin config.
type Config struct {
	sensu.PluginConfig
	Warning      int
	Critical     int
	Metrics      bool
	MetricScheme string
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "check-entropy",
			Short:    "Sensu check to monitor available system entropy",
			Keyspace: "sensu.io/plugins/check-entropy/config",
		},
	}

	options = []sensu.ConfigOption{
		&sensu.PluginConfigOption[int]{
			Path:      "warning",
			Env:       "CHECK_ENTROPY_WARNING",
			Argument:  "warning",
			Shorthand: "w",
			Default:   60,
			Usage:     "Warning threshold for available entropy",
			Value:     &plugin.Warning,
		},
		&sensu.PluginConfigOption[int]{
			Path:      "critical",
			Env:       "CHECK_ENTROPY_CRITICAL",
			Argument:  "critical",
			Shorthand: "c",
			Default:   30,
			Usage:     "Critical threshold for available entropy",
			Value:     &plugin.Critical,
		},
		&sensu.PluginConfigOption[bool]{
			Path:      "metrics",
			Env:       "CHECK_ENTROPY_METRICS",
			Argument:  "metrics",
			Shorthand: "m",
			Default:   false,
			Usage:     "Output entropy as metrics in Graphite format",
			Value:     &plugin.Metrics,
		},
		&sensu.PluginConfigOption[string]{
			Path:      "metric-scheme",
			Env:       "CHECK_ENTROPY_METRIC_SCHEME",
			Argument:  "metric-scheme",
			Shorthand: "s",
			Default:   getDefaultScheme(),
			Usage:     "Metric naming scheme when using --metrics (default: hostname.entropy)",
			Value:     &plugin.MetricScheme,
		},
	}
)

func main() {
	check := sensu.NewCheck(&plugin.PluginConfig, options, checkArgs, executeCheck, false)
	check.Execute()
}

func checkArgs(event *types.Event) (int, error) {
	if plugin.Critical < 0 || plugin.Warning < 0 {
		return sensu.CheckStateUnknown, fmt.Errorf("invalid entropy threshold")
	}
	if plugin.Critical > plugin.Warning {
		return sensu.CheckStateUnknown, fmt.Errorf("critical threshold must be less than or equal to warning threshold")
	}
	return sensu.CheckStateOK, nil
}

func executeCheck(event *types.Event) (int, error) {
	entropy, err := readEntropy()
	if err != nil {
		return sensu.CheckStateUnknown, fmt.Errorf("failed to read entropy: %v", err)
	}

	// If metrics mode is enabled, output Graphite format and exit OK
	if plugin.Metrics {
		timestamp := time.Now().Unix()
		fmt.Printf("%s %d %d\n", plugin.MetricScheme, entropy, timestamp)
		return sensu.CheckStateOK, nil
	}

	if entropy <= plugin.Critical {
		fmt.Printf("CheckEntropy CRITICAL: entropy is %d\n", entropy)
		return sensu.CheckStateCritical, nil
	}

	if entropy <= plugin.Warning {
		fmt.Printf("CheckEntropy WARNING: entropy is %d\n", entropy)
		return sensu.CheckStateWarning, nil
	}

	fmt.Printf("CheckEntropy OK: entropy is %d\n", entropy)
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
