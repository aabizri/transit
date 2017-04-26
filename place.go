package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aabizri/navitia"
	"github.com/aabizri/navitia/pretty"
	"github.com/fatih/color"
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
	ctx := context.Background()
	for i, query := range c.Args() {
		ctx, cancel := context.WithTimeout(ctx, requestTimeout)
		defer cancel()

		req := navitia.PlacesRequest{Query: query, Count: c.Uint("count")}

		pr, err := session.Places(ctx, req)
		if err != nil {
			return errors.Wrap(err, "Error while requesting places")
		}

		// Now let's print
		msg := fmt.Sprintf("\n[%d/%d] Query \"%s\" ", i, len(c.Args()), color.New(color.Underline, color.FgHiCyan).Sprint(query))
		buf := bytes.NewBuffer([]byte(msg))
		err = pretty.DefaultPlacesResultsConf.PrettyWrite(pr, buf)
		if err != nil {
			return errors.Wrapf(err, "Error while preparing result output for query #%d", i)
		}
		io.Copy(os.Stdout, buf)
	}
	return nil
}
