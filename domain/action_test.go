package domain_test

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/young-steveo/jacquard-ai/domain"
)

func TestParseAction(t *testing.T) {
	// Arrange
	expected := domain.Action{
		Name:   "search",
		Reason: "I want to know how to build a website",
		Params: map[string]string{
			"query": "How to build a website",
		},
	}
	text := "```\n" + `name = "search"
reason = "I want to know how to build a website"

[params]
query = "How to build a website"
` + "```"

	// Act
	actual, err := domain.ParseAction(text)

	// Assert
	assert.Equal(t, nil, err)
	assert.Equal(t, expected, actual)
}

func TestAction_String(t *testing.T) {
	// Arrange
	expected := "```\n" + `name = "search"
reason = "I want to know how to build a website"

[params]
query = "How to build a website"
` + "```"

	// Act
	actual := domain.Action{
		Name:   "search",
		Reason: "I want to know how to build a website",
		Params: map[string]string{
			"query": "How to build a website",
		},
	}.String()

	// Assert
	assert.Equal(t, expected, actual)
}
