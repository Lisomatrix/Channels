package pgxsql

import (
	"context"
	"fmt"
	"github.com/lib/pq"
	"github.com/lisomatrix/channels/channels/core"
	"os"
)

// Device SQL
var createDeviceSQL = `INSERT INTO "Device"("ID", "Token", "ClientID") VALUES ( $1 , $2 , $3 );`
var selectDeviceSQL = `SELECT "Token", "ClientID" FROM "Device" WHERE "ID" = $1;`
var selectClientDevicesSQL = `SELECT "ID", "Token" FROM "Device" WHERE "ClientID" = $1;`
var selectClientDeviceTokensSQL = `SELECT "Token" FROM "Device" WHERE "ClientID" = $1;`
var selectClientsDeviceTokensSQL = `SELECT "Token" FROM "Device" WHERE "ClientID" = ANY($1) LIMIT $2 ;`
var deleteClientDevicesSQL = `DELETE FROM "Device" WHERE "ClientID" = $1;`
var deleteDeviceSQL = `DELETE FROM "Device" WHERE "ID" = $1;`

// PGXDeviceRepository - SQL repository for table Device
type PGXDeviceRepository struct {
	dbHolder *PGXDatabaseStorage
	ctx      context.Context
}

// GetDevice - Get device
func (repo *PGXDeviceRepository) GetDevice(id string) (*core.Device, error) {
	row := repo.dbHolder.db.QueryRow(repo.ctx, selectDeviceSQL, id)

	var token string
	var clientID string

	err := row.Scan(&token, &clientID)

	if err != nil {

		if err.Error() == "no rows in result set" {
			return nil, nil
		}

		_, _ = fmt.Fprintf(os.Stderr, "GetDevice: query scan failed: %v\n", err)
		return nil, err
	}

	return &core.Device{
		ID:       id,
		Token:    token,
		ClientID: clientID,
	}, nil
}

// GetClientsDeviceTokens - Get all clients device tokens up to given amount
func (repo *PGXDeviceRepository) GetClientsDeviceTokens(clientIDs []string, amount int) ([]string, error) {

	rows, err := repo.dbHolder.db.Query(repo.ctx, selectClientsDeviceTokensSQL, pq.Array(&clientIDs), amount)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetClientDeviceTokens: query failed: %v\n", err)
		return nil, err
	}

	devices := make([]string, 0)

	for rows.Next() {
		var token string

		err = rows.Scan(&token)

		devices = append(devices, token)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "GetClientDeviceTokens: row scan failed: %v\n", err)
			return nil, err
		}
	}

	return devices, nil
}

// GetClientDeviceTokens - Get all client device tokens
func (repo *PGXDeviceRepository) GetClientDeviceTokens(clientID string) ([]string, error) {
	rows, err := repo.dbHolder.db.Query(repo.ctx, selectClientDeviceTokensSQL, clientID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetClientDeviceTokens: query failed: %v\n", err)
		return nil, err
	}

	devices := make([]string, 0)

	for rows.Next() {
		var token string

		err = rows.Scan(&token)

		devices = append(devices, token)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "GetClientDeviceTokens: row scan failed: %v\n", err)
			return nil, err
		}
	}

	return devices, nil
}

// GetClientDevices - Get all client devices
func (repo *PGXDeviceRepository) GetClientDevices(clientID string) ([]*core.Device, error) {
	rows, err := repo.dbHolder.db.Query(repo.ctx, selectClientDevicesSQL, clientID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetClientDevices: query failed: %v\n", err)
		return nil, err
	}

	devices := make([]*core.Device, 0)

	for rows.Next() {
		var id string
		var token string

		err = rows.Scan(&id, &token)

		devices = append(devices, &core.Device{ID: id, Token: token, ClientID: clientID})

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "GetClientDevices: row scan failed: %v\n", err)
			return nil, err
		}
	}

	return devices, nil
}

// DeleteClientDevices - Delete all client devices
func (repo *PGXDeviceRepository) DeleteClientDevices(clientID string) error {
	_, err := repo.dbHolder.db.Exec(repo.ctx, deleteClientDevicesSQL, clientID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "DeleteClientDevices: statement execution failed: %v\n", err)
		return err
	}

	return nil
}

// DeleteDevice - Delete a device
func (repo *PGXDeviceRepository) DeleteDevice(id string) error {
	_, err := repo.dbHolder.db.Exec(repo.ctx, deleteDeviceSQL, id)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "DeleteDevice: statement execution failed: %v\n", err)
		return err
	}

	return nil
}

// CreateDevice - Insert a new device
func (repo *PGXDeviceRepository) CreateDevice(id string, token string, clientID string) error {
	_, err := repo.dbHolder.db.Exec(repo.ctx, createDeviceSQL, id, token, clientID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "CreateDevice: statement execution failed: %v\n", err)
		return err
	}

	return nil
}

// NewSQLPGXDeviceRepository - Create a new instance of SQLPGXDeviceRepository
func NewSQLPGXDeviceRepository(db *PGXDatabaseStorage) *PGXDeviceRepository {
	return &PGXDeviceRepository{dbHolder: db, ctx: context.Background()}
}
