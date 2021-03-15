package pgxsql

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Channels/Channels/core"
	"github.com/jackc/pgx/v4"
)

// Channel SQL
var selectChannelClients = `SELECT "clientID" FROM "Channel_Client" WHERE "channelID" = (SELECT "ID" FROM "Channel" WHERE "ChannelID" = $1 AND "AppID" = $2 LIMIT 1);`
var createChannelSQL = `INSERT INTO "Channel"("ChannelID", "AppID", "Name", "Created_At", "IsClosed", "Extra", "Persistent", "Private", "Presence") VALUES ( $1, $2, $3, $4, $5, $6, $7, $8, $9);`
var deleteChannelSQL = `DELETE FROM "Channel" WHERE "ChannelID" = $1 AND "AppID" = $2;`
var deleteAppChannelsSQL = `DELETE FROM "Channel" WHERE "AppID" = $1;`
var joinChannelSQL = `INSERT INTO public."Channel_Client"("clientID", "channelID") VALUES ($2, (SELECT "ID" FROM "Channel" WHERE "ChannelID" = $1 AND "AppID" = $3 LIMIT 1));`
var leaveChannelSQL = `DELETE FROM "Channel_Client" WHERE "channelID" = (SELECT "ID" FROM "Channel" WHERE "ChannelID" = $1 AND "AppID" = $3 LIMIT 1) AND "clientID" = $2;`
var setCloseStatusSQL = `UPDATE "Channel" SET "IsClosed" = $1 WHERE "ChannelID" = $2 AND "AppID" = $3;`
var selectClientAllowedChannelsSQL = `SELECT "ChannelID" FROM "Channel" WHERE "Private" = false AND "ID" IN (SELECT "channelID" FROM "Channel_Client" WHERE "clientID" = $1);`
var selectClientOpenOrPrivateChannels = `SELECT "ChannelID", "AppID", "Name", "Created_At", "IsClosed", "Extra", "Persistent", "Private", "Presence" FROM "Channel" WHERE "Private" = $1 AND "ChannelID" IN (SELECT "channelID" FROM "Channel_Client" WHERE "clientID" = $2);`
var selectOpenOrPrivateAppChannels = `SELECT "ChannelID", "AppID", "Name", "Created_At", "IsClosed", "Extra", "Persistent", "Private", "Presence" FROM "Channel" WHERE "Private" = $1 AND "AppID" = $2;`
var selectAppChannels = `SELECT "ChannelID", "AppID", "Name", "Created_At", "IsClosed", "Extra", "Persistent", "Private", "Presence" FROM "Channel" WHERE "AppID" = $1;`
var selectAppChannelAmount = `SELECT COUNT("ChannelID") FROM "Channel" WHERE "AppID" = $1;`
var selectAppChannelExists = `SELECT COUNT("ChannelID") FROM "Channel" WHERE "AppID" = $1 AND "ChannelID" = $2 LIMIT 1;`
var selectAllOpenOrPrivateChannels = `SELECT "ChannelID", "AppID", "Name", "Created_At", "IsClosed", "Extra", "Persistent", "Private", "Presence" FROM "Channel" WHERE "Private" = $1;`
var selectAllChannels = `SELECT "ChannelID", "AppID", "Name", "Created_At", "IsClosed", "Extra", "Persistent", "Private", "Presence" FROM "Channel";`
var selectAllChannelsAmount = `SELECT COUNT("ChannelID") FROM "Channel"`
var selectAppChannel = `SELECT "ChannelID", "AppID", "Name", "Created_At", "IsClosed", "Extra", "Persistent", "Private", "Presence" FROM "Channel" WHERE "AppID" = $1 AND "ChannelID" = $2;`
var addChannelEventSQL = `INSERT INTO "Channel_Event"("SenderID", "EventType", "Payload", "ChannelID", "TimeStamp") VALUES ( $1 , $2 , $3 , (SELECT "ID" FROM "Channel" WHERE "ChannelID" = $4 AND "AppID" = $6 LIMIT 1) , $5 );`

var selectEventsSinceTimeStampSQL = `SELECT "SenderID", "EventType", "Payload", "TimeStamp" FROM "Channel_Event" WHERE "ChannelID" = (SELECT "ID" FROM "Channel" WHERE "ChannelID" = $1 AND "AppID" = $2) AND "TimeStamp" >= $3;`
var selectEventsBetweenTimeStampsSQL = `SELECT "SenderID", "EventType", "Payload", "TimeStamp" FROM "Channel_Event" WHERE "ChannelID" = (SELECT "ID" FROM "Channel" WHERE "ChannelID" = $1 AND "AppID" = $2) AND "TimeStamp" >= $3 AND "TimeStamp" <= $4;`

