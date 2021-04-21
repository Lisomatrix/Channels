package storagesql

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	"github.com/lisomatrix/channels/channels/core"
)

// Channel SQL
var selectChannelClients = `SELECT "clientID" FROM "Channel_Client" WHERE "channelID" = (SELECT "ID" FROM "Channel" WHERE "ChannelID" = $1 AND "AppID" = $2 LIMIT 1);`
var createChannelSQL = `INSERT INTO "Channel"("ChannelID", "AppID", "Name", "Created_At", "IsClosed", "Extra", "Persistent", "Private", "Presence", "Push") VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`
var deleteChannelSQL = `DELETE FROM "Channel" WHERE "ChannelID" = $1 AND "AppID" = $2;`
var deleteAppChannelsSQL = `DELETE FROM "Channel" WHERE "AppID" = $1;`
var joinChannelSQL = `INSERT INTO public."Channel_Client"("clientID", "channelID") VALUES ($2, (SELECT "ID" FROM "Channel" WHERE "ChannelID" = $1 AND "AppID" = $3 LIMIT 1));`
var leaveChannelSQL = `DELETE FROM "Channel_Client" WHERE "channelID" = (SELECT "ID" FROM "Channel" WHERE "ChannelID" = $1 AND "AppID" = $3 LIMIT 1) AND "clientID" = $2;`
var setCloseStatusSQL = `UPDATE "Channel" SET "IsClosed" = $1 WHERE "ChannelID" = $2 AND "AppID" = $3;`
var selectClientAllowedChannelsSQL = `SELECT "ChannelID" FROM "Channel" WHERE "ID" IN (SELECT "channelID" FROM "Channel_Client" WHERE "clientID" = $1);`
var selectClientOpenOrPrivateChannels = `SELECT "ChannelID", "AppID", "Name", "Created_At", "IsClosed", "Extra", "Persistent", "Private", "Presence", "Push" FROM "Channel" WHERE "Private" = $1 AND "ID" IN (SELECT "channelID" FROM "Channel_Client" WHERE "clientID" = $2);`
var selectOpenOrPrivateAppChannels = `SELECT "ChannelID", "AppID", "Name", "Created_At", "IsClosed", "Extra", "Persistent", "Private", "Presence", "Push" FROM "Channel" WHERE "Private" = $1 AND "AppID" = $2;`
var selectAppChannels = `SELECT "ChannelID", "AppID", "Name", "Created_At", "IsClosed", "Extra", "Persistent", "Private", "Presence", "Push" FROM "Channel" WHERE "AppID" = $1;`
var selectAppChannelAmount = `SELECT COUNT("ChannelID") FROM "Channel" WHERE "AppID" = $1;`
var selectAppChannelExists = `SELECT COUNT("ChannelID") FROM "Channel" WHERE "AppID" = $1 AND "ChannelID" = $2 LIMIT 1;`
var selectAllOpenOrPrivateChannels = `SELECT "ChannelID", "AppID", "Name", "Created_At", "IsClosed", "Extra", "Persistent", "Private", "Presence", "Push" FROM "Channel" WHERE "Private" = $1;`
var selectAllChannels = `SELECT "ChannelID", "AppID", "Name", "Created_At", "IsClosed", "Extra", "Persistent", "Private", "Presence", "Push" FROM "Channel";`
var selectAllChannelsAmount = `SELECT COUNT("ChannelID") FROM "Channel"`
var selectAppChannel = `SELECT "ChannelID", "AppID", "Name", "Created_At", "IsClosed", "Extra", "Persistent", "Private", "Presence", "Push" FROM "Channel" WHERE "AppID" = $1 AND "ChannelID" = $2;`
var addChannelEventSQL = `INSERT INTO "Channel_Event"("SenderID", "EventType", "Payload", "ChannelID", "TimeStamp") VALUES ( $1 , $2 , $3 , (SELECT "ID" FROM "Channel" WHERE "ChannelID" = $4 AND "AppID" = $6 LIMIT 1) , $5 );`

