package main

import (
	"github.com/breiting/tview"
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
	c.SetConnected(false, "not connected")
	return c
}

// SetConnected can be called to change the state
func (c *ConnectionStatus) SetConnected(status bool, msg string) {
	if status == false {
		c.SetText("[white:red]" + msg)
	} else {
		c.SetText("[green]" + msg)
	}
}