// * The new one is based on primary key since its auto incremented to it's way faster
var selectLastEventsSQL = `SELECT "SenderID", "EventType", "Payload", "TimeStamp" FROM "Channel_Event" WHERE "ChannelID" = (SELECT "ID" FROM "Channel" WHERE "ChannelID" = $1 AND "AppID" = $2) ORDER BY "ID" DESC LIMIT $3;`

var selectLastEventsSinceTimeStampSQL = `SELECT "SenderID", "EventType", "Payload", "TimeStamp" FROM (SELECT "SenderID", "EventType", "Payload", "TimeStamp" FROM "Channel_Event" WHERE "ChannelID" = (SELECT "ID" FROM "Channel" WHERE "ChannelID" = $1 AND "AppID" = $2) AND "TimeStamp" >= $3) as "t" ORDER BY "TimeStamp" DESC LIMIT $4;`

// NewSQLChannelRepository - Create a new instance of SQLChannelRepository
func NewSQLChannelRepository(db *PGXDatabaseStorage) *PGXChannelRepository {
	return &PGXChannelRepository{
		dbHolder: db,
		ctx:      context.Background(),
	}
}

// PGXChannelRepository - SQL repository for tables Channel, Channel_Client and Channel_Event
type PGXChannelRepository struct {
	dbHolder *PGXDatabaseStorage
	ctx      context.Context
}

// GetChannelClients - Get channel clients
func (repo *PGXChannelRepository) GetChannelClients(appID string, channelID string) ([]string, error) {
	rows, err := repo.dbHolder.db.Query(repo.ctx, selectChannelClients, channelID, appID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetChannelClients: query failed: %v\n", err)
		return nil, err
	}

	clientIDs := make([]string, 0)

	for rows.Next() {
		var channelID string

		err = rows.Scan(&channelID)

		if err != nil {
			fmt.Fprintf(os.Stderr, "GetChannelClients: row scan failed: %v\n", err)
			return nil, err
		}

		clientIDs = append(clientIDs, channelID)
	}

	return clientIDs, nil
}

// CreateChannel - Insert new channel row
func (repo *PGXChannelRepository) CreateChannel(id string, appID string, name string, createdAt int64, isClosed bool, extra string, persistent bool, private bool, presence bool) error {
	_, err := repo.dbHolder.db.Exec(repo.ctx, createChannelSQL, id, appID, name, createdAt, isClosed, extra, persistent, private, presence)

	if err != nil {
		fmt.Fprintf(os.Stderr, "CreateChannel: statement execution failed: %v\n", err)
		return err
	}

	return nil
}

// DeleteChannel - Remove channel row
func (repo *PGXChannelRepository) DeleteChannel(appID string, id string) error {

	_, err := repo.dbHolder.db.Exec(repo.ctx, deleteChannelSQL, id, appID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "DeleteChannel: statement execution failed: %v\n", err)
		return err
	}

	return nil
}

// DeleteAppChannels - Remove channel row
func (repo *PGXChannelRepository) DeleteAppChannels(appID string) error {
	_, err := repo.dbHolder.db.Exec(repo.ctx, deleteAppChannelsSQL, appID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "DeleteAppChannels: statement execution failed: %v\n", err)
		return err
	}

	return nil
}

// JoinClient - Add client to channel
func (repo *PGXChannelRepository) JoinClient(appID string, channelID string, clientID string) error {
	_, err := repo.dbHolder.db.Exec(repo.ctx, joinChannelSQL, channelID, clientID, appID)

	if err != nil {
		// If duplicate was found then ignore the error
		// We are already expecting this error
		// This way we can avoid query
		if !strings.Contains(err.Error(), "duplicate key value") {
			fmt.Fprintf(os.Stderr, "JoinChannel: statement execution failed: %v\n", err)
			return err
		}

		fmt.Fprintf(os.Stderr, "JoinChannel WARN: attempted to insert duplicate value: %v\n", err)
	}

	return nil
}

// LeaveClient - Remove client to channel
func (repo *PGXChannelRepository) LeaveClient(appID string, channelID string, clientID string) error {
	_, err := repo.dbHolder.db.Exec(repo.ctx, leaveChannelSQL, channelID, clientID, appID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "LeaveChannel: statement execution failed: %v\n", err)
		return err
	}

	return nil
}

// SetChannelCloseStatus - Set channel closed or open
func (repo *PGXChannelRepository) SetChannelCloseStatus(appID string, channelID string, isClosed bool) error {
	_, err := repo.dbHolder.db.Exec(repo.ctx, setCloseStatusSQL, isClosed, channelID, appID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "SetChannelCloseStatus: statement execution failed: %v\n", err)
		return err
	}

	return nil
}