var selectEventsSinceTimeStampSQL = `SELECT "SenderID", "EventType", "Payload", "TimeStamp" FROM "Channel_Event" WHERE "ChannelID" = (SELECT "ID" FROM "Channel" WHERE "ChannelID" = $1 AND "AppID" = $2) AND "TimeStamp" >= $3;`
var selectEventsBetweenTimeStampsSQL = `SELECT "SenderID", "EventType", "Payload", "TimeStamp" FROM "Channel_Event" WHERE "ChannelID" = (SELECT "ID" FROM "Channel" WHERE "ChannelID" = $1 AND "AppID" = $2) AND "TimeStamp" >= $3 AND "TimeStamp" <= $4;`

// * The new one is based on primary key since its auto incremented to it's way faster
var selectLastEventsSQL = `SELECT "SenderID", "EventType", "Payload", "TimeStamp" FROM "Channel_Event" WHERE "ChannelID" = (SELECT "ID" FROM "Channel" WHERE "ChannelID" = $1 AND "AppID" = $2) ORDER BY "ID" DESC LIMIT $3;`

var selectLastEventsSinceTimeStampSQL = `SELECT "SenderID", "EventType", "Payload", "TimeStamp" FROM (SELECT "SenderID", "EventType", "Payload", "TimeStamp" FROM "Channel_Event" WHERE "ChannelID" = (SELECT "ID" FROM "Channel" WHERE "ChannelID" = $1 AND "AppID" = $2) AND "TimeStamp" >= $3) as "t" ORDER BY "TimeStamp" ASC LIMIT $4;`
var selectLastEventsBeforeTimeStampSQL = `SELECT "SenderID", "EventType", "Payload", "TimeStamp" FROM (SELECT "SenderID", "EventType", "Payload", "TimeStamp" FROM "Channel_Event" WHERE "ChannelID" = (SELECT "ID" FROM "Channel" WHERE "ChannelID" = $1 AND "AppID" = $2) AND "TimeStamp" <= $3) as "t" ORDER BY "TimeStamp" DESC LIMIT $4;`

// NewSQLChannelRepository - Create a new instance of SQLChannelRepository
func NewSQLChannelRepository(db *DatabaseStorage) *ChannelRepository {
	return &ChannelRepository{dbHolder: db}
}

// ChannelRepository - SQL repository for tables Channel, Channel_Client and Channel_Event
type ChannelRepository struct {
	dbHolder *DatabaseStorage
}

// GetChannelClients - Get channel clients
func (repo *ChannelRepository) GetChannelClients(appID string, channelID string) ([]string, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectChannelClients)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetChannelClients: preparing statement failed: %v\n", err)
		return nil, err
	}

	rows, err := stmt.Query(channelID, appID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetChannelClients: query failed: %v\n", err)
		return nil, err
	}

	defer stmt.Close()

	clientIDs := make([]string, 0)

	for rows.Next() {
		var channelID string

		err = rows.Scan(&channelID)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "GetChannelClients: row scan failed: %v\n", err)
			return nil, err
		}

		clientIDs = append(clientIDs, channelID)
	}

	return clientIDs, nil
}

// CreateChannel - Insert new channel row
func (repo *ChannelRepository) CreateChannel(id string, appID string, name string, createdAt int64, isClosed bool, extra string, persistent bool, private bool, presence bool, push bool) error {
	stmt, err := repo.dbHolder.db.Prepare(createChannelSQL)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "CreateChannel: preparing statement failed: %v\n", err)
		return err
	}

	_, err = stmt.Exec(id, appID, name, createdAt, isClosed, extra, persistent, private, presence, push)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "CreateChannel: statement execution failed: %v\n", err)
		return err
	}

	defer stmt.Close()

	return nil
}

// DeleteChannel - Remove channel row
func (repo *ChannelRepository) DeleteChannel(appID string, id string) error {
	stmt, err := repo.dbHolder.db.Prepare(deleteChannelSQL)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "DeleteChannel: preparing statement failed: %v\n", err)
		return err
	}

	_, err = stmt.Exec(id, appID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "DeleteChannel: statement execution failed: %v\n", err)
		return err
	}

	defer stmt.Close()

	return nil
}

// DeleteAppChannels - Remove channel row
func (repo *ChannelRepository) DeleteAppChannels(appID string) error {
	stmt, err := repo.dbHolder.db.Prepare(deleteAppChannelsSQL)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "DeleteAppChannels: preparing statement failed: %v\n", err)
		return err
	}

	_, err = stmt.Exec(appID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "DeleteAppChannels: statement execution failed: %v\n", err)
		return err
	}

	defer stmt.Close()

	return nil
}

