package command

import (
	"gopkg.in/urfave/cli.v1"

	"github.com/tmessi/wallpaper/command/internal/logging"
	"github.com/tmessi/wallpaper/command/internal/next"
	"github.com/tmessi/wallpaper/command/internal/pid"
	"github.com/tmessi/wallpaper/command/internal/reload"
	"github.com/tmessi/wallpaper/command/internal/run"
	"github.com/tmessi/wallpaper/version"
)

// App returns the cli App with its subcomands and flags defined.
func App() *cli.App {
	app := cli.NewApp()
	app.Name = "wp"
	app.Version = version.Description()

	flags := append(run.Flags(), logging.Flags()...)
	flags = append(flags, pid.Flags()...)

	app.Before = run.Before(flags)
	app.Flags = flags
	app.Action = pid.HandlePIDFile(logging.HandleLogger(run.Action))

	app.Commands = []cli.Command{
		{
			Name:   "next",
			Usage:  "send singal to use the next image",
			Flags:  append(logging.Flags(), pid.Flags()...),
			Action: logging.HandleLogger(next.Action),
		},
		{
			Name:   "reload",
			Usage:  "send singal to reload the config",
			Flags:  append(logging.Flags(), pid.Flags()...),
			Action: logging.HandleLogger(reload.Action),
		},
	}
	return app
}
