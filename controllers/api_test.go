package controllers

import (
	"bytes"
	"chirp.com/app"
	"chirp.com/config"
	"chirp.com/middleware"
	"chirp.com/models"
	"chirp.com/testdata"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

/*
Struct for creating test cases
 */
type apiTestCase struct {
	tag      string
	method   string
	url      string
	body     interface{}
	status   int
	remember string //sets remember_token on http cookie
	got      map[string]interface{}
	want     interface{}
}

/*
Tests the given API endpoint
 */
func testAPI(router http.Handler, method, URL string, body interface{}, remember string) *httptest.ResponseRecorder {
	var bodyBytes []byte
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		bodyBytes = b
	}
	req, _ := http.NewRequest(method, URL, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	if remember != "" {

		cookie := http.Cookie{
			Name:     "remember_token",
			Value:    remember,
			HttpOnly: true,
		}
		req.AddCookie(&cookie)
	}
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res
}

func runAPITests(t *testing.T, router http.Handler, tests []apiTestCase) {
	for _, test := range tests {
		t.Run(test.tag, func(t *testing.T) {
			res := testAPI(router, test.method, test.url, test.body, test.remember)
			assert.Equal(t, test.status, res.Code, test.tag)
			if test.want == nil {
				return
			}
			dec := json.NewDecoder(res.Body)
			var got interface{}
			var err error
			if res.Body.String()[0] == '[' {
				array := make([]map[string]interface{}, 0)
				err = dec.Decode(&array)
				for _, v := range array {
					deleteUnwantedFields(v)
				}
				got = array
			} else {
				err = dec.Decode(&test.got)
				deleteUnwantedFields(test.got)
				got = test.got
			}
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, test.want, got, test.tag)

		})
	}
}

func toMap(obj interface{}, removeFields ...string) map[string]interface{} {
	var inInterface map[string]interface{}
	inrec, err := json.Marshal(obj)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(inrec, &inInterface)
	deleteUnwantedFields(inInterface, removeFields...)
	return inInterface
}

func deleteFields(m map[string]interface{}, fields ...string) {
	for _, field := range fields {
		delete(m, field)
	}
}

func deleteUnwantedFields(m map[string]interface{}, fields ...string) {
	s := []string{"created_at", "updated_at", "deleted_at"}
	s = append(s, fields...)
	deleteFields(m, s...)
	for _, v := range m {
		switch x := v.(type) {
		case []interface{}:
			for _, value := range x {
				v, ok := value.(map[string]interface{})
				if ok {
					deleteFields(v, s...)
				}
			}
		case map[string]interface{}:
			deleteFields(x, s...)
		default:
			// fmt.Printf("Unsupported type: %T\n", x)
		}
	}
}

func getSetup() (*models.Services, http.Handler) {
	router := app.NewRouter()
	cfg := config.TestConfig()
	services := app.Setup(cfg)
	testdata.ResetDB(cfg)
	usersAPI := NewUsers(services.User, services.Like, services.Follow, services.Tweet, nil)
	tweetsAPI := NewTweets(services.Tweet, services.Like, services.Tag, services.Tagging)
	tagsAPI := NewTags(services.Tag, services.Tagging)
	//init middleware
	userMw := middleware.NewUserMw(services.User)
	requireUserMw := middleware.NewRequireUserMw(userMw)
	ServeUserResource(router, usersAPI, &requireUserMw)
	ServeTweetResource(router, tweetsAPI, &requireUserMw)
	ServeTagResource(router, tagsAPI, &requireUserMw)
	return services, userMw.Apply(router)
}

/*
WARNING: Authentication for api tests must be reconfigured!
 */
const (
	tokenAuthTesting  = ""
	tokenUserRequired = ""
)
