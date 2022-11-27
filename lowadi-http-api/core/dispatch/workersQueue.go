package dispatch

/*
type Dispatcher interface {
	// Push takes an Event and pushes to a queue.
	Push(Event) error
	// Run spawns the workers and waits indefinitely for
	// the events to be processed.
	Run()
}

// EventDispatcher represents the datastructure for an
// EventDispatcher instance. This struct satisfies the
// Dispatcher interface.
type EventDispatcher struct {
	Opts     Options
	Queue    chan models.Notification
	Finished bool
}

// Options represent options for EventDispatcher.
type Options struct {
	MaxWorkers   int // Number of workers to spawn.
	MaxQueueSize int // Maximum length for the queue to hold events.
}

// NewEventDispatcher initialises a new event dispatcher.
func NewEventDispatcher(opts Options) Dispatcher {
	return EventDispatcher{
		Opts:     opts,
		Queue:    make(chan Event, opts.MaxQueueSize),
		Finished: false,
	}
}
*/
