package main

import (
	"os"

	"github.com/tmessi/wallpaper/command"
)

func main() {
	app := command.App()
	app.Run(os.Args)
}
