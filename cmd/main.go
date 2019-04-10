package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	ph "github.com/pyaillet/placeholder/pkg/placeholder"
	"github.com/urfave/cli"
)

var version string

func main() {
	app := cli.NewApp()
	app.Name = "placeholder"
	app.Version = version
	app.Usage = "Manage the placeholders"
	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "list placeholders in the provided files",
			Action: func(c *cli.Context) error {
				sep := ph.SeparatorFrom(c.String("start"), c.String("end"))
				placeHolders := ph.ListPlaceHoldersInFiles(c.Args(), sep)
				fmt.Printf("%s\n", strings.Join(placeHolders, "\n"))
				return nil
			},
		},
		{
			Name:    "replace",
			Aliases: []string{"rp"},
			Usage:   "replace placeholders in the provided files",
			Action: func(c *cli.Context) error {
				sep := ph.SeparatorFrom(c.String("start"), c.String("end"))
				err := ph.ReplacingPlaceHoldersInFilesFromEnv(c.Args(), sep)
				return err
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "start, s",
			Value: "%#",
			Usage: "Separator starting separator",
		},
		cli.StringFlag{
			Name:  "end, e",
			Value: "#%",
			Usage: "Separator ending separator",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
