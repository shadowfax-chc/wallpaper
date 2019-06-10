package pid

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"gopkg.in/urfave/cli.v1"
	"gopkg.in/urfave/cli.v1/altsrc"
)

func writePID(pidfile string) error {
	pid := strconv.Itoa(os.Getpid())
	return ioutil.WriteFile(pidfile, []byte(pid), 0644)
}

func rmPID(pidfile string) error {
	return os.Remove(pidfile)
}

// Get returns the recorded pid.
func Get(file string) (int, error) {
	dat, err := ioutil.ReadFile(file)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(string(dat))
}

// HandlePIDFile manages setting and removing the pidfile.
func HandlePIDFile(af cli.ActionFunc) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		pidfile := ctx.String("pid-file")
		if err := writePID(pidfile); err != nil {
			return err
		}
		defer rmPID(pidfile)
		return af(ctx)
	}
}

// Flags get the command line flags for setting the pid file
func Flags() []cli.Flag {
	return []cli.Flag{
		altsrc.NewStringFlag(cli.StringFlag{
			Name:   "pid-file",
			Usage:  "the file name were the pid is stored",
			EnvVar: "WP_PID",
			Value:  fmt.Sprintf("/var/run/user/%d/wp.pid", os.Getuid()),
		}),
	}
}
