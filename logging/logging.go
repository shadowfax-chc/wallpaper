// Package logging is used for common logging setup
package logging

import (
	"fmt"
	"io"
	"log"
	"log/syslog"
	"os"
	"path"

	"github.com/hashicorp/logutils"
)

// Config is used to configure the Writer and the log package
type Config struct {
	Name    string
	Handler string
	LogFile string // Only used if Handler is "file"
	Level   string
	Levels  []logutils.LogLevel
}

// Setup will initialize the logger.
func Setup(c *Config) error {
	var out io.Writer
	var err error
	var flags int

	logHandler := c.Handler
	switch logHandler {
	case "syslog":
		out, err = syslog.New(syslog.LOG_ALERT, path.Base(os.Args[0]))
		if err != nil {
			return err
		}
		flags = log.Llongfile
	case "stdout":
		out = os.Stdout
		flags = log.LstdFlags | log.Lmicroseconds | log.LUTC | log.Lshortfile
	default:
		fmt.Printf("Invalid log-handler: %s, Defaulting to stdout", logHandler)
		out = os.Stdout
		flags = log.LstdFlags | log.Lmicroseconds | log.LUTC | log.Llongfile
	}

	log.SetOutput(&logutils.LevelFilter{
		Levels:   c.Levels,
		MinLevel: logutils.LogLevel(c.Level),
		Writer:   out,
	})

	log.SetFlags(flags)
	return nil
}
