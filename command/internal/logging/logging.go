package logging

import (
	"github.com/hashicorp/logutils"
	"github.com/urfave/cli"

	"github.com/shadowfax-chc/wallpaper/logging"
)

func enableLogger(ctx *cli.Context) error {
	return logging.Setup(&logging.Config{
		Handler: ctx.GlobalString("log-handler"),
		Level:   ctx.GlobalString("log-level"),
		Levels:  []logutils.LogLevel{"TRACE", "DEBUG", "INFO", "WARN", "ERROR"},
	})
}

// HandleLogger is a wrapper action for setting up logging.
func HandleLogger(af cli.ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		err := enableLogger(ctx)
		if err != nil {
			return err
		}
		return af(ctx)
	}
}

// Flags get the command line flags for logging.
func Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:   "log-handler",
			Usage:  "Set where the logs will go",
			EnvVar: "WP_LOG_HANDLER",
			Value:  "stdout",
		},
		cli.StringFlag{
			Name:   "log-level",
			Usage:  "Set the logging verbosity level.",
			EnvVar: "WP_LOG_LEVEL",
			Value:  "WARN",
		},
	}
}
