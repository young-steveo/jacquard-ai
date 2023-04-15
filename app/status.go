package app

import (
	"github.com/fatih/color"
	"github.com/gobuffalo/flect"
)

type Condition string

const (
	ConditionUnknown Condition = "unknown"
	ConditionGood    Condition = "good"
	ConditionBad     Condition = "bad"
)

type StatusMsg struct {
	Text  string
	Color color.Attribute
}

type Status struct {
	ApplicationState string
	WaitingForPrompt bool
	OpenAIStatus     Condition
	GoogleStatus     Condition
}

func NewStatus(state string) *Status {
	return &Status{
		ApplicationState: state,
		WaitingForPrompt: false,
		OpenAIStatus:     ConditionUnknown,
		GoogleStatus:     ConditionUnknown,
	}
}

func (s *Status) Messages() []StatusMsg {
	msgs := []StatusMsg{}
	msgs = append(msgs, StatusMsg{Text: "State:   ", Color: color.FgWhite})
	txt := flect.New(s.ApplicationState).Humanize().String()
	switch s.ApplicationState {
	case errorConfig:
		msgs = append(msgs, StatusMsg{Text: txt + "\n", Color: color.FgRed})
	case missingConfig:
		msgs = append(msgs, StatusMsg{Text: txt + "\n", Color: color.FgRed})
	default:
		msgs = append(msgs, StatusMsg{Text: txt + "\n", Color: color.FgYellow})
	}

	msgs = append(msgs, StatusMsg{Text: "OpenAPI: ", Color: color.FgWhite})

	switch s.OpenAIStatus {
	case ConditionGood:
		msgs = append(msgs, StatusMsg{Text: "Good\n", Color: color.FgGreen})
	case ConditionBad:
		msgs = append(msgs, StatusMsg{Text: "Bad\n", Color: color.FgRed})
	default:
		msgs = append(msgs, StatusMsg{Text: "Unknown\n", Color: color.FgYellow})
	}

	msgs = append(msgs, StatusMsg{Text: "Google:  ", Color: color.FgWhite})

	switch s.GoogleStatus {
	case ConditionGood:
		msgs = append(msgs, StatusMsg{Text: "Good\n", Color: color.FgGreen})
	case ConditionBad:
		msgs = append(msgs, StatusMsg{Text: "Bad\n", Color: color.FgRed})
	default:
		msgs = append(msgs, StatusMsg{Text: "Unknown\n", Color: color.FgYellow})
	}

	return msgs
}
