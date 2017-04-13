/*
travel is a command line tool for planning your travel between two places
*/
package main

import (
	"github.com/aabizri/gonavitia"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
	"net/http"
	"os"
)

var commands = []cli.Command{
	placeCommand,
	journeyCommand,
}

var (
	apiKey  string
	session *gonavitia.Session
)

var flags = []cli.Flag{
	cli.StringFlag{
		Name:        "key",
		Usage:       "Api Key for navitia.io",
		Destination: &apiKey,
		EnvVar:      "NAVITIA_TOKEN",
	},
}

func establishSession(ctx *cli.Context) error {
	if apiKey == "" {
		return errors.Errorf("ERROR: No Api Key specified")
	}

	var err error
	session, err = gonavitia.NewCustom(apiKey, "http://api.navitia.io/v1", &http.Client{})
	return errors.Wrap(err, "Error while creating session")
}

func main() {
	app := cli.NewApp()
	app.Before = establishSession
	app.Flags = flags
	app.Commands = commands
	app.Run(os.Args)
}
