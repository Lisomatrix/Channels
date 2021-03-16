// This package holds a default implementation of the storage interface
package storagesql

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/lisomatrix/channels/channels/core"

	_ "github.com/jackc/pgx/v4/stdlib" // PostgreSQL Driver
)

// DatabaseStorage - DatabaseStorage implementation using SQL Database
type DatabaseStorage struct {
	db *sql.DB
}

var appStorage *AppRepository = nil
var clientStorage *ClientRepository = nil
var channelStorage *ChannelRepository = nil
var deviceStorage *DeviceRepository = nil

// GetAppRepository - Get SQL implementation of AppRepository
func (storage *DatabaseStorage) GetAppRepository() core.AppRepository {

	if appStorage == nil {
		appStorage = NewSQLAppRepository(storage)
	}

	return appStorage
}

// GetDeviceRepository - Get SQL implementation of DeviceRepository
func (storage *DatabaseStorage) GetDeviceRepository() core.DeviceRepository {

	if deviceStorage == nil {
		deviceStorage = NewSQLDeviceRepository(storage)
	}

	return deviceStorage
}

// GetClientRepository - Get SQL implementation of ClientRepository
func (storage *DatabaseStorage) GetClientRepository() core.ClientRepository {

	if clientStorage == nil {
		clientStorage = NewSQLClientRepository(storage)
	}

	return clientStorage
}

// GetChannelRepository - Get SQL implementation of ChannelRepository
func (storage *DatabaseStorage) GetChannelRepository() core.ChannelRepository {

	if channelStorage == nil {
		channelStorage = NewSQLChannelRepository(storage)
	}

	return channelStorage
}

var postgresDriver = ""
var user = ""
var host = ""
var port = ""
var password = ""
var dbName = ""

var dataSourceName = fmt.Sprintf("host=%s port=%s user=%s "+
	"password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)

func SetConnectionParams(dbUser string, dbPassword string, dbHost string, dbPort string, db string) {
	user = dbUser
	host = dbHost
	port = dbPort
	password = dbPassword
	dbName = db

	dataSourceName = fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable", host, port, user, password, dbName)
}

// NewSQLStorageDatabase - Create new SQLStorageDatabase instance
func NewSQLStorageDatabase() core.DatabaseStorage {

	db, err := sql.Open(postgresDriver, dataSourceName)

	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxOpenConns(5)

	return &DatabaseStorage{db: db}
}
