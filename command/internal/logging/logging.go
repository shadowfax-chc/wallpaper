package logging

import (
	"github.com/hashicorp/logutils"
	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"

	"github.com/tmessi/wallpaper/logging"
)

// ReloadLogger setups up the logger base on the cli args.
func ReloadLogger(ctx *cli.Context) error {
	return logging.Setup(&logging.Config{
		Handler: ctx.GlobalString("log-handler"),
		Level:   ctx.GlobalString("log-level"),
		Levels:  []logutils.LogLevel{"TRACE", "DEBUG", "INFO", "WARN", "ERROR"},
	})
}

// HandleLogger is a wrapper action for setting up logging.
func HandleLogger(af cli.ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		err := ReloadLogger(ctx)
		if err != nil {
			return err
		}
		return af(ctx)
	}
}

// Flags get the command line flags for logging.
func Flags() []cli.Flag {
	return []cli.Flag{
		altsrc.NewStringFlag(cli.StringFlag{
			Name:   "log-handler",
			Usage:  "Set where the logs will go",
			EnvVar: "WP_LOG_HANDLER",
			Value:  "stdout",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:   "log-level",
			Usage:  "Set the logging verbosity level.",
			EnvVar: "WP_LOG_LEVEL",
			Value:  "WARN",
		}),
	}
}
