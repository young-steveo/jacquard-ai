package app

import (
	"context"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/fatih/color"
	"github.com/looplab/fsm"
	"github.com/young-steveo/jacquard-ai/config"
	"github.com/young-steveo/jacquard-ai/tools"
)

var localConfig *config.Config

type State struct {
	g      *gocui.Gui
	fsm    *fsm.FSM
	status *Status
}

func console(g *gocui.Gui) *gocui.View {
	v, err := g.View("Console")
	if err != nil {
		panic(err)
	}
	return v
}

func NewState(g *gocui.Gui) *State {
	status := NewStatus(idle)
	return &State{
		g:      g,
		status: status,
		fsm: fsm.NewFSM(
			idle,
			fsmEventTransitions,
			fsm.Callbacks{
				"after_event": getAfterEventHandler(g, status),
				booting:       getBootingHandler(g, status),
				errorConfig:   getErrorConfigHandler(g, status),
				missingConfig: getMissingConfigHandler(g, status),
				quitting:      getQuittingHandler(g, status),

				// FSM for config wizard
				configWizard:               getConfigWizardHandler(g, status),
				configWizardReenterOpenAI:  getConfigWizardReenterOpenAIHandler(g, status),
				configWizardTestingOpenAI:  getConfigWizardTestingOpenAIHandler(g, status),
				configWizardSearchEngineID: getConfigWizardSearchEngineIDHandler(g, status),
				configWizardGoogleAPIKey:   getConfigWizardGoogleAPIKeyHander(g, status),
				configWizardTestingGoogle:  getConfigWizardTestingGoogleHandler(g, status),
				configWizardReenterGoogle:  getConfigWizardReenterGoogleHandler(g, status),
				configWizardRestarting:     getConfigWizardRestartingHander(g, status),
				configWizardFinalize:       getConfigWizardFinalizeHandler(g, status),
				configWizardAskToSave:      getConfigWizardAskToSaveHandler(g, status),
				configWizardSaving:         getConfigWizardSavingHandler(g, status),
				configWizardNotSaving:      getConfigWizardNotSavingHandler(g, status),

				// FSM for Main Menu
				mainMenu:     getMainMenuHandler(g, status),
				newProject:   getNewProjectHandler(g, status),
				listProjects: getListProjectsHandler(g, status),
			},
		),
	}
}

// Handle prompts from the user during myriad states.
func (s *State) Prompt(prompt string) error {
	s.status.ApplicationState = s.fsm.Current()
	s.status.WaitingForPrompt = false

	// no matter what, any time the user types "quit" or "exit", we should quit.
	if prompt == "quit" || prompt == "exit" {
		s.fsm.Event(context.Background(), QuitRequested)
		return showStatus(s.g, s.status)
	}

	switch s.status.ApplicationState {
	case booting:
		// if we need to check a prompt during booting, it's because something went wrong with one of the services.
		// we have asked the user if they want to start the configuration server.
		err := s.confirmPrompt(prompt, "")
		if err != nil {
			return err
		}
	case errorConfig:
		err := s.confirmPrompt(prompt, "")
		if err != nil {
			return err
		}

	case missingConfig:
		err := s.confirmPrompt(prompt, "")
		if err != nil {
			return err
		}
	case configWizard:
		restartOrNext(prompt, s.fsm, &localConfig.OpenAIAPIKey, OpenAIAPIKeyCompleted)

	case configWizardReenterOpenAI:
		restartOrNext(prompt, s.fsm, &localConfig.OpenAIAPIKey, OpenAIAPIKeyCompleted)

	case configWizardSearchEngineID:
		restartOrNext(prompt, s.fsm, &localConfig.GoogleSearchEngineID, GoogleSearchEngineIDCompleted)

	case configWizardReenterGoogle:
		restartOrNext(prompt, s.fsm, &localConfig.GoogleSearchEngineID, GoogleSearchEngineIDCompleted)

	case configWizardGoogleAPIKey:
		restartOrNext(prompt, s.fsm, &localConfig.GoogleAPIKey, GoogleAPIKeyCompleted)

	case configWizardFinalize:
		// if the user stops here for a prompt, it's because something went wrong.
		// otherwise they will already be forwarded to the next step.
		restartOrQuit(prompt, s.fsm)

	case configWizardAskToSave:
		err := s.confirmPrompt(prompt, "")
		if err != nil {
			return err
		}

	case mainMenu:
		if MenuOptionOne.MatchString(prompt) {
			s.fsm.Event(context.Background(), MenuOptionNewSelected)
		} else if MenuOptionTwo.MatchString(prompt) {
			s.fsm.Event(context.Background(), MenuOptionListSelected)
		} else if MenuOptionThree.MatchString(prompt) {
			s.fsm.Event(context.Background(), QuitRequested)
		} else if MenuOptionProject.MatchString(prompt) {
			if MenuOptionProjects.MatchString(prompt) {
				s.fsm.Event(context.Background(), MenuOptionListSelected)
			} else {
				s.fsm.Event(context.Background(), MenuOptionNewSelected)
			}
		} else {
			tools.FprintlnRed(console(s.g), "I didn't understand that. Try typing one of the numbers or words in the menu.")
		}

	default:
		// application state not implemented
		v, err := s.g.View("Console")
		if err != nil {
			return err
		}
		tools.FprintlnRed(v, "Parsing prompts during `"+s.status.ApplicationState+"` is not implemented.")
	}

	return showStatus(s.g, s.status)
}

