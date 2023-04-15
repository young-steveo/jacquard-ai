package app

import (
	"context"
	"time"

	"github.com/awesome-gocui/gocui"
	"github.com/looplab/fsm"
	"github.com/young-steveo/jacquard-ai/config"
	"github.com/young-steveo/jacquard-ai/openai"
	"github.com/young-steveo/jacquard-ai/tools"
	"github.com/young-steveo/jacquard-ai/web"
)

// config wizard states
const (
	configWizard              = "configWizard"
	configWizardRestarting    = "configWizardRestarting"
	configWizardTestingOpenAI = "configWizardTestingOpenAI"
	configWizardReenterOpenAI = "configWizardReenterOpenAI"

	configWizardSearchEngineID = "configWizardSearchEngineID"
	configWizardGoogleAPIKey   = "configWizardGoogleAPIKey"
	configWizardTestingGoogle  = "configWizardTestingGoogle"
	configWizardReenterGoogle  = "configWizardReenterGoogle"

	configWizardFinalize  = "configWizardFinalize"
	configWizardAskToSave = "configWizardAskToSave"
	configWizardSaving    = "configWizardSaving"
	configWizardNotSaving = "configWizardNotSaving"
)

// config wizard events
const (
	OpenAIAPIKeyCompleted = "openAIAPIKeyCompleted"
	OpenAIAPIKeyTestedOK  = "openAIAPIKeyTestedOK"
	OpenAIAPIKeyTestedBad = "openAIAPIKeyTestedBad"

	GoogleSearchEngineIDCompleted = "googleSearchEngineIDCompleted"
	GoogleAPIKeyCompleted         = "googleAPIKeyCompleted"
	GoogleSearchEngineTestedOK    = "googleSearchEngineTestedOK"
	GoogleSearchEngineTestedBad   = "googleSearchEngineTestedBad"

	ConfigurationWizardCompleted = "configurationWizardCompleted"
)

// config wizard transitions
func init() {
	fsmEventTransitions = append(fsmEventTransitions,
		// wizard states
		fsm.EventDesc{Src: []string{configWizard, configWizardSearchEngineID}, Name: Restarted, Dst: configWizardRestarting},
		fsm.EventDesc{Src: []string{configWizardRestarting}, Name: Restarted, Dst: configWizard},

		// wizard steps
		fsm.EventDesc{Src: []string{configWizard, configWizardReenterOpenAI}, Name: OpenAIAPIKeyCompleted, Dst: configWizardTestingOpenAI},
		fsm.EventDesc{Src: []string{configWizardTestingOpenAI}, Name: OpenAIAPIKeyTestedBad, Dst: configWizardReenterOpenAI},
		fsm.EventDesc{Src: []string{configWizardTestingOpenAI}, Name: OpenAIAPIKeyTestedOK, Dst: configWizardSearchEngineID},

		fsm.EventDesc{Src: []string{configWizardSearchEngineID, configWizardReenterGoogle}, Name: GoogleSearchEngineIDCompleted, Dst: configWizardGoogleAPIKey},
		fsm.EventDesc{Src: []string{configWizardGoogleAPIKey}, Name: GoogleAPIKeyCompleted, Dst: configWizardTestingGoogle},
		fsm.EventDesc{Src: []string{configWizardTestingGoogle}, Name: GoogleSearchEngineTestedBad, Dst: configWizardReenterGoogle},
		fsm.EventDesc{Src: []string{configWizardTestingGoogle}, Name: GoogleSearchEngineTestedOK, Dst: configWizardFinalize},

		fsm.EventDesc{Src: []string{configWizardFinalize}, Name: ConfigurationWizardCompleted, Dst: configWizardAskToSave},
		fsm.EventDesc{Src: []string{configWizardAskToSave}, Name: Confirmed, Dst: configWizardSaving},
		fsm.EventDesc{Src: []string{configWizardAskToSave}, Name: Denied, Dst: configWizardNotSaving},
		fsm.EventDesc{Src: []string{configWizardSaving}, Name: Confirmed, Dst: idle},

		fsm.EventDesc{Src: []string{configWizardNotSaving, configWizardSaving}, Name: Restarted, Dst: booting},

		// all states can quit
		fsm.EventDesc{Src: []string{
			configWizard,
			configWizardRestarting,
			configWizardTestingOpenAI,
			configWizardReenterOpenAI,
			configWizardSearchEngineID,
			configWizardGoogleAPIKey,
			configWizardTestingGoogle,
			configWizardReenterGoogle,
			configWizardFinalize,
			configWizardAskToSave,
			configWizardSaving,
			configWizardNotSaving,
		}, Name: QuitRequested, Dst: quitting},
	)
}

