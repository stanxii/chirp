package api

import (
	"net/http"

	"chirp.com/errors"
	"chirp.com/models"
	"github.com/gorilla/mux"
)

type Tags struct {
	tagS     models.TagService
	taggingS models.TaggingService
}

func NewTags(tagS models.TagService, taggingS models.TaggingService) *Tags {
	return &Tags{
		tagS:     tagS,
		taggingS: taggingS,
	}
}

// GET /tags/{:name}
func (t *Tags) Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	tag, err := t.tagS.ByName(name)
	if err != nil {
		RenderAPIError(w, errors.NotFound("Tag"))
		return
	}

	tweets, err := t.taggingS.GetTweets(tag.ID)
	if err != nil {
		RenderAPIError(w, errors.SetCustomError(err))
		return
	}
	Render(w, tweets)

}
