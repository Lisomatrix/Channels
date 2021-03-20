package redishandler

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

var messageEnqueuedQueueSuffix = "-q"
var messageProcessingQueueSuffix = "-p"

var publishConnectedMessage = "cn"

var fetchMessageTimeOut time.Duration = 10

var redisHandler = NewRedisHandler()

// GetRedisHandler - Get unique redis handler instance
func GetRedisHandler() *RedisHandler {
	return redisHandler
}

// RedisHandler responsible for handling communications with redis
// And some helping functions
type RedisHandler struct {
	rdb    *redis.Client
	pubSub *redis.PubSub
}

// Subscribe to a specific channel
func (handler *RedisHandler) Subscribe(channel string) *redis.PubSub {
	pubSubHandler := handler.rdb.Subscribe(ctx, channel)

	return pubSubHandler
}

// Unsubscribe a specific channel
func (handler *RedisHandler) Unsubscribe(channel string, pubSubHandler *redis.PubSub) {

	_ = pubSubHandler.Unsubscribe(ctx, channel)
	_ = pubSubHandler.Close()
}

// UpdateHeartBeat of a connected client
func (handler *RedisHandler) UpdateHeartBeat(clientID string) {
	now := time.Now()
	sec := now.Unix()

	handler.rdb.Set(ctx, clientID+"-hb", sec, 0)
}

var javascriptISOString = "2006-01-02T15:04:05.999Z07:00"

// GetUsersOnlineStatus from the the given clientIDs
func (handler *RedisHandler) GetUsersOnlineStatus(clientIDs []string) map[string]string {

	clientIDsSlice := make([]string, 0, len(clientIDs))

	for index, clientID := range clientIDs {
		clientIDsSlice[index] = clientID + "-hb"
	}

	//clientIDsSlice := clientIDs
	onlineResults := make(map[string]string)

	if results, err := handler.rdb.MGet(ctx, clientIDsSlice...).Result(); err != nil {
		log.Println(err)
	} else {

		for index, timeStampBin := range results {

			var lastHeartBeat time.Time
			//currentTime := time.Now()
			if timeStampBin == nil {
				onlineResults[clientIDsSlice[index]] = "unknown"
				continue
			}

			parsedTimestamp, err := strconv.ParseInt(timeStampBin.(string), 10, 64)

			if err != nil {
				log.Println(err)
				onlineResults[clientIDsSlice[index]] = "unknown"

				continue
			}

			lastHeartBeat = time.Unix(parsedTimestamp, 0)

			// If passed more than 2 minutes since last heart beat then we consider it offline

			onlineResults[clientIDsSlice[index]] = lastHeartBeat.UTC().Format(javascriptISOString)
		}
	}

	return onlineResults
}

// Publish an event
func (handler *RedisHandler) Publish(channel string, message []byte) {
	handler.rdb.Publish(ctx, channel, message)
}

// ClientConnected publish an event warning others that this client connected
func (handler *RedisHandler) ClientConnected(clientID string) {
	handler.rdb.Publish(ctx, clientID, publishConnectedMessage)
}

// DeleteUserQueues of the given clientID
func (handler *RedisHandler) DeleteUserQueues(clientID string) {
	handler.rdb.Del(ctx, clientID+messageEnqueuedQueueSuffix)
	handler.rdb.Del(ctx, clientID+messageProcessingQueueSuffix)
}

// EnqueueMessage into the waiting messages queue of this user if the user devices aren't in redis
// Then a fetch request will be sent
func (handler *RedisHandler) EnqueueMessage(clientID string, message []byte) bool {

	handler.rdb.LPush(ctx, clientID+messageEnqueuedQueueSuffix, message)
	/*
		devices, isOK := handler.GetUserDevices(clientID)

		if !isOK || len(devices) == 0 {
			devices, isOK = handler.RequestUserDevices(clientID)

			if !isOK {
				return false
			}
		}

		for _, device := range devices {
			deviceListKey := clientID + "-" + device.GetDeviceID() + messageEnqueuedQueueSuffix
			handler.rdb.LPush(ctx, deviceListKey, message)
		}*/

	return true
}

// FetchMessage from the user queue
func (handler *RedisHandler) FetchMessage(clientID string, deviceID string) []byte {
	// TODO: KEEP THE ORIGINAL
	//data, _ := handler.rdb.BRPopLPush(ctx, clientID+"-"+deviceID+messageEnqueuedQueueSuffix, clientID+"-"+deviceID+messageProcessingQueueSuffix, fetchMessageTimeOut*time.Second).Bytes()
	data, _ := handler.rdb.BRPopLPush(ctx, clientID+messageEnqueuedQueueSuffix, clientID+messageProcessingQueueSuffix, fetchMessageTimeOut*time.Second).Bytes()

	return data
}

// NoAckReceived send not ack message back to waiting queue
func (handler *RedisHandler) NoAckReceived(clientID string, deviceID string) {
	if _, err := handler.rdb.RPop(ctx, clientID+"-"+deviceID+messageProcessingQueueSuffix).Result(); err != nil {
		log.Fatal(err)
	} else {
		_, _ = handler.rdb.RPush(ctx, clientID+"-"+deviceID+messageEnqueuedQueueSuffix).Result()
	}

}

// MessageAck remove ack message from processing queue
func (handler *RedisHandler) MessageAck(clientID string, deviceID string) {
	_, _ = handler.rdb.RPop(ctx, clientID+"-"+deviceID+messageProcessingQueueSuffix).Result()
}

// NewRedisHandler instiantiate a new RedisHandler
func NewRedisHandler() *RedisHandler {

	redisHandler := new(RedisHandler)

	redisHandler.rdb = redis.NewClient(
		&redis.Options{
			Addr: "127.0.0.1:6379",
			// Password: "CulP3gnpgSAxFlbjO/JrNCR/uTKFKvTLbW7gJoVQfg1sh1BmzeNBUs5TsXy0Q7YDgGbfazSZy5LKnU3l", // no password set
			DB:       0,
			PoolSize: 5,
		})

	return redisHandler
}
