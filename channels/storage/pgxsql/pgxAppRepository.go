package pgxsql

import (
	"context"
	"fmt"
	"github.com/lisomatrix/channels/channels/core"
	"os"

)

// App SQL
var createAppSQL = `INSERT INTO "App"("AppID", "Name") VALUES ( $1 , $2 );`
var deleteAppSQL = `DELETE FROM "App" WHERE "AppID" = $1 ;`
var getAppsSQL = `SELECT "AppID", "Name" FROM "App";`
var updateAppSQL = `UPDATE "App" SET "Name" = $1 WHERE "AppID" = $2 ;`
var appExistsSQL = `SELECT COUNT("AppID") AS "EXISTS" FROM "App" WHERE "AppID" = $1 LIMIT 1;`

// PGXAppRepository - SQL repository for table App
type PGXAppRepository struct {
	dbHolder *PGXDatabaseStorage
	ctx      context.Context
}

// CreateApp - Create a new App row in the database
func (storage *PGXAppRepository) CreateApp(id string, name string) error {
	_, err := storage.dbHolder.db.Exec(storage.ctx, createAppSQL, id, name)

	if err != nil {
		fmt.Fprintf(os.Stderr, "CreateApp: statement execution failed: %v\n", err)
		return err
	}

	return nil
}

// DeleteApp - Delete App Row in the database
func (storage *PGXAppRepository) DeleteApp(id string) error {
	_, err := storage.dbHolder.db.Exec(storage.ctx, deleteAppSQL, id)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "DeleteApp: statement execution failed: %v\n", err)
		return err
	}

	return nil
}

// GetApps - Get all stored apps in the database
func (storage *PGXAppRepository) GetApps() ([]*core.App, error) {
	rows, err := storage.dbHolder.db.Query(storage.ctx, getAppsSQL)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetApps: query failed: %v\n", err)
		return nil, err
	}

	apps := make([]*core.App, 0)

	for rows.Next() {
		var appID string
		var name string

		err = rows.Scan(&appID, &name)

		apps = append(apps, &core.App{AppID: appID, Name: name})

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "GetApps: row scan failed: %v\n", err)
			return nil, err
		}
	}

	return apps, nil
}

// UpdateApp - Update App Row in the database
func (storage *PGXAppRepository) UpdateApp(id string, name string) error {
	_, err := storage.dbHolder.db.Exec(storage.ctx, updateAppSQL, name, id)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "UpdateApp: statement execution failed: %v\n", err)
		return err
	}

	return nil
}

// AppExists - Check if App already exists
func (storage *PGXAppRepository) AppExists(id string) (bool, error) {
	row := storage.dbHolder.db.QueryRow(storage.ctx, appExistsSQL, id)

	var exists int

	err := row.Scan(&exists)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "AppExists: query scan failed: %v\n", err)
		return false, err
	}

	return exists >= 1, nil
}

// NewSQLPGXAppRepository - Create a new instance of SQLPGXAppRepository
func NewSQLPGXAppRepository(db *PGXDatabaseStorage) *PGXAppRepository {
	return &PGXAppRepository{dbHolder: db, ctx: context.Background()}
}
