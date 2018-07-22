package api

import (
	"fmt"
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
		models.WithFollow(),
		models.WithLike(),
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
	runAPITests(t, userMw.Apply(router),
		ut.testCases,
	)
}

type UsersTest struct {
	// users     map[string]*models.User
	users     map[string]map[string]interface{}
	testCases []apiTestCase
}

func (ut *UsersTest) CreateUsers() {
	// ut.users = make(map[string]*models.User)

	// ut.users["samsmith"] = &models.User{
	// 	Username: "samsmith",
	// 	Name:     "Sam Smith",
	// 	Email:    "sam2018@gmail.com",
	// }
	// ut.users["kanye_west"] = &models.User{
	// 	Username: "kanye_west",
	// 	Name:     "Kanye West",
	// 	Email:    "kanye@kanye.com",
	// }
	// ut.users["duasings007"] = &models.User{
	// 	Username: "duasings",
	// 	Name:     "Dua Lipa",
	// 	Email:    "dua@lipa.com",
	// }
	// ut.users["bobbyd"] = &models.User{
	// 	Username: "bobbyd",
	// 	Name:     "bob@dylan.com",
	// 	Email:    "bob@dylan.com",
	// }
	// ut.users["vincent-xiao"] = &models.User{
	// 	Username: "vincent-xiao",
	// 	Name:     "vince",
	// 	Email:    "vincent@gmail.com",
	// }
	ut.users = make(map[string]map[string]interface{})
	ut.users["samsmith"] = map[string]interface{}{
		"username": "samsmith",
		"name":     "Sam Smith",
		"email":    "sam2018@gmail.com",
	}

	ut.users["kanye_west"] = map[string]interface{}{
		"username": "kanye_west",
		"name":     "Kanye West",
		"email":    "kanye@kanye.com",
	}
	ut.users["duasings007"] = map[string]interface{}{
		"username": "duasings",
		"name":     "Dua Lipa",
		"email":    "dua@lipa.com",
	}
	ut.users["bobbyd"] = map[string]interface{}{
		"username": "bobbyd",
		"name":     "Bob Dylan",
		"email":    "bob@dylan.com",
	}
	ut.users["vincent-xiao"] = map[string]interface{}{
		"username": "vincent-xiao",
		"name":     "vince",
		"email":    "vincent@gmail.com",
	}
}

func (ut *UsersTest) CreateTestCases() {
	ut.testCases = make([]apiTestCase, 0)

	getUser := apiTestCase{
		tag:    "Get user's info",
		method: "GET",
		url:    "/kanye_west",
		status: http.StatusOK,
		want:   ut.users["kanye_west"],
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
		want:   ut.users["vincent-xiao"],
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
		want:   ut.users["vincent-xiao"],
	}
	logoutUser := apiTestCase{
		tag:      "logout user",
		method:   "POST",
		url:      "/logout",
		status:   http.StatusOK,
		loggedIn: true,
	}
	userWithFollowers := copyMap(ut.users["bobbyd"])
	userWithFollowers["followers"] = []interface{}{ut.users["samsmith"],
		ut.users["kanye_west"]}
	getFollowers := apiTestCase{
		tag:    "get user's followers",
		method: "GET",
		url:    "/bobbyd/followers",
		status: http.StatusOK,
		want:   userWithFollowers,
	}
	fmt.Println(ut.users["bobbyd"])

	ut.testCases = append(ut.testCases, getUser, signUpUser, loginUser, logoutUser, getFollowers)

}

func sanitize(u *models.User) *models.User {
	u.CreatedAt = time.Time{}
	u.UpdatedAt = time.Time{}
	u.DeletedAt = nil
	return u
}

func copyMap(original map[string]interface{}) map[string]interface{} {
	newMap := map[string]interface{}{}
	for k, v := range original {
		newMap[k] = v
	}
	return newMap
}
