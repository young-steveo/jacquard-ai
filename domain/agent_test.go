package domain_test

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/young-steveo/jacquard-ai/domain"
)

func TestNewAgent(t *testing.T) {
	// Arrange
	expectedID := domain.AgentID(1)
	expectedName := "Agent"
	expectedRole := domain.Role("task_coordinator")
	expectedObjective := "Build a website"
	expectedActions := []domain.ActionDefinition{
		{
			Name:        "search",
			Description: "Use a search engine",
		},
	}

	// Act
	actual := domain.NewAgent(expectedID, expectedName, expectedRole, expectedObjective, expectedActions)

	// Assert
	assert.Equal(t, expectedID, actual.ID)
	assert.Equal(t, expectedName, actual.Name)
	assert.Equal(t, expectedRole, actual.Role)
	assert.Equal(t, expectedObjective, actual.Objective)
	assert.Equal(t, expectedActions, actual.Actions)
}

func TestAgent_AddAction(t *testing.T) {
	// Arrange
	agent := domain.NewAgent(domain.AgentID(1), "Agent", domain.Role("task_coordinator"), "Build a website", []domain.ActionDefinition{})
	expected := domain.ActionDefinition{
		Name:        "search",
		Description: "Use a search engine",
	}

	// Act
	agent.AddAction(expected)

	// Assert
	assert.Equal(t, 1, len(agent.Actions))
	assert.Equal(t, expected, agent.Actions[0])
}

func TestAgent_RemoveAction(t *testing.T) {
	// Arrange
	expected := []domain.ActionDefinition{
		{
			Name:        "search",
			Description: "Use a search engine",
		},
		{
			Name:        "loadURL",
			Description: "Load a URL, e.g. a website, formatted as a markdown document",
		},
	}
	agent := domain.NewAgent(domain.AgentID(1), "Agent", domain.Role("task_coordinator"), "Build a website", expected)

	// Act
	agent.RemoveAction(expected[0].Name)

	// Assert
	assert.Equal(t, 1, len(agent.Actions))
}

func TestAgent_HasAction(t *testing.T) {
	// Arrange
	expected := []domain.ActionDefinition{
		{
			Name:        "search",
			Description: "Use a search engine",
		},
		{
			Name:        "loadURL",
			Description: "Load a URL, e.g. a website, formatted as a markdown document",
		},
	}
	agent := domain.NewAgent(domain.AgentID(1), "Agent", domain.Role("task_coordinator"), "Build a website", expected)

	// Act
	hasAction := agent.HasAction(expected[0].Name)
	doesntHaveAction := agent.HasAction("doesntExist")

	// Assert
	assert.Equal(t, true, hasAction)
	assert.Equal(t, false, doesntHaveAction)
}
