package api

import (
	"net/http"

	"chirp.com/errors"
	"chirp.com/internal/utils"
	"chirp.com/middleware"
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

func ServeTagResource(r *mux.Router, t *Tags, m *middleware.RequireUser) {
	r.HandleFunc("/tags/{name}", t.Show).Methods("GET")
}

// GET /tags/{:name}
func (t *Tags) Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	tag, err := t.tagS.ByName(name)
	if err != nil {
		utils.RenderAPIError(w, errors.NotFound("Tag"))
		return
	}

	tweets, err := t.taggingS.GetTweets(tag.ID)
	if err != nil {
		utils.RenderAPIError(w, errors.SetCustomError(err))
		return
	}
	utils.Render(w, tweets)

}
