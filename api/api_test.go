package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"chirp.com/models"
	"github.com/stretchr/testify/assert"
)

type apiTestCase struct {
	tag    string
	method string
	url    string
	body   interface{}
	status int
	got    interface{}
	want   interface{}
}

func testAPI(router http.Handler, method, URL string, body interface{}) *httptest.ResponseRecorder {
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
	res := httptest.NewRecorder()
	// Drop a cookie on the recorder.
	SetPreferencesCookie(res, &Preferences{Colour: "Blue"})
	router.ServeHTTP(res, req)
	return res
}

func runAPITests(t *testing.T, router http.Handler, tests []apiTestCase) {
	for _, test := range tests {
		t.Run(test.tag, func(t *testing.T) {

			res := testAPI(router, test.method, test.url, test.body)
			fmt.Println("body\n", res.Body)

			dec := json.NewDecoder(res.Body)

			err := dec.Decode(&test.got)
			test.got = reflect.ValueOf(test.got).Elem().Interface() //convert back to its original type
			switch v := test.got.(type) {
			case models.User:
				test.got = sanitize(&v)
			default:
				test.got = &test.got
			}

			if err != nil {
				fmt.Printf("%+v\n", res)
				t.Fatal(err)
			}

			assert.Equal(t, test.status, res.Code, test.tag)
			assert.Equal(t, test.want, test.got, test.tag)
		})
	}
}

// func TestAPI(t *testing.T) {
// 	router := app.NewRouter()
// 	router.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
// 		sample := &sample{Msg: "hello"}
// 		utils.Render(w, sample)
// 	})

// 	want := sample{
// 		Msg: "hello",
// 	}
// 	runAPITests(t, router, []apiTestCase{
// 		{"t1 - say hi", "GET", "/hi", "", http.StatusOK, &sample{}, want},
// 	})
// }

// type sample struct {
// 	Msg string
// }

type Preferences struct {
	Colour string
}

func SetPreferencesCookie(w http.ResponseWriter, prefs *Preferences) error {
	// data, err := json.Marshal(prefs)
	// if err != nil {
	// 	return err
	// }
	http.SetCookie(w, &http.Cookie{
		Name:  "test",
		Value: "hi",
	})
	return nil
}
