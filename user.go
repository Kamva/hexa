package hexa

import (
	"errors"
	"fmt"

	"github.com/kamva/gutil"
	"github.com/kamva/tracer"
)

// UserType is type of a user. possible values is :
// guest: Use for guest users.
// regular: Use for regular users of app (real registered users)
// service: Use for service users (microservices,...)
type UserType string

const (
	UserTypeGuest   UserType = "_guest"
	UserTypeRegular UserType = "_regular"
	UserTypeService UserType = "_service"
)

// User meta keys.
const (
	UserMetaKeyUserType   = "_user_type"
	UserMetaKeyIdentifier = "_user_identifier"
	UserMetaKeyEmail      = "_user_email"
	UserMetaKeyPhone      = "_user_phone"
	UserMetaKeyName       = "_user_name"
	UserMetaKeyUsername   = "_user_username"
	UserMetaKeyIsActive   = "_user_is_active"
	UserMetaKeyRoles      = "_user_roles"
)

var userMetaKeys = []string{
	UserMetaKeyUserType,
	UserMetaKeyIdentifier,
	UserMetaKeyEmail,
	UserMetaKeyPhone,
	UserMetaKeyName,
	UserMetaKeyUsername,
	UserMetaKeyIsActive,
	UserMetaKeyRoles,
}

// guestUserID is the guest user's id
const guestUserID = "_guest_id"

// User who sends request to the app (can be guest,regular user,service user,...)
type User interface {
	// Type specifies user's type (guest,regular,...)
	Type() UserType

	// Identifier returns user's identifier
	Identifier() string

	// Email returns the user's email.
	// This value can be empty.
	Email() string

	// Phone returns the user's phone number.
	// This value can be empty.
	Phone() string

	// Name returns the user name.
	Name() string

	// Username can be unique username,email,phone number or
	// everything else which can be used as username.
	Username() string

	// IsActive specify that user is active or no.
	IsActive() bool

	// Roles returns user's roles.
	Roles() []string

	Meta(key string) (val interface{}, exists bool)

	SetMeta(key string, val interface{}) (User, error)

	// User must be able be export and import using this meta data.
	// Meta data must be json serializable.
	MetaData() Map
}

// user is default implementation of hexa User for real users.
type user struct {
	meta map[string]interface{}
}

func (u *user) Meta(key string) (val interface{}, exists bool) {
	val, exists = u.meta[key]
	return
}

func (u *user) copyMeta() (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if err := gutil.UnmarshalStruct(u.meta, &m); err != nil {
		return nil, tracer.Trace(err)
	}
	return m, tracer.Trace(userMetaInterfaceToTrueTypedMeta(m))
}

func (u *user) SetMeta(key string, val interface{}) (User, error) {
	m, err := u.copyMeta()
	if err != nil {
		return nil, tracer.Trace(err)
	}

	m[key] = val

	if err := validateUserMetaData(m); err != nil {
		return nil, tracer.Trace(err)
	}

	return NewUserFromMeta(m)
}

func (u *user) MetaData() Map {
	m := make(map[string]interface{})
	for k, v := range u.meta {
		m[k] = v
	}
	return m
}

func (u *user) metaString(k string) string {
	return u.meta[k].(string)
}

func (u *user) Type() UserType {
	return u.meta[UserMetaKeyUserType].(UserType)
}

func (u *user) Identifier() string {
	return u.metaString(UserMetaKeyIdentifier)
}

func (u *user) Email() string {
	return u.metaString(UserMetaKeyEmail)
}

func (u *user) Phone() string {
	return u.metaString(UserMetaKeyPhone)
}

func (u *user) Name() string {
	return u.metaString(UserMetaKeyName)
}

func (u *user) Username() string {
	return u.metaString(UserMetaKeyUsername)
}

func (u *user) IsActive() bool {
	return u.meta[UserMetaKeyIsActive].(bool)
}

func (u *user) Roles() []string {
	return u.meta[UserMetaKeyRoles].([]string)
}

// NewUserFromMeta creates new user from the meta keys.
func NewUserFromMeta(meta Map) (User, error) {
	if err := validateUserMetaData(meta); err != nil {
		return nil, tracer.Trace(err)
	}

	return &user{meta: meta}, nil
}

func MustNewUserFromMeta(meta Map) User {
	u, err := NewUserFromMeta(meta)
	if err != nil {
		panic(err)
	}
	return u
}

type UserParams struct {
	Id       string
	Type     UserType
	Email    string
	Phone    string
	Name     string
	UserName string
	IsActive bool
	Roles    []string
}

// NewUser returns new hexa user instance.
func NewUser(p UserParams) User {
	meta := map[string]interface{}{
		UserMetaKeyIdentifier: p.Id,
		UserMetaKeyUserType:   p.Type,
		UserMetaKeyEmail:      p.Email,
		UserMetaKeyPhone:      p.Phone,
		UserMetaKeyName:       p.Name,
		UserMetaKeyUsername:   p.UserName,
		UserMetaKeyIsActive:   p.IsActive,
		UserMetaKeyRoles:      p.Roles,
	}
	return MustNewUserFromMeta(meta)
}

// NewGuest returns new instance of guest user.
func NewGuest() User {
	return NewUser(UserParams{
		Id:       guestUserID,
		Type:     UserTypeGuest,
		Email:    "",
		Phone:    "",
		Name:     "_guest",
		UserName: "_guest_username",
		IsActive: false,
		Roles:    []string{},
	})
}

// NewServiceUser returns new instance of Service user.
func NewServiceUser(id, name string, isActive bool, roles []string) User {
	return NewUser(UserParams{
		Id:       id,
		Type:     UserTypeService,
		Email:    "",
		Phone:    "",
		Name:     name,
		UserName: "_service_username",
		IsActive: isActive,
		Roles:    roles,
	})
}

func validateUserMetaData(meta map[string]interface{}) error {
	// validate meta keys: all required meta keys must exists.
	for _, k := range userMetaKeys {
		if _, ok := meta[k]; !ok {
			errStr := fmt.Sprintf("key %s not found in user's meta keys", k)
			return tracer.Trace(errors.New(errStr))
		}
	}

	// Validate userType field
	if _, ok := meta[UserMetaKeyUserType].(UserType); !ok {
		return tracer.Trace(errors.New("invalid type for usertype field in user's meta data"))
	}

	// Validate IsActive field
	if _, ok := meta[UserMetaKeyIsActive].(bool); !ok {
		return tracer.Trace(errors.New("invalid type for isActive field in user's meta data"))
	}

	// Validate roles field:
	if _, ok := meta[UserMetaKeyRoles].([]string); !ok {
		return tracer.Trace(errors.New("invalid type for roles field in user's meta data"))
	}

	return nil
}

func WithUserRole(u User, role string) User {
	roles := append(u.Roles(), role)
	return gutil.Must(u.SetMeta(UserMetaKeyRoles, roles)).(User)
}

func userMetaInterfaceToTrueTypedMeta(meta map[string]interface{}) error {

	meta[UserMetaKeyUserType] = UserType(meta[UserMetaKeyUserType].(string))

	// Convert user roles from []interface{} to []string:
	roles := make([]string, 0)
	err := gutil.UnmarshalStruct(meta[UserMetaKeyRoles], &roles)
	if err != nil {
		return tracer.Trace(err)
	}
	meta[UserMetaKeyRoles] = roles
	return nil
}

// Assertion
var _ User = &user{}
