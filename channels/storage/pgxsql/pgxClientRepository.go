package pgxsql

import (
	"context"
	"fmt"
	"github.com/lisomatrix/channels/channels/core"
	"os"
)

// Client SQL
var createClientSQL = `INSERT INTO "Client"("ID", "Username", "AppID", "Extra") VALUES ( $1 , $2 , $3 , $4 );`
var deleteClientByIDSQL = `DELETE FROM "Client" WHERE "ID" = $1;`
var deleteClientByAppIDSQL = `DELETE FROM "Client" WHERE "AppID" = $1;`
var updateClientSQL = `UPDATE "Client" SET "Username" = $1, "Extra" = $2 WHERE "ID" = $3 ;`
var updateClientUsernameSQL = `UPDATE "Client" SET "Username" = $1 WHERE "ID" = $2 ;`
var updateClientExtraSQL = `UPDATE "Client" SET "Extra" = $1 WHERE "ID" = $2;`
var selectAppClientsSQL = `SELECT "ID", "Username", "AppID", "Extra" FROM "Client" WHERE "AppID" = $1;`
var selectAppClientSQL = `SELECT "ID", "Username", "AppID", "Extra" FROM "Client" WHERE "AppID" = $1 AND "ID" = $2;`
var selectAppClientExistsSQL = `SELECT COUNT("ID") FROM "Client" WHERE "AppID" = $1 AND "ID" = $2 LIMIT 1;`
var selectAppClientsAmountSQL = `SELECT COUNT("ID") FROM "Client" WHERE "AppID" = $1;`
var selectAllClientsSQL = `SELECT "ID", "Username", "AppID", "Extra" FROM "Client";`
var selectAllClientsAmountSQL = `SELECT COUNT("ID") AS "AMOUNT" FROM "Client";`
var selectClientExtraSQL = `SELECT "Extra" FROM "Client" WHERE "ID" = $1;`

// PGXClientRepository - SQL implementation of client repository
type PGXClientRepository struct {
	dbHolder *PGXDatabaseStorage
	ctx      context.Context
}

// ExistsAppClient - Get app client
func (repo *PGXClientRepository) ExistsAppClient(AppID string, ClientID string) (bool, error) {
	row := repo.dbHolder.db.QueryRow(repo.ctx, selectAppClientExistsSQL, AppID, ClientID)

	var found int64

	err := row.Scan(&found)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetAppClient: row scan failed: %v\n", err)
		return false, err
	}

	return found == 1, nil
}

// GetAppClient - Get app client
func (repo *PGXClientRepository) GetAppClient(AppID string, ClientID string) (*core.Client, error) {
	row := repo.dbHolder.db.QueryRow(repo.ctx, selectAppClientSQL, AppID, ClientID)

	var id string
	var username string
	var appID string
	var extra string

	err := row.Scan(&id, &username, &appID, &extra)

	client := &core.Client{
		ID:       id,
		Username: username,
		AppID:    appID,
		Extra:    extra,
	}

	if err != nil {
		// TODO: Need a better fix for this
		if err.Error() == "no rows in result set" {
			return nil, nil
		}

		_, _ = fmt.Fprintf(os.Stderr, "GetAppClient: row scan failed: %v\n", err)
		return nil, err
	}

	return client, nil
}

// CreateClient - Insert new client row
func (repo *PGXClientRepository) CreateClient(id string, username string, appID string, extra string) error {
	_, err := repo.dbHolder.db.Exec(repo.ctx, createClientSQL, id, username, appID, extra)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "CreateClient: statement execution failed: %v\n", err)
		return err
	}

	return nil
}

// GetClientExtra - Get column Extra of the client
func (repo *PGXClientRepository) GetClientExtra(id string) (string, error) {
	row := repo.dbHolder.db.QueryRow(repo.ctx, selectClientExtraSQL, id)

	var extra string

	err := row.Scan(&extra)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetClientExtra: query scan failed: %v\n", err)
		return "", err
	}

	return extra, nil
}

// DeleteClient - Remove client row with given id
func (repo *PGXClientRepository) DeleteClient(id string) error {
	_, err := repo.dbHolder.db.Exec(repo.ctx, deleteClientByIDSQL, id)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "DeleteClient: statement execution failed: %v\n", err)
		return err
	}

	return nil
}

