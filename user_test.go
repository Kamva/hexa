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

func TestUser_SetMeta(t *testing.T) {
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

	newUser, err := u.SetMeta(UserMetaKeyIdentifier, "123")
	if assert.Nil(t, err) {
		assert.Equal(t, u.Identifier(), "abc")
		assert.Equal(t, newUser.Identifier(), "123")
	}
}

func TestUser_SetCopyMeta(t *testing.T) {
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

	newUser, err := u.SetMeta(UserMetaKeyIdentifier, "new_identifier")
	if !assert.Nil(t, err) {
		return
	}

	roles := u.Roles()
	newRoles := newUser.Roles()

	newRoles[0] = "new_role"
	assert.Equal(t, "a", roles[0])
	assert.NotEqual(t, roles, newRoles)
}

func TestWithUserRole(t *testing.T) {
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
	newUser := WithUserRole(u, "abc")

	assert.Equal(t, u.Roles(), []string{"a", "b"})
	assert.Equal(t, newUser.Roles(), []string{"a", "b","abc"})
}
