package testdata

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"chirp.com/app"
	_ "github.com/lib/pq" // initialize posgresql for test
)

var TestConfig = app.TestConfig()

func getSQLFile() string {
	if _, err := os.Stat("testdata/db.sql"); err == nil {
		return "testdata/db.sql"
	}
	return "../testdata/db.sql"
}

// ResetDB re-create the database schema and re-populate the initial data using the SQL statements in db.sql.
func ResetDB(cfg app.Config) {
	dbCfg := cfg.Database
	out, err := exec.Command("psql", "-U", dbCfg.User, "-d", dbCfg.Name, "-f", getSQLFile()).CombinedOutput()
	if err != nil {
		log.Fatalf("Error executing query.\nCommand Output: %s\n%s", out, err.Error())
	}
	fmt.Println(string(out))

}
