package app

import (
	"flag"
	"fmt"

	"chirp.com/errors"
	"github.com/gorilla/mux"
)

func Init() Config {
	boolPtr := flag.Bool("prod", false, "Provide this flag in production. This ensures that a .config file is provided before the application starts.")
	flag.Parse()

	cfg := LoadConfig(*boolPtr)

	// load error messages
	if err := errors.LoadMessages("config/errors.yaml"); err != nil {
		panic(fmt.Errorf("Failed to read the error message file: %s", err))
	}

	return cfg

}

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	return router
}
