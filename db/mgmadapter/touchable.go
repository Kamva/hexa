//go:generate msgp

package mgmadapter

import "time"

type Touchable struct {
	VersionTouchable `bson:",inline"`
	TimeTouchable    `bson:",inline"`
}

type VersionTouchable struct {
	Version   int64 `json:"v" bson:"v"`
	isTouched bool
}

type TimeTouchable struct {
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// Touch is idempotent.
func (t *Touchable) Touch() {
	t.VersionTouchable.Touch()
	t.TimeTouchable.Touch()
}

// TouchAt is idempotent.
func (t *Touchable) TouchAt(at time.Time) {
	t.VersionTouchable.Touch()
	t.TimeTouchable.TouchAt(at)
}

// Touch is idempotent, so you can call it multiple times.
func (t *VersionTouchable) Touch() {
	if t.isTouched {
		return
	}

	t.Version++
	t.isTouched = true
}

func (t *TimeTouchable) Touch() {
	t.TouchAt(time.Now())
}

// TouchAt is idempotent, so you can call it multiple time.
func (t *TimeTouchable) TouchAt(at time.Time) {
	if t.CreatedAt.IsZero() {
		t.CreatedAt = at
	}
	t.UpdatedAt = at
}
