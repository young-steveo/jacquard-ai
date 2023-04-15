package domain

type Project struct {
	Name   string  `json:"name"`
	Goal   string  `json:"goal"`
	Agents []Agent `json:"agents"`
}
