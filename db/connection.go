package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/apex/log"
	"github.com/geniusmonkey/gander/creds"
	"github.com/geniusmonkey/gander/env"
	"github.com/geniusmonkey/gander/migration"
	"github.com/geniusmonkey/gander/project"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

var sqlDb *sql.DB

func Setup(project project.Project, env env.Environment, cred creds.Credentials) {
	var err error
	if err = migration.SetDialect(project.Driver); err != nil {
		log.Fatalf("failed to set dialect %s, %v", project.Driver, err)
	}

	var driver string
	var dsn string
	switch project.Driver {
	case "redshift", "cockroach":
		driver = "postgres"
		dsn = buildPostgres(env, cred)
	case "mysql":
		driver = "mysql"
		dsn = buildMysql(env, cred)
	case "postgres":
		driver = "postgres"
		dsn = buildPostgres(env, cred)
	}

	sqlDb, err = sql.Open(driver, dsn)
	if err != nil {
		log.Fatalf("failed to open connection %s", err)
	}
}

func buildMysql(environment env.Environment, cred creds.Credentials) string {
	sb := new(strings.Builder)

	if cred.Username != "" {
		fmt.Fprintf(sb, "%v:%v@", cred.Username, cred.Password)
	}

	fmt.Fprintf(sb, "%v(%v:%v)/%v?", environment.Protocol, environment.Host, environment.Port, environment.Schema)

	for k, v := range environment.Paramas {
		fmt.Fprintf(sb, "%v=%v&", k, v)
	}

	s := sb.String()
	fmt.Println(s)
	return s
}

func buildPostgres(environment env.Environment, cred creds.Credentials) string {
	dsn := map[string]string{
		"host":   environment.Host,
		"port":   fmt.Sprintf("%v", environment.Port),
		"dbname": environment.Schema,
	}

	if cred.Username != "" {
		dsn["user"] = cred.Username
		dsn["password"] = cred.Password
	}

	for k, v := range environment.Paramas {
		dsn[k] = v
	}

	dsnStr := ""
	for k, v := range dsn {
		dsnStr += fmt.Sprintf("%v=%v ", k, v)
	}

	return dsnStr
}

func Close() {
	_ = sqlDb.Close()
}

func Get() *sql.DB {
	return sqlDb
}
