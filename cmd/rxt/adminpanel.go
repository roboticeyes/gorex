package main

import (
	"github.com/breiting/tview"
	"github.com/gdamore/tcell"
)

// AdminPanel shows some admin tools
type AdminPanel struct {
	*tview.Box
}

// NewAdminPanel creates a admin panel
func NewAdminPanel() *AdminPanel {
	p := &AdminPanel{
		Box: tview.NewBox(),
	}
	p.SetTitle("Admin")
	p.SetBorder(true)

	return p
}

func (p *AdminPanel) name() string {
	return "AdminPanel"
}

func (p *AdminPanel) key() tcell.Key {
	return tcell.KeyF2
}

func (p *AdminPanel) content() tview.Primitive {
	return p
}

func (p *AdminPanel) update() error {
	return nil
}
