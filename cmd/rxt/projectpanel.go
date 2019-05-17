package main

import (
	"strconv"

	"github.com/breiting/tview"
	"github.com/gdamore/tcell"
	"github.com/roboticeyes/gorex/http/rexos/listing"
)

// ProjectPanel shows a table of rexOS projects
type ProjectPanel struct {
	*tview.Table
	controller  *ViewController
	tableHeader []string
	alignment   []int
}

// NewProjectPanel creates a new UI component
func NewProjectPanel(c *ViewController) *ProjectPanel {

	p := &ProjectPanel{
		Table: tview.NewTable().SetFixed(1, 1).SetSelectable(true, false),
	}
	p.controller = c

	p.tableHeader = []string{
		"Name",
		"Type",
		"Size (kb)",
	}
	p.alignment = []int{
		tview.AlignLeft,
		tview.AlignCenter,
		tview.AlignRight,
	}

	p.SetTitle("Project")
	p.SetBorder(true)

	p.SetSelectedFunc(func(row, column int) {
		// p.GetProjectFiles()
		// p.SwitchToPage("newproject")
		// fmt.Println("selected", row, column)
	})

	p.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'n' {
			// p.SwitchToPage("newproject")
		}
		return event
	})

	p.SetProjectFiles([]listing.ProjectFile{})
	return p
}

func (pp *ProjectPanel) name() string {
	return "ProjectFile"
}

func (pp *ProjectPanel) key() tcell.Key {
	// user cannot switch to this pane with a key binding
	return KeySwitchToProjectFiles // TODO
}

func (pp *ProjectPanel) content() tview.Primitive {
	return pp
}

func (pp *ProjectPanel) update() error {
	p, err := pp.controller.GetProjectFiles()
	if err != nil {
		return err
	}
	pp.SetProjectFiles(p)
	return nil
}

// SetProjectFiles sets the project files to view
func (pp *ProjectPanel) SetProjectFiles(projects []listing.ProjectFile) {

	for i, h := range pp.tableHeader {

		pp.SetCell(0, i, &tview.TableCell{
			Text:          h,
			Color:         tcell.ColorAzure,
			Align:         pp.alignment[i],
			NotSelectable: true,
			Expansion:     1,
		})
	}

	for row, p := range projects {

		var color tcell.Color = tcell.ColorYellow
		pp.SetCell(row+1, 0, &tview.TableCell{
			Text:          p.Name,
			Color:         color,
			Align:         pp.alignment[0],
			NotSelectable: false,
		})

		pp.SetCell(row+1, 1, &tview.TableCell{
			Text:          p.Type,
			Color:         color,
			Align:         pp.alignment[1],
			NotSelectable: false,
			Expansion:     1,
		})

		pp.SetCell(row+1, 2, &tview.TableCell{
			Text:          strconv.Itoa(p.FileSize / 1024),
			Color:         color,
			Align:         pp.alignment[2],
			NotSelectable: false,
		})

		pp.ScrollToBeginning()
	}
}
