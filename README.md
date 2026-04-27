# dsmr

An idiomatic Go implementation of a DSMR 5.0.2 P1 telegram parser and Prometheus metrics exporter.

Connects to the P1 serial port on a Dutch smart meter, parses each telegram, and exposes all measurements as Prometheus gauges.

## Requirements

- Go 1.26+
- A P1 USB cable connected to your smart meter's P1 port

## Development

```sh
make build            # compile
make test             # unit tests
make test-integration # integration tests (reads examples/ fixtures)
make test-coverage    # unit tests with coverage report
make lint             # golangci-lint
make clean            # remove coverage artifacts
```

CI runs lint, unit tests, and integration tests automatically on every push and pull request to `main`. Dependabot keeps Go modules and GitHub Actions up to date.

## Build

```sh
go build -o dsmr .
```

## Run

```sh
./dsmr [flags]
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `-port` | auto-detect | Serial port path (e.g. `/dev/ttyUSB0`) |
| `-baud` | `115200` | Serial port baud rate |
| `-listen` | `:8080` | Prometheus metrics listen address |
| `-debug` | `false` | Enable debug logging |

### Port auto-detection

When `-port` is not set, the process scans available serial ports for names matching known P1 cable patterns:

| Pattern | Adapter type |
|---------|--------------|
| `/dev/ttyUSB*` | Linux — FTDI / CP210x |
| `/dev/ttyACM*` | Linux — CDC ACM |
| `/dev/cu.usbserial*` | macOS — FTDI |
| `/dev/cu.usbmodem*` | macOS — CDC ACM |

If no pattern matches, the first available port is used.

### Graceful shutdown

The process handles `SIGINT` and `SIGTERM`. On shutdown it:
1. Closes the serial port (stops the parser)
2. Waits up to 5 s for in-flight telegrams to drain
3. Shuts down the HTTP server with a 5 s timeout

## Prometheus metrics

Metrics are served at `http://<listen>/metrics` (redirected from `/`).

### Info metrics (GaugeVec, value always 1)

| Metric | Label | Source OBIS | Description |
|--------|-------|-------------|-------------|
| `dsmr_info` | `version` | `1-3:0.2.8` | DSMR protocol version |
| `electricity_equipment_info` | `identifier` | `0-0:96.1.1` | Meter equipment identifier |

### M-Bus metrics (GaugeVec)

| Metric | Labels | Source OBIS | Description |
|--------|--------|-------------|-------------|
| `mbus_last_value` | `channel`, `device_type` | `0-n:24.2.1` | Last 5-minute metered value from connected M-Bus device |
| `mbus_valve_state` | `channel` | `0-n:24.4.0` | Valve/switch position (0=disconnected, 1=connected, 2=ready) |
| `mbus_equipment_info` | `channel`, `identifier` | `0-n:96.1.0` | M-Bus equipment identifier (value always 1) |

### Electricity gauges

| Metric | Unit | OBIS | Description |
|--------|------|------|-------------|
| `electricity_delivered_to_client_tariff1_kwh` | kWh | `1-0:1.8.1` | Cumulative energy import tariff 1 |
| `electricity_delivered_to_client_tariff2_kwh` | kWh | `1-0:1.8.2` | Cumulative energy import tariff 2 |
| `electricity_delivered_by_client_tariff1_kwh` | kWh | `1-0:2.8.1` | Cumulative energy export tariff 1 |
| `electricity_delivered_by_client_tariff2_kwh` | kWh | `1-0:2.8.2` | Cumulative energy export tariff 2 |
| `electricity_tariff` | — | `0-0:96.14.0` | Active tariff indicator (1 or 2) |
| `actual_electricity_power_delivered_kw` | kW | `1-0:1.7.0` | Actual power import |
| `actual_electricity_power_received_kw` | kW | `1-0:2.7.0` | Actual power export |
| `power_failures_in_any_phase` | count | `0-0:96.7.21` | Total power failure count |
| `long_power_failures_in_any_phase` | count | `0-0:96.7.9` | Total long power failure count |
| `voltage_sags_phase_l1` | count | `1-0:32.32.0` | Voltage sags L1 |
| `voltage_sags_phase_l2` | count | `1-0:52.32.0` | Voltage sags L2 ¹ |
| `voltage_sags_phase_l3` | count | `1-0:72.32.0` | Voltage sags L3 ¹ |
| `voltage_swells_phase_l1` | count | `1-0:32.36.0` | Voltage swells L1 |
| `voltage_swells_phase_l2` | count | `1-0:52.36.0` | Voltage swells L2 ¹ |
| `voltage_swells_phase_l3` | count | `1-0:72.36.0` | Voltage swells L3 ¹ |
| `instant_voltage_l1` | V | `1-0:32.7.0` | Instantaneous voltage L1 |
| `instant_voltage_l2` | V | `1-0:52.7.0` | Instantaneous voltage L2 ¹ |
| `instant_voltage_l3` | V | `1-0:72.7.0` | Instantaneous voltage L3 ¹ |
| `instant_current_l1` | A | `1-0:31.7.0` | Instantaneous current L1 |
| `instant_current_l2` | A | `1-0:51.7.0` | Instantaneous current L2 ¹ |
| `instant_current_l3` | A | `1-0:71.7.0` | Instantaneous current L3 ¹ |
| `instant_active_power_delivered_l1` | kW | `1-0:21.7.0` | Active power import L1 |
| `instant_active_power_delivered_l2` | kW | `1-0:41.7.0` | Active power import L2 ¹ |
| `instant_active_power_delivered_l3` | kW | `1-0:61.7.0` | Active power import L3 ¹ |
| `instant_active_power_received_l1` | kW | `1-0:22.7.0` | Active power export L1 |
| `instant_active_power_received_l2` | kW | `1-0:42.7.0` | Active power export L2 ¹ |
| `instant_active_power_received_l3` | kW | `1-0:62.7.0` | Active power export L3 ¹ |

