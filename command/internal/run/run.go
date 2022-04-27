/*
Package run defines the main run command
*/
package run

import (
	"context"
	"log"
	"os"
	"os/signal"
	"path"
	"syscall"

	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"

	"github.com/tmessi/wallpaper/command/internal/logging"
	"github.com/tmessi/wallpaper/directory"
	"github.com/tmessi/wallpaper/wallpaper"
)

// Before is used as a cli.BeforeFunc that is called before Action.
func Before(flags []cli.Flag) cli.BeforeFunc {
	return altsrc.InitInputSourceWithContext(flags, altsrc.NewYamlSourceFromFlagFunc("config"))
}

// Action is the action function for the command.
func Action(c *cli.Context) error {
	config := &directory.Config{
		Root:    c.String("directory"),
		Shuffle: c.Bool("shuffle"),
	}
	r, err := directory.NewRepository(config)
	if err != nil {
		return err
	}
	r.Load()

	uc := &wallpaper.UpdaterConfig{
		Mode:       wallpaper.Mode(c.String("mode")),
		Repository: r,
		Frequency:  c.Duration("update-frequency"),
	}
	updater := wallpaper.NewUpdater(uc)

	sigs := make(chan os.Signal, 1)
	signal.Notify(
		sigs,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGHUP,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)
	ctx, cancel := context.WithCancel(context.Background())

	// Signal handling
	go func() {
		for {
			s := <-sigs
			log.Printf("[DEBUG] signal %s", s)
			switch s {
			case syscall.SIGUSR1:
				updater.Next()
			case syscall.SIGHUP:
				Before(c.App.Flags)(c)
				logging.ReloadLogger(c)
				rc := &wallpaper.ReloadConfig{
					Mode:      wallpaper.Mode(c.String("mode")),
					Location:  c.String("directory"),
					Shuffle:   c.Bool("shuffle"),
					Frequency: c.Duration("update-frequency"),
				}
				updater.Reload(rc)
			case syscall.SIGTERM:
				fallthrough
			case syscall.SIGINT:
				cancel()
				return
			}
		}
	}()

	return updater.Run(ctx)
}

// Flags gets the command line flags for this action.
func Flags() []cli.Flag {
	return []cli.Flag{
		cli.StringFlag{
			Name:  "config",
			Value: path.Join(os.Getenv("HOME"), ".wp.yaml"),
		},
		altsrc.NewStringFlag(cli.StringFlag{
			Name:   "directory",
			Usage:  "The root directory that contains background images",
			EnvVar: "WP_DIRECTORY,WALLPAPERS",
			Value:  path.Join(os.Getenv("HOME"), ".wallpaper"),
		}),
		altsrc.NewBoolFlag(cli.BoolFlag{
			Name:   "shuffle",
			Usage:  "Randomize the images",
			EnvVar: "WP_SHUFFLE",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:   "mode",
			Usage:  "The background mode to apply to the image so that it fits the screen",
			EnvVar: "WP_MODE",
			Value:  "fill",
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:   "update-frequency",
			Usage:  "How often to update the background image",
			EnvVar: "WP_UPDATE_FREQ",
			Value:  "5m",
		}),
	}
}
