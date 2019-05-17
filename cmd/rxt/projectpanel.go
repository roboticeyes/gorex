package main

import (
	"fmt"
	"strconv"

	"github.com/breiting/tview"
	"github.com/gdamore/tcell"
	"github.com/roboticeyes/gorex/http/rexos/listing"
)

var (
	tableHeader []string
	alignment   []int
)

func init() {

	tableHeader = []string{
		"Urn",
		"Name",
		"#Files",
		"Size (kb)",
		"Public",
	}
	alignment = []int{
		tview.AlignCenter,
		tview.AlignLeft,
		tview.AlignRight,
		tview.AlignRight,
		tview.AlignLeft,
	}
}

// ProjectPanel shows a table of rexOS projects
type ProjectPanel struct {
	*tview.Table
	controller *ViewController
}

// NewProjectPanel creates a new UI component
func NewProjectPanel(c *ViewController) *ProjectPanel {
	p := &ProjectPanel{
		Table:      tview.NewTable().SetFixed(1, 1).SetSelectable(true, false),
		controller: c,
	}
	p.SetTitle("Projects")
	p.SetBorder(true)

	p.Table.SetSelectedFunc(func(row, column int) {

		fmt.Println("selected", row, column)
	})

	p.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == 'n' {
			// modal := tview.NewModal().
			// 	SetText("Do you want to quit the application?").
			// 	AddButtons([]string{"Quit", "Cancel"}).
			// 	SetDoneFunc(func(buttonIndex int, buttonLabel string) {
			// 		if buttonLabel == "Quit" {
			// 		}
			// 	})
		}
		return event
	})

	p.SetProjects("", []listing.Project{})
	return p
}

// // InputHandler returns the handler of the primitive
// func (pp *ProjectPanel) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
// 	return pp.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
// 		fmt.Println(event)
//
// 	})
// }

func (pp *ProjectPanel) name() string {
	return "ProjectPanel"
}

func (pp *ProjectPanel) key() tcell.Key {
	return tcell.KeyF1
}

func (pp *ProjectPanel) content() tview.Primitive {
	return pp
}

func (pp *ProjectPanel) update() error {
	p, err := pp.controller.GetProjects()
	if err != nil {
		return err
	}
	pp.SetProjects(pp.controller.GetUserID(), p)
	return nil
}

// SetProjects sets the projects for this view. The owner is used to color shared project differently
func (pp *ProjectPanel) SetProjects(owner string, projects []listing.Project) {

	for i, h := range tableHeader {

		pp.SetCell(0, i, &tview.TableCell{
			Text:          h,
			Color:         tcell.ColorAzure,
			Align:         alignment[i],
			NotSelectable: true,
			Expansion:     1,
		})
	}

	for row, p := range projects {

		var color tcell.Color = tcell.ColorYellow
		if p.Owner != owner {
			color = tcell.ColorGray
		}

		pp.SetCell(row+1, 0, &tview.TableCell{
			Text:          p.Urn,
			Color:         color,
			Align:         alignment[0],
			NotSelectable: false,
		})

		pp.SetCell(row+1, 1, &tview.TableCell{
			Text:          p.Name,
			Color:         color,
			Align:         alignment[1],
			NotSelectable: false,
			Expansion:     1,
		})

		pp.SetCell(row+1, 2, &tview.TableCell{
			Text:          strconv.Itoa(p.NumberOfProjectFiles),
			Color:         color,
			Align:         alignment[2],
			NotSelectable: false,
		})
		pp.SetCell(row+1, 3, &tview.TableCell{
			Text:          strconv.Itoa(p.TotalProjectFileSize / 1024),
			Color:         color,
			Align:         alignment[3],
			NotSelectable: false,
		})
		pp.SetCell(row+1, 4, &tview.TableCell{
			Text:          strconv.FormatBool(p.Public),
			Color:         color,
			Align:         alignment[4],
			NotSelectable: false,
		})

		pp.ScrollToBeginning()
	}
}
