package telegram

// Sink consumes parsed telegrams and forwards them to a destination such as
// Prometheus, a database, or a log file.
//
// Implementations are responsible for their own thread safety and resource
// cleanup; sinks needing graceful shutdown should additionally implement
// io.Closer and be closed by the caller that owns them.
type Sink interface {
	Write(t *Telegram) error
}
