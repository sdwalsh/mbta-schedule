package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sdwalsh/mbta-v3-go/mbta"
	"os"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

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

// Update is called to update the Tea model
// https://github.com/charmbracelet/bubbletea/tree/master/tutorials/commands/
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case error:
		m.err = msg
	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		m.routeList.SetSize(msg.Width-left-right, msg.Height-top-bottom)
		m.stopList.SetSize(msg.Width-left-right, msg.Height-top-bottom)
		m.directionList.SetSize(msg.Width-left-right, msg.Height-top-bottom)
		m.timetableList.SetSize(msg.Width-left-right, msg.Height-top-bottom)

	// Callback from Route API Request
	case []*mbta.Route:
		m.Routes = msg
		m.routeList.SetItems(processRoutes(msg))
	// Callback from Stop API Request
	case []*mbta.Stop:
		m.Stops = msg
		m.stopList.SetItems(processStops(msg))
		// Directions are provided by Routes. When we have a selected Route we can populate the direction list
		m.directionList.SetItems(m.processDirections())
		m.stopList.Title = "Select a Stop"
	case []*mbta.Prediction:
		m.Timetable = msg
		m.timetableList.SetItems(m.processTimetable())

	// Handle user input
	case tea.KeyMsg:
		// Don't match any of the keys below if we're actively filtering.
		if m.routeList.FilterState() == list.Filtering ||
			m.stopList.FilterState() == list.Filtering ||
			m.directionList.FilterState() == list.Filtering ||
			m.timetableList.FilterState() == list.Filtering {
			break
		}
		switch keypress := msg.String(); keypress {
		// Listen for exit command
		case "ctrl+c":
			return m, tea.Quit
		// Selecting
		case "enter":
			switch m.step {
			case menuStep:
				// TODO add main menu
			case routeStep:
				i, ok := m.routeList.SelectedItem().(item)
				sr := selectedRoute{index: i.index, name: i.id}
				if ok {
					m.selectedRoute = sr
				}
				m.step = stopStep
				return m, m.fetchStopsCmd()
			case stopStep:
				i, ok := m.stopList.SelectedItem().(item)
				if ok {
					m.selectedStop = i.id
				}
				m.step = directionStep
				return m, nil // we don't need to send a command here since direction are already present in our data
			case directionStep:
				i, ok := m.directionList.SelectedItem().(item)
				if ok {
					m.selectedDirection = i.id
				}
				m.step = timetableStep
				return m, m.fetchTimetableCmd()
			case timetableStep:
				// TODO add focus or option to restart
			}
			return m, nil
		}
	}
	// Handle Bubbles: TODO cleanup / look to simplify
	var cmd tea.Cmd
	switch m.step {
	case routeStep:
		m.routeList, cmd = m.routeList.Update(msg)
	case stopStep:
		m.stopList, cmd = m.stopList.Update(msg)
	case directionStep:
		m.directionList, cmd = m.directionList.Update(msg)
	case timetableStep:
		m.timetableList, cmd = m.timetableList.Update(msg)
	}
	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}

	// Render based on step
	switch m.step {
	case routeStep:
		return docStyle.Render(m.routeList.View())
	case stopStep:
		return docStyle.Render(m.stopList.View())
	case directionStep:
		return docStyle.Render(m.directionList.View())
	case timetableStep:
		return docStyle.Render(m.timetableList.View())
	}

	// Hopefully we don't reach this
	return docStyle.Render("Nothing to display!")
}

func main() {
	m := model{}
	m.step = routeStep
	// TODO import via env (seems to work w/o an API key, maybe limited?)
	m.Client = *mbta.NewClient(mbta.ClientConfig{APIKey: ""})

	// init lists
	m.routeList = list.NewModel([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	m.stopList = list.NewModel([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	m.directionList = list.NewModel([]list.Item{}, list.NewDefaultDelegate(), 0, 0)
	m.timetableList = list.NewModel([]list.Item{}, list.NewDefaultDelegate(), 0, 0)

	m.routeList.Title = "Select a Route"
	m.stopList.Title = "Select a Stop"
	m.directionList.Title = "Select a Direction"
	m.timetableList.Title = "Timetable"

	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
