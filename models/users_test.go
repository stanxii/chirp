package models

import (
	"fmt"
	"testing"

	"chirp.com/pkg/hash"
	"chirp.com/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserDB(t *testing.T) {
	cfg := testdata.TestConfig
	dbCfg := cfg.Database
	testdata.ResetDB(cfg)

	services, err := NewServices(
		WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		WithLogMode(!cfg.IsProd()),
		// WithUser(cfg.Pepper, cfg.HMACKey),
	)
	if err != nil {
		t.Error(err)
	}
	defer services.Close()
	userDB := &userGorm{services.db}
	user, err := userDB.ByEmail("sam2018@gmail.com")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(user)
}

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
	// callsSpy []string
	mock.Mock
}

func (u *userDBMock) ByID(id uint) (*User, error) {
	// u.callsSpy = append(u.callsSpy, byID)
	return u.user, nil
}
func (u *userDBMock) ByEmail(email string) (*User, error) {
	// u.callsSpy = append(u.callsSpy, byEmail)
	u.user.Email = email
	u.Called()
	return u.user, nil
}
func (u *userDBMock) ByUsername(username string) (*User, error) {
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
	cfg := testdata.TestConfig
	hmac := hash.NewHMAC(cfg.HMACKey)
	return newUserValidator(mockDB, hmac, cfg.Pepper)
}

func TestByEmailValidator(t *testing.T) {
	mockDB := newUserDBMock(nil)
	uv := newUserValidatorWithMock(mockDB)
	dbCalls := 0

	tests := map[string]struct {
		email string
		want  interface{}
	}{
		"no uppercase":    {email: "Hey@GMAIL.coM", want: "hey@gmail.com"},
		"no spaces":       {email: "  nospace @yahoo.com ", want: "nospace@yahoo.com"},
		"no empty string": {email: "", want: ErrEmailRequired},
	}
	for name, test := range tests {
		var got interface{}
		t.Run(name, func(t *testing.T) {
			mockDB.On("ByEmail")
			user, err := uv.ByEmail(test.email)
			if err != nil {
				got = err
			} else {
				dbCalls++
				got = user.Email
			}
			assert.Equal(t, test.want, got)
		})
	}
	mockDB.AssertNumberOfCalls(t, "ByEmail", dbCalls)
}