// GetClientAllowedChannels - Get all allowed channels for the given client, including public and private
func (repo *PGXChannelRepository) GetClientAllowedChannels(clientID string) ([]string, error) {
	rows, err := repo.dbHolder.db.Query(repo.ctx, selectClientAllowedChannelsSQL, clientID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetClientAllowedChannels: query failed: %v\n", err)
		return nil, err
	}

	channelIDs := make([]string, 0)

	for rows.Next() {
		var channelID string

		err = rows.Scan(&channelID)

		if err != nil {
			fmt.Fprintf(os.Stderr, "GetClientAllowedChannels: row scan failed: %v\n", err)
			return nil, err
		}

		channelIDs = append(channelIDs, channelID)
	}

	return channelIDs, nil
}

// GetClientPublicChannels - Get all client public channels
func (repo *PGXChannelRepository) GetClientPublicChannels(clientID string) ([]*core.Channel, error) {
	rows, err := repo.dbHolder.db.Query(repo.ctx, selectClientOpenOrPrivateChannels, false, clientID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetClientPublicChannels: query failed: %v\n", err)
		return nil, err
	}

	channels := make([]*core.Channel, 0)

	for rows.Next() {
		chann, err := repo.rowToChannel(rows)

		if err != nil {
			fmt.Fprintf(os.Stderr, "GetClientPublicChannels: row scan failed: %v\n", err)
			return nil, err
		}

		channels = append(channels, chann)
	}

	return channels, nil
}

// GetClientPrivateChannels - Get all client private channels
func (repo *PGXChannelRepository) GetClientPrivateChannels(clientID string) ([]*core.Channel, error) {
	rows, err := repo.dbHolder.db.Query(repo.ctx, selectClientOpenOrPrivateChannels, true, clientID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetClientPrivateChannels: query failed: %v\n", err)
		return nil, err
	}

	channels := make([]*core.Channel, 0)

	for rows.Next() {
		chann, err := repo.rowToChannel(rows)

		if err != nil {
			fmt.Fprintf(os.Stderr, "GetClientPrivateChannels: row scan failed: %v\n", err)
			return nil, err
		}

		channels = append(channels, chann)
	}

	return channels, nil
}

// GetAllPrivateChannels - Get all private channels without joined users
func (repo *PGXChannelRepository) GetAllPrivateChannels() ([]*core.Channel, error) {
	rows, err := repo.dbHolder.db.Query(repo.ctx, selectAllOpenOrPrivateChannels, true)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetAllPrivateChannels: query failed: %v\n", err)
		return nil, err
	}

	channels := make([]*core.Channel, 0)

	for rows.Next() {
		chann, err := repo.rowToChannel(rows)

		if err != nil {
			fmt.Fprintf(os.Stderr, "GetAllPrivateChannels: row scan failed: %v\n", err)
			return nil, err
		}

		channels = append(channels, chann)
	}

	return channels, nil
}

// GetAllPublicChannels - Get all public channels without joined users
func (repo *PGXChannelRepository) GetAllPublicChannels() ([]*core.Channel, error) {
	rows, err := repo.dbHolder.db.Query(repo.ctx, selectAllOpenOrPrivateChannels, false)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetAllPublicChannels: query failed: %v\n", err)
		return nil, err
	}

	channels := make([]*core.Channel, 0)

	for rows.Next() {
		chann, err := repo.rowToChannel(rows)

		if err != nil {
			fmt.Fprintf(os.Stderr, "GetAllPublicChannels: row scan failed: %v\n", err)
			return nil, err
		}

		channels = append(channels, chann)
	}

	return channels, nil
}

// GetAppPrivateChannels - Get all app private channels without joined users
func (repo *PGXChannelRepository) GetAppPrivateChannels(appID string) ([]*core.Channel, error) {
	rows, err := repo.dbHolder.db.Query(repo.ctx, selectOpenOrPrivateAppChannels, true, appID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetAppPrivateChannels: query failed: %v\n", err)
		return nil, err
	}

	channels := make([]*core.Channel, 0)

	for rows.Next() {
		chann, err := repo.rowToChannel(rows)

		if err != nil {
			fmt.Fprintf(os.Stderr, "GetAppPrivateChannels: row scan failed: %v\n", err)
			return nil, err
		}

		channels = append(channels, chann)
	}

	return channels, nil
}

