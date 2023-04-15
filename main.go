package main

import (
	"errors"
	"log"

	"github.com/awesome-gocui/gocui"
	"github.com/young-steveo/jacquard-ai/app"
	"github.com/young-steveo/jacquard-ai/layout"
)

func main() {
	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	state := app.NewState(g)

	mgr := layout.NewManager(g, state)

	g.SetManager(mgr)

	if err := mgr.SetKeybindings(); err != nil {
		log.Panicln(err)
	}

	g.Update(func(g *gocui.Gui) error {
		if err := state.ShowStatus(); err != nil {
			return err
		}
		return state.Event(app.BootRequested)
	})

	if err := g.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		log.Panicln(err)
	}
}
