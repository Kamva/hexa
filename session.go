package hexa

import (
	"context"
	"time"
)

//--------------------------------
// Important: Session is prototype.
//--------------------------------

const SessionContextKey = "_ctx_session"

type SessionProvider interface {
	// Get returns a session if it exists.
	// Sometimes we need to single store per session (e.g., cookieSessions), and
	// sometimes single store per all sessions (e.g., redis), by this pattern
	// (providing store at creation time) we support both.
	Get(store SessionStore, id string) (Session, error)

	// GetOrNew will returns old session or creates a new one with the provided ID.
	// id could be empty if you want to get just a new one session.
	GetOrNew(store SessionStore, id string) (Session, error)

	// Copy copies the session with a new ID.
	Copy() (Session, error)
}

type SessionStore interface {
	Save(ctx context.Context, sess ...Session) error
}

type Session interface {
	// SessionID returns the session id.
	SessionID() string

	// Expiry returns the session expiry time.
	Expiry() time.Time

	// SetExpiry sets the session expiry time. to delete a
	// session just set its expiry in before and save it.
	SetExpiry(at time.Time) error

	// Get returns the session value. it must be nil
	// if isn't found.
	Get(key string) (interface{}, error)

	// Set sets the key and value in the session. nil
	// value meaning delete the key from the session.
	Set(key string, val interface{}) error

	Save(ctx context.Context) error

	// TODO: Add flash methods too (set and delete
	//  after one read), like gorilla sessions.
}

// SessionFromContext extracts the session from the context.
// It will be nil if not found in the context.
func SessionFromContext(ctx context.Context) Session {
	sess, _ := ctx.Value(SessionContextKey).(Session)
	return sess
}
