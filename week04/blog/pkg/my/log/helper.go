package log

import "fmt"

// Helper is a logger helper.
type Helper struct {
	debug Logger
	info  Logger
	warn  Logger
	err   Logger
}

// NewHelper new a logger helper.
func NewHelper(name string, logger Logger) *Helper {
	const LevelKey = "level"
	logger = With(logger, "module", name)
	return &Helper{
		debug: With(logger, LevelKey, "Debug"),
		info:  With(logger, LevelKey, "Info" ),
		warn:  With(logger, LevelKey, "Warn"),
		err:   With(logger, LevelKey, "Error"),
	}
}

// With with logger fields.
func With(l Logger, kv ...interface{}) Logger {
	if c, ok := l.(*context); ok {
		kvs := make([]interface{}, 0, len(c.prefix)+len(kv))
		kvs = append(kvs, kv...)
		kvs = append(kvs, c.prefix...)
		return &context{
			logs:      c.logs,
			prefix:    kvs,
			hasValuer: containsValuer(kvs),
		}
	}
	return &context{logs: []Logger{l}, prefix: kv, hasValuer: containsValuer(kv)}
}

func containsValuer(keyvals []interface{}) bool {
	for i := 1; i < len(keyvals); i += 2 {
		if _, ok := keyvals[i].(Valuer); ok {
			return true
		}
	}
	return false
}

func (h *Helper) Debug(a ...interface{}) {
	h.debug.Log("msg", fmt.Sprint(a...))
}
func (h *Helper) Debugf(format string, a ...interface{}) {
	h.debug.Log("msg", fmt.Sprintf(format, a...))
}

func (h *Helper) Info(a ...interface{}) {
	h.info.Log("msg", fmt.Sprint(a...))
}
func (h *Helper) Infof(format string, a ...interface{}) {
	h.info.Log("msg", fmt.Sprintf(format, a...))
}
func (h *Helper) Infow(kv ...interface{}) {
	h.info.Log(kv...)
}

func (h *Helper) Warnf(format string, a ...interface{}) {
	h.warn.Log("msg", fmt.Sprintf(format, a...))
}

func (h *Helper) Errorf(format string, a ...interface{}) {
	h.err.Log("msg", fmt.Sprintf(format, a...))
}
func (h *Helper) Error(a ...interface{}) {
	h.err.Log("msg", fmt.Sprint(a...))
}
