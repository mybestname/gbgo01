package config

import "my/log"

// Option is config option.
type Option func(*options)

type options struct {
	sources []Source
	decoder Decoder
	logger  log.Logger
}

// Source is config source.
type Source interface {
	Load() ([]*KeyValue, error)
	Watch() (Watcher, error)
}

// Watcher watches a source for changes.
type Watcher interface {
	Next() ([]*KeyValue, error)
	Stop() error
}

// KeyValue is config key value.
type KeyValue struct {
	Key      string
	Value    []byte
	Metadata map[string]string
}

// WithSource with config source.
func WithSource(s ...Source) Option {
	return func(o *options) {
		o.sources = s
	}
}

// Decoder is config decoder.
type Decoder func(*KeyValue, map[string]interface{}) error

// WithDecoder with config decoder.
func WithDecoder(d Decoder) Option {
	return func(o *options) {
		o.decoder = d
	}
}

// WithLogger with config loogger.
func WithLogger(l log.Logger) Option {
	return func(o *options) {
		o.logger = l
	}
}