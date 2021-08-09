package gormsql

import (
	"log"

	"github.com/lisomatrix/channels/channels/core"
	"gorm.io/gorm"
)

type ChannelsApp struct {
	AppID string `gorm:"column:app_id;primaryKey;not null"`
	Name  string `gorm:"column:name;not null"`
}

type GormAppRepository struct {
	gormDB *gorm.DB
}

func (repo *GormAppRepository) Migrate() error {
	return repo.gormDB.AutoMigrate(&ChannelsApp{})
}

func (repo *GormAppRepository) CreateApp(id string, name string) error {
	return repo.gormDB.Create(&ChannelsApp{
		AppID: id,
		Name:  name,
	}).Error
}

func (repo *GormAppRepository) DeleteApp(id string) error {
	return repo.gormDB.Delete(&ChannelsApp{}, id).Error
}

func (repo *GormAppRepository) GetApps() ([]*core.App, error) {
	apps := make([]ChannelsApp, 0)

	tx := repo.gormDB.Find(&apps)

	if tx.Error != nil {
		log.Println(tx.Error)
		return nil, tx.Error
	}

	coreApps := make([]*core.App, 0, len(apps))

	for _, a := range apps {
		coreApps = append(coreApps, &core.App{
			AppID: a.AppID,
			Name:  a.Name,
		})
	}

	return coreApps, nil
}

func (repo *GormAppRepository) GetApp(id string) (*core.App, error) {
	var channelApp ChannelsApp

	tx := repo.gormDB.First(&channelApp)

	if tx.Error != nil {
		return nil, tx.Error
	}

	return &core.App{
		AppID: id,
		Name:  channelApp.Name,
	}, nil
}

func (repo *GormAppRepository) UpdateApp(id string, name string) error {
	return repo.gormDB.Save(&ChannelsApp{
		AppID: id,
		Name:  name,
	}).Error
}
