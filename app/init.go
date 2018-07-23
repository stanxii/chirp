package app

import (
	"fmt"

	"chirp.com/errors"
	"chirp.com/internal/utils"
	"chirp.com/models"
	"github.com/gorilla/mux"
)

func Setup(cfg Config) *models.Services {
	fmt.Println("========================== LOOK HERE ==========================")

	dbCfg := cfg.Database
	services, err := models.NewServices(
		models.WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		models.WithLogMode(!cfg.IsProd()),
		models.WithUser(cfg.Pepper, cfg.HMACKey),
		models.WithTweet(),
		models.WithTag(),
		models.WithTagging(),
		models.WithLike(),
		models.WithFollow(),
	)
	utils.Must(err)
	services.AutoMigrate()

	// load error messages
	if err := errors.LoadMessages("../config/errors.yaml"); err != nil {
		panic(fmt.Errorf("Failed to read the error message file: %s", err))
	}
	return services
}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	return router
}
