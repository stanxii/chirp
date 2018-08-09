package models

import (
	"testing"

	"chirp.com/config"
	"chirp.com/pkg/hash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	byID       = "byId"
	byEmail    = "byEmail"
	byUsername = "byUsername"
	byRemember = "byRemember"
	create     = "create"
	update     = "update"
	delete     = "delete"
)

// type UserDB
type userDBMock struct {
	user *User
	mock.Mock
}

func (u *userDBMock) ByID(id uint) (*User, error) {
	u.user.ID = id
	u.Called()
	return u.user, nil
}
func (u *userDBMock) ByEmail(email string) (*User, error) {
	u.user.Email = email
	u.Called()
	return u.user, nil
}
func (u *userDBMock) ByUsername(username string) (*User, error) {
	u.user.Username = username
	u.Called()
	return u.user, nil
}
func (u *userDBMock) ByRemember(token string) (*User, error) {
	return u.user, nil
}
func (u *userDBMock) Create(user *User) error {
	return nil
}
func (u *userDBMock) Update(user *User) error {
	return nil
}
func (u *userDBMock) Delete(id uint) error {
	return nil
}

func newUserDBMock(u *User) *userDBMock {
	var user *User
	if user != nil {
		user = u
	} else {
		user = &User{}
	}
	return &userDBMock{
		user: user,
	}

}
func newUserValidatorWithMock(mockDB *userDBMock) *userValidator {
	cfg := config.TestConfig()
	hmac := hash.NewHMAC(cfg.HMACKey)
	return newUserValidator(mockDB, hmac, cfg.Pepper)
}

type validatorTestCase struct {
	tag   string //short description of test
	input interface{}
	got   interface{}
	want  interface{}
}

func runValidatorTestCases(t *testing.T, testCases ...validatorTestCase) {
	mockDB := newUserDBMock(nil)
	uv := newUserValidatorWithMock(mockDB)
	dbCalls := 0
	for _, test := range testCases {
		t.Run(test.tag, func(t *testing.T) {
			mockDB.On("ByEmail")
			user, err := uv.ByEmail(test.input.(string))
			if err != nil {
				test.got = err
			} else {
				test.got = user.Email
				dbCalls++
			}
			assert.Equal(t, test.want, test.got)
		})
	}
	mockDB.AssertNumberOfCalls(t, "ByEmail", dbCalls)
}

func TestByEmailValidator(t *testing.T) {
	tests := []validatorTestCase{
		{tag: "no uppercase", input: "Hey@GMAIL.coM", want: "hey@gmail.com"},
		{tag: "no spaces", input: "  nospace @yahoo.com ", want: "nospace@yahoo.com"},
		{tag: "no empty string", input: "", want: ErrEmailRequired},
	}
	runValidatorTestCases(t, tests...)
}

func TestByUsernameValidator(t *testing.T) {
	tests := []validatorTestCase{
		{tag: "no uppercase", input: "Hey@GMAIL.coM", want: "hey@gmail.com"},
		{tag: "no whitespace", input: "Hey@GMAIL.coM", want: "hey@gmail.com"},
		{tag: "no whitespace", input: "Hey@GMAIL.coM", want: "hey@gmail.com"},
		{tag: "no whitespace", input: "Hey@GMAIL.coM", want: "hey@gmail.com"},
	}
	runValidatorTestCases(t, tests...)
}
