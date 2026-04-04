# DSMR 5.0.2 P1 Implementation Plan

## Current State

The project parses DSMR telegrams from a P1 serial port and exports Prometheus metrics. The `cosem_typing` branch is restructuring packages under `pkg/`. Core tokenizer, parser, and basic COSEM types exist but several spec requirements are incomplete.

---

## Phase 1: Finish Package Restructuring

### 1.1 Move parser into `pkg/`
- Relocate `telegram/parser/` to `pkg/telegram/parser/`
- Fix the `tokenizer2` import alias artifact
- Update all import paths in `main.go`
- Remove old `telegram/` directory once fully migrated

### 1.2 Clean up stale references
- Remove placeholder/easter-egg functions from `main.go`
- Ensure all `pkg/` packages compile cleanly with `go build ./...`
- Run `go vet ./...` and fix any issues

---

## Phase 2: Complete COSEM Type System

### 2.1 Implement missing types
Each type should follow the existing option-function pattern (`New()` + `WithX()` options).

| Type | File | Notes |
|------|------|-------|
| `Integer` | `pkg/cosem/integer/integer.go` | Validates `In` format (e.g. `I4` → 4-digit integer) |
| `OctetString` | `pkg/cosem/octet_string/octet_string.go` | Hex-encoded string, 2 hex chars per octet, decode to bytes |
| `Timestamp` | `pkg/cosem/timestamp/timestamp.go` | Parse `YYMMDDhhmmssX` where `X` is `S` (summer/DST) or `W` (winter), expose `time.Time` |
| `Enum` | `pkg/cosem/enum/enum.go` | Integer-backed enumeration (tag 22) |

### 2.2 Enrich existing types
- `FloatingPoint`: add `Fn(x,y)` format metadata so callers can express e.g. `F9(3,3)` constraints directly
- `String`: distinguish `VisibleString` (printable ASCII) from `OctetString` (hex-encoded bytes)

### 2.3 Unit type
- Move unit constants (kWh, kW, V, A, m3, GJ, s) into `pkg/cosem/unit/unit.go`
- Reference units from OBIS format definitions instead of raw strings

---

## Phase 3: Complete OBIS Reference Registry

### 3.1 Verify and fill missing OBIS codes
Cross-reference with the spec table. Currently missing or incomplete:

| OBIS Code | Description | Status |
|-----------|-------------|--------|
| `1-0:99.97.0` | Power failure event log (profile generic) | **Missing** |
| `0-0:96.13.0` | Text message (max 2048 chars) | **Missing** |
| `0-n:24.1.0` | M-Bus Device-Type (channels 1-4) | **Missing** |
| `0-n:96.1.0` | M-Bus equipment identifier (channels 1-4) | **Missing** |
| `0-n:24.2.1` | M-Bus last 5-min value (channels 1-4) | **Partial** (only channel 1) |

For M-Bus codes, register all 4 channels (`n` = 1–4).

### 3.2 Add format metadata to each reference
Each `Reference` should carry its exact COSEM format spec:
- Tag (COSEM data type)
- Format string (e.g. `F9(3,3)`, `S96`, `TST`, `I4`)
- Unit
- Class and attribute

This enables the parser to validate values at parse time.

### 3.3 Separate poly-phase references
Mark L2/L3 OBIS codes as poly-phase only so single-phase telegrams can be validated correctly.

---

## Phase 4: Parser Enhancements

### 4.1 CRC16 validation
- Implement CRC-16 with polynomial `0xA001` (reflected), no XOR-in/out
- Calculate over all bytes from `/` through `!` (inclusive)
- Compare against the 4-hex-char footer value
- Reject telegrams with CRC mismatch (log warning, skip telegram)

### 4.2 Power failure event log parsing
Format: `(count)(0-0:96.7.19)(TST1)(duration1*s)(TST2)(duration2*s)...`
- Parse the count
- Extract up to 10 `(timestamp, duration)` pairs
- Store as structured data on the `Data` object (e.g. `[]PowerFailureEvent`)

