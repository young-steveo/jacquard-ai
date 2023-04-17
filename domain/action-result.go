package domain

type ActionResult struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
