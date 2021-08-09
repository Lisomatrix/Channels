package gormsql

import (
	"errors"

	"github.com/lisomatrix/channels/channels/core"
	"gorm.io/gorm"
)

type ChannelsChannel struct {
	ID         string           `gorm:"column:id;primaryKey;not null"`
	AppID      string           `gorm:"column:app_id;not null"`
	Name       string           `gorm:"column:name;not null"`
	CreatedAt  int64            `gorm:"column:created_at;not null"`
	IsClosed   bool             `gorm:"column:is_closed;not null"`
	Extra      string           `gorm:"column:extra;not null"`
	Persistent bool             `gorm:"column:persistent;not null"`
	Private    bool             `gorm:"column:private;not null"`
	Presence   bool             `gorm:"column:presence;not null"`
	Push       bool             `gorm:"column:push;not null"`
	Clients    []ChannelsClient `gorm:"many2many:channel_client;"`
}

func (c *ChannelsChannel) TableName() string {
	return "channel"
}

type GormChannelRepository struct {
	gormDB *gorm.DB
}

func (repo *GormChannelRepository) Migrate() error {
	return repo.gormDB.AutoMigrate(&ChannelsChannel{})
}

func (repo *GormChannelRepository) CreateChannel(id string, appID string, name string, createdAt int64, isClosed bool, extra string, persistent bool, private bool, presence bool, push bool) error {
	return repo.gormDB.Create(&ChannelsChannel{
		ID:         id,
		AppID:      appID,
		Name:       name,
		CreatedAt:  createdAt,
		IsClosed:   isClosed,
		Extra:      extra,
		Persistent: persistent,
		Private:    private,
		Presence:   presence,
		Push:       push,
	}).Error
}

func (repo *GormChannelRepository) GetChannelClients(appID string, channelID string) ([]string, error) {
	var channel ChannelsChannel

	tx := repo.gormDB.Where(map[string]interface{}{"app_id": appID, "id": channelID}).Preload("Clients").First(&channel)

	if tx.Error != nil {
		return nil, tx.Error
	}

	clientIDs := make([]string, 0, len(channel.Clients))

	for _, c := range channel.Clients {
		clientIDs = append(clientIDs, c.ID)
	}

	return clientIDs, nil
}

func (repo *GormChannelRepository) DeleteChannel(appID string, id string) error {
	return repo.gormDB.Delete(&ChannelsClient{}, id).Error
}

func (repo *GormChannelRepository) DeleteAppChannels(appID string) error {
	return repo.gormDB.Where(map[string]interface{}{"app_id": appID}).Delete(&ChannelsChannel{}).Error
}

func (repo *GormChannelRepository) JoinClient(appID string, channelID string, clientID string) error {
	channel := ChannelsChannel{ID: channelID, AppID: appID}
	return repo.gormDB.Model(&channel).Association("Clients").Append(&ChannelsClient{ID: clientID})
}

func (repo *GormChannelRepository) LeaveClient(appID string, channelID string, clientID string) error {
	channel := ChannelsChannel{ID: channelID, AppID: appID}
	return repo.gormDB.Model(&channel).Association("Clients").Delete(&ChannelsClient{ID: clientID})
}

func (repo *GormChannelRepository) SetChannelCloseStatus(appID string, channelID string, isClosed bool) error {
	return repo.gormDB.Model(&ChannelsChannel{ID: channelID, AppID: appID}).UpdateColumn("is_closed", isClosed).Error
}

func (repo *GormChannelRepository) GetAppChannel(appID string, channelID string) (*core.Channel, error) {
	var channel ChannelsChannel

	tx := repo.gormDB.First(&channel, channelID)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, tx.Error
	}

	return &core.Channel{
		ID:         channelID,
		AppID:      appID,
		Name:       channel.Name,
		CreatedAt:  channel.CreatedAt,
		IsClosed:   channel.IsClosed,
		Extra:      channel.Extra,
		Persistent: channel.Persistent,
		Private:    channel.Private,
		Presence:   channel.Presence,
		Push:       channel.Push,
	}, nil
}

