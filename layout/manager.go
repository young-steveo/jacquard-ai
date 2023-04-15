package layout

import (
	"errors"
	"fmt"

	"github.com/awesome-gocui/gocui"
	"github.com/young-steveo/jacquard-ai/app"
	"github.com/young-steveo/jacquard-ai/tools"
)

type Manager struct {
	g     *gocui.Gui
	state *app.State
}

func NewManager(g *gocui.Gui, state *app.State) *Manager {
	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	return &Manager{g: g, state: state}
}

func (m *Manager) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	var v *gocui.View
	var err error

	if maxX >= 100 {
		// show logo on large screens.
		if v, err = g.SetView("Logo", 0, 0, maxX/2-2, 7, 0); err != nil {
			if err != nil && !errors.Is(err, gocui.ErrUnknownView) {
				return err
			}
		}
		v.Clear()
		fmt.Fprint(v, logoASCII)

		// Show status in large box top right.
		if _, err := g.SetView("Status", maxX/2-1, 0, maxX/2+maxX/4, 7, 0); err != nil {
			if !errors.Is(err, gocui.ErrUnknownView) {
				return err
			}
		}

		// show large console
		if v, err := g.SetView("Console", 0, 8, maxX/2+maxX/4, maxY-4, 0); err != nil {
			if !errors.Is(err, gocui.ErrUnknownView) {
				return err
			}
			v.Title = " Console "
			v.Wrap = true
		}

		// show backlog on the right
		if v, err := g.SetView("Backlog", maxX/2+maxX/4+1, 0, maxX-2, maxY-20, 0); err != nil {
			if !errors.Is(err, gocui.ErrUnknownView) {
				return err
			}
			v.Title = " Backlog "
			v.Wrap = true
		}

		if v, err := g.SetView("Prompt", 0, maxY-3, maxX/2+maxX/4, maxY-1, 0); err != nil {
			if !errors.Is(err, gocui.ErrUnknownView) {
				return err
			}
			v.Title = " Prompt "
			v.Editable = true
		}

		// Agents are bottom right on large screens
		if v, err := g.SetView("Agents", maxX/2+maxX/4+1, maxY-19, maxX-2, maxY-1, 0); err != nil {
			if !errors.Is(err, gocui.ErrUnknownView) {
				return err
			}
			v.Title = " Agents "
			v.Editable = true
		}

	} else {
		// show small logo on small screens.
		if v, err = g.SetView("Logo", 0, 0, maxX-2, 2, 0); err != nil {
			if !errors.Is(err, gocui.ErrUnknownView) {
				return err
			}
		}
		v.Clear()
		fmt.Fprint(v, "Jacquard AI")

		// show status in small box top.
		if _, err := g.SetView("Status", 0, 3, maxX-2, 7, 0); err != nil {
			if !errors.Is(err, gocui.ErrUnknownView) {
				return err
			}
		}

		// show small console
		if v, err := g.SetView("Console", 0, 8, maxX-2, 25, 0); err != nil {
			if !errors.Is(err, gocui.ErrUnknownView) {
				return err
			}
			v.Title = " Console "
			v.Wrap = true
		}

		if v, err := g.SetView("Prompt", 0, maxY-3, maxX-2, maxY-1, 0); err != nil {
			if !errors.Is(err, gocui.ErrUnknownView) {
				return err
			}
			v.Title = " Prompt "
			v.Editable = true

		}

		// backlog and ageents have a second breakpoint at 50
		if maxX < 50 {
			// tiny screen
			if v, err := g.SetView("Backlog", 0, 26, maxX-2, maxY-14, 0); err != nil {
				if !errors.Is(err, gocui.ErrUnknownView) {
					return err
				}
				v.Title = " Backlog "
				v.Wrap = true
			}
			// agents are after backlog on small screens
			if v, err := g.SetView("Agents", 0, maxY-13, maxX-2, maxY-4, 0); err != nil {
				if !errors.Is(err, gocui.ErrUnknownView) {
					return err
				}
				v.Title = " Agents "
				v.Editable = true
			}

		} else {
			// show backlog after console
			if v, err := g.SetView("Backlog", 0, 26, maxX/2, maxY-4, 0); err != nil {
				if !errors.Is(err, gocui.ErrUnknownView) {
					return err
				}
				v.Title = " Backlog "
				v.Wrap = true
			}
			// agents are after backlog on small screens
			if v, err := g.SetView("Agents", maxX/2+1, 26, maxX-2, maxY-4, 0); err != nil {
				if !errors.Is(err, gocui.ErrUnknownView) {
					return err
				}
				v.Title = " Agents " + fmt.Sprintf("%d", maxX)
				v.Editable = true
			}
		}
	}

	if _, err := g.SetCurrentView("Prompt"); err != nil {
		return err
	}
	if _, err := g.SetViewOnTop("Prompt"); err != nil {
		return err
	}

	return nil
}

func (m *Manager) SetKeybindings() error {
	if err := m.g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, m.quit); err != nil {
		return err
	}
	if err := m.g.SetKeybinding("Prompt", gocui.KeyEnter, gocui.ModNone, m.handlePrompt); err != nil {
		return err
	}
	return nil
}

func (m *Manager) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (m *Manager) handlePrompt(g *gocui.Gui, v *gocui.View) error {
	if !m.state.CanPrompt() {
		return nil
	}

	prompt := v.Buffer()

	if len(prompt) == 0 {
		return nil
	}

	v.Clear()

	v, err := g.View("Console")

	if err != nil {
		return err
	}

	tools.FprintWhite(v, "> ")
	tools.FprintlnCyan(v, prompt)

	return m.state.Prompt(prompt)
}