func getConfigWizardHandler(g *gocui.Gui, status *Status) fsm.Callback {
	return func(_ context.Context, e *fsm.Event) {
		// start with a new blank config
		localConfig = &config.Config{}
		status.OpenAIStatus = ConditionUnknown
		status.GoogleStatus = ConditionUnknown

		v := console(g)
		v.Clear()
		printConfigWizardTitle(v)

		tools.FprintlnWhite(v, "Jacquard connects to OpenAI's API to communicate with ChatGPT and other models.")
		tools.FprintWhite(v, "Please enter your ")
		tools.FprintGreen(v, "OpenAI API key")
		tools.FprintlnWhite(v, " in the prompt below.")
		tools.FprintlnMagenta(v, "(You can obtain an API key at https://platform.openai.com/account/api-keys)")

		printQuitRestartMessage(v)

		tools.FprintlnBlue(v, "Anything else will be interpreted as your API key.")

		status.WaitingForPrompt = true
	}
}

func getConfigWizardReenterOpenAIHandler(g *gocui.Gui, status *Status) fsm.Callback {
	return func(_ context.Context, e *fsm.Event) {
		localConfig.OpenAIAPIKey = ""
		v := console(g)
		v.Clear()
		printConfigWizardTitle(v)

		tools.FprintlnWhite(v, "Uh Oh, it looks like the API key you entered was invalid, or perhaps the OpenAI API is down.")
		tools.FprintWhite(v, "Reason Given: ")
		tools.FprintlnRed(v, e.Args[0].(string))
		tools.FprintlnWhite(v, " ")

		tools.FprintWhite(v, "Please reenter your ")
		tools.FprintGreen(v, "OpenAI API key")
		tools.FprintlnWhite(v, " in the prompt below.")

		printQuitRestartMessage(v)

		tools.FprintlnBlue(v, "Anything else will be interpreted as your API key.")

		status.WaitingForPrompt = true
	}
}

func getConfigWizardTestingOpenAIHandler(g *gocui.Gui, status *Status) fsm.Callback {
	return func(_ context.Context, e *fsm.Event) {
		// if we are testing the connection, the condition is unknown
		status.OpenAIStatus = ConditionUnknown

		client := openai.NewClient(localConfig.OpenAIAPIKey)

		go func() {
			err := client.TestConnection()
			if err != nil {
				status.OpenAIStatus = ConditionBad

				<-time.NewTimer(time.Second).C
				g.Update(func(g *gocui.Gui) error {
					return e.FSM.Event(context.Background(), OpenAIAPIKeyTestedBad, err.Error())
				})

				return
			}

			status.OpenAIStatus = ConditionGood
			<-time.NewTimer(time.Second).C
			g.Update(func(g *gocui.Gui) error {
				return e.FSM.Event(context.Background(), OpenAIAPIKeyTestedOK)
			})
		}()
	}
}

