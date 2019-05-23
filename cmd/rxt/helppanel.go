package main

import (
	"github.com/breiting/tview"
	"github.com/gdamore/tcell"
)

// HelpPanel shows the main help
type HelpPanel struct {
	*tview.TextView
}

// NewHelpPanel creates a admin panel
func NewHelpPanel() *HelpPanel {
	p := &HelpPanel{
		TextView: tview.NewTextView(),
	}
	p.SetTitle("Help")
	p.SetBorder(true)
	p.setDefaultText()

	return p
}

func (p *HelpPanel) name() string {
	return "Help"
}

func (p *HelpPanel) key() tcell.Key {
	return KeySwitchToHelp
}

func (p *HelpPanel) content() tview.Primitive {
	return p
}

func (p *HelpPanel) update() error {
	return nil
}

func (p *HelpPanel) setDefaultText() {
	p.SetText(`

	rxt provides the following key bindings:

	F1 | Help
	F2 | Project
	F3 | Admin

	ESC quits the application.
	`)
}
