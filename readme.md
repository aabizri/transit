# transit is a public transit planner and information tool.
[![Build Status](https://travis-ci.org/aabizri/transit.svg?branch=dev)](https://travis-ci.org/aabizri/transit)
[![Go Report Card](https://goreportcard.com/badge/github.com/aabizri/transit)](https://goreportcard.com/report/github.com/aabizri/transit)
[![Coverage](https://gocover.io/_badge/github.com/aabizri/transit?0)](http://gocover.io/github.com/aabizri/transit)

## Install
Run `go get -u transit`

## Use

### Search places (Works !)
For finding places (not public transport places) use: `transit place "rue de vanves" "avenue de gaulle" "rue du caire"`

### Plan a journey (WIP but has rudimentary working !)
To plan a journey from one point to the next use: `transit journey from "11 Avenue Victor Hugo, Paris" to "Notre-dame-des-champs"`

Options:
- `-d`, `-duration`: Limit the maximum time in transit
- `-m`, `-mode`: Limit the transit to a single mode

### See the departures or arrivals (WIP)
Use `transit departures "Etoile"` for departures, `...arrivals "Etoile"` for arrivals.

### Isochrones (WIP)
Use `transit isochrone from "Etoile"` for isochrone from that place, `...to etoile` for an isochrone towards that place.


