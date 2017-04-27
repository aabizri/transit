package maps

import (
	"fmt"
	"image"
	"image/color"

	"github.com/aabizri/navitia/types"
	"github.com/flopp/go-staticmaps"
	"github.com/golang/geo/s2"
	"github.com/pkg/errors"
)

// DrawPlaces draws a slice of coordinates-enabled containers
func DrawPlaces(places []types.Container) (image.Image, error) {
	smctx := sm.NewContext()
	smctx.SetTileProvider(sm.NewTileProviderCartoLight())
	smctx.SetSize(420, 420)
	for i, c := range places {
		pos, err := posContainer(c)
		if err != nil {
			return nil, errors.Wrapf(err, "error while getting latitude/longitude for place #%d", i)
		}

		marker := &sm.Marker{
			Position:   pos,
			Color:      color.RGBA{0xff, 0, 0, 0xff},
			Size:       16.0,
			Label:      fmt.Sprintf("%d. (%s)", i, c.Name),
			LabelColor: color.RGBA{0, 0, 0, 0xff},
		}
		smctx.AddMarker(marker)
	}

	// Render
	img, err := smctx.Render()
	if err != nil {
		return nil, errors.Wrap(err, "DrawPlaces: error while rendering image")
	}

	return img, nil
}

// DrawRegions draws multiple regions
func DrawRegions(regions []types.Region) (image.Image, error) {
	smctx := sm.NewContext()
	smctx.SetTileProvider(sm.NewTileProviderCartoLight())
	smctx.SetSize(420, 420)
	for _, r := range regions {
		area := &sm.Area{}
		for _, l1 := range r.Shape.Coords() {
			for _, l2 := range l1 {
				for _, l3 := range l2 {
					area.Positions = append(area.Positions, s2.LatLngFromDegrees(l3[0], l3[0]))
				}
			}
		}
		smctx.AddArea(area)
	}

	// Render
	img, err := smctx.Render()
	if err != nil {
		return nil, errors.Wrap(err, "DrawRegions: error while rendering image")
	}

	return img, nil
}

func posContainer(cont types.Container) (s2.LatLng, error) {
	obj, err := cont.Object()
	if err != nil {
		return s2.LatLng{}, err
	}

	// Switch between types
	var coord types.Coordinates
	switch t := obj.(type) {
	case *types.Address:
		coord = t.Coord
	case *types.StopArea:
		coord = t.Coord
	case *types.StopPoint:
		coord = t.Coord
	case *types.Admin:
		coord = t.Coord
	case *types.POI:
		err = errors.Errorf("POI has no latitude/longitude")
	default:
		err = errors.Errorf("Cannot find a matching type for container marked with embedded type: %s", cont.EmbeddedType)
	}
	if err != nil {
		return s2.LatLng{}, err
	}
	pos := s2.LatLngFromDegrees(coord.Latitude, coord.Longitude)
	return pos, nil
}
