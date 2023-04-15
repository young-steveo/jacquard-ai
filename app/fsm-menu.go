package app

import (
	"context"
	"regexp"

	"github.com/awesome-gocui/gocui"
	"github.com/looplab/fsm"
	"github.com/young-steveo/jacquard-ai/tools"
)

// Main Menu regex
var (
	MenuOptionOne      = regexp.MustCompile(`1|start|new|first|one`)
	MenuOptionTwo      = regexp.MustCompile(`2|view|previous|second|two`)
	MenuOptionThree    = regexp.MustCompile(`3|quit|third|three`)
	MenuOptionProject  = regexp.MustCompile(`project`)
	MenuOptionProjects = regexp.MustCompile(`projects`)
)

// Main Menu states
const (
	mainMenu     = "mainMenu"
	newProject   = "newProject"
	listProjects = "listProjects"
)

// Main Menu events
const (
	MenuOptionNewSelected  = "menuOptionNewSelected"
	MenuOptionListSelected = "menuOptionListSelected"
	MainMenuRequested      = "mainMenuRequested"
)

// Main Menu transitions
func init() {
	fsmEventTransitions = append(fsmEventTransitions,
		fsm.EventDesc{Src: []string{mainMenu}, Name: MenuOptionNewSelected, Dst: newProject},
		fsm.EventDesc{Src: []string{mainMenu}, Name: MenuOptionListSelected, Dst: listProjects},
		fsm.EventDesc{Src: []string{mainMenu}, Name: MainMenuRequested, Dst: mainMenu},

		// all states can quit
		fsm.EventDesc{Src: []string{mainMenu, newProject, listProjects}, Name: QuitRequested, Dst: quitting},
	)
}

// Main menu handlers
func getMainMenuHandler(g *gocui.Gui, status *Status) fsm.Callback {
	return func(_ context.Context, e *fsm.Event) {
		v := console(g)
		v.Clear()

		tools.FprintlnWhite(v, "Jacquard AI Main Menu")
		tools.FprintlnWhite(v, "---------------------")

		tools.FprintlnWhite(v, " ")
		tools.FprintlnWhite(v, "What would you like to do today?")
		tools.FprintlnWhite(v, " ")
		tools.FprintlnWhite(v, "1. Start a new project")
		tools.FprintlnWhite(v, "2. View previous projects")
		tools.FprintlnWhite(v, "3. Quit")

		status.WaitingForPrompt = true
	}
}

func getNewProjectHandler(g *gocui.Gui, status *Status) fsm.Callback {
	return func(_ context.Context, e *fsm.Event) {
	}
}

func getListProjectsHandler(g *gocui.Gui, status *Status) fsm.Callback {
	return func(_ context.Context, e *fsm.Event) {
	}
}
