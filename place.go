package main

import (
	"fmt"
	"github.com/aabizri/gonavitia"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var placeCommand = cli.Command{
	Name:    "place",
	Aliases: []string{"p"},
	Usage:   "Search for places",
	Action:  placeAction,
}

func placeAction(c *cli.Context) error {
	var response string
	for _, query := range c.Args() {
		req := gonavitia.PlacesRequest{Query: query}

		res, err := session.Places(req)
		if err != nil {
			return errors.Wrap(err, "Error while requesting places")
		}
		response += fmt.Sprintf("\nQuery \"%s\" (%d results):\n%s", query, len(res.Places), res.String())
	}
	fmt.Print(response)
	return nil
}
