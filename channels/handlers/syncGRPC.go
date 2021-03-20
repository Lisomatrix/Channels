package handlers

import (
	"context"
	"fmt"
	"github.com/lisomatrix/channels/channels/auth"
	"github.com/lisomatrix/channels/channels/core"
	"os"
)

type SyncServer struct {}

func (srv *SyncServer) GetMessagesBetween(ctx context.Context, req *GetMessagesBetweenRequest) (*GetChannelEventsResponse, error) {
	// Check if is admin, and validate token
	identity, isOK := auth.VerifyToken(req.Token)

	// If not valid return
	if !isOK || !identity.CanUseAppID(req.AppID) {
		return nil, &AuthError{reason: "Invalid token"}
	}

	// Check if channel exists
	exists, err := core.GetEngine().GetChannelRepository().ExistsAppChannel(req.AppID, req.ChannelID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GRPC Get messages between timestamps: failed to check app channel existence %v\n", err)
		return nil, &InternalError{reason: err.Error()}
	}

	if !exists {
		return nil, &NotFoundError{}
	}

	// Get events
	events, err := core.GetEngine().GetChannelRepository().GetChannelEventsAfterAndBefore(req.AppID, req.ChannelID, req.FromTimeStamp, req.ToTimeStamp)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GRPC Get messages between timestamps: failed fetch events %v\n", err)
		return nil, &InternalError{reason: err.Error()}
	}

	outEvents := make([]*OutChannelEvent, len(events))

	for _, event := range events {
		outEvents = append(outEvents, &OutChannelEvent{
			SenderID:      event.SenderID,
			EventType:     event.EventType,
			Payload:       event.Payload,
			ChannelID:     event.ChannelID,
			Timestamp:     event.Timestamp,
		})
	}

	response := GetChannelEventsResponse{Events: outEvents}

	return &response, nil
}

func (srv *SyncServer) GetMessagesSince(ctx context.Context, req *GetMessagesSinceRequest) (*GetChannelEventsResponse, error) {
	// Check if is admin, and validate token
	identity, isOK := auth.VerifyToken(req.Token)

	// If not valid return
	if !isOK || !identity.CanUseAppID(req.AppID) {
		return nil, &AuthError{reason: "Invalid token"}
	}

	// Check if channel exists
	exists, err := core.GetEngine().GetChannelRepository().ExistsAppChannel(req.AppID, req.ChannelID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GRPC Get messages between timestamps: failed to check app channel existence %v\n", err)
		return nil, &InternalError{reason: err.Error()}
	}

	if !exists {
		return nil, &NotFoundError{}
	}

	// Get events
	events, err := core.GetEngine().GetChannelRepository().GetChannelEventsAfter(req.AppID, req.ChannelID, req.FromTimeStamp)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GRPC Get messages since timestamp: failed fetch events %v\n", err)
		return nil, &InternalError{reason: err.Error()}
	}

	outEvents := make([]*OutChannelEvent, len(events))

	for _, event := range events {
		outEvents = append(outEvents, &OutChannelEvent{
			SenderID:      event.SenderID,
			EventType:     event.EventType,
			Payload:       event.Payload,
			ChannelID:     event.ChannelID,
			Timestamp:     event.Timestamp,
		})
	}

	response := GetChannelEventsResponse{Events: outEvents}

	return &response, nil
}

func (srv *SyncServer) GetLastMessagesSince(ctx context.Context, req *GetLastMessagesSinceRequest) (*GetChannelEventsResponse, error) {
	// Check if is admin, and validate token
	identity, isOK := auth.VerifyToken(req.Token)

	// If not valid return
	if !isOK || !identity.CanUseAppID(req.AppID) {
		return nil, &AuthError{reason: "Invalid token"}
	}

	// Check if channel exists
	exists, err := core.GetEngine().GetChannelRepository().ExistsAppChannel(req.AppID, req.ChannelID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GRPC Get messages between timestamps: failed to check app channel existence %v\n", err)
		return nil, &InternalError{reason: err.Error()}
	}

	if !exists {
		return nil, &NotFoundError{}
	}

	events, err := core.GetEngine().GetChannelRepository().GetChannelLastEventsAfter(req.AppID, req.ChannelID, req.Amount, req.FromTimeStamp)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GRPC Get last messages since timestamp: failed fetch events %v\n", err)
		return nil, &InternalError{reason: err.Error()}
	}

	outEvents := make([]*OutChannelEvent, len(events))

	for _, event := range events {
		outEvents = append(outEvents, &OutChannelEvent{
			SenderID:      event.SenderID,
			EventType:     event.EventType,
			Payload:       event.Payload,
			ChannelID:     event.ChannelID,
			Timestamp:     event.Timestamp,
		})
	}

	response := GetChannelEventsResponse{Events: outEvents}

	return &response, nil
}

func (srv *SyncServer) GetLastMessages(ctx context.Context, req *GetLastMessagesRequest) (*GetChannelEventsResponse, error) {
	// Check if is admin, and validate token
	identity, isOK := auth.VerifyToken(req.Token)

	// If not valid return
	if !isOK || !identity.CanUseAppID(req.AppID) {
		return nil, &AuthError{reason: "Invalid token"}
	}

	// Check if channel exists
	exists, err := core.GetEngine().GetChannelRepository().ExistsAppChannel(req.AppID, req.ChannelID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GRPC Get messages between timestamps: failed to check app channel existence %v\n", err)
		return nil, &InternalError{reason: err.Error()}
	}

	if !exists {
		return nil, &NotFoundError{}
	}

	var events []*core.ChannelEvent

	if req.Amount <= core.CacheQueueSize {
		size := core.GetEngine().GetCacheStorage().GetChannelEventsSize(req.ChannelID, req.AppID)

		if size >= uint64(req.Amount) {
			events = core.GetEngine().GetCacheStorage().GetChannelEvents(req.ChannelID, req.AppID, req.Amount)
		} else {
			events, err = core.GetEngine().GetChannelRepository().GetChannelLastEvents(req.AppID, req.ChannelID, req.Amount)

			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "GRPC Get last messages: failed fetch events %v\n", err)
				return nil, &InternalError{reason: err.Error()}
			}
		}

	} else {
		// Get events
		events, err = core.GetEngine().GetChannelRepository().GetChannelLastEvents(req.AppID, req.ChannelID, req.Amount)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "HTTP Get last messages: failed fetch events %v\n", err)
			return nil, &InternalError{reason: err.Error()}
		}
	}

	outEvents := make([]*OutChannelEvent, len(events))

	for _, event := range events {
		outEvents = append(outEvents, &OutChannelEvent{
			SenderID:      event.SenderID,
			EventType:     event.EventType,
			Payload:       event.Payload,
			ChannelID:     event.ChannelID,
			Timestamp:     event.Timestamp,
		})
	}

	response := GetChannelEventsResponse{Events: outEvents}

	return &response, nil
}

func (srv *SyncServer) mustEmbedUnimplementedSyncServiceServer() {

}