// JoinClient - Add client to channel
func (repo *ChannelRepository) JoinClient(appID string, channelID string, clientID string) error {
	stmt, err := repo.dbHolder.db.Prepare(joinChannelSQL)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "JoinChannel: preparing statement failed: %v\n", err)
		return err
	}

	_, err = stmt.Exec(channelID, clientID, appID)

	if err != nil {
		// If duplicate was found then ignore the error
		// We are already expecting this error
		// This way we can avoid query
		if !strings.Contains(err.Error(), "duplicate key value") {
			_, _ = fmt.Fprintf(os.Stderr, "JoinChannel: statement execution failed: %v\n", err)
			return err
		}

		_, _ = fmt.Fprintf(os.Stderr, "JoinChannel WARN: attempted to insert duplicate value: %v\n", err)
	}

	defer stmt.Close()

	return nil
}

// LeaveClient - Remove client to channel
func (repo *ChannelRepository) LeaveClient(appID string, channelID string, clientID string) error {
	stmt, err := repo.dbHolder.db.Prepare(leaveChannelSQL)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "LeaveChannel: preparing statement failed: %v\n", err)
		return err
	}

	_, err = stmt.Exec(channelID, clientID, appID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "LeaveChannel: statement execution failed: %v\n", err)
		return err
	}

	defer stmt.Close()

	return nil
}

// SetChannelCloseStatus - Set channel closed or open
func (repo *ChannelRepository) SetChannelCloseStatus(appID string, channelID string, isClosed bool) error {
	stmt, err := repo.dbHolder.db.Prepare(setCloseStatusSQL)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "SetChannelCloseStatus: preparing statement failed: %v\n", err)
		return err
	}

	_, err = stmt.Exec(isClosed, channelID, appID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "SetChannelCloseStatus: statement execution failed: %v\n", err)
		return err
	}

	defer stmt.Close()

	return nil
}

// GetClientAllowedChannels - Get all allowed channels for the given client, including public and private
func (repo *ChannelRepository) GetClientAllowedChannels(clientID string) ([]string, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectClientAllowedChannelsSQL)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetClientAllowedChannels: preparing statement failed: %v\n", err)
		return nil, err
	}

	rows, err := stmt.Query(clientID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetClientAllowedChannels: query failed: %v\n", err)
		return nil, err
	}

	defer stmt.Close()

	channelIDs := make([]string, 0)

	for rows.Next() {
		var channelID string

		err = rows.Scan(&channelID)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "GetClientAllowedChannels: row scan failed: %v\n", err)
			return nil, err
		}

		channelIDs = append(channelIDs, channelID)
	}

	return channelIDs, nil
}

// GetClientPublicChannels - Get all client public channels
func (repo *ChannelRepository) GetClientPublicChannels(clientID string) ([]*core.Channel, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectClientOpenOrPrivateChannels)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetClientPublicChannels: preparing statement failed: %v\n", err)
		return nil, err
	}

	rows, err := stmt.Query(false, clientID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetClientPublicChannels: query failed: %v\n", err)
		return nil, err
	}

	defer stmt.Close()

	channels := make([]*core.Channel, 0)

	for rows.Next() {
		chann, err := repo.rowToChannel(rows)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "GetClientPublicChannels: row scan failed: %v\n", err)
			return nil, err
		}

		channels = append(channels, chann)
	}

	return channels, nil
}

// GetClientPrivateChannels - Get all client private channels
func (repo *ChannelRepository) GetClientPrivateChannels(clientID string) ([]*core.Channel, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectClientOpenOrPrivateChannels)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetClientPrivateChannels: preparing statement failed: %v\n", err)
		return nil, err
	}

	rows, err := stmt.Query(true, clientID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetClientPrivateChannels: query failed: %v\n", err)
		return nil, err
	}

	defer stmt.Close()

	channels := make([]*core.Channel, 0)

	for rows.Next() {
		chann, err := repo.rowToChannel(rows)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "GetClientPrivateChannels: row scan failed: %v\n", err)
			return nil, err
		}

		channels = append(channels, chann)
	}

	return channels, nil
}

