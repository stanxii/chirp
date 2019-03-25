package models

import (
	"chirp.com/config"
	"chirp.com/pkg/hash"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
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
	u.user.RememberHash = token
	u.Called()
	return u.user, nil
}
func (u *userDBMock) Create(user *User) error {
	u.user = user
	u.Called()
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

type testCase struct {
	tag     string //short description of test
	gotErr  error
	wantErr error
}

type stringTestCase struct {
	tc    testCase
	input string
	got   string
	want  string
}

type userTestCase struct {
	tc    testCase
	input *User
	got   User
	want  User
}

func setupUserValTests(user *User) (*userDBMock, *userValidator) {
	mockDB := newUserDBMock(nil)
	uv := newUserValidatorWithMock(mockDB)
	return mockDB, uv
}

func TestByEmailValidator(t *testing.T) {
	tests := []stringTestCase{
		{tc: testCase{tag: "no uppercase"},
			input: "Hey@GMAIL.coM",
			want:  "hey@gmail.com"},
		{tc: testCase{tag: "no spaces"},
			input: " nospace @yahoo.com ",
			want:  "nospace@yahoo.com"},
		{tc: testCase{tag: "no empty string", wantErr: ErrEmailRequired}, input: ""},
	}

	mockDB, uv := setupUserValTests(nil)
	dbCalls := 0
	for _, test := range tests {
		t.Run(test.tc.tag, func(t *testing.T) {
			mockDB.On(ByEmail)
			user, err := uv.ByEmail(test.input)
			if err != nil {
				test.tc.gotErr = err
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
	tests := []stringTestCase{
		{tc: testCase{tag: "no uppercase"}, input: "KOBE4EVER", want: "kobe4ever"},
		{tc: testCase{tag: "no spaces"}, input: " le bron fan  ", want: "lebronfan"},
		{tc: testCase{tag: "doesnt begin with letter", wantErr: ErrUsernameDoesntBeginWithLetter}, input: "5guys"},
		{tc: testCase{tag: "not enough characters", wantErr: ErrCharMin.customCharLimitError(3, "username")}, input: "ab"},
		{tc: testCase{tag: "max characters exceeded", wantErr: ErrCharMax.customCharLimitError(25, "username")}, input: "abcdefghijklmonpqrustuvwxyz2018"},
	}
	mockDB, uv := setupUserValTests(nil)
	dbCalls := 0
	for _, test := range tests {
		t.Run(test.tc.tag, func(t *testing.T) {
			mockDB.On(ByUsername)
			user, err := uv.ByUsername(test.input)
			if err != nil {
				test.tc.gotErr = err
				assert.Equal(t, test.tc.wantErr, test.tc.gotErr)
			} else {
				test.got = user.Username
				dbCalls++
				assert.Equal(t, test.want, test.got)
			}
		})
	}
	mockDB.AssertNumberOfCalls(t, ByUsername, dbCalls)
}

func TestByRememberValidator(t *testing.T) {
	tests := []stringTestCase{
		{tc: testCase{tag: "remember token gets hashed"},
			input: "fakeInputToken_123", want: "BHQIyJPodP2exJnuAIukgXRkqENICrNJr7bQc1JBO9I="},
	}
	mockDB, uv := setupUserValTests(nil)
	dbCalls := 0
	for _, test := range tests {
		t.Run(test.tc.tag, func(t *testing.T) {
			mockDB.On(ByRemember)
			user, err := uv.ByRemember(test.input)
			if err != nil {
				test.tc.gotErr = err
				assert.Equal(t, test.tc.wantErr, test.tc.gotErr)
			} else {
				test.got = user.RememberHash
				dbCalls++
				assert.Equal(t, test.want, test.got)
			}
		})
	}
	mockDB.AssertNumberOfCalls(t, ByRemember, dbCalls)
}

func newDefaultTestUser() *User {
	return &User{
		Name:         "George Hill",
		Username:     "hill_9000",
		Email:        "george@hill.com",
		Password:     "12345678",
		Remember:     "mockmockfakerememberhash12345678fakefakefakefake",
		RememberHash: "A5VkvtUNRxYLNWDQIr_L1IzEA6S7k2Xrf3kVjkRJcxU=",
	}
}

func TestCreateUserValidator(t *testing.T) {
	tests := []userTestCase{}
	mockDB, uv := setupUserValTests(nil)
	dbCalls := 0

	mockDB.On(ByEmail).Return(newDefaultTestUser())

	noPw := newDefaultTestUser()
	noPw.Password = ""
	tests = append(tests, userTestCase{
		tc: testCase{
			tag:     "password required",
			wantErr: ErrPasswordRequired,
		},
		input: noPw,
	})

	shortPw := newDefaultTestUser()
	shortPw.Password = "123ab"
	tests = append(tests, userTestCase{
		tc: testCase{
			tag:     "password required",
			wantErr: ErrPasswordTooShort},
		input: shortPw,
	})

	for _, test := range tests {
		t.Run(test.tc.tag, func(t *testing.T) {
			mockDB.On(Create)

			err := uv.Create(test.input)
			if err != nil {
				test.tc.gotErr = err
				assert.Equal(t, test.tc.wantErr, test.tc.gotErr)
			} else {
				fmt.Printf("struct: %+v\n", test.input)
				test.got = *test.input
				test.want.PasswordHash = test.got.PasswordHash
				dbCalls++
				assert.Equal(t, test.want, test.got)
			}
		})
	}

	mockDB.AssertNumberOfCalls(t, Create, dbCalls)
}
