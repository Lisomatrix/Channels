package storagesql

import (
	"fmt"
	"os"

	"github.com/channelserver/channelserver/core"
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
var selectAllClientsAmountSQL = `SELECT COUNT("ID") FROM "Client";`
var selectClientExtraSQL = `SELECT "Extra" FROM "Client" WHERE "ID" = $1;`

// ClientRepository - SQL implementation of client repository
type ClientRepository struct {
	dbHolder *DatabaseStorage
}

// ExistsAppClient - Get app client
func (repo *ClientRepository) ExistsAppClient(AppID string, ClientID string) (bool, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectAppClientExistsSQL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ExistsAppClient: preparing statement failed: %v\n", err)
		return false, err
	}

	row := stmt.QueryRow(AppID, ClientID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ExistsAppClient: query failed: %v\n", err)
		return false, err
	}

	defer stmt.Close()

	var found int64

	err = row.Scan(&found)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetAppClient: row scan failed: %v\n", err)
		return false, err
	}

	return found == 1, nil
}

// GetAppClient - Get app client
func (repo *ClientRepository) GetAppClient(AppID string, ClientID string) (*core.Client, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectAppClientSQL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetAppClient: preparing statement failed: %v\n", err)
		return nil, err
	}

	row := stmt.QueryRow(AppID, ClientID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetAppClient: query failed: %v\n", err)
		return nil, err
	}

	defer stmt.Close()

	var id string
	var username string
	var appID string
	var extra string

	err = row.Scan(&id, &username, &appID, &extra)

	client := &core.Client{
		ID:       id,
		Username: username,
		AppID:    appID,
		Extra:    extra,
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetAppClient: row scan failed: %v\n", err)
		return nil, err
	}

	return client, nil
}

// CreateClient - Insert new client row
func (repo *ClientRepository) CreateClient(id string, username string, appID string, extra string) error {
	stmt, err := repo.dbHolder.db.Prepare(createClientSQL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "CreateClient: preparing statement failed: %v\n", err)
		return err
	}

	_, err = stmt.Exec(id, username, appID, extra)

	if err != nil {
		fmt.Fprintf(os.Stderr, "CreateClient: statement execution failed: %v\n", err)
		return err
	}

	defer stmt.Close()

	return nil
}

// GetClientExtra - Get column Extra of the client
func (repo *ClientRepository) GetClientExtra(id string) (string, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectClientExtraSQL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetClientExtra: preparing statement failed: %v\n", err)
		return "", err
	}

	row := stmt.QueryRow(id)

	defer stmt.Close()

	var extra string

	err = row.Scan(&extra)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetClientExtra: query scan failed: %v\n", err)
		return "", err
	}

	return extra, nil
}

// DeleteClient - Remove client row with given id
func (repo *ClientRepository) DeleteClient(id string) error {
	stmt, err := repo.dbHolder.db.Prepare(deleteClientByIDSQL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "DeleteClient: preparing statement failed: %v\n", err)
		return err
	}

	_, err = stmt.Exec(id)

	if err != nil {
		fmt.Fprintf(os.Stderr, "DeleteClient: statement execution failed: %v\n", err)
		return err
	}

	defer stmt.Close()

	return nil
}

// DeleteAppClients - Remove all clients that belong to the given appID
func (repo *ClientRepository) DeleteAppClients(appID string) error {
	stmt, err := repo.dbHolder.db.Prepare(deleteClientByAppIDSQL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "DeleteAppClients: preparing statement failed: %v\n", err)
		return err
	}

	_, err = stmt.Exec(appID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "DeleteAppClients: statement execution failed: %v\n", err)
		return err
	}

	defer stmt.Close()

	return nil
}

// UpdateClient - Update Client Username and Extra
func (repo *ClientRepository) UpdateClient(id string, username string, extra string) error {
	stmt, err := repo.dbHolder.db.Prepare(updateClientSQL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "UpdateClient: preparing statement failed: %v\n", err)
		return err
	}

	_, err = stmt.Exec(username, extra, id)

	if err != nil {
		fmt.Fprintf(os.Stderr, "UpdateClient: statement execution failed: %v\n", err)
		return err
	}

	defer stmt.Close()

	return nil
}

// UpdateClientUsername - Update Client Username
func (repo *ClientRepository) UpdateClientUsername(id string, username string) error {
	stmt, err := repo.dbHolder.db.Prepare(updateClientUsernameSQL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "UpdateClientUsername: preparing statement failed: %v\n", err)
		return err
	}

	_, err = stmt.Exec(username, id)

	if err != nil {
		fmt.Fprintf(os.Stderr, "UpdateClientUsername: statement execution failed: %v\n", err)
		return err
	}

	defer stmt.Close()

	return nil
}

// UpdateClientExtra - Update Client Extra
func (repo *ClientRepository) UpdateClientExtra(id string, extra string) error {
	stmt, err := repo.dbHolder.db.Prepare(updateClientUsernameSQL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "UpdateClientExtra: preparing statement failed: %v\n", err)
		return err
	}

	_, err = stmt.Exec(extra, id)

	if err != nil {
		fmt.Fprintf(os.Stderr, "UpdateClientExtra: statement execution failed: %v\n", err)
		return err
	}

	defer stmt.Close()

	return nil
}

// GetAppClients - Get all App clients
func (repo *ClientRepository) GetAppClients(appID string) ([]*core.Client, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectAppClientsSQL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetAppClients: preparing statement failed: %v\n", err)
		return nil, err
	}

	rows, err := stmt.Query(appID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetAppClients: query failed: %v\n", err)
		return nil, err
	}

	defer stmt.Close()

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
			fmt.Fprintf(os.Stderr, "GetAppClients: row scan failed: %v\n", err)
			return nil, err
		}
	}

	return clients, nil
}

// GetAppClientsCount - Get how much clients an App has
func (repo *ClientRepository) GetAppClientsCount(appID string) (uint64, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectAppClientsAmountSQL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetAppClientsAmount: preparing statement failed: %v\n", err)
		return 0, err
	}

	row := stmt.QueryRow(appID)

	defer stmt.Close()

	var amount uint64

	err = row.Scan(&amount)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetAppClients: row scan failed: %v\n", err)
		return 0, err
	}

	return amount, nil
}

// GetAllClients - Get all Apps clients
func (repo *ClientRepository) GetAllClients() ([]*core.Client, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectAllClientsSQL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetAppClients: preparing statement failed: %v\n", err)
		return nil, err
	}

	rows, err := stmt.Query()

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetAppClients: query failed: %v\n", err)
		return nil, err
	}

	defer stmt.Close()

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
			fmt.Fprintf(os.Stderr, "GetAppClients: row scan failed: %v\n", err)
			return nil, err
		}
	}

	return clients, nil
}

// GetAllClientsCount - Get how much clients an App has
func (repo *ClientRepository) GetAllClientsCount() (uint64, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectAllClientsAmountSQL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetAllClientsAmount: preparing statement failed: %v\n", err)
		return 0, err
	}

	row := stmt.QueryRow()

	defer stmt.Close()

	var amount uint64

	err = row.Scan(&amount)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetAllClients: row scan failed: %v\n", err)
		return 0, err
	}

	return amount, nil
}

// NewSQLClientRepository - Create a new instance of SQLClientRepository
func NewSQLClientRepository(db *DatabaseStorage) *ClientRepository {
	return &ClientRepository{dbHolder: db}
}
