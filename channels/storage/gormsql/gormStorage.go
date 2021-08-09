package gormsql

import (
	"log"

	"github.com/lisomatrix/channels/channels/core"
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
	if err := storage.GetAppRepository().(*GormAppRepository).Migrate(); err != nil {
		log.Fatal(err)
	}

	if err := storage.GetDeviceRepository().(*GormDeviceRepository).Migrate(); err != nil {
		log.Fatal(err)
	}

	if err := storage.GetClientRepository().(*GormClientRepository).Migrate(); err != nil {
		log.Fatal(err)
	}

	if err := storage.GetChannelRepository().(*GormChannelRepository).Migrate(); err != nil {
		log.Fatal(err)
	}
}

func (storage *GormDatabaseStorage) GetDeviceRepository() core.DeviceRepository {
	if deviceStorage == nil {
		deviceStorage = &GormDeviceRepository{gormDB: storage.gormDB}
	}

	return deviceStorage
}

func (storage *GormDatabaseStorage) GetAppRepository() core.AppRepository {
	if appStorage == nil {
		appStorage = &GormAppRepository{gormDB: storage.gormDB}
	}

	return appStorage
}

func (storage *GormDatabaseStorage) GetClientRepository() core.ClientRepository {
	if clientStorage == nil {
		clientStorage = &GormClientRepository{gormDB: storage.gormDB}
	}

	return clientStorage
}

func (storage *GormDatabaseStorage) GetChannelRepository() core.ChannelRepository {
	if channelStorage == nil {
		channelStorage = &GormChannelRepository{gormDB: storage.gormDB}
	}

	return channelStorage
}
