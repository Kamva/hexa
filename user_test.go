package hexa

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUser(t *testing.T) {
	p := UserParams{
		Id:       "abc",
		Type:     UserTypeRegular,
		Email:    "a@b.com",
		Phone:    "+123",
		Name:     "def",
		UserName: "abcd",
		IsActive: true,
		Roles: []string{
			"a", "b",
		},
	}
	u := NewUser(p)
	if !assert.NotNil(t, u) {
		return
	}

	assert.Equal(t, u.Identifier(), p.Id)
	assert.Equal(t, u.Type(), p.Type)
	assert.Equal(t, u.Email(), p.Email)
	assert.Equal(t, u.Phone(), p.Phone)
	assert.Equal(t, u.Name(), p.Name)
	assert.Equal(t, u.Username(), p.UserName)
	assert.Equal(t, u.IsActive(), p.IsActive)
	assert.Equal(t, u.Roles(), p.Roles)
}