// DeleteAppClients - Remove all clients that belong to the given appID
func (repo *PGXClientRepository) DeleteAppClients(appID string) error {
	_, err := repo.dbHolder.db.Exec(repo.ctx, deleteClientByAppIDSQL, appID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "DeleteAppClients: statement execution failed: %v\n", err)
		return err
	}

	return nil
}

// UpdateClient - Update Client Username and Extra
func (repo *PGXClientRepository) UpdateClient(id string, username string, extra string) error {
	_, err := repo.dbHolder.db.Exec(repo.ctx, updateClientSQL, username, extra, id)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "UpdateClient: statement execution failed: %v\n", err)
		return err
	}

	return nil
}

// UpdateClientUsername - Update Client Username
func (repo *PGXClientRepository) UpdateClientUsername(id string, username string) error {
	_, err := repo.dbHolder.db.Exec(repo.ctx, updateClientUsernameSQL, username, id)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "UpdateClientUsername: statement execution failed: %v\n", err)
		return err
	}

	return nil
}

// UpdateClientExtra - Update Client Extra
func (repo *PGXClientRepository) UpdateClientExtra(id string, extra string) error {
	_, err := repo.dbHolder.db.Exec(repo.ctx, updateClientUsernameSQL, extra, id)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "UpdateClientExtra: statement execution failed: %v\n", err)
		return err
	}

	return nil
}

// GetAppClients - Get all App clients
func (repo *PGXClientRepository) GetAppClients(appID string) ([]*core.Client, error) {
	rows, err := repo.dbHolder.db.Query(repo.ctx, selectAppClientsSQL, appID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetAppClients: query failed: %v\n", err)
		return nil, err
	}

	clients := make([]*core.Client, 0)

	for rows.Next() {
		var id string
		var username string
		var appID string
		var extra string

		err = rows.Scan(&id, &username, &appID, &extra)

		clients = append(clients, &core.Client{
			ID:       id,
			Username: username,
			AppID:    appID,
			Extra:    extra,
		})

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "GetAppClients: row scan failed: %v\n", err)
			return nil, err
		}
	}

	return clients, nil
}

// GetAppClientsCount - Get how much clients an App has
func (repo *PGXClientRepository) GetAppClientsCount(appID string) (uint64, error) {
	row := repo.dbHolder.db.QueryRow(repo.ctx, selectAppClientsAmountSQL, appID)

	var amount uint64

	err := row.Scan(&amount)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetAppClients: row scan failed: %v\n", err)
		return 0, err
	}

	return amount, nil
}

// GetAllClients - Get all Apps clients
func (repo *PGXClientRepository) GetAllClients() ([]*core.Client, error) {
	rows, err := repo.dbHolder.db.Query(repo.ctx, selectAllClientsSQL)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetAppClients: query failed: %v\n", err)
		return nil, err
	}

	clients := make([]*core.Client, 0)

	for rows.Next() {
		var id string
		var username string
		var appID string
		var extra string

		err = rows.Scan(&id, &username, &appID, &extra)

		clients = append(clients, &core.Client{
			ID:       id,
			Username: username,
			AppID:    appID,
			Extra:    extra,
		})

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "GetAppClients: row scan failed: %v\n", err)
			return nil, err
		}
	}

	return clients, nil
}

// GetAllClientsCount - Get how much clients an App has
func (repo *PGXClientRepository) GetAllClientsCount() (uint64, error) {
	row, err := repo.dbHolder.db.Query(repo.ctx, selectAllClientsAmountSQL)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetAllClients: query failed: %v\n", err)
		return 0, err
	}

	var amount uint64

	err = row.Scan(&amount)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetAllClients: row scan failed: %v\n", err)
		return 0, err
	}

	return amount, nil
}

// NewSQLPGXClientRepository - Create a new instance of SQLPGXClientRepository
func NewSQLPGXClientRepository(db *PGXDatabaseStorage) *PGXClientRepository {
	return &PGXClientRepository{dbHolder: db, ctx: context.Background()}
}
