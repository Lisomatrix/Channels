package gormsql

import (
	"log"

	"gorm.io/gorm"
)

var appStorage *GormAppRepository = nil
var clientStorage *GormClientRepository = nil
var deviceStorage *GormDeviceRepository = nil
var channelStorage *GormChannelRepository = nil

type GormDatabaseStorage struct {
	gormDB *gorm.DB
}

// NewGormDatabaseStorage - Pass your pointer so you can specify your own database config
func NewGormDatabaseStorage(gormDB *gorm.DB) *GormDatabaseStorage {
	return &GormDatabaseStorage{gormDB: gormDB}
}

func (storage *GormDatabaseStorage) Migrate() {

	// App Table
	if err := storage.GetAppRepository().Migrate(); err != nil {
		log.Fatal(err)
	}
}

func (storage *GormDatabaseStorage) GetDeviceRepository() *GormDeviceRepository {
	if deviceStorage == nil {
		deviceStorage = &GormDeviceRepository{gormDB: storage.gormDB}
	}

	return deviceStorage
}

func (storage *GormDatabaseStorage) GetAppRepository() *GormAppRepository {
	if appStorage == nil {
		appStorage = &GormAppRepository{gormDB: storage.gormDB}
	}

	return appStorage
}

func (storage *GormDatabaseStorage) GetClientRepository() *GormClientRepository {
	if clientStorage == nil {
		clientStorage = &GormClientRepository{gormDB: storage.gormDB}
	}

	return clientStorage
}

func (storage *GormDatabaseStorage) GetChannelRepository() *GormChannelRepository {
	if channelStorage == nil {
		channelStorage = &GormChannelRepository{gormDB: storage.gormDB}
	}

	return channelStorage
}
