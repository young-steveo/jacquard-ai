package domain

import (
	"bytes"
	"strings"

	"github.com/BurntSushi/toml"
)

type Action struct {
	Name   string            `json:"name" toml:"name"`
	Reason string            `json:"reason" toml:"reason"`
	Params map[string]string `json:"params" toml:"params"`
}

func ParseAction(text string) (Action, error) {
	text = strings.Trim(text, "`")

	var a Action
	_, err := toml.Decode(text, &a)
	return a, err
}

func (a Action) String() string {
	buf := new(bytes.Buffer)
	encoder := toml.NewEncoder(buf)
	encoder.Indent = ""
	encoder.Encode(a)

	return "```\n" + buf.String() + "```"
}
