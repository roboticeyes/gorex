package main

import (
	"github.com/gdamore/tcell"
)

// const defines the default key bindings
const (
	KeySwitchToHelp    = tcell.KeyF1
	KeySwitchToProject = tcell.KeyF2
	KeySwitchToAdmin   = tcell.KeyF3
)

var keyString = map[tcell.Key](string){
	tcell.KeyF1: "F1",
	tcell.KeyF2: "F2",
	tcell.KeyF3: "F3",
}
