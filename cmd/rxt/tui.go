package main

import (
	// "fmt"
	// "strings"

	"github.com/breiting/tview"
	"github.com/gdamore/tcell"
	"github.com/roboticeyes/gorex/http/rexos"
)

// UIRunner wrapps the function to run the UI
type UIRunner interface {
	Run() error
}

// ViewModel is the main user interface
type ViewModel struct {
	app          *tview.Application
	cmd          *tview.InputField
	status       *ConnectionStatus
	projectsView *ProjectsView
	statusBar    *tview.TextView
	root         *tview.Flex

	controller Controller
}

// Controller is an interface for providing business logic
type Controller interface {
	Connect() (string, error)
	GetConfiguration() *Configuration
	GetAllProjects() ([]rexos.ProjectSimple, error)
}

// NewTui creates a new TUI
func NewTui(c Controller) UIRunner {

	view := ViewModel{
		controller: c,
	}

	view.status = NewConnectionStatus()
	view.projectsView = NewProjectView()

	titleBar := tview.NewGrid().
		SetRows(1).
		SetColumns(0, -1, 0).
		AddItem(tview.NewTextView().SetDynamicColors(true).SetText("[yellow]rexOS terminal"), 0, 0, 1, 1, 0, 0, false).
		AddItem(tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter).SetText("[blue]"+c.GetConfiguration().Default), 0, 1, 1, 1, 0, 0, false).
		AddItem(view.status, 0, 2, 1, 1, 0, 0, false)

	// status bar
	status := tview.NewTextView()
	status.SetText("F1 Connect")

	// main area use pages
	main := tview.NewPages()
	main.AddPage("projects", view.projectsView, true, true)
	main.AddPage("leanbim", tview.NewBox().SetBorder(true).SetTitle("LeanBIM"), true, false)

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
		} else if event.Key() == tcell.KeyF2 {
			p, err := view.controller.GetAllProjects()
			if err != nil {
				view.status.SetConnected(false, err.Error())
			}
			view.projectsView.SetProjects(p)
			main.SwitchToPage("projects")
		} else if event.Key() == tcell.KeyF3 {
			main.SwitchToPage("leanbim")
		}
		return event
	})

	return &view
}

// Run starts the user interface
func (v *ViewModel) Run() error {
	return v.app.Run()
}
