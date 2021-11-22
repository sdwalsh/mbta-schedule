# MBTA Scheduler
...in your terminal!

## Design
### ELM / React-style Architecture
#### Model <> Update <> View

## Libraries Used
- Built using TUI libraries
  - charmbracelet/bubbletea
  - charmbracelet/bubbles
  - charmbracelet/lipgloss
- MBTA API Library / Typings
  - (currently a forked version on my account, needs minor fixes)

## Next steps
- Testing is currently missing. Need to add automated tests.
- The update loop is crowded. Refactor?
- Forked the MBTA API library, issues with marshaling Time fields
- Custom styling (lipgloss library)

## Compile and Run
- Built with: go1.17.3 linux/amd64 on Ubuntu 21.10
- `go run main.go` or `go build` and `./mbta`