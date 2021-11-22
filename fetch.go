package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sdwalsh/mbta-v3-go/mbta"
)

// fetchRoutes returns the raw data structure from the mbta package
func fetchRoutes(api mbta.Client) tea.Msg {
	requestConfig := &mbta.GetAllRoutesRequestConfig{
		FilterRouteTypes: []mbta.RouteType{mbta.RouteTypeLightRail, mbta.RouteTypeHeavyRail},
	}
	routes, _, err := api.Routes.GetAllRoutes(requestConfig)
	if err != nil {
		return err
	}

	return routes
}

// fetchStops returns the raw data structure from the mbta package
func fetchStops(api mbta.Client, route string) tea.Msg {
	requestConfig := &mbta.GetAllStopsRequestConfig{
		FilterRouteIDs: []string{route},
	}
	stops, _, err := api.Stops.GetAllStops(requestConfig)
	if err != nil {
		return err
	}

	return stops
}

// fetchTimetable returns the raw data structure from the mbta package
func fetchTimetable(api mbta.Client, stop string, route string, direction string) tea.Msg {
	requestConfig := &mbta.GetAllPredictionsRequestConfig{
		FilterStopIDs:     []string{stop},
		FilterRouteIDs:    []string{route},
		FilterDirectionID: direction,
		Sort: mbta.PredictionsSortByDepartureTimeAscending,
	}
	timetable, _, err := api.Predictions.GetAllPredictions(requestConfig)
	if err != nil {
		return err
	}

	return timetable
}

// processRoutes converts the raw mbta structure to a slice of list.Item in order to populate the terminal list
func processRoutes(m []*mbta.Route) []list.Item {
	var routeItems []list.Item
	for i, r := range m {
		routeItems = append(routeItems, item{id: r.ID, index: i, title: r.LongName,
			desc: fmt.Sprintf("%s <--> %s", r.DirectionDestinations[1], r.DirectionDestinations[0])})
	}

	return routeItems
}

// processStops converts the raw mbta structure to a slice of list.Item in order to populate the terminal list
func processStops(stops []*mbta.Stop) []list.Item {
	var stopItems []list.Item
	for _, s := range stops {
		// We won't display a description unless the mbta provides one (otherwise we dereference nil pointer / panic)
		d := ""
		if s.Address != nil {
			d = *s.Address
		}
		stopItems = append(stopItems, item{id: s.ID, title: s.Name, desc: d})
	}

	return stopItems
}
