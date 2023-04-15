package app

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/looplab/fsm"
	"github.com/young-steveo/jacquard-ai/config"
	"github.com/young-steveo/jacquard-ai/openai"
	"github.com/young-steveo/jacquard-ai/tools"
	"github.com/young-steveo/jacquard-ai/web"
)

var fsmEventTransitions = fsm.Events{}

// boot-level and global application states
const (
	idle          = "idle"
	booting       = "booting"
	errorConfig   = "errorConfig"
	missingConfig = "missingConfig"
	quitting      = "quitting"
)

// boot-level and global events
const (
	ConfigNotFound     = "configNotFound"
	ErrorReadingConfig = "errorReadingConfig"
	Confirmed          = "confirmed"
	Denied             = "denied"
	Restarted          = "restarted"
	QuitRequested      = "quitRequested"
	BootRequested      = "bootRequested"
	BootSucceeded      = "bootSucceeded"
)

// boot-level and global transitions
func init() {
	fsmEventTransitions = append(fsmEventTransitions,
		// startup
		fsm.EventDesc{Src: []string{idle}, Name: BootRequested, Dst: booting},

		// canot find config file during boot
		fsm.EventDesc{Src: []string{booting}, Name: ConfigNotFound, Dst: missingConfig},
		fsm.EventDesc{Src: []string{booting}, Name: ErrorReadingConfig, Dst: errorConfig},

		// make a new config?
		fsm.EventDesc{Src: []string{booting, missingConfig, errorConfig}, Name: Confirmed, Dst: configWizard},

		fsm.EventDesc{Src: []string{booting}, Name: BootSucceeded, Dst: mainMenu},

		// quit if the user does not want to configure the application
		fsm.EventDesc{Src: []string{missingConfig, errorConfig}, Name: Denied, Dst: quitting},

		// all states can quit
		fsm.EventDesc{Src: []string{idle, booting, missingConfig, errorConfig}, Name: QuitRequested, Dst: quitting},
	)
}

func getAfterEventHandler(g *gocui.Gui, status *Status) fsm.Callback {
	return func(_ context.Context, e *fsm.Event) {
		// after all events, update the status window.
		status.ApplicationState = e.Dst
		if err := showStatus(g, status); err != nil {
			panic(err)
		}
	}
}

func getBootingHandler(g *gocui.Gui, status *Status) fsm.Callback {
	return func(_ context.Context, e *fsm.Event) {
		v := console(g)
		v.Clear()
		tools.FprintlnWhite(v, "Welcome to Jacquard AI!")
		tools.FprintlnWhite(v, "-----------------------")

		// is this the first time the user has run the application?
		if conf, err := config.Read(); err != nil {
			g.Update(func(g *gocui.Gui) error {
				return e.FSM.Event(context.Background(), ErrorReadingConfig, err.Error())
			})
		} else if conf == nil {
			g.Update(func(g *gocui.Gui) error {
				return e.FSM.Event(context.Background(), ConfigNotFound)
			})
		} else {
			localConfig = conf

			tools.FprintlnWhite(v, "Configuration loaded.")

			go func() {
				if status.OpenAIStatus != ConditionGood {
					tools.FprintWhite(v, "Testing connectivity to OpenAI. ")
					client := openai.NewClient(localConfig.OpenAIAPIKey)
					if err := client.TestConnection(); err != nil {
						tools.FprintlnRed(v, "Failed.")
					} else {
						status.OpenAIStatus = ConditionGood
						tools.FprintlnGreen(v, "OK.")
					}
				}

				if status.GoogleStatus != ConditionGood {
					tools.FprintWhite(v, "Testing connectivity to Google. ")
					client, err := web.NewSearcher(context.Background(), localConfig.GoogleAPIKey, localConfig.GoogleSearchEngineID)
					if err != nil {
						tools.FprintlnRed(v, "Failed.")
					}
					if err := client.TestConnection(); err != nil {
						tools.FprintlnRed(v, "Failed.")
					} else {
						status.GoogleStatus = ConditionGood
						tools.FprintlnGreen(v, "OK.")
					}
				}

				if status.OpenAIStatus != ConditionGood || status.GoogleStatus != ConditionGood {
					tools.FprintlnWhite(v, "It looks like there was an error trying to connect to one or more services.")
					tools.FprintlnWhite(v, "Would you like to run the configuration wizard? (y/n)")
					status.WaitingForPrompt = true
				} else {
					g.Update(func(g *gocui.Gui) error {
						return e.FSM.Event(context.Background(), BootSucceeded)
					})
				}
			}()
		}
	}
}

func getErrorConfigHandler(g *gocui.Gui, status *Status) fsm.Callback {
	return func(_ context.Context, e *fsm.Event) {
		v := console(g)
		tools.FprintlnRed(v, e.Args[0].(string))
		tools.FprintlnWhite(v, "It looks like there was an error trying to read the Jacquard configuration file.")
		tools.FprintlnWhite(v, "Would you like to create a new configuration file? (y/n)")
		status.WaitingForPrompt = true
	}
}

func getMissingConfigHandler(g *gocui.Gui, status *Status) fsm.Callback {
	return func(_ context.Context, e *fsm.Event) {
		v := console(g)
		tools.FprintlnWhite(v, "It looks like this is your first time running the program.")
		tools.FprintlnWhite(v, "We'll need to configure Jacquard with API keys before you can get started.")
		tools.FprintWhite(v, "Would you like to do that now? ")
		tools.FprintlnMagenta(v, "(type your answer in the prompt below)")
		status.WaitingForPrompt = true
	}
}

func getQuittingHandler(g *gocui.Gui, status *Status) fsm.Callback {
	return func(_ context.Context, e *fsm.Event) {
		v := console(g)
		tools.FprintlnWhite(v, "Have a nice day!")

		timer := time.NewTimer(2 * time.Second)
		go func() {
			<-timer.C
			g.Close()
			fmt.Println("Goodbye!")
			os.Exit(0)
		}()
	}
}
