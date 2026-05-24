package hexa

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserPropagator_RoundTrip(t *testing.T) {
	p := NewUserPropagator()
	u := NewUser(UserParams{
		Id:       "1",
		Type:     UserTypeRegular,
		Email:    "a@b.com",
		Phone:    "+1",
		Name:     "n",
		UserName: "un",
		IsActive: true,
		Roles:    []string{"r1", "r2"},
	})

	b, err := p.ToBytes(u)
	require.NoError(t, err)

	got, err := p.FromBytes(b)
	require.NoError(t, err)

	assert.Equal(t, "1", got.Identifier())
	assert.Equal(t, UserTypeRegular, got.Type())
	assert.Equal(t, "a@b.com", got.Email())
	assert.True(t, got.IsActive())
	assert.Equal(t, []string{"r1", "r2"}, got.Roles())
}

func TestUserPropagator_FromInvalidBytes(t *testing.T) {
	_, err := NewUserPropagator().FromBytes([]byte("not json"))
	assert.Error(t, err)
}

func TestSessionFromContext_Absent(t *testing.T) {
	assert.Nil(t, SessionFromContext(context.Background()))
}