func (repo *GormChannelRepository) GetClientAllowedChannels(clientID string) ([]string, error) {
	return clientStorage.GetClientAllowedChannels(clientID)
}

func (repo *GormChannelRepository) GetClientPrivateChannels(clientID string) ([]*core.Channel, error) {
	client, err := clientStorage.GetClientWithChannels(clientID)

	if err != nil {
		return nil, err
	} else if client == nil {
		return nil, nil
	}

	channels := make([]*core.Channel, 0, len(client.Channels)/2)

	for _, c := range client.Channels {
		if c.Private {
			channels = append(channels, &core.Channel{
				ID:         c.ID,
				AppID:      c.AppID,
				Name:       c.Name,
				CreatedAt:  c.CreatedAt,
				IsClosed:   c.IsClosed,
				Extra:      c.Extra,
				Persistent: c.Persistent,
				Private:    c.Private,
				Presence:   c.Presence,
				Push:       c.Push,
			})
		}
	}

	return channels, nil
}

func (repo *GormChannelRepository) GetClientPublicChannels(clientID string) ([]*core.Channel, error) {
	client, err := clientStorage.GetClientWithChannels(clientID)

	if err != nil {
		return nil, err
	} else if client == nil {
		return nil, nil
	}

	channels := make([]*core.Channel, 0, len(client.Channels)/2)

	for _, c := range client.Channels {
		if !c.Private {
			channels = append(channels, &core.Channel{
				ID:         c.ID,
				AppID:      c.AppID,
				Name:       c.Name,
				CreatedAt:  c.CreatedAt,
				IsClosed:   c.IsClosed,
				Extra:      c.Extra,
				Persistent: c.Persistent,
				Private:    c.Private,
				Presence:   c.Presence,
				Push:       c.Push,
			})
		}
	}

	return channels, nil
}

func (repo *GormChannelRepository) GetAppPrivateChannels(appID string) ([]*core.Channel, error) {
	channels := make([]ChannelsChannel, 0)
	tx := repo.gormDB.Where(map[string]interface{}{"app_id": appID}).Find(&channels)

	if tx.Error != nil {
		return nil, tx.Error
	}

	privChannels := make([]*core.Channel, 0, len(channels))

	for _, c := range channels {
		if c.Private {
			privChannels = append(privChannels, &core.Channel{
				ID:         c.ID,
				AppID:      c.AppID,
				Name:       c.Name,
				CreatedAt:  c.CreatedAt,
				IsClosed:   c.IsClosed,
				Extra:      c.Extra,
				Persistent: c.Persistent,
				Private:    c.Private,
				Presence:   c.Presence,
				Push:       c.Push,
			})
		}
	}

	return privChannels, nil
}

func (repo *GormChannelRepository) GetAppPublicChannels(appID string) ([]*core.Channel, error) {
	channels := make([]ChannelsChannel, 0)
	tx := repo.gormDB.Where(map[string]interface{}{"app_id": appID}).Find(&channels)

	if tx.Error != nil {
		return nil, tx.Error
	}

	privChannels := make([]*core.Channel, 0, len(channels))

	for _, c := range channels {
		if !c.Private {
			privChannels = append(privChannels, &core.Channel{
				ID:         c.ID,
				AppID:      c.AppID,
				Name:       c.Name,
				CreatedAt:  c.CreatedAt,
				IsClosed:   c.IsClosed,
				Extra:      c.Extra,
				Persistent: c.Persistent,
				Private:    c.Private,
				Presence:   c.Presence,
				Push:       c.Push,
			})
		}
	}

	return privChannels, nil
}

func (repo *GormChannelRepository) ExistsAppChannel(appID string, channelID string) (bool, error) {
	channel, err := repo.GetAppChannel(appID, channelID)

	if err != nil {
		return false, err
	}

	if channel == nil {
		return false, nil
	}

	return true, nil
}