### 4.3 M-Bus multi-value lines
Format: `0-n:24.2.1(timestamp)(value*unit)`
- Parse the capture timestamp as first value group
- Parse the metered value as second value group
- Identify device type from `0-n:24.1.0` to determine unit and format

### 4.4 Typed value parsing
Instead of storing all values as `[]string`, parse into concrete COSEM types:
- Look up the OBIS reference format
- Instantiate the correct COSEM type (`FloatingPoint`, `Integer`, `Timestamp`, `OctetString`, `String`)
- Validate against format constraints (decimal places, length, range)
- Expose a `Value() interface{}` or typed accessor on `Data`

### 4.5 Header validation
- Validate the `5` baud rate identifier in position 4 of the header
- Extract manufacturer (3 chars), identification (variable length) per IEC 62056-21

---

## Phase 5: Prometheus Integration ✓

### 5.1 Reconnect metric export ✓
- `obis.Register()` returns `*Metrics` (gauges + GaugeVecs)
- `parser.WithMetrics(m)` injects metrics into the parser
- `parser.updateMetrics()` iterates `DataMap` on each complete telegram
- `toFloat64()` converts `*FloatingPoint`, `*Integer`, and raw strings

### 5.2 M-Bus metrics ✓
- Single `mbus_last_value{channel="n"}` GaugeVec (channel label "1"–"4")
- Metric names removed from `0-n:24.2.1` references; handled in parser

### 5.3 Metadata metrics ✓
- `dsmr_info{version="vv"} 1` — DSMR version info metric (GaugeVec)
- `electricity_equipment_info{identifier="..."} 1` — equipment info metric (GaugeVec)
- `electricity_tariff` — tariff indicator gauge (S4 string parsed as float64)

---

## Phase 6: Configuration & Runtime

### 6.1 CLI flags and environment variables
- Serial port path (default `/dev/ttyUSB0`)
- Baud rate (default `115200`)
- Prometheus listen address (default `:8080`)
- Log level

### 6.2 Graceful shutdown
- Handle `SIGINT` / `SIGTERM`
- Close serial port cleanly
- Drain in-flight telegram processing
- Shut down HTTP server with timeout

### 6.3 Serial port auto-detection
- Optionally scan for known USB-to-serial VID/PID pairs used by P1 cables

---

## Phase 7: Testing

### 7.1 Unit tests for new types
- `Integer`: valid/invalid formats, boundary values
- `OctetString`: hex encoding/decoding, odd-length rejection
- `Timestamp`: DST flag parsing, time zone handling, edge dates
- `Enum`: valid range, unknown values

### 7.2 CRC16 tests
- Verify against the example telegram in `examples/telegram_v5_0_2.txt`
- Test with known CRC values from the spec
- Test rejection of corrupted telegrams

### 7.3 Integration tests
- Parse the full example telegram end-to-end
- Verify all 40+ OBIS values are extracted correctly
- Verify M-Bus data extraction
- Verify power failure log extraction

### 7.4 Edge cases
- Empty telegram (header + footer, no data)
- Single-phase vs poly-phase telegrams
- Maximum length text messages (1024 / 2048 chars)
- M-Bus device on different channels
- Missing or malformed data lines (parser should skip, not crash)

---

## Phase 8: Documentation

### 8.1 README
- Usage instructions (build, run, configure)
- Prometheus metrics reference
- Supported OBIS codes table
- Example Grafana dashboard or PromQL queries

### 8.2 Go doc comments
- Public types and functions should have doc comments
- Package-level doc comments explaining each package's role

---

## Priority Order

| Priority | Phase | Rationale |
|----------|-------|-----------|
| **P0** | 1 (Restructuring) | Unblocks all other work |
| **P0** | 4.1 (CRC16) | Core protocol correctness |
| **P1** | 2 (COSEM types) | Foundation for typed parsing |
| **P1** | 3 (OBIS registry) | Completeness against spec |
| **P1** | 4.2–4.4 (Parser) | Full telegram interpretation |
| **P2** | 5 (Prometheus) | Usable metrics output |
| **P2** | 6 (Config/runtime) | Production readiness |
| **P3** | 7 (Testing) | Ongoing alongside each phase |
| **P3** | 8 (Documentation) | After stabilization |
