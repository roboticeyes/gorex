package main

import (
	"github.com/breiting/tview"
)

// StatusCode should be used to set the status
type StatusCode int

const (
	// Success everything good
	Success StatusCode = iota
	// Error everything bad
	Error StatusCode = iota
	// Info just for info
	Info StatusCode = iota
	// Warning raises a warning
	Warning StatusCode = iota
)

// ConnectionStatus is a UI component for handling rexOS connection
type ConnectionStatus struct {
	*tview.TextView
}

// NewConnectionStatus creates a new UI component
func NewConnectionStatus() *ConnectionStatus {
	c := &ConnectionStatus{
		TextView: tview.NewTextView(),
	}
	c.SetTextAlign(tview.AlignRight).
		SetDynamicColors(true)
	c.SetStatus(Info, "not connected")
	return c
}

// SetStatus can be called to change the state
func (c *ConnectionStatus) SetStatus(code StatusCode, msg string) {
	c.SetText(getColor(code) + msg)
}

func getColor(s StatusCode) string {

	colorCodes := []string{
		"[green]",
		"[white:red]",
		"[yellow]",
		"[white:orange]",
	}
	return colorCodes[s]
}
