package next

import (
	"fmt"
	"os"
	"syscall"

	"github.com/shadowfax-chc/wallpaper/command/internal/pid"
	"gopkg.in/urfave/cli.v1"
)

// Action is the action function for the next command.
func Action(c *cli.Context) error {
	p, err := pid.Get(c.String("pid-file"))
	if err != nil {
		fmt.Printf("Could not get pid, is wp running?\n")
		return err
	}
	proc, err := os.FindProcess(p)
	if err != nil {
		fmt.Printf("Could not get process for pid %d, is wp running?\n", p)
		return err
	}
	return proc.Signal(syscall.SIGUSR1)
}
