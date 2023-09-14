package assets

import (
	"os"

	"decred.org/dcrwallet/errors"
	"github.com/decred/slog"
	"github.com/jrick/logrotate/rotator"
)

// logWriter implements an io.Writer that outputs to both standard output and
// the write-end pipe of an initialized log rotator.
type logWriter struct {
	logRotator *rotator.Rotator
}

func (lw *logWriter) Write(p []byte) (n int, err error) {
	os.Stdout.Write(p)
	lw.logRotator.Write(p)
	return len(p), nil
}

// initLogRotator initializes the logging rotater to write logs to logFile and
// create roll files in the same directory.  It must be called before the
// package-global log rotater variables are used.
func initLogRotator(logFile string) (*slog.Backend, error) {
	r, err := rotator.New(logFile, 10*1024, false, 3)
	if err != nil {
		return nil, errors.Errorf("failed to create file rotator: %v", err)
	}

	return slog.NewBackend(&logWriter{logRotator: r}), nil
}
