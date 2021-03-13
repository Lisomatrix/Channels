package storagesql

import (
	"database/sql"
	"fmt"
	"os"

	"lisomatrix.pt/channelserver/channelserver/core"
)

// Device SQL
var createDeviceSQL = `INSERT INTO "Device"("ID", "Token", "ClientID") VALUES ( $1 , $2 , $3 );`
var selectDeviceSQL = `SELECT "Token", "ClientID" FROM "Device" WHERE "ID" = $1;`
var selectClientDevicesSQL = `SELECT "ID", "Token" FROM "Device" WHERE "ClientID" = $1;`
var selectClientDeviceTokensSQL = `SELECT "Token" FROM "Device" WHERE "ClientID" = $1;`
var deleteClientDevicesSQL = `DELETE FROM "Device" WHERE "ClientID" = $1;`
var deleteDeviceSQL = `DELETE FROM "Device" WHERE "ID" = $1;`

// DeviceRepository - SQL repository for table Device
type DeviceRepository struct {
	dbHolder *DatabaseStorage
}

// GetDevice - Get device
func (repo *DeviceRepository) GetDevice(id string) (*core.Device, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectDeviceSQL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetDevice: preparing statement failed: %v\n", err)
		return nil, err
	}

	row := stmt.QueryRow(id)

	defer stmt.Close()

	var token string
	var clientID string

	err = row.Scan(&token, &clientID)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetDevice: query scan failed: %v\n", err)
		return nil, err
	}

	return &core.Device{
		ID:       id,
		Token:    token,
		ClientID: clientID,
	}, nil
}

// GetClientDeviceTokens - Get all client device tokens
func (repo *DeviceRepository) GetClientDeviceTokens(clientID string) ([]string, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectClientDeviceTokensSQL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetClientDeviceTokens: preparing statement failed: %v\n", err)
		return nil, err
	}

	rows, err := stmt.Query(clientID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetClientDeviceTokens: query failed: %v\n", err)
		return nil, err
	}

	defer stmt.Close()

	devices := make([]string, 0)

	for rows.Next() {
		var token string

		err = rows.Scan(&token)

		devices = append(devices, token)

		if err != nil {
			fmt.Fprintf(os.Stderr, "GetClientDeviceTokens: row scan failed: %v\n", err)
			return nil, err
		}
	}

	defer stmt.Close()

	return devices, nil
}

// GetClientDevices - Get all client devices
func (repo *DeviceRepository) GetClientDevices(clientID string) ([]*core.Device, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectClientDevicesSQL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetClientDevices: preparing statement failed: %v\n", err)
		return nil, err
	}

	rows, err := stmt.Query(clientID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetClientDevices: query failed: %v\n", err)
		return nil, err
	}

	defer stmt.Close()

	devices := make([]*core.Device, 0)

	for rows.Next() {
		var id string
		var token string

		err = rows.Scan(&id, &token)

		devices = append(devices, &core.Device{ID: id, Token: token, ClientID: clientID})

		if err != nil {
			fmt.Fprintf(os.Stderr, "GetClientDevices: row scan failed: %v\n", err)
			return nil, err
		}
	}

	defer stmt.Close()

	return devices, nil
}

// DeleteClientDevices - Delete all client devices
func (repo *DeviceRepository) DeleteClientDevices(clientID string) error {
	stmt, err := repo.dbHolder.db.Prepare(deleteClientDevicesSQL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "DeleteClientDevices: preparing statement failed: %v\n", err)
		return err
	}

	_, err = stmt.Exec(clientID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "DeleteClientDevices: statement execution failed: %v\n", err)
		return err
	}

	defer stmt.Close()

	return nil
}

// DeleteDevice - Delete a device
func (repo *DeviceRepository) DeleteDevice(id string) error {
	stmt, err := repo.dbHolder.db.Prepare(deleteDeviceSQL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "DeleteDevice: preparing statement failed: %v\n", err)
		return err
	}

	_, err = stmt.Exec(id)

	if err != nil {
		fmt.Fprintf(os.Stderr, "DeleteDevice: statement execution failed: %v\n", err)
		return err
	}

	defer stmt.Close()

	return nil
}

// CreateDevice - Insert a new device
func (repo *DeviceRepository) CreateDevice(id string, token string, clientID string) error {
	stmt, err := repo.dbHolder.db.Prepare(createDeviceSQL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "CreateDevice: preparing statement failed: %v\n", err)
		return err
	}

	_, err = stmt.Exec(id, token, clientID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "CreateDevice: statement execution failed: %v\n", err)
		return err
	}

	defer stmt.Close()

	return nil
}

// NewSQLDeviceRepository - Create a new instance of SQLDeviceRepository
func NewSQLDeviceRepository(db *DatabaseStorage) *DeviceRepository {
	return &DeviceRepository{dbHolder: db}
}
