package hexa

import (
	"errors"
	"github.com/Kamva/gutil"
	"github.com/Kamva/tracer"
)

type (
	// IDGenerator can generate fresh ID.
	IDGenerator func() ID

	// UserType is type of a user. possible values is :
	// guest: Use for guest users.
	// regular: Use for regular users of app (real registered users)
	// service: Use for service users (microservices,...)
	UserType string

	// User who sends request to the app (can be guest,regular user,service user,...)
	User interface {
		// Type specifies user's type (guest,regular,...)
		Type() UserType

		// Identifier returns user's identifier
		Identifier() ID

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
	}

	// user is default implementation of hexa User for real users.
	user struct {
		id       ID
		userType UserType
		email    string
		phone    string
		name     string
		username string
		isActive bool
		roles    []string
	}

	// UserExporterImporter export a user to json and then import it.
	exportedUser struct {
		ID       interface{} `json:"id"`
		Type     UserType    `json:"type"`
		Email    string      `json:"email"`
		Phone    string      `json:"phone"`
		Name     string      `json:"name"`
		Username string      `json:"username"`
		IsActive bool        `json:"is_active"`
		Roles    []string    `json:"roles"`
	}

	// UserExporterImporter export a user to json and then import it.
	UserExporterImporter interface {
		Export(user User) (Map, error)
		Import(exportedMap Map) (User, error)
	}
	// userExporterImporter implements the UserExporterImporter interface.
	userExporterImporter struct {
		idGenerator IDGenerator
	}

	// UserSDK is the user's kit to import, export and generate guest.
	UserSDK interface {
		UserExporterImporter
		// GenerateGuest returns new Guest User.
		NewGuest() User
	}

	// userSDK implements the UserSDK.
	userSDK struct {
		UserExporterImporter
	}
)

const (
	UserTypeGuest   UserType = "__guest__"
	UserTypeRegular UserType = "__regular__"
	UserTypeService UserType = "__service__"
)

// guestUserID is the guest user's id
var guestUserID = "__guest_id__"

func (u *user) Type() UserType {
	return u.userType
}

func (u *user) Identifier() ID {
	return u.id
}

func (u *user) Email() string {
	return u.email
}

func (u *user) Phone() string {
	return u.phone
}

func (u *user) Name() string {
	return u.name
}

func (u *user) Username() string {
	return u.email
}

func (u *user) IsActive() bool {
	return u.isActive
}

func (u *user) Roles() []string {
	return u.roles
}

// Export method export a user to map.
func (e *userExporterImporter) Export(user User) (Map, error) {
	if user == nil {
		return nil, tracer.Trace(errors.New("user can not be nil"))
	}
	return gutil.StructToMap(exportedUser{
		ID:       user.Identifier().Val(),
		Type:     user.Type(),
		Email:    user.Email(),
		Phone:    user.Phone(),
		Name:     user.Name(),
		Username: user.Username(),
		IsActive: user.IsActive(),
		Roles:    user.Roles(),
	}), nil
}

// Import method a user from map.
func (e *userExporterImporter) Import(exportedMap Map) (User, error) {
	eu := exportedUser{}
	err := gutil.MapToStruct(exportedMap, &eu)
	if err != nil {
		return nil, err
	}

	id := e.idGenerator()

	// By convention guest and service users use stringID type,
	// but for user we should generate, because each app can
	// uses specific database with specific implemented userID
	//type.
	if eu.Type == UserTypeGuest || eu.Type == UserTypeService {
		id = NewStringID(eu.ID.(string))
	} else {
		if err := id.From(eu.ID); err != nil {
			return nil, err
		}
	}

	user := NewUser(id, eu.Type, eu.Email, eu.Phone, eu.Name, eu.Username, eu.IsActive, eu.Roles)

	return user, nil
}

func (u *userSDK) NewGuest() User {
	return NewGuest()
}

// NewUser returns new hexa user instance.
func NewUser(id ID, utype UserType, email, phone, name, username string, isActive bool, roles []string) User {
	return &user{
		id:       id,
		userType: utype,
		email:    email,
		phone:    phone,
		name:     name,
		username: username,
		isActive: isActive,
		roles:    roles,
	}
}

// NewGuest returns new instance of guest user.
func NewGuest() User {
	email := ""
	phone := ""
	name := "__guest__"
	username := "__guest__username__"
	return NewUser(NewStringID(guestUserID), UserTypeGuest, email, phone, name, username, false, []string{})
}

// NewServiceUser returns new instance of Service user.
func NewServiceUser(id, name string, isActive bool, roles []string) User {
	email := ""
	phone := ""
	username := "__service_username__"
	return NewUser(NewStringID(id), UserTypeService, email, phone, name, username, isActive, roles)
}

// NewUserExporterImporter returns new instance of user exporter.
func NewUserExporterImporter(idGenerator IDGenerator) UserExporterImporter {
	return &userExporterImporter{idGenerator}
}

// NewUserSDK returns new instance of the user SDK.
func NewUserSDK(ei UserExporterImporter) UserSDK {
	return &userSDK{ei}
}

// Assertion
var _ User = &user{}
var _ UserExporterImporter = &userExporterImporter{}
var _ UserSDK = &userSDK{}
