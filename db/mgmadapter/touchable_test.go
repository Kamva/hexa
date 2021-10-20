package mgmadapter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type UserWithVersion struct {
	Entity           `bson:",inline"`
	VersionTouchable `bson:",inline"`
}

type UserWithTime struct {
	Entity        `bson:",inline"`
	TimeTouchable `bson:",inline"`
}

type User struct {
	Entity    `bson:",inline"`
	Touchable `bson:",inline"`
}

func TestVersionTouchable_Touch(t *testing.T) {
	u := &UserWithVersion{}
	assert.Equal(t, int64(0), u.Version)

	u.Touch()
	assert.Equal(t, int64(1), u.Version)

	u.Touch()
	assert.Equal(t, int64(1), u.Version)
}

func TestTimeTouchable_TouchAt(t *testing.T) {
	u := &UserWithTime{}
	assert.Equal(t, time.Time{}.UnixNano(), u.CreatedAt.UnixNano())
	assert.Equal(t, time.Time{}.UnixNano(), u.UpdatedAt.UnixNano())

	now := time.Now()
	u.TouchAt(now)
	assert.Equal(t, now.UnixNano(), u.CreatedAt.UnixNano())
	assert.Equal(t, now.UnixNano(), u.UpdatedAt.UnixNano())

	nowAgain := time.Now()
	u.TouchAt(nowAgain)
	assert.Equal(t, now.UnixNano(), u.CreatedAt.UnixNano())
	assert.Equal(t, nowAgain.UnixNano(), u.UpdatedAt.UnixNano())
}

func TestTouchable_TouchAt(t *testing.T) {
	u := &User{}
	assert.Equal(t, int64(0), u.Version)
	assert.Equal(t, time.Time{}.UnixNano(), u.CreatedAt.UnixNano())
	assert.Equal(t, time.Time{}.UnixNano(), u.UpdatedAt.UnixNano())

	now := time.Now()
	u.TouchAt(now)
	assert.Equal(t, int64(1), u.Version)
	assert.Equal(t, now.UnixNano(), u.CreatedAt.UnixNano())
	assert.Equal(t, now.UnixNano(), u.UpdatedAt.UnixNano())

	nowAgain := time.Now()
	u.TouchAt(nowAgain)
	assert.Equal(t, int64(1), u.Version)
	assert.Equal(t, now.UnixNano(), u.CreatedAt.UnixNano())
	assert.Equal(t, nowAgain.UnixNano(), u.UpdatedAt.UnixNano())
}

func TestTouchable_Touch(t *testing.T) {
	u := &User{}
	assert.Equal(t, int64(0), u.Version)
	assert.Equal(t, time.Time{}.UnixNano(), u.CreatedAt.UnixNano())
	assert.Equal(t, time.Time{}.UnixNano(), u.UpdatedAt.UnixNano())

	now := time.Now()
	u.Touch()
	assert.Equal(t, int64(1), u.Version)
	assert.True(t, now.UnixNano() <= u.CreatedAt.UnixNano())
	assert.True(t, now.UnixNano() <= u.UpdatedAt.UnixNano())

	nowAgain := time.Now()
	u.TouchAt(nowAgain)
	assert.Equal(t, int64(1), u.Version)
	assert.True(t, now.UnixNano() <= u.CreatedAt.UnixNano())
	assert.True(t, nowAgain.UnixNano() >= u.CreatedAt.UnixNano())
	assert.True(t, nowAgain.UnixNano() <= u.UpdatedAt.UnixNano())
}
