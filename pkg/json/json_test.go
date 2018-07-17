package json

import (
	"reflect"
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	type form struct {
		ID   int
		Post string
	}
	var got form
	want := form{
		ID:   5,
		Post: "hello",
	}
	data := `{"Post":"hello","ID":5}`
	err := Decode(&got, strings.NewReader(data))
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %+v got %+v", want, got)
	}
}
