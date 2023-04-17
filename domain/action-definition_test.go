package domain_test

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/young-steveo/jacquard-ai/domain"
)

func TestActionDefinition_GetAction(t *testing.T) {
	// Arrange
	expectedName := "search"
	expectedParams := map[string]string{
		"query": "How to build a website",
	}
	actionDefinition := domain.ActionDefinition{
		Name:        expectedName,
		Description: "Use a search engine",
		Params: map[string]string{
			"query": "string",
		},
	}

	// Act
	actual, err := actionDefinition.GetAction(expectedParams)

	// Assert
	assert.Equal(t, nil, err)
	assert.Equal(t, expectedName, actual.Name)
	assert.Equal(t, "", actual.Reason)
	assert.Equal(t, expectedParams, actual.Params)
}

func TestActionDefinition_GetAction_MissingParam(t *testing.T) {
	// Arrange
	actionDefinition := domain.ActionDefinition{
		Name:        "search",
		Description: "Use a search engine",
		Params: map[string]string{
			"query": "string",
		},
	}

	// Act
	_, err := actionDefinition.GetAction(map[string]string{})

	// Assert
	assert.Equal(t, "missing required parameter query", err.Error())
}

func TestActionDefinition_GetAction_UnexpectedParam(t *testing.T) {
	// Arrange
	actionDefinition := domain.ActionDefinition{
		Name:        "search",
		Description: "Use a search engine",
		Params: map[string]string{
			"query": "string",
		},
	}

	// Act
	_, err := actionDefinition.GetAction(map[string]string{
		"query": "How to build a website",
		"foo":   "bar",
	})

	// Assert
	assert.Equal(t, "unexpected parameter foo", err.Error())
}

func TestParseActionDefinition(t *testing.T) {
	// Arrange
	expected := domain.ActionDefinition{
		Name:        "search",
		Description: "Use a search engine",
		Params: map[string]string{
			"query": "string",
		},
	}
	text := "```\n" + `name = "search"
description = "Use a search engine"

[params]
query = "string"
` + "```"

	// Act
	actual, err := domain.ParseActionDefinition(text)

	// Assert
	assert.Equal(t, nil, err)
	assert.Equal(t, expected, actual)
}

func TestActionDefinition_String(t *testing.T) {
	// Arrange
	expected := "```\n" + `name = "search"
description = "Use a search engine"

[params]
query = "string"
` + "```"

	// Act
	actual := domain.ActionDefinition{
		Name:        "search",
		Description: "Use a search engine",
		Params: map[string]string{
			"query": "string",
		},
	}.String()

	// Assert
	assert.Equal(t, expected, actual)
}
