//--------------------------------
// Worker contains interfaces to
// define workers to process background
// jobs.
//--------------------------------

package kitty

type (

	// JobHandlerFunc is the handler of each job in the worker.
	JobHandlerFunc func(Context, Payload) error

	// Job is the new to push to the queue by Jobs interface
	Job struct {
		Payload
		Name  string
		Queue string
		// Retry specify retry counts of the job.
		// 0: means that throw job away (and dont push to dead queue) on first fail.
		// -1: means that push job to the dead queue on first fail.
		Retry int
	}

	Payload struct {
		Header Map `json:"header"`
		Data   Map `json:"data"`
	}

	// Worker is the background jobs worker
	Worker interface {
		// Register handler for new job
		Register(name string, handlerFunc JobHandlerFunc) error

		// Set worker concurrency
		Concurrency(c int) error

		// start process on some queues
		Process(queues ...string) error
	}

	// Jobs pushes jobs to process by worker.
	Jobs interface {
		// Push push job to the default queue
		Push(Context, *Job) error
	}
)
