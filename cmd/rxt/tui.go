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
	app          *tview.Application
	cmd          *tview.InputField
	status       *ConnectionStatus
	projectsView *ProjectsView
	statusBar    *tview.TextView
	root         *tview.Flex
	main         *tview.Pages
	controller   *ViewController
}

// NewTui creates a new TUI
func NewTui(c *ViewController) UIRunner {

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
	status.SetText("F1 Projects")

	// main area use pages
	view.main = tview.NewPages()
	view.main.AddPage("projects", view.projectsView, true, true)
	view.main.AddPage("leanbim", tview.NewBox().SetBorder(true).SetTitle("LeanBIM"), true, false)

	view.root = tview.NewFlex().SetDirection(tview.FlexRow)
	view.root.AddItem(titleBar, 1, 0, false)
	view.root.AddItem(view.main, 0, 1, true)

	view.root.AddItem(status, 1, 0, false)

	view.app = tview.NewApplication().SetRoot(view.root, true)

	view.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc || event.Rune() == 'q' {
			view.app.Stop()
		} else if event.Key() == tcell.KeyF1 {
			view.projects()
		} else if event.Key() == tcell.KeyF2 {
			view.main.SwitchToPage("leanbim")
		}
		return event
	})

	// Auto-connect
	view.status.SetStatus(Info, "connecting")
	go view.connect()
	return &view
}

// Run starts the user interface
func (v *ViewModel) Run() error {
	return v.app.Run()
}

func (v *ViewModel) connect() {
	if username, err := v.controller.Connect(); err == nil {
		v.app.QueueUpdateDraw(func() {
			v.status.SetStatus(Success, username)
			v.projects()
		})
	} else {
		v.app.QueueUpdateDraw(func() {
			v.status.SetStatus(Error, err.Error())
		})
	}
}

func (v *ViewModel) projects() {
	p, err := v.controller.GetProjects()
	if err != nil {
		v.status.SetStatus(Error, err.Error())
	}
	v.projectsView.SetProjects(v.controller.GetUserID(), p)
	v.main.SwitchToPage("projects")
}
