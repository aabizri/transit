package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/aabizri/navitia"
	"github.com/aabizri/navitia/pretty"
	"github.com/aabizri/transit/maps"
	"github.com/fogleman/gg"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var coverageCommand = cli.Command{
	Name:    "coverage",
	Aliases: []string{"c"},
	Usage:   "List coverage",
	Action:  coverageAction,
}

func coverageAction(c *cli.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
	defer cancel()

	req := navitia.RegionRequest{
		Geo: true,
	}

	res, err := session.Regions(ctx, req)
	if err != nil {
		return errors.Wrap(err, "coverageAction: error while retrieving list of regions")
	}

	err = pretty.DefaultRegionResultsConf.PrettyWrite(res, os.Stdout)
	if err != nil {
		return errors.Wrap(err, "error while pretty-printing")
	}

	// If enabled, draw !
	if path := c.GlobalString("map"); path != "" {
		if !filepath.IsAbs(path) {
			wd, err := os.Getwd()
			if err != nil {
				return errors.Wrap(err, "couldn't retrieve working directory, consider giving an absolute path")
			}

			path = filepath.Join(wd, path)
		}

		img, err := maps.DrawRegions(res.Regions)
		if err != nil {
			return errors.Wrap(err, "error while drawing places")
		}

		err = gg.SavePNG(path, img)
		if err != nil {
			return errors.Wrapf(err, "error while saving image to %s", path)
		}
	}

	return nil
}
