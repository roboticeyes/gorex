package main

import (
	"github.com/breiting/tview"
	"github.com/gdamore/tcell"
)

// UIRunner wrapps the function to run the UI
type UIRunner interface {
	Run() error
}

// ViewModel is the main user interface
type ViewModel struct {
	app           *tview.Application
	cmd           *tview.InputField
	status        *ConnectionStatus
	statusBar     *tview.TextView
	listNotebooks *tview.List
	listNotes     *tview.List
	preview       *tview.TextView
	root          *tview.Flex

	controller Controller
}

// Controller is an interface for providing business logic
type Controller interface {
	Connect() (string, error)
	GetConfiguration() *Configuration
}

// NewTui creates a new TUI
func NewTui(c Controller) UIRunner {

	view := ViewModel{
		controller: c,
	}

	// title bar
	view.status = NewConnectionStatus()
	titleBar := tview.NewGrid().
		SetRows(1).
		SetColumns(0, -1, 0).
		AddItem(tview.NewTextView().SetDynamicColors(true).SetText("[yellow]rexOS terminal"), 0, 0, 1, 1, 0, 0, false).
		AddItem(tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter).SetText("[blue]"+c.GetConfiguration().Default), 0, 1, 1, 1, 0, 0, false).
		AddItem(view.status, 0, 2, 1, 1, 0, 0, false)

	// main area
	main := tview.NewBox().SetBorder(true)

	// status bar
	status := tview.NewTextView()
	status.SetText("F1 Connect")

	view.root = tview.NewFlex().SetDirection(tview.FlexRow)
	view.root.AddItem(titleBar, 1, 0, false)
	view.root.AddItem(main, 0, 1, true)
	view.root.AddItem(status, 1, 0, false)

	view.app = tview.NewApplication().SetRoot(view.root, true)

	view.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc || event.Rune() == 'q' {
			view.app.Stop()
		} else if event.Key() == tcell.KeyF1 {
			if username, err := view.controller.Connect(); err == nil {
				view.status.SetConnected(true, username)
			} else {
				view.status.SetConnected(false, err.Error())
			}
		}
		return event
	})

	return &view
}

// Run starts the user interface
func (v *ViewModel) Run() error {
	return v.app.Run()
}
