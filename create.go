package gander

import (
	"fmt"
	"github.com/spf13/afero"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

var fs = afero.NewOsFs()

// Create writes a new blank migration file.
func CreateWithTemplate(dir string, migrationTemplate *template.Template, name, migrationType string) error {
	version := time.Now().Format(timestampFormat)
	filename := fmt.Sprintf("%v_%v.%v", version, name, migrationType)

	fpath := filepath.Join(dir, filename)

	tmpl := sqlMigrationTemplate
	if migrationType == "go" {
		tmpl = goSQLMigrationTemplate
	}

	if migrationTemplate != nil {
		tmpl = migrationTemplate
	}

	path, err := writeTemplateToFile(fpath, tmpl, version)
	if err != nil {
		return err
	}

	log.Infof("Created new file: %s\n", path)
	return nil
}

// Create writes a new blank migration file.
func Create(dir, name, migrationType string) error {
	return CreateWithTemplate(dir, nil, name, migrationType)
}

func writeTemplateToFile(path string, t *template.Template, version string) (string, error) {
	if _, err := fs.Stat(path); !os.IsNotExist(err) {
		return "", fmt.Errorf("failed to create file: %v already exists", path)
	}

	f, err := fs.Create(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	err = t.Execute(f, version)
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}

var sqlMigrationTemplate = template.Must(template.New("goose.sql-migration").Parse(`-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
`))

var goSQLMigrationTemplate = template.Must(template.New("goose.go-migration").Parse(`package migration

import (
	"database/sql"
	"github.com/geniusmonkey/gander"
)

func init() {
	goose.AddMigration(Up{{.}}, Down{{.}})
}

func Up{{.}}(tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return nil
}

func Down{{.}}(tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
`))
