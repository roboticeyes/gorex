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
	c.SetConnected(false)
	return c
}

// SetConnected can be called to change the state
func (c *ConnectionStatus) SetConnected(status bool) {
	if status == false {
		c.SetText("[red]not connected")
	} else {
		c.SetText("[green]connected")
	}
}

// SetError sets the error message
func (c *ConnectionStatus) SetError(e error) {
	c.SetText("[white:red]" + e.Error())
}
