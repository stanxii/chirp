package models

import (
	"fmt"
	"testing"

	"chirp.com/config"
	"chirp.com/testdata"
)

func TestUserDB(t *testing.T) {
	cfg := config.TestConfig()
	dbCfg := cfg.Database
	testdata.ResetDB(cfg)

	services, err := NewServices(
		WithGorm(dbCfg.Dialect(), dbCfg.ConnectionInfo()),
		WithLogMode(!cfg.IsProd()),
		// WithUser(cfg.Pepper, cfg.HMACKey),
	)
	if err != nil {
		t.Error(err)
	}
	defer services.Close()
	userDB := &userGorm{services.db}
	user, err := userDB.ByEmail("sam2018@gmail.com")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(user)

}
