package models

import (
	"testing"

	"chirp.com/config"
	"chirp.com/pkg/hash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	ByID       = "ByID"
	ByEmail    = "ByEmail"
	ByUsername = "ByUsername"
	ByRemember = "ByRemember"
	Create     = "Create"
	Update     = "Update"
	Delete     = "Delete"
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
	tag     string //short description of test
	input   interface{}
	got     interface{}
	want    interface{}
	wantErr error
}

func setupUserValTests(user *User) (*userDBMock, *userValidator) {
	mockDB := newUserDBMock(nil)
	uv := newUserValidatorWithMock(mockDB)
	return mockDB, uv

}

// func runValidatorTestCases(t *testing.T, valFunc userValFunc, testCases ...validatorTestCase) {
// 	mockDB, uv := setupUserValTests(nil)
// 	dbCalls := 0
// 	for _, test := range testCases {
// 		t.Run(test.tag, func(t *testing.T) {
// 			mockDB.On(ByEmail)
// 			user, err := uv.ByEmail(test.input.(string))
// 			if err != nil {
// 				test.got = err
// 			} else {
// 				test.got = user.Email
// 				dbCalls++
// 			}
// 			assert.Equal(t, test.want, test.got)
// 		})
// 	}
// 	mockDB.AssertNumberOfCalls(t, ByEmail, dbCalls)
// }

func TestByEmailValidator(t *testing.T) {
	tests := []validatorTestCase{
		{tag: "no uppercase", input: "Hey@GMAIL.coM", want: "hey@gmail.com"},
		{tag: "no spaces", input: "  nospace @yahoo.com ", want: "nospace@yahoo.com"},
		{tag: "no empty string", input: "", want: ErrEmailRequired},
	}

	mockDB, uv := setupUserValTests(nil)
	dbCalls := 0
	for _, test := range tests {
		t.Run(test.tag, func(t *testing.T) {
			mockDB.On(ByEmail)
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
	mockDB.AssertNumberOfCalls(t, ByEmail, dbCalls)
}

func TestByUsernameValidator(t *testing.T) {
	tests := []validatorTestCase{
		{tag: "no uppercase", input: "KOBE4EVER", want: "kobe4ever"},
		{tag: "no spaces", input: " le bron fan  ", want: "lebronfan"},
		{tag: "doesnt being with non-letter", input: "5guys", wantErr: ErrUsernameNoLetter},
		{tag: "not enough characters", input: "ab", wantErr: ErrCharMin.customCharLimitError(3, "username")},
		{tag: "max characters exceeded", input: "abcdefghijklmonpqrustuvwxyz", wantErr: ErrCharMax.customCharLimitError(25, "username")},
		// {tag: "", input: "5g", err: ErrCharMin.customCharLimitError(3, "username")},

	}
	mockDB, uv := setupUserValTests(nil)
	dbCalls := 0
	for _, test := range tests {
		t.Run(test.tag, func(t *testing.T) {
			mockDB.On(ByUsername)
			user, err := uv.ByUsername(test.input.(string))
			if err != nil {
				test.got = err
			} else {
				test.got = user.Username
				dbCalls++
			}

			if test.wantErr != nil {
				assert.Equal(t, test.wantErr, test.got)
			} else {
				assert.Equal(t, test.want, test.got)
			}
		})
	}
	mockDB.AssertNumberOfCalls(t, ByUsername, dbCalls)
}
