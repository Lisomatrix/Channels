package core

// App - Database representation of a App
type App struct {
	AppID string
	Name  string
}

// AppRepository - Repository for handling App table
type AppRepository interface {
	CreateApp(id string, name string) error
	DeleteApp(id string) error
	GetApps() ([]*App, error)
	UpdateApp(id string, name string) error
	AppExists(id string) (bool, error)
}

// Client - Database representation of a Client
type Client struct {
	ID       string
	Username string
	AppID    string
	Extra    string
}

// ClientRepository - Repository for handling Client table
type ClientRepository interface {
	CreateClient(id string, username string, appID string, extra string) error
	ExistsAppClient(AppID string, ClientID string) (bool, error)
	GetClientExtra(id string) (string, error)
	GetAppClient(AppID string, ClientID string) (*Client, error)
	DeleteClient(id string) error
	DeleteAppClients(appID string) error
	UpdateClient(id string, username string, extra string) error
	UpdateClientUsername(id string, username string) error
	UpdateClientExtra(id string, extra string) error
	GetAppClients(appID string) ([]*Client, error)
	GetAppClientsCount(appID string) (uint64, error)
	GetAllClients() ([]*Client, error)
	GetAllClientsCount() (uint64, error)
}

// Device - Database representation of a device
type Device struct {
	ID       string
	Token    string
	ClientID string
}

// DeviceRepository - Repository for handling client devices
type DeviceRepository interface {
	CreateDevice(id string, token string, clientID string) error
	GetDevice(id string) (*Device, error)
	DeleteDevice(id string) error
	DeleteClientDevices(clientID string) error
	GetClientDevices(clientID string) ([]*Device, error)
	GetClientDeviceTokens(clientID string) ([]string, error)
}

// Channel - Database representation of a channel
type Channel struct {
	ID         string `json:"id"`
	AppID      string `json:"appID"`
	Name       string `json:"name"`
	CreatedAt  int64  `json:"createdAt"`
	IsClosed   bool   `json:"isClosed"`
	Extra      string `json:"extra"`
	Persistent bool   `json:"isPersistent"`
	Private    bool   `json:"isPrivate"`
	Presence   bool   `json:"isPresence"`
}

// ChannelRepository - Repository for handling Channel, Channel_Event and Channel_Client tables
type ChannelRepository interface {
	CreateChannel(id string, appID string, name string, createdAt int64, isClosed bool, extra string, persistent bool, private bool, presence bool) error

	GetChannelClients(appID string, channelID string) ([]string, error)
	DeleteChannel(appID string, id string) error
	DeleteAppChannels(appID string) error

	JoinClient(appID string, channelID string, clientID string) error
	LeaveClient(appID string, channelID string, clientID string) error

	SetChannelCloseStatus(appID string, channelID string, isClosed bool) error

	GetClientAllowedChannels(clientID string) ([]string, error)
	GetClientPrivateChannels(clientID string) ([]*Channel, error)
	GetClientPublicChannels(clientID string) ([]*Channel, error)

	GetAppPrivateChannels(appID string) ([]*Channel, error)
	GetAppPublicChannels(appID string) ([]*Channel, error)

	ExistsAppChannel(appID string, channelID string) (bool, error)
	GetAppChannel(appID string, channelID string) (*Channel, error)

	GetAppChannels(appID string) ([]*Channel, error)
	GetAppChannelsCount(appID string) (uint64, error)

	GetAllPublicChannels() ([]*Channel, error)
	GetAllPrivateChannels() ([]*Channel, error)
	GetAllChannels() ([]*Channel, error)
	GetAllChannelsCount() (uint64, error)

	AddChannelEvent(appID string, channelID string, event *ChannelEvent) error
	AddChannelEvents(items []InsertItem) error

	GetChannelEventsAfter(appID string, channelID string, timestamp int64) ([]*ChannelEvent, error)
	GetChannelEventsAfterAndBefore(appID string, channelID string, timestampAfter int64, timestampBefore int64) ([]*ChannelEvent, error)
	GetChannelLastEvents(appID string, channelID string, amount int64) ([]*ChannelEvent, error)
	GetChannelLastEventsAfter(appID string, channelID string, amount int64, timestamp int64) ([]*ChannelEvent, error)
}

// DatabaseStorage - Persistent database storage interface
type DatabaseStorage interface {
	GetAppRepository() AppRepository
	GetClientRepository() ClientRepository
	GetChannelRepository() ChannelRepository
	GetDeviceRepository() DeviceRepository
}
