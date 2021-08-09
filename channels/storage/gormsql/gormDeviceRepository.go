package gormsql

import (
	"errors"

	"github.com/lisomatrix/channels/channels/core"
	"gorm.io/gorm"
)

type ChannelsDevice struct {
	ID       string `gorm:"column:id;primaryKey;not null"`
	Token    string `gorm:"column:token;not null"`
	ClientID string `gorm:"column:client_id;not null"`
}

type ChannelsDeviceToken struct {
	Token string `gorm:"column:token;not null"`
}

type GormDeviceRepository struct {
	gormDB *gorm.DB
}

func (repo *GormDeviceRepository) CreateDevice(id, token, clientID string) error {
	return repo.gormDB.Create(ChannelsDevice{
		ID:       id,
		Token:    token,
		ClientID: clientID,
	}).Error
}

func (repo *GormDeviceRepository) GetDevice(id string) (*core.Device, error) {
	var device ChannelsDevice

	tx := repo.gormDB.First(&device, id)

	if tx.Error != nil {
		if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, tx.Error
	}

	return &core.Device{
		ID:       device.ID,
		Token:    device.Token,
		ClientID: device.ClientID,
	}, nil
}

func (repo *GormDeviceRepository) DeleteDevice(id string) error {
	return repo.gormDB.Delete(&ChannelsDevice{}, id).Error
}

func (repo *GormDeviceRepository) DeleteClientDevices(clientID string) error {
	return repo.gormDB.Delete(&ChannelsDevice{}, "client_id = ?", clientID).Error
}

func (repo *GormDeviceRepository) GetClientDevices(clientID string) ([]*core.Device, error) {
	devices := make([]ChannelsDevice, 0)

	tx := repo.gormDB.Where(map[string]interface{}{"client_id": clientID}).Find(&devices)

	if tx.Error != nil {
		return nil, tx.Error
	}

	coreDevices := make([]*core.Device, 0, len(devices))

	for _, d := range devices {
		coreDevices = append(coreDevices, &core.Device{
			ID:       d.ID,
			Token:    d.Token,
			ClientID: d.ClientID,
		})
	}

	return coreDevices, nil
}

func (repo *GormDeviceRepository) GetClientsDeviceTokens(clientIDs []string, amount int) ([]string, error) {
	devicesTokens := make([]ChannelsDeviceToken, 0)
	tx := repo.gormDB.Select("token").Where("client_id in (?)", clientIDs)

	if tx.Error != nil {
		return nil, tx.Error
	}

	tokens := make([]string, 0, len(devicesTokens))

	for _, t := range devicesTokens {
		tokens = append(tokens, t.Token)
	}

	return tokens, nil
}
