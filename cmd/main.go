package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/pyaillet/placeholder/pkg/placeholder"
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
			Usage:   "list place holders on the provided files",
			Action: func(c *cli.Context) error {
				placeHolders := placeholder.ListPlaceHoldersInFiles(c.Args())
				fmt.Printf("%s\n", strings.Join(placeHolders, "\n"))
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
