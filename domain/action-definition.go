package domain

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/BurntSushi/toml"
)

type ActionDefinition struct {
	Name        string `json:"name" toml:"name"`
	Description string `json:"description" toml:"description"`

	// Params is a key-value map that describes the parameters
	// that can be passed to the action.
	Params map[string]string `json:"params" toml:"params"`
}

func (ad ActionDefinition) GetAction(params map[string]string) (Action, error) {
	for param := range ad.Params {
		if _, ok := params[param]; !ok {
			return Action{}, fmt.Errorf("missing required parameter %s", param)
		}
	}

	for param := range params {
		if _, ok := ad.Params[param]; !ok {
			return Action{}, fmt.Errorf("unexpected parameter %s", param)
		}
	}

	return Action{
		Name:   ad.Name,
		Reason: "",
		Params: params,
	}, nil
}

func ParseActionDefinition(text string) (ActionDefinition, error) {
	text = strings.Trim(text, "`")

	var ad ActionDefinition
	_, err := toml.Decode(text, &ad)
	return ad, err
}

func (a ActionDefinition) String() string {
	buf := new(bytes.Buffer)
	encoder := toml.NewEncoder(buf)
	encoder.Indent = ""
	encoder.Encode(a)
	return "```\n" + buf.String() + "```"
}
