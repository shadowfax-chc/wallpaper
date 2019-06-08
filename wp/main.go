package main

import (
	"os"

	"github.com/shadowfax-chc/wallpaper/command"
)

func main() {
	app := command.App()
	app.Run(os.Args)
}
