package main

import (
	"github.com/breiting/tview"
	"github.com/gdamore/tcell"
)

// Panel is a simple interface for updating the content when switching to the panel
type Panel interface {
	name() string
	key() tcell.Key
	content() tview.Primitive
	update() error
}
