package main

import tea "github.com/charmbracelet/bubbletea"

func (m model) fetchStopsCmd() tea.Cmd {
	return func() tea.Msg {
		return fetchStops(m.Client, m.selectedRoute.name)
	}
}

func (m model) fetchTimetableCmd() tea.Cmd {
	return func() tea.Msg {
		return fetchTimetable(m.Client, m.selectedStop, m.selectedRoute.name, m.selectedDirection)
	}
}
