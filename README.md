# Sensu Check Entropy

[![Go Report Card](https://goreportcard.com/badge/github.com/nmollerup/sensu-check-entropy)](https://goreportcard.com/report/github.com/nmollerup/sensu-check-entropy)
[![GoDoc](https://godoc.org/github.com/nmollerup/sensu-check-entropy?status.svg)](https://godoc.org/github.com/nmollerup/sensu-check-entropy)

## Description

This plugin provides native entropy instrumentation for monitoring and metrics collection of available system entropy on Linux systems.

The plugin includes:

- **check-entropy**: Monitor entropy levels with configurable warning and critical thresholds
- **metrics-entropy**: Collect entropy metrics in Graphite format for time-series databases

Both commands read the available entropy from `/proc/sys/kernel/random/entropy_avail`.

## Installation

### Binary Releases

Download the latest release from the [releases page](https://github.com/nmollerup/sensu-check-entropy/releases).

### From Source

```bash
go install github.com/nmollerup/sensu-check-entropy/cmd/check-entropy@latest
go install github.com/nmollerup/sensu-check-entropy/cmd/metrics-entropy@latest
```

## Configuration

### Asset Registration

Assets are the recommended way to install plugins in Sensu. The asset can be configured using sensuctl:

```bash
sensuctl asset add nmollerup/sensu-check-entropy
```

### Check Definition

Example check definition:

```yaml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: check-entropy
spec:
  command: check-entropy --warning 60 --critical 30
  subscriptions:
    - linux
  runtime_assets:
    - nmollerup/sensu-check-entropy
  interval: 60
  publish: true
```

## Usage

### check-entropy

```bash
Sensu check to monitor available system entropy

Usage:
  check-entropy [flags]

Flags:
  -c, --critical int        Critical threshold for available entropy (default 30)
  -h, --help                help for check-entropy
  -m, --metrics             Output entropy as metrics in Graphite format
  -s, --metric-scheme       Metric naming scheme when using --metrics (default: hostname.entropy)
  -w, --warning int         Warning threshold for available entropy (default 60)
```

### Examples

Check with default thresholds (warning: 60, critical: 30):

```bash
check-entropy
```

Check with custom thresholds:

```bash
check-entropy --warning 100 --critical 50
```

Check using environment variables:

```bash
CHECK_ENTROPY_WARNING=100 CHECK_ENTROPY_CRITICAL=50 check-entropy
```

Output metrics in Graphite format:

```bash
check-entropy --metrics
# Output: hostname.entropy 256 1738147200
```

Output metrics with custom scheme:

```bash
check-entropy --metrics --metric-scheme myserver.system.entropy
# Output: myserver.system.entropy 256 1738147200
```

## Output

The check will output one of the following states:

- **OK**: Entropy is above the warning threshold
- **WARNING**: Entropy is at or below the warning threshold but above the critical threshold
- **CRITICAL**: Entropy is at or below the critical threshold
- **UNKNOWN**: Unable to read entropy or invalid thresholds

Example outputs:

```bash
CheckEntropy OK: entropy is 256
CheckEntropy WARNING: entropy is 55
CheckEntropy CRITICAL: entropy is 25
```

### metrics-entropy

```bash
Sensu metric plugin to collect system entropy metrics

Usage:
  metrics-entropy [flags]

Flags:
  -h, --help           help for metrics-entropy
  -s, --scheme string  Metric naming scheme, text to prepend to metric (default: hostname.entropy)
```

#### Examples

Collect entropy metrics with default scheme:

```bash
metrics-entropy
# Output: hostname.entropy 256 1738147200
```

Collect entropy metrics with custom scheme:

```bash
metrics-entropy --scheme production.server01.entropy
# Output: production.server01.entropy 256 1738147200
```

Example metric check definition:

```yaml
---
type: CheckConfig
api_version: core/v2
metadata:
  name: metrics-entropy
spec:
  command: metrics-entropy
  subscriptions:
    - linux
  runtime_assets:
    - nmollerup/sensu-check-entropy
  interval: 60
  publish: true
  output_metric_format: graphite_plaintext
  output_metric_handlers:
    - influxdb
```

## Platform Support

This plugin is designed for Linux systems only, as it reads from `/proc/sys/kernel/random/entropy_avail`.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

MIT License - see LICENSE file for details.
