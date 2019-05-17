package main

import (
	"strconv"

	"github.com/breiting/tview"
	"github.com/gdamore/tcell"
	"github.com/roboticeyes/gorex/http/rexos/listing"
)

// ProjectsPanel shows a table of rexOS projects
type ProjectsPanel struct {
	*tview.Pages
	projectTable *tview.Table
	controller   *ViewController
	tableHeader  []string
	alignment    []int
}

// NewProjectsPanel creates a new UI component
func NewProjectsPanel(c *ViewController) *ProjectsPanel {

	p := &ProjectsPanel{
		Pages: tview.NewPages(),
	}
	p.projectTable = tview.NewTable().SetFixed(1, 1).SetSelectable(true, false)
	p.controller = c

	p.tableHeader = []string{
		"Urn",
		"Name",
		"#Files",
		"Size (kb)",
		"Public",
	}
	p.alignment = []int{
		tview.AlignCenter,
		tview.AlignLeft,
		tview.AlignRight,
		tview.AlignRight,
		tview.AlignLeft,
	}

	p.AddPage("projects", p.projectTable, true, true)
	p.AddPage("newproject", tview.NewForm().
		AddDropDown("Title", []string{"Mr.", "Ms.", "Mrs.", "Dr.", "Prof."}, 0, nil).
		AddInputField("First name", "", 20, nil, nil).
		AddInputField("Last name", "", 20, nil, nil).
		AddCheckbox("Age 18+", false, nil).
		AddPasswordField("Password", "", 10, '*', nil).
		AddButton("Save", nil).
		AddButton("Quit", func() {
		}), true, false)

	p.SetTitle("Project")
	p.SetBorder(true)

	p.projectTable.SetSelectedFunc(func(row, column int) {
		// p.GetProjectFiles(p.controller.GetProjects()[0])
		// fmt.Println("selected", row, column)
	})

	p.projectTable.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyEsc {
			p.SwitchToPage("projects")
		}
		if event.Rune() == 'n' {
			p.SwitchToPage("newproject")
		}
		return event
	})

	p.SetProjects("", []listing.Project{})
	return p
}

// // InputHandler returns the handler of the primitive
// func (pp *ProjectsPanel) InputHandler() func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
// 	return pp.WrapInputHandler(func(event *tcell.EventKey, setFocus func(p tview.Primitive)) {
// 		fmt.Println(event)
//
// 	})
// }

func (pp *ProjectsPanel) name() string {
	return "Project"
}

func (pp *ProjectsPanel) key() tcell.Key {
	return KeySwitchToProject
}

func (pp *ProjectsPanel) content() tview.Primitive {
	return pp
}

func (pp *ProjectsPanel) update() error {
	p, err := pp.controller.GetProjects()
	if err != nil {
		return err
	}
	pp.SetProjects(pp.controller.GetUserID(), p)
	return nil
}

// SetProjects sets the projects for this view. The owner is used to color shared project differently
func (pp *ProjectsPanel) SetProjects(owner string, projects []listing.Project) {

	for i, h := range pp.tableHeader {

		pp.projectTable.SetCell(0, i, &tview.TableCell{
			Text:          h,
			Color:         tcell.ColorAzure,
			Align:         pp.alignment[i],
			NotSelectable: true,
			Expansion:     1,
		})
	}

	for row, p := range projects {

		var color tcell.Color = tcell.ColorYellow
		if p.Owner != owner {
			color = tcell.ColorGray
		}

		pp.projectTable.SetCell(row+1, 0, &tview.TableCell{
			Text:          p.Urn,
			Color:         color,
			Align:         pp.alignment[0],
			NotSelectable: false,
		})

		pp.projectTable.SetCell(row+1, 1, &tview.TableCell{
			Text:          p.Name,
			Color:         color,
			Align:         pp.alignment[1],
			NotSelectable: false,
			Expansion:     1,
		})

		pp.projectTable.SetCell(row+1, 2, &tview.TableCell{
			Text:          strconv.Itoa(p.NumberOfProjectFiles),
			Color:         color,
			Align:         pp.alignment[2],
			NotSelectable: false,
		})
		pp.projectTable.SetCell(row+1, 3, &tview.TableCell{
			Text:          strconv.Itoa(p.TotalProjectFileSize / 1024),
			Color:         color,
			Align:         pp.alignment[3],
			NotSelectable: false,
		})
		pp.projectTable.SetCell(row+1, 4, &tview.TableCell{
			Text:          strconv.FormatBool(p.Public),
			Color:         color,
			Align:         pp.alignment[4],
			NotSelectable: false,
		})

		pp.projectTable.ScrollToBeginning()
	}
}

func (pp *ProjectsPanel) GetProjectFiles() {
	pp.SwitchToPage("projectfile")
}