// GetAllPrivateChannels - Get all private channels without joined users
func (repo *ChannelRepository) GetAllPrivateChannels() ([]*core.Channel, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectAllOpenOrPrivateChannels)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetAllPrivateChannels: preparing statement failed: %v\n", err)
		return nil, err
	}

	rows, err := stmt.Query(true)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetAllPrivateChannels: query failed: %v\n", err)
		return nil, err
	}

	defer stmt.Close()

	channels := make([]*core.Channel, 0)

	for rows.Next() {
		chann, err := repo.rowToChannel(rows)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "GetAllPrivateChannels: row scan failed: %v\n", err)
			return nil, err
		}

		channels = append(channels, chann)
	}

	return channels, nil
}

// GetAllPublicChannels - Get all public channels without joined users
func (repo *ChannelRepository) GetAllPublicChannels() ([]*core.Channel, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectAllOpenOrPrivateChannels)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetAllPublicChannels: preparing statement failed: %v\n", err)
		return nil, err
	}

	rows, err := stmt.Query(false)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetAllPublicChannels: query failed: %v\n", err)
		return nil, err
	}

	defer stmt.Close()

	channels := make([]*core.Channel, 0)

	for rows.Next() {
		chann, err := repo.rowToChannel(rows)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "GetAllPublicChannels: row scan failed: %v\n", err)
			return nil, err
		}

		channels = append(channels, chann)
	}

	return channels, nil
}

// GetAppPrivateChannels - Get all app private channels without joined users
func (repo *ChannelRepository) GetAppPrivateChannels(appID string) ([]*core.Channel, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectOpenOrPrivateAppChannels)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetAppPrivateChannels: preparing statement failed: %v\n", err)
		return nil, err
	}

	rows, err := stmt.Query(true, appID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetAppPrivateChannels: query failed: %v\n", err)
		return nil, err
	}

	defer stmt.Close()

	channels := make([]*core.Channel, 0)

	for rows.Next() {
		chann, err := repo.rowToChannel(rows)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "GetAppPrivateChannels: row scan failed: %v\n", err)
			return nil, err
		}

		channels = append(channels, chann)
	}

	return channels, nil
}

// GetAppPublicChannels - Get all app public channels without joined users
func (repo *ChannelRepository) GetAppPublicChannels(appID string) ([]*core.Channel, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectOpenOrPrivateAppChannels)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetAppPublicChannels: preparing statement failed: %v\n", err)
		return nil, err
	}

	rows, err := stmt.Query(false, appID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetAppPublicChannels: query failed: %v\n", err)
		return nil, err
	}

	defer stmt.Close()

	channels := make([]*core.Channel, 0)

	for rows.Next() {
		chann, err := repo.rowToChannel(rows)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "GetAppPublicChannels: row scan failed: %v\n", err)
			return nil, err
		}

		channels = append(channels, chann)
	}

	return channels, nil
}

// GetAppChannel - Get channel with given AppID and ChannelID
func (repo *ChannelRepository) GetAppChannel(appID string, channelID string) (*core.Channel, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectAppChannel)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetAppChannel: preparing statement failed: %v\n", err)
		return nil, err
	}

	row := stmt.QueryRow(appID, channelID)

	defer stmt.Close()

	chann, err := repo.singleRowToChannel(row)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetAppChannel: row scan failed: %v\n", err)
		return nil, err
	}

	return chann, nil
}

// GetAppChannels - Get all app channels without joined users
func (repo *ChannelRepository) GetAppChannels(appID string) ([]*core.Channel, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectAppChannels)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetAppChannels: preparing statement failed: %v\n", err)
		return nil, err
	}

	rows, err := stmt.Query(appID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetAppChannels: query failed: %v\n", err)
		return nil, err
	}

	defer stmt.Close()

	channels := make([]*core.Channel, 0)

	for rows.Next() {
		chann, err := repo.rowToChannel(rows)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "GetAppChannels: row scan failed: %v\n", err)
			return nil, err
		}

		channels = append(channels, chann)
	}

	return channels, nil
}

// ExistsAppChannel - Check if app channel exists
func (repo *ChannelRepository) ExistsAppChannel(appID string, channelID string) (bool, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectAppChannelExists)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ExistsAppChannel: preparing statement failed: %v\n", err)
		return false, err
	}

	row := stmt.QueryRow(appID, channelID)

	defer stmt.Close()

	var amount uint64

	err = row.Scan(&amount)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ExistsAppChannel: row scan failed: %v\n", err)
		return false, err
	}

	return amount == 1, nil
}

