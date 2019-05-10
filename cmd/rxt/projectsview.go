package main

import (
	"strconv"

	"github.com/breiting/tview"
	"github.com/gdamore/tcell"
	"github.com/roboticeyes/gorex/http/rexos"
)

var (
	tableHeader []string
	alignment   []int
)

func init() {

	tableHeader = []string{
		"ID",
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

	var tp rexos.ProjectComplexList
	p.SetProjects("", tp)
	return p
}

// SetProjects sets the projects for this view
func (v *ProjectsView) SetProjects(owner string, projectList rexos.ProjectComplexList) {

	for i, h := range tableHeader {

		v.SetCell(0, i, &tview.TableCell{
			Text:          h,
			Color:         tcell.ColorAzure,
			Align:         alignment[i],
			NotSelectable: true,
			Expansion:     1,
		})
	}

	for row, p := range projectList.Embedded.Projects {

		var color tcell.Color = tcell.ColorYellow
		if p.Owner != owner {
			color = tcell.ColorGray
		}

		v.SetCell(row+1, 0, &tview.TableCell{
			Text:          p.ID,
			Color:         color,
			Align:         alignment[0],
			NotSelectable: false,
		})

		v.SetCell(row+1, 1, &tview.TableCell{
			Text:          p.Name,
			Color:         color,
			Align:         alignment[1],
			NotSelectable: false,
			Expansion:     1,
		})

		v.SetCell(row+1, 2, &tview.TableCell{
			Text:          strconv.Itoa(p.NumberOfProjectFiles),
			Color:         color,
			Align:         alignment[2],
			NotSelectable: false,
		})
		v.SetCell(row+1, 3, &tview.TableCell{
			Text:          strconv.Itoa(p.TotalProjectFileSize / 1024),
			Color:         color,
			Align:         alignment[3],
			NotSelectable: false,
		})
		v.SetCell(row+1, 4, &tview.TableCell{
			Text:          strconv.FormatBool(p.Public),
			Color:         color,
			Align:         alignment[4],
			NotSelectable: false,
		})

		v.ScrollToBeginning()
	}
}