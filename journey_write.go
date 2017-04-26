package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/aabizri/navitia"
	"github.com/aabizri/navitia/types"
	"github.com/fatih/color"
)

func journeyResultsWrite(jr *navitia.JourneyResults, out io.Writer) {
	// Buffers to line-up the reads, sequentially
	buffers := make([]io.Reader, jr.Count())
	// Waitgroup for each goroutine
	wg := sync.WaitGroup{}

	// Iterate through the journeys, printing them
	for i, j := range jr.Journeys {
		buf := &bytes.Buffer{}
		buffers[i] = buf

		// Increment the WaitGroup
		wg.Add(1)

		// Launch !
		go func(j types.Journey) {
			defer wg.Done()
			journeyWrite(&j, buf)
		}(j)
	}

	// Create the reader
	reader := io.MultiReader(buffers...)

	// Wait for completion
	wg.Wait()

	// Copy the new reader to the given output
	_, err := io.Copy(out, reader)
	if err != nil {
		panic(0)
	}
}

const timeLayout = "15:04"

func journeyWrite(j *types.Journey, out io.Writer) {
	// Build the envellope
	msg := fmt.Sprintf("%s ➡️ %s | %s\n", color.RedString(j.Departure.Format(timeLayout)), color.RedString(j.Arrival.Format(timeLayout)), color.MagentaString(j.Duration.String()))
	// Buffers to line-up the reads, sequentially
	buffers := make([]io.Reader, len(j.Sections))
	// Waitgroup for each goroutine
	wg := sync.WaitGroup{}

	// Iterate through the journeys, printing them
	for i, s := range j.Sections {
		buf := &bytes.Buffer{}
		buffers[i] = buf

		// Increment the WaitGroup
		wg.Add(1)

		// Launch !
		go func(s types.Section) {
			defer wg.Done()
			sectionWrite(&s, buf)
		}(s)
	}

	// Create the reader
	readers := append([]io.Reader{strings.NewReader(msg)}, buffers...)
	reader := io.MultiReader(readers...)

	// Wait for completion
	wg.Wait()

	// Copy the new reader to the given output
	_, err := io.Copy(out, reader)
	if err != nil {
		panic(0)
	}
}

var modeEmoji = map[string]string{
	string(types.PhysicalModeAir):               "✈️",
	string(types.PhysicalModeBoat):              "⛴️",
	string(types.PhysicalModeBus):               "🚍",
	string(types.PhysicalModeBusRapidTransit):   "🚍",
	string(types.PhysicalModeCoach):             "🚍",
	string(types.PhysicalModeFerry):             "⛴️",
	string(types.PhysicalModeFunicular):         "🚞",
	string(types.PhysicalModeLocalTrain):        "🚆",
	string(types.PhysicalModeLongDistanceTrain): "🚆",
	string(types.PhysicalModeMetro):             "🚇",
	string(types.PhysicalModeRapidTransit):      "🚍",
	string(types.PhysicalModeShuttle):           "🚐",
	string(types.PhysicalModeTaxi):              "🚖",
	string(types.PhysicalModeTrain):             "🚆",
	string(types.PhysicalModeTramway):           "🚊",

	// Because the API doesn't always return predictable returns, we have aliases
	"Métro": "🚇",
	"Bus":   "🚍",

	// Classic Modes: Walking, biking or bikesharing
	string(types.ModeWalking):   "🚶",
	string(types.ModeBike):      "🚴",
	string(types.ModeBikeShare): "🚴",
}

func sectionWrite(s *types.Section, out io.Writer) {
	// if there's no from or no to, finish now
	if s.From.Name == "" || s.To.Name == "" {
		return
	}

	var middle string
	switch {
	case s.Mode != "":
		middle = modeEmoji[string(s.Mode)]
	case s.Display.PhysicalMode != "":
		middle = modeEmoji[string(s.Display.PhysicalMode)] + s.Display.Label
	}
	msg := fmt.Sprintf("\t%s (%s)\t%s➡️%s\n", color.GreenString(middle), color.MagentaString(s.Duration.String()), color.BlueString(s.From.Name), color.BlueString(s.To.Name))

	_, err := out.Write([]byte(msg))
	if err != nil {
		panic(err)
	}
}
