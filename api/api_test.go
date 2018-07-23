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

const (
	tokenAuthTesting  = "CxblHnxYhskXGkVmbwbwPF8_i4fKXH4JWHY-qKzgLfE="
	tokenUserRequired = "UoQNgRSTJckrZFFAMlNXValG5c3IVmiL9oNspreQykY="
)

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

func testAPI(router http.Handler, method, URL string, body interface{}, remember string) *httptest.ResponseRecorder {
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

			if test.want != nil {
				dec := json.NewDecoder(res.Body)
				err := dec.Decode(&test.got)
				if err != nil {
					t.Error(err)
				}
				//delete all json fields we want to ignore
				deleteUnwantedFields(test.got)
				assert.Equal(t, test.want, test.got, test.tag)
			}

			assert.Equal(t, test.status, res.Code, test.tag)
		})
	}
}

func toMap(d interface{}) map[string]interface{} {
	var inInterface map[string]interface{}
	inrec, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(inrec, &inInterface)
	deleteUnwantedFields(inInterface)
	fmt.Printf("map: %+v\n", inInterface)
	return inInterface
}

func deleteFields(m map[string]interface{}) {
	delete(m, "created_at")
	delete(m, "updated_at")
	delete(m, "deleted_at")
}

func deleteUnwantedFields(m map[string]interface{}) {
	deleteFields(m)
	for k, v := range m {
		switch x := v.(type) {
		case []interface{}:
			for _, value := range x {
				v, ok := value.(map[string]interface{})
				if ok {
					deleteFields(v)
				}
			}
		case map[string]interface{}:
			deleteFields(x)
		default:
			fmt.Printf("Unsupported type: %T\n", x)
		}

		fmt.Printf("\nkey[%s] value[%s]\n", k, v)
	}
}
