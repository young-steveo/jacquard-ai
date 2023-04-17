package domain

type Project struct {
	Name             string  `json:"name"`
	Goal             string  `json:"goal"`
	WorkingDirectory string  `json:"working_directory"`
	Agents           []Agent `json:"agents"`
}
