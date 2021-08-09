package storagesql

import (
	"fmt"
	"os"

	"github.com/lisomatrix/channels/channels/core"
)

// App SQL
var createAppSQL = `INSERT INTO "App"("AppID", "Name") VALUES ( $1 , $2 );`
var deleteAppSQL = `DELETE FROM "App" WHERE "AppID" = $1 ;`
var getAppsSQL = `SELECT "AppID", "Name" FROM "App";`
var getAppSQL = `SELECT "AppID", "Name" FROM "App" WHERE "AppID" = $1 ;`
var updateAppSQL = `UPDATE "App" SET "Name" = $1 WHERE "AppID" = $2 ;`
var appExistsSQL = `SELECT COUNT("AppID") AS "EXISTS" FROM "App" WHERE "AppID" = $1 LIMIT 1;`

// AppRepository - SQL repository for table App
type AppRepository struct {
	dbHolder *DatabaseStorage
}

// CreateApp - Create a new App row in the database
func (storage *AppRepository) CreateApp(id string, name string) error {
	stmt, err := storage.dbHolder.db.Prepare(createAppSQL)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "CreateApp: preparing statement failed: %v\n", err)
		return err
	}

	_, err = stmt.Exec(id, name)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "CreateApp: statement execution failed: %v\n", err)
		return err
	}

	defer stmt.Close()

	return nil
}

// DeleteApp - Delete App Row in the database
func (storage *AppRepository) DeleteApp(id string) error {
	stmt, err := storage.dbHolder.db.Prepare(deleteAppSQL)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "DeleteApp: preparing statement failed: %v\n", err)
		return err
	}

	_, err = stmt.Exec(id)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "DeleteApp: statement execution failed: %v\n", err)
		return err
	}

	defer stmt.Close()

	return nil
}

func (storage *AppRepository) GetApp(id string) (*core.App, error) {
	stmt, err := storage.dbHolder.db.Prepare(getAppSQL)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetApp: preparing statement failed: %v\n", err)
		return nil, err
	}

	row := stmt.QueryRow()

	var app core.App

	err = row.Scan(&app)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetApp: query scan failed: %v\n", err)
		return nil, err
	}

	return &app, nil
}

// GetApps - Get all stored apps in the database
func (storage *AppRepository) GetApps() ([]*core.App, error) {
	stmt, err := storage.dbHolder.db.Prepare(getAppsSQL)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetApps: preparing statement failed: %v\n", err)
		return nil, err
	}

	rows, err := stmt.Query()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetApps: query failed: %v\n", err)
		return nil, err
	}

	defer stmt.Close()

	apps := make([]*core.App, 0)

	for rows.Next() {
		var appID string
		var name string

		err = rows.Scan(&appID, &name)

		apps = append(apps, &core.App{AppID: appID, Name: name})

		if err != nil {
			fmt.Fprintf(os.Stderr, "GetApps: row scan failed: %v\n", err)
			return nil, err
		}
	}

	defer stmt.Close()

	return apps, nil
}

// UpdateApp - Update App Row in the database
func (storage *AppRepository) UpdateApp(id string, name string) error {
	stmt, err := storage.dbHolder.db.Prepare(updateAppSQL)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "UpdateApp: preparing statement failed: %v\n", err)
		return err
	}

	_, err = stmt.Exec(name, id)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "UpdateApp: statement execution failed: %v\n", err)
		return err
	}

	defer stmt.Close()

	return nil
}

// AppExists - Check if App already exists
func (storage *AppRepository) AppExists(id string) (bool, error) {
	stmt, err := storage.dbHolder.db.Prepare(appExistsSQL)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "AppExists: preparing statement failed: %v\n", err)
		return false, err
	}

	row := stmt.QueryRow(id)

	defer stmt.Close()

	var exists bool

	err = row.Scan(&exists)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "AppExists: query scan failed: %v\n", err)
		return false, err
	}

	return exists, nil
}

// NewSQLAppRepository - Create a new instance of SQLAppRepository
func NewSQLAppRepository(db *DatabaseStorage) *AppRepository {
	return &AppRepository{dbHolder: db}
}
