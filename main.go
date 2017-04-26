/*
travel is a command line tool for planning your travel between two places
*/
package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aabizri/navitia"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

const standardRequestTimeout = 10 * time.Second

var commands = []cli.Command{
	placeCommand,
	journeyCommand,
}

var (
	apiKey         string
	requestTimeout time.Duration
	session        *navitia.Session
)

var flags = []cli.Flag{
	cli.StringFlag{
		Name:        "key",
		Usage:       "Api Key for navitia.io",
		Destination: &apiKey,
		EnvVar:      "NAVITIA_TOKEN",
	},
	cli.DurationFlag{
		Name:        "req-timeout, rt",
		Usage:       "Request timeout value",
		Value:       standardRequestTimeout,
		Destination: &requestTimeout,
	},
	cli.StringFlag{
		Name:  "in",
		Usage: "Coverage to use [NOT IMPLEMENTED]",
	},
}

var authors = []cli.Author{
	cli.Author{
		Name:  "Alexandre A. Bizri",
		Email: "alexandre@bizri.fr",
	},
}

const description = "transit is a tool for planning, monitoring and searching public transit information"

func establishSession(ctx *cli.Context) error {
	if apiKey == "" {
		return errors.Errorf("ERROR: No Api Key specified")
	}

	var err error
	session, err = navitia.NewCustom(apiKey, "http://api.navitia.io/v1", &http.Client{})
	return errors.Wrap(err, "Error while creating session")
}

func main() {
	app := cli.NewApp()
	app.Version = "-dev"
	app.Description = description
	app.Name = "transit"
	app.Usage = "The public transit tool for the CLI"
	app.Authors = authors
	app.Before = establishSession
	app.Flags = flags
	app.Commands = commands
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