func restartOrNext(prompt string, fsm *fsm.FSM, configKey *string, nextEvent string) {
	if prompt == "restart" {
		fsm.Event(context.Background(), Restarted)
		return
	}
	*configKey = prompt
	fsm.Event(context.Background(), nextEvent)

}

func restartOrQuit(prompt string, fsm *fsm.FSM) {
	if prompt == "restart" {
		fsm.Event(context.Background(), Restarted)
	} else if prompt == "quit" {
		fsm.Event(context.Background(), QuitRequested)
	}
}

func (s *State) confirmPrompt(prompt string, followup string) error {
	switch sentiment, err := tools.GetSentiment(prompt); err == nil {
	case sentiment == tools.Positive:
		go func() {
			<-time.NewTimer(time.Second).C
			s.g.Update(func(g *gocui.Gui) error {
				return s.fsm.Event(context.Background(), Confirmed)
			})
		}()
	case sentiment == tools.Negative:
		go func() {
			<-time.NewTimer(time.Second).C
			s.g.Update(func(g *gocui.Gui) error {
				return s.fsm.Event(context.Background(), Denied)
			})
		}()
	case sentiment == tools.Ambiguous:
		v, err := s.g.View("Console")
		if err != nil {
			return err
		}
		tools.FprintWhite(v, "I'm sorry, '")
		tools.FprintCyan(v, prompt)
		tools.FprintWhite(v, "' is too ambiguous; I don't want to do the wrong thing, so did you mean ")
		tools.FprintGreen(v, "yes")
		tools.FprintWhite(v, " or ")
		tools.FprintRed(v, "no")
		tools.FprintlnWhite(v, "?")
		if len(followup) > 0 {
			tools.FprintlnWhite(v, followup)
		}
		s.status.WaitingForPrompt = true
	}
	return nil
}

func (s *State) ShowStatus() error {
	return showStatus(s.g, s.status)
}

func showStatus(g *gocui.Gui, status *Status) error {
	// write the State to the Status view
	v, err := g.View("Status")
	if err != nil {
		return err
	}
	v.Clear()
	for _, msg := range status.Messages() {
		color.New(msg.Color).Fprint(v, msg.Text)
	}

	// update Console subtitle status
	v, err = g.View("Console")
	if err != nil {
		return err
	}

	if status.WaitingForPrompt {
		v.TitleColor = gocui.ColorGreen
		v.Subtitle = " waiting for input "
	} else {
		v.TitleColor = gocui.ColorWhite
		v.Subtitle = ""
	}
	return err
}

func (s *State) Event(event string, args ...interface{}) error {
	return s.fsm.Event(context.Background(), event, args...)
}

func (s *State) CanPrompt() bool {
	return s.status.WaitingForPrompt
}
