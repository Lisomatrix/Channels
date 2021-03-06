// This package holds the PGX driver implementation of the storage interface
package pgxsql

import (
	"context"
	"fmt"
	"github.com/lisomatrix/channels/channels/core"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

var user = ""
var host = ""
var port = ""
var password = ""
var dbName = ""

var dataSourceName = fmt.Sprintf("host=%s port=%s user=%s "+
	"password=%s dbname=%s sslmode=prefer", host, port, user, password, dbName)

// PGXDatabaseStorage - DatabaseStorage implementation with postgres driver
type PGXDatabaseStorage struct {
	db *pgxpool.Pool
}

var appStorage *PGXAppRepository = nil
var clientStorage *PGXClientRepository = nil
var channelStorage *PGXChannelRepository = nil
var deviceStorage *PGXDeviceRepository = nil

func PGXSetConnectionParams(dbUser string, dbPassword string, dbHost string, dbPort string, db string) {
	user = dbUser
	host = dbHost
	port = dbPort
	password = dbPassword
	dbName = db

	dataSourceName = fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=prefer", host, port, user, password, dbName)
}

func (storage *PGXDatabaseStorage) GetDB() *pgxpool.Pool {
	return storage.db
}

// GetAppRepository - Get SQL implementation of AppRepository
func (storage *PGXDatabaseStorage) GetAppRepository() core.AppRepository {

	if appStorage == nil {
		appStorage = NewSQLPGXAppRepository(storage)
	}

	return appStorage
}

// GetDeviceRepository - Get SQL implementation of DeviceRepository
func (storage *PGXDatabaseStorage) GetDeviceRepository() core.DeviceRepository {

	if deviceStorage == nil {
		deviceStorage = NewSQLPGXDeviceRepository(storage)
	}

	return deviceStorage
}

// GetClientRepository - Get SQL implementation of ClientRepository
func (storage *PGXDatabaseStorage) GetClientRepository() core.ClientRepository {

	if clientStorage == nil {
		clientStorage = NewSQLPGXClientRepository(storage)
	}

	return clientStorage
}

// GetChannelRepository - Get SQL implementation of ChannelRepository
func (storage *PGXDatabaseStorage) GetChannelRepository() core.ChannelRepository {

	if channelStorage == nil {
		channelStorage = NewSQLChannelRepository(storage)
	}

	return channelStorage
}

// NewSQLStorageDatabase - Create new SQLStorageDatabase implementation with postgre specific driver, it also works with YugaByteDB tested it
func NewSQLStorageDatabase() *PGXDatabaseStorage {

	poolConfig, err := pgxpool.ParseConfig(dataSourceName)

	if err != nil {
		log.Fatal("Unable to parse DATABASE_URL", "error", err)
	}
	poolConfig.MaxConns = 5

	conn, err := pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		log.Fatal("Unable to create connection pool", "error", err)
	}

	return &PGXDatabaseStorage{db: conn}
}
