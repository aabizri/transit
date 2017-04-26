package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aabizri/navitia"
	"github.com/aabizri/navitia/types"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

// Add departure, arrival, mode, etc.
var journeyFlags = []cli.Flag{
	cli.UintFlag{
		Name:  "count, c",
		Usage: "Amount of journey results to display",
	},
	cli.StringFlag{
		Name:  "departure, dep",
		Usage: "Departure date/time (NOT IMPLEMENTED)",
	},
	cli.StringFlag{
		Name:  "arrival, arr",
		Usage: "Arrival date/time (NOT IMPLEMENTED)",
	},
	cli.UintFlag{
		Name:  "max-transfers",
		Usage: "Maximum number of transfers",
	},
	cli.DurationFlag{
		Name:  "max-duration, max",
		Usage: "Maximum duration of journey",
	},
}

var journeyCommand = cli.Command{
	Name:    "journey",
	Aliases: []string{"j"},
	Usage:   "Build journey propositions",
	Action:  journeyAction,
	Flags:   journeyFlags,
}

func parseJourneyArgs(args []string) (from string, to string, err error) {
	if len(args) > 4 || len(args) == 0 {
		return "", "", errors.Errorf("Number of arguments for journey incorect (%d<%d<%d)", 2, len(args), 4)
	}

	for len(args) >= 2 {
		switch args[0] {
		case "from":
			from = args[1]
			args = args[2:]
		case "to":
			to = args[1]
			args = args[2:]
		}
	}

	return
}

/*
journeyAction works like that:
	- First, retrieve the from and to arguments
	- For each of them, call a goroutine which will retrieve the most likely result
	- Then query
*/
func journeyAction(c *cli.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	fromQuery, toQuery, err := parseJourneyArgs(c.Args())
	if err != nil {
		fmt.Print(err)
		return err
	}

	fromChan := make(chan types.Container)
	toChan := make(chan types.Container)
	getPlace := func(query string, c chan types.Container) {

		req := navitia.PlacesRequest{Query: query, Count: 1}

		res, err := session.Places(ctx, req)
		if err != nil {
			panic(errors.Wrap(err, "Error while requesting places"))
		} else if len(res.Places) == 0 {
			panic("Not enough responses")
		}

		c <- res.Places[0]
	}

	if fromQuery != "" {
		go getPlace(fromQuery, fromChan)
	}

	if toQuery != "" {
		go getPlace(toQuery, toChan)
	}

	var (
		from types.ID
		to   types.ID
	)
	// While both haven't returned, wait
	for (fromQuery != "" && from == "") || (toQuery != "" && to == "") {
		select {
		case rec := <-fromChan:
			from = rec.ID
		case rec := <-toChan:
			to = rec.ID
		}
	}

	// Build query
	req := navitia.JourneyRequest{
		Count:        c.Uint("count"),
		MaxDuration:  c.Duration("max-duration"),
		MaxTransfers: c.Uint("max-transfers"),
	}
	if from != "" {
		req.From = from
	}
	if to != "" {
		req.To = to
	}

	// Send it
	res, err := session.Journeys(ctx, req)
	fmt.Printf("Got %d journeys\n", res.Count())
	journeyResultsWrite(res, os.Stdout)
	if err != nil {
		fmt.Printf("Got an error: %v", err)
		return err
	}

	return nil
}
