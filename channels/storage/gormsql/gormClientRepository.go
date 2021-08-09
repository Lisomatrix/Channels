package gormsql

import (
	"errors"

	"github.com/lisomatrix/channels/channels/core"
	"gorm.io/gorm"
)

type ChannelsClient struct {
	ID       string            `gorm:"column:id;primaryKey;not null"`
	Username string            `gorm:"column:username;not null"`
	AppID    string            `gorm:"column:app_id;not null"`
	Extra    string            `gorm:"column:extra;"`
	Channels []ChannelsChannel `gorm:"many2many:channel_client;"`
}

func (c *ChannelsClient) TableName() string {
	return "client"
}

type GormClientRepository struct {
	gormDB *gorm.DB
}

func (repo *GormClientRepository) Migrate() error {
	return repo.gormDB.AutoMigrate(&ChannelsClient{})
}

func (repo *GormClientRepository) CreateClient(id, username, appID, extra string) error {
	return repo.gormDB.Create(&ChannelsClient{
		ID:       id,
		Username: username,
		AppID:    appID,
		Extra:    extra,
	}).Error
}

func (repo *GormClientRepository) GetAppClient(AppID string, ClientID string) (*core.Client, error) {
	var client ChannelsClient

	tx := repo.gormDB.Where(&ChannelsClient{ID: ClientID, AppID: AppID}).First(&client)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, tx.Error
		}

		return nil, tx.Error
	}

	return &core.Client{
		ID:       ClientID,
		Username: client.Username,
		AppID:    AppID,
		Extra:    client.Extra,
	}, nil
}

func (repo *GormClientRepository) ExistsAppClient(AppID, ClientID string) (bool, error) {
	client, err := repo.GetAppClient(AppID, ClientID)

	if err != nil {
		return false, err
	}

	if client == nil {
		return false, nil
	}

	return true, nil
}

func (repo *GormClientRepository) DeleteClient(id string) error {
	return repo.gormDB.Delete(&ChannelsClient{}, id).Error
}

func (repo *GormClientRepository) DeleteAppClients(appID string) error {
	return repo.gormDB.Delete(ChannelsClient{}, "app_id = ?", appID).Error
}

func (repo *GormClientRepository) UpdateClient(id string, username string, extra string) error {
	return repo.gormDB.Save(&ChannelsClient{
		ID:       id,
		Username: username,
		AppID:    id,
		Extra:    extra,
	}).Error
}

func (repo *GormClientRepository) GetAppClients(appID string) ([]*core.Client, error) {
	clients := make([]ChannelsClient, 0)
	tx := repo.gormDB.Where(map[string]interface{}{"app_id": appID}).Find(&clients)

	if tx.Error != nil {
		return nil, tx.Error
	}

	coreClients := make([]*core.Client, 0, len(clients))

	for _, c := range clients {
		coreClients = append(coreClients, &core.Client{
			ID:       c.ID,
			Username: c.Username,
			AppID:    appID,
			Extra:    c.Extra,
		})
	}

	return coreClients, nil
}

func (repo *GormClientRepository) GetClientWithChannels(clientID string) (*ChannelsClient, error) {
	var client ChannelsClient

	tx := repo.gormDB.Select("Channels.ID").Preload("Channels").First(&client)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, tx.Error
	}

	return &client, nil
}

func (repo *GormClientRepository) GetClientAllowedChannels(clientID string) ([]string, error) {
	var client ChannelsClient

	tx := repo.gormDB.Select("Channels.ID").Preload("Channels").First(&client)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, tx.Error
	}

	channelIDs := make([]string, 0, len(client.Channels))

	for _, c := range client.Channels {
		channelIDs = append(channelIDs, c.ID)
	}

	return channelIDs, nil
}

func (repo *GormClientRepository) GetAllClients() ([]*core.Client, error) {
	clients := make([]ChannelsClient, 0)
	tx := repo.gormDB.Find(&clients)

	if tx.Error != nil {
		return nil, tx.Error
	}

	coreClients := make([]*core.Client, 0, len(clients))

	for _, c := range clients {
		coreClients = append(coreClients, &core.Client{
			ID:       c.ID,
			Username: c.Username,
			AppID:    c.AppID,
			Extra:    c.Extra,
		})
	}

	return coreClients, nil
}
