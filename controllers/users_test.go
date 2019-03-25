package controllers

import (
	"net/http"
	"testing"

	"chirp.com/models"
)

//usernames
const (
	samsmith       = "samsmith"
	kanye_west     = "kanye_west"
	duasings       = "duasings"
	bobbyd         = "bobbyd"
	tommyTesterton = "tommytesterton"
	vinceTester    = "vincetester"
	vincentXiao    = "vincent-xiao"
)

func TestUsers(t *testing.T) {
	services, router := getSetup()
	defer services.Close()

	ut := newUsersTester()
	testCases := ut.createTestCases()
	runAPITests(t, router, testCases)
}

type usersTester struct {
	users map[string]*models.User
}

func newUsersTester() *usersTester {
	ut := &usersTester{}
	ut.createUsers()
	return ut

}

func (ut *usersTester) createUsers() {
	ut.users = make(map[string]*models.User)
	ut.users[samsmith] = &models.User{
		Username: samsmith,
		Name:     "Sam Smith",
		Email:    "sam2018@gmail.com",
	}

	ut.users[kanye_west] = &models.User{
		Username: kanye_west,
		Name:     "Kanye West",
		Email:    "kanye@kanye.com",
	}
	ut.users[duasings] = &models.User{
		Username: duasings,
		Name:     "Dua Lipa",
		Email:    "dua@lipa.com",
	}
	ut.users[bobbyd] = &models.User{
		Username: bobbyd,
		Name:     "Bob Dylan",
		Email:    "bob@dylan.com",
	}
	ut.users["vincent-xiao"] = &models.User{
		Username: vincentXiao,
		Name:     "vince",
		Email:    "vincent@gmail.com",
	}
	ut.users[vinceTester] = &models.User{
		Username: vinceTester,
		Name:     "Vince Main",
		Email:    "vtester@gmail.com",
	}
}

func (ut *usersTester) createTestCases() (testCases []apiTestCase) {
	getUser := apiTestCase{
		tag:    "Get user's info",
		method: "GET",
		url:    "/kanye_west",
		status: http.StatusOK,
		want:   toMap(ut.users[kanye_west]),
	}

	signUpUser := apiTestCase{
		tag:    "sign up user",
		method: "POST",
		body: SignUpForm{
			Name:     "vince",
			Username: vincentXiao,
			Email:    "vincent@gmail.com",
			Password: "super-secret-password",
		},
		url:    "/signup",
		status: http.StatusOK,
		want:   toMap(ut.users[vincentXiao]),
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
		want:   toMap(ut.users[vincentXiao]),
	}
	logoutUser := apiTestCase{
		tag:      "logout user",
		method:   "POST",
		url:      "/logout",
		status:   http.StatusOK,
		remember: tokenAuthTesting,
	}

	userWithFollowers := *ut.users[bobbyd]
	userWithFollowers.Followers = []models.User{*ut.users[samsmith], *ut.users[kanye_west], *ut.users[vinceTester]}
	getFollowers := apiTestCase{
		tag:    "get user's followers",
		method: "GET",
		url:    "/bobbyd/followers",
		status: http.StatusOK,
		want:   toMap(userWithFollowers),
	}

	userWithFollowing := *ut.users[duasings]
	userWithFollowing.Following = []models.User{*ut.users[kanye_west]}
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
				User:       ut.users[duasings],
			},
		),
		remember: tokenUserRequired,
	}

	deleteFollow := apiTestCase{
		tag:      "unfollow user",
		method:   "POST",
		url:      "/bobbyd/follow/delete",
		status:   http.StatusOK,
		want:     toMap(ut.users[bobbyd]),
		remember: tokenUserRequired,
	}

	testCases = append(testCases,
		getUser,
		signUpUser,
		loginUser,
		logoutUser,
		getFollowers,
		getFollowing,
		followUser,
		deleteFollow,
	)
	return testCases
}