// GetAppPublicChannels - Get all app public channels without joined users
func (repo *PGXChannelRepository) GetAppPublicChannels(appID string) ([]*core.Channel, error) {
	rows, err := repo.dbHolder.db.Query(repo.ctx, selectOpenOrPrivateAppChannels, false, appID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetAppPublicChannels: query failed: %v\n", err)
		return nil, err
	}

	channels := make([]*core.Channel, 0)

	for rows.Next() {
		chann, err := repo.rowToChannel(rows)

		if err != nil {
			fmt.Fprintf(os.Stderr, "GetAppPublicChannels: row scan failed: %v\n", err)
			return nil, err
		}

		channels = append(channels, chann)
	}

	return channels, nil
}

// GetAppChannel - Get channel with given AppID and ChannelID
func (repo *PGXChannelRepository) GetAppChannel(appID string, channelID string) (*core.Channel, error) {
	row := repo.dbHolder.db.QueryRow(repo.ctx, selectAppChannel, appID, channelID)

	chann, err := repo.singleRowToChannel(row)

	if err != nil {

		// TODO: Need a better fix for this
		if err.Error() == "no rows in result set" {
			return nil, nil
		}

		fmt.Fprintf(os.Stderr, "GetAppChannel: row scan failed: %v\n", err)
		return nil, err
	}

	return chann, nil
}

// GetAppChannels - Get all app channels without joined users
func (repo *PGXChannelRepository) GetAppChannels(appID string) ([]*core.Channel, error) {
	rows, err := repo.dbHolder.db.Query(repo.ctx, selectAppChannels, appID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetAppChannels: query failed: %v\n", err)
		return nil, err
	}

	channels := make([]*core.Channel, 0)

	for rows.Next() {
		chann, err := repo.rowToChannel(rows)

		if err != nil {
			fmt.Fprintf(os.Stderr, "GetAppChannels: row scan failed: %v\n", err)
			return nil, err
		}

		channels = append(channels, chann)
	}

	return channels, nil
}

// ExistsAppChannel - Check if app channel exists
func (repo *PGXChannelRepository) ExistsAppChannel(appID string, channelID string) (bool, error) {
	row := repo.dbHolder.db.QueryRow(repo.ctx, selectAppChannelExists, appID, channelID)

	var amount uint64

	err := row.Scan(&amount)

	if err != nil {
		fmt.Fprintf(os.Stderr, "ExistsAppChannel: row scan failed: %v\n", err)
		return false, err
	}

	return amount == 1, nil
}

// GetAppChannelsCount - Get how much clients an App has
func (repo *PGXChannelRepository) GetAppChannelsCount(appID string) (uint64, error) {
	row := repo.dbHolder.db.QueryRow(repo.ctx, selectAppChannelAmount, appID)

	var amount uint64

	err := row.Scan(&amount)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetAppChannelsCount: row scan failed: %v\n", err)
		return 0, err
	}

	return amount, nil
}

// GetAllChannels - Get all channels without joined users
func (repo *PGXChannelRepository) GetAllChannels() ([]*core.Channel, error) {
	rows, err := repo.dbHolder.db.Query(repo.ctx, selectAllChannels)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetAllChannels: query failed: %v\n", err)
		return nil, err
	}

	channels := make([]*core.Channel, 0)

	for rows.Next() {
		chann, err := repo.rowToChannel(rows)

		if err != nil {
			fmt.Fprintf(os.Stderr, "GetAllChannels: row scan failed: %v\n", err)
			return nil, err
		}

		channels = append(channels, chann)
	}

	return channels, nil
}

// GetAllChannelsCount - Get how much clients an App has
func (repo *PGXChannelRepository) GetAllChannelsCount() (uint64, error) {
	row := repo.dbHolder.db.QueryRow(repo.ctx, selectAllChannelsAmount)

	var amount uint64

	err := row.Scan(&amount)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetAllChannelsCount: row scan failed: %v\n", err)
		return 0, err
	}

	return amount, nil
}

// AddChannelEvent - Add event to given channel
func (repo *PGXChannelRepository) AddChannelEvent(appID string, channelID string, event *core.ChannelEvent) error {

	_, err := repo.dbHolder.db.Exec(repo.ctx, addChannelEventSQL, event.SenderID, event.EventType, event.Payload, channelID, event.Timestamp, appID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "AddChannelEvent: statement execution failed: %v\n", err)
		return err
	}

	return nil
}

