package domain

type Agent struct {
	ID        uint32 `json:"id"`
	Name      string `json:"name"`
	Role      Role   `json:"role"`
	Objective string `json:"objective"`

	// can the agennt search google?
	CanSearch bool `json:"can_search"`

	// can the agent access URLs?
	CanBrowse bool `json:"can_browse"`

	// can the agent read files?
	CanReadFiles bool `json:"can_read_files"`

	// can the agent write files?
	CanWriteFiles bool `json:"can_write_files"`
}
