package domain

type AgentID uint32

type Agent struct {
	ID        AgentID            `json:"id"`
	Name      string             `json:"name"`
	Role      Role               `json:"role"`
	Objective string             `json:"objective"`
	Actions   []ActionDefinition `json:"actions"`
}

func NewAgent(id AgentID, name string, role Role, objective string, actions []ActionDefinition) Agent {
	return Agent{
		ID:        id,
		Name:      name,
		Role:      role,
		Objective: objective,
		Actions:   actions,
	}
}

func (a *Agent) AddAction(action ActionDefinition) {
	a.Actions = append(a.Actions, action)
}

func (a *Agent) RemoveAction(name string) {
	for i, action := range a.Actions {
		if action.Name == name {
			a.Actions = append(a.Actions[:i], a.Actions[i+1:]...)
			return
		}
	}
}

func (a *Agent) HasAction(name string) bool {
	for _, action := range a.Actions {
		if action.Name == name {
			return true
		}
	}
	return false
}