// AddChannelEvents - Add a batch of events
func (repo *PGXChannelRepository) AddChannelEvents(items []core.InsertItem) error {

	batch := &pgx.Batch{}

	for _, item := range items {
		event := item.Event
		batch.Queue(addChannelEventSQL, event.SenderID, event.EventType, event.Payload, event.ChannelID, event.Timestamp, item.AppID)
	}
	conn, err := repo.dbHolder.db.Acquire(repo.ctx)

	if err != nil {
		fmt.Fprintf(os.Stderr, "AddChannelEvents: failed to acquire connection: %v\n", err)
		return err
	}

	br := conn.SendBatch(repo.ctx, batch)

	// br := repo.dbHolder.db.SendBatch(repo.ctx, batch)

	_, err = br.Exec()
	conn.Release()

	if err != nil {
		fmt.Fprintf(os.Stderr, "AddChannelEvents: batch execution failed: %v\n", err)
		return err
	}

	return nil
}

// GetChannelEventsAfter - Get all events after given timestamp
func (repo *PGXChannelRepository) GetChannelEventsAfter(appID string, channelID string, timestamp int64) ([]*core.ChannelEvent, error) {
	rows, err := repo.dbHolder.db.Query(repo.ctx, selectEventsSinceTimeStampSQL, appID, channelID, timestamp)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetChannelEventsAfter: query failed: %v\n", err)
		return nil, err
	}

	channelEvents := make([]*core.ChannelEvent, 0)

	for rows.Next() {
		event, err := repo.rowToChannelEvent(channelID, rows)

		if err != nil {
			fmt.Fprintf(os.Stderr, "GetChannelEventsAfter: row scan failed: %v\n", err)
			return nil, err
		}

		channelEvents = append(channelEvents, event)
	}

	return channelEvents, nil
}

// GetChannelEventsAfterAndBefore - Get all events between given timestamps
func (repo *PGXChannelRepository) GetChannelEventsAfterAndBefore(appID string, channelID string, timestampAfter int64, timestampBefore int64) ([]*core.ChannelEvent, error) {
	rows, err := repo.dbHolder.db.Query(repo.ctx, selectEventsBetweenTimeStampsSQL, channelID, appID, timestampAfter, timestampBefore)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetChannelEventsAfterAndBefore: query failed: %v\n", err)
		return nil, err
	}

	channelEvents := make([]*core.ChannelEvent, 0)

	for rows.Next() {
		event, err := repo.rowToChannelEvent(channelID, rows)

		if err != nil {
			fmt.Fprintf(os.Stderr, "GetChannelEventsAfterAndBefore: row scan failed: %v\n", err)
			return nil, err
		}

		channelEvents = append(channelEvents, event)
	}

	return channelEvents, nil
}

// GetChannelLastEventsAfter - Get an given amount events after given timestamp
func (repo *PGXChannelRepository) GetChannelLastEventsAfter(appID string, channelID string, amount int64, timestamp int64) ([]*core.ChannelEvent, error) {
	rows, err := repo.dbHolder.db.Query(repo.ctx, selectLastEventsSinceTimeStampSQL, channelID, appID, timestamp, amount)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetChannelLastEventsAfter: query failed: %v\n", err)
		return nil, err
	}

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
func (repo *PGXChannelRepository) GetChannelLastEvents(appID string, channelID string, amount int64) ([]*core.ChannelEvent, error) {
	rows, err := repo.dbHolder.db.Query(repo.ctx, selectLastEventsSQL, channelID, appID, amount)

	if err != nil {
		fmt.Fprintf(os.Stderr, "GetChannelLastEvents: query failed: %v\n", err)
		return nil, err
	}

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
func (repo *PGXChannelRepository) rowToChannelEvent(channelID string, rows pgx.Rows) (*core.ChannelEvent, error) {
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

	/*
		channEvent := &core.ChannelEvent{
			ID:        id,
			SenderID:  senderID,
			EventType: eventType,
			Payload:   payload,
			ChannelID: channelID,
			Timestamp: timestamp,
		}*/

	return channEvent, err
}

// rowToChannel - Small helper to keep code cleaner
func (repo *PGXChannelRepository) rowToChannel(rows pgx.Rows) (*core.Channel, error) {
	var id string
	var appID string
	var name string
	var createdAt int64
	var isClosed bool
	var extra string
	var persistent bool
	var private bool
	var presence bool

	err := rows.Scan(&id, &appID, &name, &createdAt, &isClosed, &extra, &persistent, &private, &presence)

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
	}

	return chann, err
}

func (repo *PGXChannelRepository) singleRowToChannel(rows pgx.Row) (*core.Channel, error) {
	var id string
	var appID string
	var name string
	var createdAt int64
	var isClosed bool
	var extra string
	var persistent bool
	var private bool
	var presence bool

	err := rows.Scan(&id, &appID, &name, &createdAt, &isClosed, &extra, &persistent, &private, &presence)

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
	}

	return chann, err
}