func getConfigWizardSearchEngineIDHandler(g *gocui.Gui, status *Status) fsm.Callback {
	return func(_ context.Context, e *fsm.Event) {
		v := console(g)
		v.Clear()
		printConfigWizardTitle(v)

		tools.FprintlnWhite(v, "Jacquard uses a search engine to help complete tasks.")
		tools.FprintWhite(v, "Google provides a nifty solution called ")
		tools.FprintMagenta(v, "Programmable Search Engine")
		tools.FprintlnWhite(v, ", which allows Jacquard to search without a browser.")
		tools.FprintlnWhite(v, " ")

		tools.FprintlnWhite(v, "You must first enable Google's Custom Search API in your Google Cloud account.")
		tools.FprintlnMagenta(v, "(You can enable it at https://console.cloud.google.com/apis/library/customsearch.googleapis.com)")
		tools.FprintlnWhite(v, " ")

		tools.FprintWhite(v, "Once that's done, create a new Programmable Search Engine and enter your ")
		tools.FprintGreen(v, "Programmable Search Engine ID")
		tools.FprintlnWhite(v, " in the prompt below.")

		tools.FprintlnMagenta(v, "(You can create a Programmable Search Engine at https://programmablesearchengine.google.com/controlpanel/create)")

		printQuitRestartMessage(v)

		tools.FprintlnBlue(v, "Anything else will be interpreted as your Programmable Search Engine ID.")

		status.WaitingForPrompt = true
	}
}

func getConfigWizardGoogleAPIKeyHander(g *gocui.Gui, status *Status) fsm.Callback {
	return func(_ context.Context, e *fsm.Event) {
		v := console(g)
		v.Clear()
		printConfigWizardTitle(v)

		tools.FprintlnWhite(v, "Now we need an API key from Google in order to access the search engine.")

		tools.FprintWhite(v, "Please enter your ")
		tools.FprintGreen(v, "Google API Key")
		tools.FprintlnWhite(v, " in the prompt below.")

		tools.FprintlnMagenta(v, "(You can create an API key at https://console.cloud.google.com/apis/credentials)")

		printQuitRestartMessage(v)

		tools.FprintlnBlue(v, "Anything else will be interpreted as your Google API Key.")

		status.WaitingForPrompt = true
	}
}

func getConfigWizardTestingGoogleHandler(g *gocui.Gui, status *Status) fsm.Callback {
	return func(_ context.Context, e *fsm.Event) {
		// if we are testing the connection, the condition is unknown
		status.GoogleStatus = ConditionUnknown

		searcher, err := web.NewSearcher(context.Background(), localConfig.GoogleAPIKey, localConfig.GoogleSearchEngineID)
		if err != nil {
			status.GoogleStatus = ConditionBad
			go func() {
				<-time.NewTimer(time.Second).C
				g.Update(func(g *gocui.Gui) error {
					return e.FSM.Event(context.Background(), GoogleSearchEngineTestedBad, err.Error())
				})
			}()
			return
		}

		go func() {
			err = searcher.TestConnection()

			if err != nil {
				status.GoogleStatus = ConditionBad
				reason := err.Error()

				<-time.NewTimer(time.Second).C
				g.Update(func(g *gocui.Gui) error {
					return e.FSM.Event(context.Background(), GoogleSearchEngineTestedBad, reason)
				})

				return
			}

			status.GoogleStatus = ConditionGood
			<-time.NewTimer(time.Second).C
			g.Update(func(g *gocui.Gui) error {
				return e.FSM.Event(context.Background(), GoogleSearchEngineTestedOK)
			})
		}()
	}
}

func getConfigWizardReenterGoogleHandler(g *gocui.Gui, status *Status) fsm.Callback {
	return func(_ context.Context, e *fsm.Event) {
		v := console(g)
		v.Clear()
		printConfigWizardTitle(v)

		tools.FprintlnWhite(v, "Uh Oh, we were not able to successfully search Google.")
		tools.FprintWhite(v, "Reason Given: ")
		tools.FprintlnRed(v, e.Args[0].(string))
		tools.FprintlnWhite(v, " ")
		tools.FprintWhite(v, "Please reenter your ")
		tools.FprintGreen(v, "Programmable Search Engine ID")
		tools.FprintlnWhite(v, " in the prompt below.")

		printQuitRestartMessage(v)

		tools.FprintlnBlue(v, "Anything else will be interpreted as your Programmable Search Engine ID.")

		status.WaitingForPrompt = true
	}
}

