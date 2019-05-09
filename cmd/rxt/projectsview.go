package main

import (
	// "strconv"

	"github.com/breiting/tview"
	"github.com/gdamore/tcell"
	"github.com/roboticeyes/gorex/http/rexos"
)

var (
	tableHeader []string
)

func init() {

	tableHeader = []string{"ID", "Name", "Owner"}
}

// ProjectsView shows a table of rexOS projects
type ProjectsView struct {
	*tview.Table
}

// NewProjectView creates a new UI component
func NewProjectView() *ProjectsView {
	p := &ProjectsView{
		Table: tview.NewTable().SetFixed(1, 1).SetSelectable(true, false),
	}
	p.SetTitle("Projects")
	p.SetBorder(true)

	var tp []rexos.ProjectSimple
	// for i := 0; i < 20; i++ {
	// 	tp = append(tp, rexos.ProjectSimple{
	// 		ID:    "000-0000-" + strconv.Itoa(i),
	// 		Name:  "TestProject",
	// 		Owner: "me",
	// 	})
	// }
	p.SetProjects(tp)
	return p
}

// SetProjects sets the projects for this view
func (v *ProjectsView) SetProjects(projects []rexos.ProjectSimple) {

	for i, h := range tableHeader {

		v.SetCell(0, i, &tview.TableCell{
			Text:          h,
			Color:         tcell.ColorYellow,
			Align:         tview.AlignCenter,
			NotSelectable: true,
			Expansion:     1,
		})
	}

	for row, p := range projects {

		v.SetCell(row+1, 0, &tview.TableCell{
			Text:          p.ID,
			Color:         tcell.ColorWhite,
			Align:         tview.AlignCenter,
			NotSelectable: false,
		})

		v.SetCell(row+1, 1, &tview.TableCell{
			Text:          p.Name,
			Color:         tcell.ColorWhite,
			Align:         tview.AlignLeft,
			NotSelectable: false,
			Expansion:     1,
		})

		v.SetCell(row+1, 2, &tview.TableCell{
			Text:          p.Owner,
			Color:         tcell.ColorWhite,
			Align:         tview.AlignCenter,
			NotSelectable: false,
			// Expansion:     1,
		})
	}
}
