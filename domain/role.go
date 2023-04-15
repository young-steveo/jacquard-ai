package domain

type Role string

const (

	// task coordinators are responsible for prioritizing the project backlog
	// there can be only one task coordinator per project
	TaskCoordinator Role = "task_coordinator"

	// execution specialists are responsible for executing the tasks at the top of the project backlog
	// there can be several execution specialists per project
	ExecutionSpecialist Role = "execution_specialist"

	// product managers are responsible for evaluating the work completed by the execution specialists
	// and generating new tasks for the task coordinators
	// there can be several product managers per project
	ProductManager Role = "product_manager"
)
