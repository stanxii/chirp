package api

import (
	"net/http"
	"testing"
	"time"

	"chirp.com/app"
	"chirp.com/middleware"
	"chirp.com/models"
	"chirp.com/testdata"
)

func TestUser(t *testing.T) {
	cfg := testdata.TestConfig
	dbCfg := cfg.Database
	testdata.ResetDB(cfg)
	services, err := models.NewServices(
		models.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		models.WithLogMode(!cfg.IsProd()),
		models.WithUser(cfg.Pepper, cfg.HMACKey),
	)
	if err != nil {
		t.Error(err)
	}
	defer services.Close()

	router := app.NewRouter()
	usersAPI := NewUsers(services.User, services.Like, services.Follow, nil)
	//init middleware
	userMw := middleware.NewUserMw(services.User)
	requireUserMw := middleware.NewRequireUserMw(userMw)

	ServeUserResource(router, usersAPI, &requireUserMw)
	ut := &UsersTest{}
	ut.CreateUsers()
	ut.CreateTestCases()
	runAPITests(t, router,

		ut.testCases,
	)
}

type UsersTest struct {
	users     map[int]*models.User
	testCases []apiTestCase
}

func (ut *UsersTest) CreateUsers() {
	ut.users = make(map[int]*models.User)

	ut.users[0] = &models.User{
		Username: "samsmith",
		Name:     "Sam Smith",
		Email:    "sam2018@gmail.com",
	}
	ut.users[1] = &models.User{
		Username: "kanye_west",
		Name:     "Kanye West",
		Email:    "kanye@kanye.com",
	}
	ut.users[2] = &models.User{
		Username: "duasings007",
		Name:     "Dua Lipa",
		Email:    "dua@lipa.com",
	}
	ut.users[3] = &models.User{
		Username: "bobbyd",
		Name:     "bob@dylan.com",
		Email:    "bob@dylan.com",
	}
	ut.users[4] = &models.User{
		Username: "vincent-xiao",
		Name:     "vince",
		Email:    "vincent@gmail.com",
	}
}

func (ut *UsersTest) CreateTestCases() {
	ut.testCases = make([]apiTestCase, 0)

	getUser := apiTestCase{
		tag:    "Get user's info",
		method: "GET",
		url:    "/kanye_west",
		status: http.StatusOK,
		got:    &models.User{},
		want:   ut.users[1],
	}

	signUpUser := apiTestCase{
		tag:    "sign up user",
		method: "POST",
		body: SignUpForm{
			Name:     "vince",
			Username: "vincent-xiao",
			Email:    "vincent@gmail.com",
			Password: "super-secret-password",
		},
		url:    "/signup",
		status: http.StatusOK,
		got:    &models.User{},
		want:   ut.users[4],
	}

	loginUser := apiTestCase{
		tag:    "login user",
		method: "POST",
		body: SignUpForm{
			Email:    "vincent@gmail.com",
			Password: "super-secret-password",
		},
		url:    "/login",
		status: http.StatusOK,
		got:    &models.User{},
		want:   ut.users[4],
	}

	ut.testCases = append(ut.testCases, getUser, signUpUser, loginUser)

}

func sanitize(u *models.User) *models.User {
	u.CreatedAt = time.Time{}
	u.UpdatedAt = time.Time{}
	u.DeletedAt = nil
	return u
}
