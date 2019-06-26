package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pyaillet/placeholder/pkg/placeholder"
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
				start := c.GlobalString("start")
				end := c.GlobalString("end")
				sep := placeholder.SeparatorFrom(start, end)
				placeHolders := placeholder.ListPlaceHoldersInFiles(c.Args(), sep)
				fmt.Printf("%s", strings.Join(placeHolders, "\n"))
				return nil
			},
		},
		{
			Name:    "replace",
			Aliases: []string{"rp"},
			Usage:   "replace placeholders in the provided files",
			Action: func(c *cli.Context) error {
				start := c.GlobalString("start")
				end := c.GlobalString("end")
				sep := placeholder.SeparatorFrom(start, end)
				input := c.String("input")
				var provider placeholder.ValuesProvider
				if len(input) == 0 {
					provider = placeholder.EnvProvider{}
				} else {
					var err error
					provider, err = placeholder.NewFileProvider(input)
					if err != nil {
						panic(err)
					}
				}
				err := placeholder.ReplacingPlaceHoldersInFiles(c.Args(), sep, provider)
				return err
			},
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "input, i",
					Usage: "input file (json or yaml)",
				},
			},
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "start, s",
			Value: "${",
			Usage: "Separator starting separator",
		},
		cli.StringFlag{
			Name:  "end, e",
			Value: "}",
			Usage: "Separator ending separator",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
