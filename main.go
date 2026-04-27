// Command dsmr reads DSMR P1 telegrams from a serial port and exposes parsed
// values as Prometheus metrics over HTTP.
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
	"time"

	"github.com/hyperized/dsmr/pkg/obis"
	"github.com/hyperized/dsmr/pkg/sink/prom"
	"github.com/hyperized/dsmr/pkg/telegram"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.bug.st/serial"
)

func main() {
	portFlag := flag.String("port", "", "serial port path (default: auto-detect)")
	baudFlag := flag.Int("baud", 115200, "serial port baud rate")
	listenFlag := flag.String("listen", ":8080", "Prometheus metrics listen address")
	debugFlag := flag.Bool("debug", false, "enable debug logging")
	flag.Parse()

	logLevel := slog.LevelInfo
	if *debugFlag {
		logLevel = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: logLevel})))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// HTTP server
	mux := http.NewServeMux()
	mux.Handle("/", http.RedirectHandler("/metrics", http.StatusFound))
	mux.Handle("/metrics", promhttp.Handler())
	srv := &http.Server{Addr: *listenFlag, Handler: mux}

	go func() {
		slog.Info("starting prometheus webserver", "addr", *listenFlag)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("http server error", "err", err)
			stop()
		}
	}()

	// Serial port
	portPath := *portFlag
	if portPath == "" {
		var err error
		portPath, err = detectPort()
		if err != nil {
			slog.Error("serial port auto-detection failed", "err", err)
			os.Exit(1)
		}
		slog.Info("auto-detected serial port", "port", portPath)
	}

	mode := &serial.Mode{
		BaudRate: *baudFlag,
		DataBits: 8,
		Parity:   serial.NoParity,
		StopBits: serial.OneStopBit,
	}

	slog.Info("connecting to serial port", "port", portPath)
	port, err := serial.Open(portPath, mode)
	if err != nil {
		slog.Error("failed to open serial port", "port", portPath, "err", err)
		os.Exit(1)
	}

	m := obis.Register(prometheus.DefaultRegisterer)
	p := telegram.NewParser(port, telegram.WithSink(prom.New(m)))

	done := make(chan struct{})
	go func() {
		defer close(done)
		p.ParseStream()
	}()

	// Wait for shutdown signal or parser exit.
	select {
	case <-ctx.Done():
		slog.Info("shutdown signal received")
	case <-done:
		slog.Info("parser exited")
	}

	stop() // release signal resources

	// Close serial port so the scanner unblocks if still running.
	if err := port.Close(); err != nil {
		slog.Warn("error closing serial port", "err", err)
	}

	// Wait for parser to drain in-flight telegrams.
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		slog.Warn("timed out waiting for parser to stop")
	}

	// Graceful HTTP server shutdown.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Warn("http server shutdown error", "err", err)
	}

	slog.Info("shutdown complete")
}

// knownP1Prefixes lists USB-to-serial port name patterns used by common P1 cables.
var knownP1Prefixes = []string{
	"/dev/ttyUSB",       // Linux: FTDI, CP210x
	"/dev/ttyACM",       // Linux: CDC ACM
	"/dev/cu.usbserial", // macOS: FTDI
	"/dev/cu.usbmodem",  // macOS: CDC ACM
}

// detectPort returns the first serial port that looks like a P1 cable.
// Falls back to the first available port if none match the known prefixes.
func detectPort() (string, error) {
	ports, err := serial.GetPortsList()
	if err != nil {
		return "", fmt.Errorf("listing serial ports: %w", err)
	}
	if len(ports) == 0 {
		return "", fmt.Errorf("no serial ports found")
	}

	for _, prefix := range knownP1Prefixes {
		if i := slices.IndexFunc(ports, func(p string) bool {
			return strings.HasPrefix(p, prefix)
		}); i >= 0 {
			return ports[i], nil
		}
	}

	slog.Debug("no known P1 port prefix matched, using first available port", "ports", ports)
	return ports[0], nil
}