// GetAppChannelsCount - Get how much clients an App has
func (repo *ChannelRepository) GetAppChannelsCount(appID string) (uint64, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectAppChannelAmount)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetAppChannelsCount: preparing statement failed: %v\n", err)
		return 0, err
	}

	row := stmt.QueryRow(appID)

	defer stmt.Close()

	var amount uint64

	err = row.Scan(&amount)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetAppChannelsCount: row scan failed: %v\n", err)
		return 0, err
	}

	return amount, nil
}

// GetAllChannels - Get all channels without joined users
func (repo *ChannelRepository) GetAllChannels() ([]*core.Channel, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectAllChannels)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetAllChannels: preparing statement failed: %v\n", err)
		return nil, err
	}

	rows, err := stmt.Query()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetAllChannels: query failed: %v\n", err)
		return nil, err
	}

	defer stmt.Close()

	channels := make([]*core.Channel, 0)

	for rows.Next() {
		chann, err := repo.rowToChannel(rows)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "GetAllChannels: row scan failed: %v\n", err)
			return nil, err
		}

		channels = append(channels, chann)
	}

	return channels, nil
}

// GetAllChannelsCount - Get how much clients an App has
func (repo *ChannelRepository) GetAllChannelsCount() (uint64, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectAllChannelsAmount)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetAllChannelsCount: preparing statement failed: %v\n", err)
		return 0, err
	}

	row := stmt.QueryRow()

	defer stmt.Close()

	var amount uint64

	err = row.Scan(&amount)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetAllChannelsCount: row scan failed: %v\n", err)
		return 0, err
	}

	return amount, nil
}

// AddChannelEvents - UNSUPPORTED
func (repo *ChannelRepository) AddChannelEvents(items []core.InsertItem) error {
	return nil
}

// AddChannelEvent - Add event to given channel
func (repo *ChannelRepository) AddChannelEvent(appID string, channelID string, event *core.ChannelEvent) error {

	stmt, err := repo.dbHolder.db.Prepare(addChannelEventSQL)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "AddChannelEvent: preparing statement failed: %v\n", err)
		return err
	}

	_, err = stmt.Exec(event.SenderID, event.EventType, event.Payload, channelID, event.Timestamp, appID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "AddChannelEvent: statement execution failed: %v\n", err)
		return err
	}

	defer stmt.Close()

	return nil
}

// GetChannelEventsAfter - Get all events after given timestamp
func (repo *ChannelRepository) GetChannelEventsAfter(appID string, channelID string, timestamp int64) ([]*core.ChannelEvent, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectEventsSinceTimeStampSQL)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetChannelEventsAfter: preparing statement failed: %v\n", err)
		return nil, err
	}

	rows, err := stmt.Query(appID, channelID, timestamp)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetChannelEventsAfter: query failed: %v\n", err)
		return nil, err
	}

	defer stmt.Close()

	channelEvents := make([]*core.ChannelEvent, 0)

	for rows.Next() {
		event, err := repo.rowToChannelEvent(channelID, rows)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "GetChannelEventsAfter: row scan failed: %v\n", err)
			return nil, err
		}

		channelEvents = append(channelEvents, event)
	}

	return channelEvents, nil
}

// GetChannelEventsAfterAndBefore - Get all events between given timestamps
func (repo *ChannelRepository) GetChannelEventsAfterAndBefore(appID string, channelID string, timestampAfter int64, timestampBefore int64) ([]*core.ChannelEvent, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectEventsBetweenTimeStampsSQL)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetChannelEventsAfterAndBefore: preparing statement failed: %v\n", err)
		return nil, err
	}

	rows, err := stmt.Query(channelID, appID, timestampAfter, timestampBefore)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetChannelEventsAfterAndBefore: query failed: %v\n", err)
		return nil, err
	}

	defer stmt.Close()

	channelEvents := make([]*core.ChannelEvent, 0)

	for rows.Next() {
		event, err := repo.rowToChannelEvent(channelID, rows)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "GetChannelEventsAfterAndBefore: row scan failed: %v\n", err)
			return nil, err
		}

		channelEvents = append(channelEvents, event)
	}

	return channelEvents, nil
}

