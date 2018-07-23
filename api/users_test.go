package api

import (
	"net/http"
	"testing"

	"chirp.com/app"
	"chirp.com/middleware"
	"chirp.com/models"
	"chirp.com/testdata"
	"github.com/gorilla/mux"
)

func setUpTests(router *mux.Router) *models.Services {

	cfg := testdata.TestConfig
	services := app.Init(cfg)
	testdata.ResetDB(cfg)
	usersAPI := NewUsers(services.User, services.Like, services.Follow, nil)
	//init middleware
	userMw := middleware.NewUserMw(services.User)
	requireUserMw := middleware.NewRequireUserMw(userMw)
	ServeUserResource(router, usersAPI, &requireUserMw)
	return services
}

func TestUser(t *testing.T) {
	router := app.NewRouter()
	services := setUpTests(router)
	defer services.Close()
	userMw := middleware.NewUserMw(services.User)

	ut := &UsersTest{}
	ut.CreateUsers()
	ut.CreateTestCases()
	runAPITests(t, userMw.Apply(router),
		ut.testCases,
	)
}

type UsersTest struct {
	users     map[string]*models.User
	testCases []apiTestCase
}

func (ut *UsersTest) CreateUsers() {
	ut.users = make(map[string]*models.User)
	ut.users["samsmith"] = &models.User{
		Username: "samsmith",
		Name:     "Sam Smith",
		Email:    "sam2018@gmail.com",
	}

	ut.users["kanye_west"] = &models.User{
		Username: "kanye_west",
		Name:     "Kanye West",
		Email:    "kanye@kanye.com",
	}
	ut.users["duasings"] = &models.User{
		Username: "duasings",
		Name:     "Dua Lipa",
		Email:    "dua@lipa.com",
	}
	ut.users["bobbyd"] = &models.User{
		Username: "bobbyd",
		Name:     "Bob Dylan",
		Email:    "bob@dylan.com",
	}
	ut.users["vincent-xiao"] = &models.User{
		Username: "vincent-xiao",
		Name:     "vince",
		Email:    "vincent@gmail.com",
	}
	ut.users["mainTester"] = &models.User{
		Username: "mainTester",
		Name:     "Main Tester",
		Email:    "mtester@gmail.com",
	}
}

func (ut *UsersTest) CreateTestCases() {
	ut.testCases = make([]apiTestCase, 0)

	getUser := apiTestCase{
		tag:    "Get user's info",
		method: "GET",
		url:    "/kanye_west",
		status: http.StatusOK,
		want:   toMap(ut.users["kanye_west"]),
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
		want:   toMap(ut.users["vincent-xiao"]),
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
		want:   toMap(ut.users["vincent-xiao"]),
	}
	logoutUser := apiTestCase{
		tag:      "logout user",
		method:   "POST",
		url:      "/logout",
		status:   http.StatusOK,
		remember: tokenAuthTesting,
	}

	userWithFollowers := *ut.users["bobbyd"]
	userWithFollowers.Followers = []models.User{*ut.users["samsmith"], *ut.users["kanye_west"], *ut.users["mainTester"]}

	getFollowers := apiTestCase{
		tag:    "get user's followers",
		method: "GET",
		url:    "/bobbyd/followers",
		status: http.StatusOK,
		want:   toMap(userWithFollowers),
	}

	// userWithFollowing := copyMap(ut.users["duasings"])
	// userWithFollowing["following"] = []interface{}{ut.users["kanye_west"]}
	userWithFollowing := *ut.users["duasings"]
	userWithFollowing.Following = []models.User{*ut.users["kanye_west"]}
	getFollowing := apiTestCase{
		tag:    "get users that use is following",
		method: "GET",
		url:    "/duasings/following",
		status: http.StatusOK,
		want:   toMap(userWithFollowing),
	}

	followUser := apiTestCase{
		tag:    "follow user",
		method: "POST",
		url:    "/duasings/follow",
		status: http.StatusOK,
		want: toMap(
			models.Follow{
				FollowerID: 6,
				UserID:     3,
				User:       ut.users["duasings"],
			},
		),
		remember: tokenUserRequired,
	}

	deleteFollow := apiTestCase{
		tag:      "unfollow user",
		method:   "POST",
		url:      "/bobbyd/follow/delete",
		status:   http.StatusOK,
		want:     toMap(ut.users["bobbyd"]),
		remember: tokenUserRequired,
	}

	ut.testCases = append(ut.testCases,
		getUser,
		signUpUser,
		loginUser,
		logoutUser,
		getFollowers,
		getFollowing,
		followUser,
		deleteFollow,
	)

}
