package command

import (
	"github.com/urfave/cli"

	"github.com/shadowfax-chc/wallpaper/command/internal/logging"
	"github.com/shadowfax-chc/wallpaper/command/internal/run"
	"github.com/shadowfax-chc/wallpaper/version"
)

// App returns the cli App with its subcomands and flags defined.
func App() *cli.App {
	app := cli.NewApp()
	app.Name = "wp"
	app.Version = version.Description()
	app.Flags = append(run.Flags(), logging.Flags()...)
	app.Action = logging.HandleLogger(run.Action)
	return app
}