¹ Poly-phase meters only.

### Example PromQL queries

```promql
# Net energy consumption (import minus export, both tariffs)
(
  electricity_delivered_to_client_tariff1_kwh
  + electricity_delivered_to_client_tariff2_kwh
) - (
  electricity_delivered_by_client_tariff1_kwh
  + electricity_delivered_by_client_tariff2_kwh
)

# Current power balance (positive = importing, negative = exporting)
actual_electricity_power_delivered_kw - actual_electricity_power_received_kw

# Total three-phase import power
instant_active_power_delivered_l1
  + instant_active_power_delivered_l2
  + instant_active_power_delivered_l3

# Gas usage from M-Bus channel 1
mbus_last_value{channel="1"}
```

## Supported OBIS codes

| OBIS Code | Description |
|-----------|-------------|
| `1-3:0.2.8` | DSMR version |
| `0-0:1.0.0` | P1 message timestamp |
| `0-0:96.1.1` | Electricity meter equipment identifier |
| `1-0:1.8.1` | Meter reading electricity delivered to client — tariff 1 |
| `1-0:1.8.2` | Meter reading electricity delivered to client — tariff 2 |
| `1-0:2.8.1` | Meter reading electricity delivered by client — tariff 1 |
| `1-0:2.8.2` | Meter reading electricity delivered by client — tariff 2 |
| `0-0:96.14.0` | Tariff indicator electricity |
| `1-0:1.7.0` | Actual electricity power delivered (+P) |
| `1-0:2.7.0` | Actual electricity power received (−P) |
| `0-0:96.7.21` | Number of power failures in any phase |
| `0-0:96.7.9` | Number of long power failures in any phase |
| `1-0:99.97.0` | Power failure event log |
| `1-0:32.32.0` | Number of voltage sags in phase L1 |
| `1-0:52.32.0` | Number of voltage sags in phase L2 |
| `1-0:72.32.0` | Number of voltage sags in phase L3 |
| `1-0:32.36.0` | Number of voltage swells in phase L1 |
| `1-0:52.36.0` | Number of voltage swells in phase L2 |
| `1-0:72.36.0` | Number of voltage swells in phase L3 |
| `0-0:96.13.0` | Text message (up to 2048 characters, hex-encoded) |
| `1-0:32.7.0` | Instantaneous voltage L1 |
| `1-0:52.7.0` | Instantaneous voltage L2 |
| `1-0:72.7.0` | Instantaneous voltage L3 |
| `1-0:31.7.0` | Instantaneous current L1 |
| `1-0:51.7.0` | Instantaneous current L2 |
| `1-0:71.7.0` | Instantaneous current L3 |
| `1-0:21.7.0` | Instantaneous active power L1 (+P) |
| `1-0:41.7.0` | Instantaneous active power L2 (+P) |
| `1-0:61.7.0` | Instantaneous active power L3 (+P) |
| `1-0:22.7.0` | Instantaneous active power L1 (−P) |
| `1-0:42.7.0` | Instantaneous active power L2 (−P) |
| `1-0:62.7.0` | Instantaneous active power L3 (−P) |
| `0-1:24.1.0` | M-Bus device type — channel 1 |
| `0-1:96.1.0` | M-Bus equipment identifier — channel 1 |
| `0-1:24.2.1` | M-Bus last 5-minute value — channel 1 |
| `0-1:24.4.0` | M-Bus valve/switch position — channel 1 |
| `0-2:24.1.0` | M-Bus device type — channel 2 |
| `0-2:96.1.0` | M-Bus equipment identifier — channel 2 |
| `0-2:24.2.1` | M-Bus last 5-minute value — channel 2 |
| `0-2:24.4.0` | M-Bus valve/switch position — channel 2 |
| `0-3:24.1.0` | M-Bus device type — channel 3 |
| `0-3:96.1.0` | M-Bus equipment identifier — channel 3 |
| `0-3:24.2.1` | M-Bus last 5-minute value — channel 3 |
| `0-3:24.4.0` | M-Bus valve/switch position — channel 3 |
| `0-4:24.1.0` | M-Bus device type — channel 4 |
| `0-4:96.1.0` | M-Bus equipment identifier — channel 4 |
| `0-4:24.2.1` | M-Bus last 5-minute value — channel 4 |
| `0-4:24.4.0` | M-Bus valve/switch position — channel 4 |

## Package architecture

```
pkg/
├── cosem/     # COSEM type system (IEC 62056): Enum, FloatingPoint, Integer,
│              # OctetString, String, Timestamp, plus SI unit constants and
│              # data-tag / class / attribute identifiers.
├── obis/      # OBIS reference registry and Prometheus metrics registration.
└── telegram/  # P1 telegram parsing pipeline: line tokenizer, header / data /
               # footer parsers, CRC-16/IBM validator, and the streaming
               # Parser that drives them and feeds Prometheus.
```
