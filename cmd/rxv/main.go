package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/roboticeyes/gorex/loader/rex"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/util/application"
	"github.com/g3n/engine/util/logger"
	"github.com/g3n/engine/window"
)

const (
	selectionText = "Selected object: "
)

var (
	app         *application.Application
	rc          *core.Raycaster
	rexGroup    *core.Node
	wireframe   bool
	picking     bool
	gridSize    float32
	headLight   *light.Directional
	pointLight  *light.Point
	scaleFactor float32

	labelSelectedObject *gui.Label
)

func init() {
	picking = false
	wireframe = false
	gridSize = 10
	scaleFactor = 1
}

func toggleWireframe() {
	wireframe = !wireframe
	for _, obj := range rexGroup.Children() {
		for _, n := range obj.GetNode().Children() {
			ig, ok := n.(graphic.IGraphic)
			if !ok {
				continue
			}
			gr := ig.GetGraphic()
			imat := gr.GetMaterial(0).GetMaterial()
			imat.SetWireframe(wireframe)
		}
	}
}

func onKey(evname string, ev interface{}) {

	kev := ev.(*window.KeyEvent)
	if kev.Action == window.Release {
		return
	}
	switch kev.Keycode {
	case window.KeyQ:
		app.Quit()
	case window.KeyC:
		bbox := app.Scene().BoundingBox()
		center := bbox.Center(nil)
		app.Scene().SetPosition(-center.X, 0, -center.Z)
	case window.KeyJ:
		scaleFactor *= 0.1
		rexGroup.SetScale(scaleFactor, scaleFactor, scaleFactor)
	case window.KeyK:
		scaleFactor /= 0.1
		rexGroup.SetScale(scaleFactor, scaleFactor, scaleFactor)
	case window.KeyW:
		toggleWireframe()
	case window.KeyP:
		picking = !picking
	}
}

func onMouse(ev interface{}) {
	if !picking {
		return
	}

	// Convert mouse coordinates to normalized device coordinates
	mev := ev.(*window.MouseEvent)
	width, height := app.Window().Size()
	x := 2*(mev.Xpos/float32(width)) - 1
	y := -2*(mev.Ypos/float32(height)) + 1

	// Set the raycaster from the current camera and mouse coordinates
	app.Camera().SetRaycaster(rc, x, y)

	// Checks intersection with all objects in the scene
	intersects := rc.IntersectObjects(rexGroup.Children(), true)
	// fmt.Printf("intersects:%+v\n", intersects)
	if len(intersects) == 0 {
		return
	}

	// Get first intersection
	obj := intersects[0].Object
	// Convert INode to IGraphic
	ig, ok := obj.(graphic.IGraphic)
	if !ok {
		app.Log().Debug("Not graphic:%T", obj)
		return
	}
	// app.Log().Debug("Selected node %s", obj.GetNode().Name())
	labelSelectedObject.SetText(selectionText + obj.GetNode().Name())

	// Get graphic object
	gr := ig.GetGraphic()
	imat := gr.GetMaterial(0)

	type matI interface {
		EmissiveColor() math32.Color
		SetEmissiveColor(*math32.Color)
	}

	if v, ok := imat.(matI); ok {
		if em := v.EmissiveColor(); em.R == 1 && em.G == 0 && em.B == 0 {
			v.SetEmissiveColor(&math32.Color{R: 0, G: 0, B: 0})
		} else {
			v.SetEmissiveColor(&math32.Color{R: 1, G: 0, B: 0})
		}
	}
}

func addLights() {
	headLight = light.NewDirectional(&math32.Color{R: 1, G: 1, B: 1}, 0.1)
	headLight.SetPosition(0, 0, 10)
	app.Scene().Add(headLight)

	// Adds white directional top light
	l2 := light.NewDirectional(&math32.Color{R: 1, G: 1, B: 1}, 1.0)
	l2.SetPosition(0, 100, 0)
	// app.Scene().Add(l2)

	// Adds white directional right light
	l3 := light.NewDirectional(&math32.Color{R: 1, G: 1, B: 1}, 1.0)
	l3.SetPosition(10, 0, 0)
	// app.Scene().Add(l3)

	// Add lights to the scene
	ambientLight := light.NewAmbient(&math32.Color{R: 1, G: 1, B: 1}, 0.5)
	app.Scene().Add(ambientLight)

	// pointLight = light.NewPoint(&math32.Color{R: 1, G: 1, B: 1}, 5.0)
	// pointLight.SetPosition(3, 50, 4)
	// app.Scene().Add(pointLight)
}

func addGui() {
	// Button 1
	b1 := gui.NewButton("Wireframe")
	b1.SetPosition(10, 10)
	b1.Subscribe(gui.OnClick, func(name string, ev interface{}) {
		// app.Log().Info("button 1 OnClick")
		toggleWireframe()
	})
	app.Gui().GetPanel().Add(b1)

	// Label
	labelSelectedObject = gui.NewLabel("pick with key p")
	labelSelectedObject.SetPosition(100, 10)
	app.Gui().GetPanel().Add(labelSelectedObject)

}

// addRex can take a single filename or a file pattern
func addRex(pattern string) {

	rexGroup = core.NewNode()
	rexGroup.SetName("rexgroup")
	app.Scene().Add(rexGroup)

	matches, err := filepath.Glob(pattern)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, i := range matches {
		dec, err := rex.NewDecoder(i)
		if err != nil {
			panic(err)
		}

		rexBlock, err := dec.NewGroup()
		if err != nil {
			fmt.Println("Cannot create rex group: ", err)
		}
		rexGroup.Add(rexBlock)
	}
	rexGroup.SetScale(scaleFactor, scaleFactor, scaleFactor)
}

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please specify a REX filename")
		os.Exit(1)
	}

	app, _ = application.Create(application.Options{
		Title:    "GO REX Viewer",
		Width:    800,
		Height:   600,
		LogLevel: logger.DEBUG,
	})

	// Registers shaders and program
	app.Renderer().AddShader("shaderGSDemoVertex", sourceGSDemoVertex)
	app.Renderer().AddShader("shaderGSDemoGeometry", sourceGSDemoGeometry)
	app.Renderer().AddShader("shaderGSDemoFrag", sourceGSDemoFrag)
	app.Renderer().AddProgram("progGSDemo", "shaderGSDemoVertex", "shaderGSDemoFrag", "shaderGSDemoGeometry")

	// Attach light to camera
	app.Subscribe(application.OnBeforeRender, func(evname string, ev interface{}) {
		camPos := app.CameraPersp().Position()
		if headLight != nil {
			headLight.SetPosition(camPos.X, camPos.Y, camPos.Z)
		}
		if pointLight != nil {
			pointLight.SetPosition(camPos.X, camPos.Y, camPos.Z)
		}
	})

	rc = core.NewRaycaster(&math32.Vector3{}, &math32.Vector3{})

	addLights()
	addGui()
	addRex(os.Args[1])

	app.Window().Subscribe(window.OnKeyDown, onKey)
	app.Window().Subscribe(window.OnMouseDown, func(evname string, ev interface{}) {
		onMouse(ev)
	})

	grid := graphic.NewGridHelper(gridSize, 1, &math32.Color{R: 0.4, G: 0.4, B: 0.4})
	app.Scene().Add(grid)
	axis := graphic.NewAxisHelper(2)
	app.Scene().Add(axis)

	app.CameraPersp().SetPosition(-10, 10, 10)
	app.CameraPersp().LookAt(math32.NewVector3(0.0, 0.0, 0.0))
	app.Run()
}
