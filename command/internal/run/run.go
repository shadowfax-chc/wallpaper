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
	"time"

	"github.com/shadowfax-chc/wallpaper/directory"
	"github.com/shadowfax-chc/wallpaper/wallpaper"
	"github.com/urfave/cli"
)

// Action is the action function for the command.
func Action(c *cli.Context) error {
	config := &directory.Config{
		Root:      c.String("directory"),
		Recursive: c.Bool("recursive"),
		Shuffle:   c.Bool("shuffle"),
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
				updater.Reload()
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
			Name:   "directory,d",
			Usage:  "The root directory that contains background images",
			EnvVar: "WP_DIRECTORY,WALLPAPERS",
			Value:  path.Join(os.Getenv("HOME"), ".wallpaper"),
		},
		cli.BoolFlag{
			Name:   "recursive,r",
			Usage:  "Scan recursively down the directory for images",
			EnvVar: "WP_RECURSIVE",
		},
		cli.BoolFlag{
			Name:   "shuffle,s",
			Usage:  "Randomize the images",
			EnvVar: "WP_SHUFFLE",
		},
		cli.StringFlag{
			Name:   "mode",
			Usage:  "The background mode to apply to the image so that it fits the screen",
			EnvVar: "WP_MODE",
			Value:  "fill",
		},
		cli.DurationFlag{
			Name:   "update-frequency,u",
			Usage:  "How often to update the background image",
			EnvVar: "WP_UPDATE_FREQ",
			Value:  5 * time.Minute,
		},
	}
}
