package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	ph "github.com/pyaillet/placeholder/pkg/placeholder"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "placeholder"
	app.Usage = "Manage the placeholders"
	app.Commands = []cli.Command{
		{
			Name:    "list",
			Aliases: []string{"ls"},
			Usage:   "list placeholders in the provided files",
			Action: func(c *cli.Context) error {
				placeHolders := ph.ListPlaceHoldersInFiles(c.Args(), ph.DefaultSeparator())
				fmt.Printf("%s\n", strings.Join(placeHolders, "\n"))
				return nil
			},
		},
		{
			Name:    "replace",
			Aliases: []string{"rp"},
			Usage:   "replace placeholders in the provided files",
			Action: func(c *cli.Context) error {
				err := ph.ReplacingPlaceHoldersInFilesFromEnv(c.Args(), ph.DefaultSeparator())
				return err
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
