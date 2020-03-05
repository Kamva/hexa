//--------------------------------
// Event file contains the Event
// interfaces to send and receive
// messages.
//--------------------------------
package kitty

type (
	Emitter interface {
		Emit(event Event) error
		Close() error
	}

	Event struct{
		// TODO: set event props.
	}


)