func getConfigWizardRestartingHander(g *gocui.Gui, status *Status) fsm.Callback {
	return func(_ context.Context, e *fsm.Event) {
		v := console(g)
		tools.FprintlnWhite(v, "Restarting Config Wizard.")
		timer := time.NewTimer(time.Second)
		go func() {
			<-timer.C
			g.Update(func(g *gocui.Gui) error {
				return e.FSM.Event(context.Background(), Restarted)
			})
		}()
	}
}

func getConfigWizardFinalizeHandler(g *gocui.Gui, status *Status) fsm.Callback {
	return func(_ context.Context, e *fsm.Event) {
		v := console(g)
		v.Clear()

		if status.GoogleStatus != ConditionGood || status.OpenAIStatus != ConditionGood {
			printConfigWizardTitle(v)
			tools.FprintlnRed(v, "Uh Oh, we were not able to successfully connect to all of the services.")
			tools.FprintWhite(v, "Would you like to `")
			tools.FprintCyan(v, "restart")
			tools.FprintWhite(v, "` the Configuration Wizard, or `")
			tools.FprintCyan(v, "quit")
			tools.FprintlnWhite(v, "` Jacquard?")
			tools.FprintlnWhite(v, " ")
			tools.FprintlnMagenta(v, "(If you are still having trouble after several attempts, file an issue at https://github.com/young-steveo/jacquard/issues )")

			status.WaitingForPrompt = true
			return
		}
		g.Update(func(g *gocui.Gui) error {
			return e.FSM.Event(context.Background(), ConfigurationWizardCompleted)
		})
	}
}

func getConfigWizardAskToSaveHandler(g *gocui.Gui, status *Status) fsm.Callback {
	return func(_ context.Context, e *fsm.Event) {
		v := console(g)
		v.Clear()
		printConfigWizardTitle(v)
		tools.FprintlnWhite(v, "You have successfully configured Jacquard!")
		tools.FprintlnWhite(v, " ")
		tools.FprintlnWhite(v, "Would you like to save this configuration for future use?")
		status.WaitingForPrompt = true
	}
}

func getConfigWizardSavingHandler(g *gocui.Gui, status *Status) fsm.Callback {
	return func(_ context.Context, e *fsm.Event) {
		v := console(g)
		v.Clear()
		printConfigWizardTitle(v)
		tools.FprintlnWhite(v, "OK! Saving configuration for future use and restarting Jacquard.")

		go func() {
			err := config.Write(localConfig)
			if err != nil {
				tools.FprintlnRed(v, "Error writing config file: "+err.Error())
			}
			timer := time.NewTimer(5 * time.Second)
			go func() {
				<-timer.C
				g.Update(func(g *gocui.Gui) error {
					return e.FSM.Event(context.Background(), Restarted)
				})
			}()
		}()
	}
}

func getConfigWizardNotSavingHandler(g *gocui.Gui, status *Status) fsm.Callback {
	return func(_ context.Context, e *fsm.Event) {
		v := console(g)
		v.Clear()
		printConfigWizardTitle(v)
		tools.FprintlnWhite(v, "OK, we'll only use this configuration for the current session. Please wait a few seconds while we reload.")

		timer := time.NewTimer(5 * time.Second)
		go func() {
			<-timer.C
			g.Update(func(g *gocui.Gui) error {
				return e.FSM.Event(context.Background(), Restarted)
			})
		}()
	}
}

func printQuitRestartMessage(v *gocui.View) {
	tools.FprintlnWhite(v, " ")
	tools.FprintBlue(v, "If you'd like to start the Config Wizard over, just type `")
	tools.FprintCyan(v, "restart")
	tools.FprintlnBlue(v, "` and press enter.")

	tools.FprintBlue(v, "If you'd like to quit Jacquard, type `")
	tools.FprintCyan(v, "quit")
	tools.FprintlnBlue(v, "`.")
}

func printConfigWizardTitle(v *gocui.View) {
	tools.FprintWhite(v, "Jacquard AI")
	tools.FprintlnWhite(v, " Configuration Wizard")
	tools.FprintlnWhite(v, "--------------------------------")
}
