package main

import (
	"fmt"
	"github.com/aabizri/gonavitia"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var placeFlags = []cli.Flag{
	cli.UintFlag{
		Name:  "count, c",
		Value: 6,
	},
}

var placeCommand = cli.Command{
	Name:    "place",
	Aliases: []string{"p"},
	Usage:   "Search for places",
	Action:  placeAction,
	Flags:   placeFlags,
}

func placeAction(c *cli.Context) error {
	for i, query := range c.Args() {
		req := gonavitia.PlacesRequest{Query: query, Count: c.Uint("count")}

		res, err := session.Places(req)
		if err != nil {
			return errors.Wrap(err, "Error while requesting places")
		}
		fmt.Printf("\n[%d/%d] Query \"%s\" (%d results):\n%s", i, len(c.Args()), query, len(res.Places), res.String())
	}
	return nil
}
