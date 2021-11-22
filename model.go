package main

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/sdwalsh/mbta-v3-go/mbta"
	"strconv"
)

// item is the data structure used to display data to the terminal
type item struct {
	id    string
	index int
	title string
	desc  string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// Select Menu -> Route -> Select Stop -> Select Direction -> Display Timetable
type step int

// TODO We aren't using menu yet, but I would really like to load the user into a neutral space
const (
	menuStep step = iota
	routeStep
	stopStep
	directionStep
	timetableStep
)

type selectedRoute struct {
	index int
	name  string
}

// model is our application state
type model struct {
	Client            mbta.Client
	Routes            []*mbta.Route
	Stops             []*mbta.Stop
	Timetable         []*mbta.Prediction
	routeList         list.Model
	stopList          list.Model
	directionList     list.Model
	timetableList     list.Model
	step              step
	loading           bool // TODO
	selectedRoute     selectedRoute
	selectedStop      string
	selectedDirection string
	err               error
}

// Init is automatically called by BubbleTea. Start off by calling the first fetch (async)
func (m model) Init() tea.Cmd {
	return func() tea.Msg {
		return fetchRoutes(m.Client)
	}
}

func (m *model) processDirections() []list.Item {
	var directionItems []list.Item
	// TODO brittle code
	for i, d := range m.Routes[m.selectedRoute.index].DirectionNames {
		directionItems = append(directionItems, item{index: i, id: strconv.Itoa(i), title: d,
			desc: m.Routes[m.selectedRoute.index].DirectionDestinations[i]})
	}

	return directionItems
}

func (m *model) processTimetable() []list.Item {
	var timetableItems []list.Item
	for _, t := range m.Timetable {
		// DepartureTime can be empty if we're at the end of the line and there are no more stops on the run
		// we're just going to handle that by dropping them for the time being. Future: Maybe check and notify user?
		desc := ""
		if t.Status != nil {
			desc = *t.Status
		}
		if t.DepartureTime != "" {
			timetableItems = append(timetableItems, item{title: t.DepartureTime, desc: desc})
		}
	}

	return timetableItems
}