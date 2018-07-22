package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type apiTestCase struct {
	tag      string
	method   string
	url      string
	body     interface{}
	status   int
	loggedIn bool //sets remember_token on http cookie
	got      map[string]interface{}
	want     interface{}
}

func testAPI(router http.Handler, method, URL string, body interface{}, loggedIn bool) *httptest.ResponseRecorder {
	var bodyBytes []byte
	if body != nil {
		b, err := json.Marshal(body)
		fmt.Println(string(b))

		if err != nil {
			panic(err)
		}
		bodyBytes = b
	}
	req, _ := http.NewRequest(method, URL, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	if loggedIn {
		cookie := http.Cookie{
			Name:     "remember_token",
			Value:    "ke3kO2KwD4HjC2lqhYWDD17T3aKXanDN1qiMLQLq1LI=",
			HttpOnly: true,
		}
		req.AddCookie(&cookie)
	}
	fmt.Printf("req: %+v\n", req)

	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res
}

func runAPITests(t *testing.T, router http.Handler, tests []apiTestCase) {
	for _, test := range tests {
		t.Run(test.tag, func(t *testing.T) {

			res := testAPI(router, test.method, test.url, test.body, test.loggedIn)

			if test.want != nil {
				dec := json.NewDecoder(res.Body)
				err := dec.Decode(&test.got)
				if err != nil {
					t.Error(err)
				}
				//delete all json fields we want to ignore
				delete(test.got, "created_at")
				delete(test.got, "updated_at")
				delete(test.got, "deleted_at")
				for k, v := range test.got {
					switch x := v.(type) {
					case []interface{}:
						for _, value := range x {
							v, ok := value.(map[string]interface{})
							if ok {
								delete(v, "created_at")
								delete(v, "updated_at")
								delete(v, "deleted_at")
							}
						}
					default:
						fmt.Printf("Unsupported type: %T\n", x)
					}

					fmt.Printf("\nkey[%s] value[%s]\n", k, v)
				}

				assert.Equal(t, test.want, test.got, test.tag)
			}

			assert.Equal(t, test.status, res.Code, test.tag)
		})
	}
}