// GetChannelLastEventsAfter - Get an given amount events after given timestamp
func (repo *ChannelRepository) GetChannelLastEventsBefore(appID string, channelID string, amount int64, timestamp int64) ([]*core.ChannelEvent, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectLastEventsBeforeTimeStampSQL)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetChannelLastEventsAfter: preparing statement failed: %v\n", err)
		return nil, err
	}

	rows, err := stmt.Query(channelID, appID, timestamp, amount)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetChannelLastEventsAfter: query failed: %v\n", err)
		return nil, err
	}

	defer stmt.Close()

	channelEvents := make([]*core.ChannelEvent, 0)

	for rows.Next() {
		event, err := repo.rowToChannelEvent(channelID, rows)

		if err != nil {
			fmt.Fprintf(os.Stderr, "GetChannelLastEventsAfter: row scan failed: %v\n", err)
			return nil, err
		}

		channelEvents = append(channelEvents, event)
	}

	return channelEvents, nil
}

// GetChannelLastEventsAfter - Get an given amount events after given timestamp
func (repo *ChannelRepository) GetChannelLastEventsAfter(appID string, channelID string, amount int64, timestamp int64) ([]*core.ChannelEvent, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectLastEventsSinceTimeStampSQL)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GetChannelLastEventsAfter: preparing statement failed: %v\n", err)
		return nil, err
	}

	rows, err := stmt.Query(channelID, appID, timestamp, amount)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetChannelLastEventsAfter: query failed: %v\n", err)
		return nil, err
	}

	defer stmt.Close()

	channelEvents := make([]*core.ChannelEvent, 0)

	for rows.Next() {
		event, err := repo.rowToChannelEvent(channelID, rows)

		if err != nil {
			fmt.Fprintf(os.Stderr, "GetChannelLastEventsAfter: row scan failed: %v\n", err)
			return nil, err
		}

		channelEvents = append(channelEvents, event)
	}

	return channelEvents, nil
}

// GetChannelLastEvents - Get last events
func (repo *ChannelRepository) GetChannelLastEvents(appID string, channelID string, amount int64) ([]*core.ChannelEvent, error) {
	stmt, err := repo.dbHolder.db.Prepare(selectLastEventsSQL)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetChannelLastEvents: preparing statement failed: %v\n", err)
		return nil, err
	}

	rows, err := stmt.Query(channelID, appID, amount)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetChannelLastEvents: query failed: %v\n", err)
		return nil, err
	}

	defer stmt.Close()

	channelEvents := make([]*core.ChannelEvent, 0)

	for rows.Next() {
		event, err := repo.rowToChannelEvent(channelID, rows)

		if err != nil {
			fmt.Fprintf(os.Stderr, "GetChannelLastEvents: row scan failed: %v\n", err)
			return nil, err
		}

		channelEvents = append(channelEvents, event)
	}

	return channelEvents, nil
}

// rowToChannelEvent - Small helper to keep code cleaner
func (repo *ChannelRepository) rowToChannelEvent(channelID string, rows *sql.Rows) (*core.ChannelEvent, error) {
	//var id string
	var senderID string
	var eventType string
	var payload string
	//var channelID string
	var timestamp int64

	//err := rows.Scan(&id, &senderID, &eventType, &payload, &channelID, &timestamp)
	err := rows.Scan(&senderID, &eventType, &payload, &timestamp)

	channEvent := &core.ChannelEvent{
		SenderID:  senderID,
		EventType: eventType,
		Payload:   payload,
		ChannelID: channelID,
		Timestamp: timestamp,
	}

	return channEvent, err
}

// rowToChannel - Small helper to keep code cleaner
func (repo *ChannelRepository) rowToChannel(rows *sql.Rows) (*core.Channel, error) {
	var id string
	var appID string
	var name string
	var createdAt int64
	var isClosed bool
	var extra string
	var persistent bool
	var private bool
	var presence bool
	var push bool

	err := rows.Scan(&id, &appID, &name, &createdAt, &isClosed, &extra, &persistent, &private, &presence, &push)

	chann := &core.Channel{
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
	}

	return chann, err
}

func (repo *ChannelRepository) singleRowToChannel(rows *sql.Row) (*core.Channel, error) {
	var id string
	var appID string
	var name string
	var createdAt int64
	var isClosed bool
	var extra string
	var persistent bool
	var private bool
	var presence bool
	var push bool

	err := rows.Scan(&id, &appID, &name, &createdAt, &isClosed, &extra, &persistent, &private, &presence, &push)

	chann := &core.Channel{
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
	}

	return chann, err
}
