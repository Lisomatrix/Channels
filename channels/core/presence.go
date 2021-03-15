package core

// // InitialPresenceStatus - Represents current presence status of a channel
// type InitialPresenceStatus struct {
// 	ChannelID    string           `json:"channelID"`
// 	ClientStatus map[string]int64 `json:"clientStatus"`
// }

// LastDevicePresence - Represents last client device heart beat
type LastDevicePresence struct {
	ClientID  string `json:"clientID"`
	DeviceID  string `json:"deviceID"`
	Timestamp int64  `json:"timestamp"`
}

// // ClientJoin - Event for user joining channel
// type ClientJoin struct {
// 	ChannelID string `json:"channelID"`
// 	ClientID  string `json:"clientID"`
// }

// // ClientLeave - Event for user leaving channel
// type ClientLeave struct {
// 	ChannelID string `json:"channelID"`
// 	ClientID  string `json:"clientID"`
// }

// // OnlineStatusUpdate - Event for user subscribing/unsubscribing channel
// type OnlineStatusUpdate struct {
// 	ChannelID string `json:"channelID"`
// 	ClientID  string `json:"clientID"`
// 	Status    bool   `json:"status"`
// 	TimeStamp int64  `json:"timestamp"`
// }

// PresenceHandler - Handle client online status
type PresenceHandler interface {

	// Channel Presence
	GetChannelClientsPresence(appID string, channelID string) map[string]int64
	AddOnlineChannelDevice(appID string, channelID string, clientID string, deviceID string)
	RemoveOnlineChannelDevice(appID string, channelID string, clientID string, deviceID string)
	GetChannelAmountOfClientDevices(appID string, channelID string, clientID string) int64
	IsClientDeviceConnectToChannel(appID string, channelID string, clientID string, deviceID string) bool

	// This Instant Online Status
	SetDeviceOnline(clientID string, deviceID string)
	SetDeviceOffline(clientID string, deviceID string)
	GetClientOnlineDevices(clientID string) ([]string, error)

	IsOnline(clientID string) bool
	AddDevice(clientID string, deviceID string)
	RemoveDevice(clientID string, deviceID string)

	// Last timestamps
	GetClientDevicesPresences(clientID string) ([]*LastDevicePresence, error)
	UpdateDeviceTimestamp(clientID string, deviceID string)
}
