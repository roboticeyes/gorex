package main

import (
	"github.com/breiting/tview"
	"github.com/gdamore/tcell"
)

// UIRunner wrapps the function to run the UI
type UIRunner interface {
	Run() error
}

// App is the top level UI application
type App struct {
	app        *tview.Application
	cmd        *tview.InputField
	status     *ConnectionStatus
	main       *tview.Pages
	controller *ViewController
	panels     []Panel
}

// NewApp creates the application
func NewApp(c *ViewController) UIRunner {

	app := App{
		controller: c,
	}

	app.status = NewConnectionStatus()

	titleBar := tview.NewGrid().
		SetRows(1).
		SetColumns(0, -1, 0).
		AddItem(tview.NewTextView().SetDynamicColors(true).SetText("[yellow]rexOS terminal"), 0, 0, 1, 1, 0, 0, false).
		AddItem(tview.NewTextView().SetDynamicColors(true).SetTextAlign(tview.AlignCenter).SetText("[blue]"+c.GetConfiguration().Default), 0, 1, 1, 1, 0, 0, false).
		AddItem(app.status, 0, 2, 1, 1, 0, 0, false)

	// toolbar
	toolbar := tview.NewTextView()
	toolbar.SetText("F1 Projects")

	// create all panels
	app.createPanels()

	root := tview.NewFlex().SetDirection(tview.FlexRow)
	root.AddItem(titleBar, 1, 0, false)
	root.AddItem(app.main, 0, 1, true)
	root.AddItem(toolbar, 1, 0, false)
	app.app = tview.NewApplication().SetRoot(root, true)

	app.app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc || event.Rune() == 'q' {
			app.app.Stop()
		}
		for _, p := range app.panels {
			if event.Key() == p.key() {
				app.main.SwitchToPage(p.name())
				err := p.update()
				if err != nil {
					app.status.SetStatus(Error, err.Error())
				}

				return event
			}
		}
		return event
	})

	// Auto-connect
	app.status.SetStatus(Info, "connecting")
	go app.connect()
	return &app
}

// Run starts the user interface
func (a *App) Run() error {
	return a.app.Run()
}

func (a *App) connect() {
	if username, err := a.controller.Connect(); err == nil {
		a.app.QueueUpdateDraw(func() {
			a.status.SetStatus(Success, username)
		})
	} else {
		a.app.QueueUpdateDraw(func() {
			a.status.SetStatus(Error, err.Error())
		})
	}
}

func (a *App) createPanels() {
	a.main = tview.NewPages()

	a.panels = append(a.panels, NewProjectPanel(a.controller))
	a.panels = append(a.panels, NewAdminPanel())

	for i, p := range a.panels {
		visible := false
		if i == 0 {
			visible = true
		}
		a.main.AddPage(p.name(), p.content(), true, visible)
	}
}
