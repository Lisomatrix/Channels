package handlers

import (
	"context"
	"fmt"
	"github.com/lisomatrix/channels/channels/auth"
	"github.com/lisomatrix/channels/channels/core"
	"os"
	"time"
)

type ChannelServer struct {}

func (srv *ChannelServer) mustEmbedUnimplementedChannelServiceServer() {

}

func (srv *ChannelServer) PublishEvent(ctx context.Context, req *ChannelPublishRequest) (*DefaultResponse, error) {

	// Check if is admin, and validate token
	identity, isOK := auth.AuthenticateAdmin(req.Token)

	// If not valid return
	if !isOK || !identity.CanUseAppID(req.AppID) {
		return nil, &AuthError{reason: "Invalid token"}
	}

	channel := core.GetEngine().GetCacheStorage().GetChannel(req.AppID, req.ChannelID)

	if channel == nil {
		var err error
		channel, err = core.GetEngine().GetChannelRepository().GetAppChannel(req.AppID, req.ChannelID)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "GRPC Publish Channel: failed to get app channel %v\n", err)
			return nil, &InternalError{reason: "Failed to get app channel"}
		}

		if channel == nil /*!exists*/ {
			return nil, &AuthError{reason: "Invalid token"}
		}

		core.GetEngine().GetCacheStorage().StoreChannel(req.AppID, req.ChannelID, channel)
	}

	if channel.IsClosed {
		return &DefaultResponse{IsOK: false}, nil
	}

	event := &core.ChannelEvent{
		SenderID:  identity.ClientID,
		EventType: req.EventType,
		Payload:   string(req.Payload),
		ChannelID: req.ChannelID,
		Timestamp: time.Now().Unix(),
	}

	if channel.Persistent {
		core.GetEngine().GetCacheStorage().StoreChannelEvent(req.ChannelID, req.AppID, event)
		core.GetEngine().StoreEvent(req.AppID, event)
	}

	// If no hub exists then we don't have clients from this hub
	hub := core.GetEngine().HubsHandler.ContainsHub(req.AppID)

	if hub == nil {
		return &DefaultResponse{IsOK: true}, nil
	}

	// If inside the hub there are no clients listening for channel, then ignore
	hubChannel := hub.ContainsChannel(req.ChannelID)

	if hubChannel == nil {
		return &DefaultResponse{IsOK: true}, nil
	}

	// If there client listening then send to them
	hubChannel.ExternalPublish(event)

	return &DefaultResponse{IsOK: true}, nil
}

func (srv *ChannelServer) CreateChannel(ctx context.Context, req *CreateChannelRequest) (*DefaultResponse, error) {
	// Check if is admin, and validate token
	identity, isOK := auth.AuthenticateAdmin(req.Token)

	// If not valid return
	if !isOK || !identity.CanUseAppID(req.AppID) {
		return nil, &AuthError{reason: "Invalid token"}
	}

	createdAt := time.Now().Unix()

	newChannel := core.Channel{
		ID:         req.ChannelID,
		AppID:      req.AppID,
		Name:       req.Name,
		CreatedAt:  createdAt,
		IsClosed:   false,
		Extra:      req.Extra,
		Persistent: req.Persistent,
		Private:    req.Private,
		Presence:   req.Presence,
	}

	if isOK, err := core.CreateChannel(req.AppID, &newChannel); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GRPC Create Channel failed %v\n", err)
		return nil, &InternalError{reason: "GRPC Create Channel failed"}
	} else if isOK {
		return &DefaultResponse{IsOK: true}, nil
	} else {
		return &DefaultResponse{IsOK: false}, nil
	}
}

func (srv *ChannelServer) JoinChannel(ctx context.Context, req *JoinClientRequest) (*DefaultResponse, error) {
	// Check if is admin, and validate token
	identity, isOK := auth.AuthenticateAdmin(req.Token)

	// If not valid return
	if !isOK || !identity.CanUseAppID(req.AppID) {
		return nil, &AuthError{reason: "Invalid token"}
	}

	if isOK, err := core.JoinChannel(req.AppID, req.ChannelID, req.ClientID); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GRPC Join channel: failed to join client %v\n", err)
		return nil, &InternalError{reason: "GRPC Join Channel failed"}
	} else if isOK {
		return &DefaultResponse{IsOK: true}, nil
	} else {
		return &DefaultResponse{IsOK: false}, nil
	}
}

func (srv *ChannelServer) LeaveChannel(ctx context.Context, req *LeaveClientRequest) (*DefaultResponse, error) {
	// Check if is admin, and validate token
	identity, isOK := auth.AuthenticateAdmin(req.Token)

	// If not valid return
	if !isOK || !identity.CanUseAppID(req.AppID) {
		return nil, &AuthError{reason: "Invalid token"}
	}

	if isOK, err := core.LeaveChannel(req.AppID, req.ChannelID, req.ClientID); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GRPC Leave channel: remove client from channel %v\n", err)
		return nil, &InternalError{reason: "GRPC Leave Channel failed"}
	} else if isOK {
		return &DefaultResponse{IsOK: true}, nil
	} else {
		return &DefaultResponse{IsOK: false}, nil
	}
}

func (srv *ChannelServer) DeleteChannel(ctx context.Context, req *DeleteChannelRequest) (*DefaultResponse, error) {
	// Check if is admin, and validate token
	identity, isOK := auth.AuthenticateAdmin(req.Token)

	// If not valid return
	if !isOK || !identity.CanUseAppID(req.AppID) {
		return nil, &AuthError{reason: "Invalid token"}
	}

	if isOK, err := core.DeleteChannel(req.AppID, req.ChannelID); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GRPC Delete channel: failed to delete channel %v\n", err)
		return nil, &InternalError{reason: "GRPC delete Channel failed"}
	} else if isOK {
		return &DefaultResponse{IsOK: true}, nil
	} else {
		return &DefaultResponse{IsOK: false}, nil
	}
}

func (srv *ChannelServer) CloseChannel(ctx context.Context, req *CloseChannelRequest) (*DefaultResponse, error) {
	// Check if is admin, and validate token
	identity, isOK := auth.AuthenticateAdmin(req.Token)

	// If not valid return
	if !isOK || !identity.CanUseAppID(req.AppID) {
		return nil, &AuthError{reason: "Invalid token"}
	}

	if isOK, err := core.SetChannelCloseStatus(req.AppID, req.ChannelID, true); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GRPC Close channel: failed to close channel %v\n", err)
		return nil, &InternalError{reason: "GRPC close Channel failed"}
	} else if isOK {
		return &DefaultResponse{IsOK: true}, nil
	} else {
		return &DefaultResponse{IsOK: false}, nil
	}
}

func (srv *ChannelServer) OpenChannel(ctx context.Context, req *OpenChannelRequest) (*DefaultResponse, error) {
	// Check if is admin, and validate token
	identity, isOK := auth.AuthenticateAdmin(req.Token)

	// If not valid return
	if !isOK || !identity.CanUseAppID(req.AppID) {
		return nil, &AuthError{reason: "Invalid token"}
	}

	if isOK, err := core.SetChannelCloseStatus(req.AppID, req.ChannelID, false); err != nil {
		fmt.Fprintf(os.Stderr, "GRPC Open channel: failed to open channel %v\n", err)
		return nil, &InternalError{reason: "GRPC Open Channel failed"}
	} else if isOK {
		return &DefaultResponse{IsOK: true}, nil
	} else {
		return &DefaultResponse{IsOK: false}, nil
	}
}

func (srv *ChannelServer) GetChannels(ctx context.Context, req *DefaultRequest) (*GetChannelsResponse, error) {
	// Check if is admin, and validate token
	identity, isOK := auth.AuthenticateAdmin(req.Token)

	// If not valid return
	if !isOK || !identity.CanUseAppID(identity.AppID) {
		return nil, &AuthError{reason: "Invalid token"}
	}

	var response GetChannelsResponse

	// If AppID given, and user is authorized to use AppID
	// Then fetch App specific channels
	if identity.AppID != "" && identity.CanUseAppID(identity.AppID) {

		if channels, err := core.GetEngine().GetChannelRepository().GetAppPublicChannels(identity.AppID); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "GRPC Get channels channels: failed to get app open channels %v\n", err)
			return nil, &InternalError{reason: "GRPC Get channels channels failed"}
		} else {

			outChannels := make([]*Channel, len(channels))
			for _, c := range channels {
				outChannels = append(outChannels, &Channel{
					AppID: c.AppID,
					Name: c.Name,
					Extra: c.Extra,
					CreatedAt: c.CreatedAt,
					IsClosed: c.IsClosed,
					Persistent: c.Persistent,
					Presence: c.Presence,
					Private: c.Private,
				})
			}

			response.Channels = outChannels
		}

		// Else if no given appID and user is SuperAdmin
		// Fetch all open channels from all apps
	} else if identity.AppID == "" && identity.IsSuperAdmin() {

		if channels, err := core.GetEngine().GetChannelRepository().GetAllPublicChannels(); err != nil {
			fmt.Fprintf(os.Stderr, "HTTP Get open channels: failed to get all open channels %v\n", err)
			return nil, &InternalError{reason: "GRPC Get channels channels failed"}
		} else {

			outChannels := make([]*Channel, len(channels))
			for _, c := range channels {
				outChannels = append(outChannels, &Channel{
					AppID: c.AppID,
					Name: c.Name,
					Extra: c.Extra,
					CreatedAt: c.CreatedAt,
					IsClosed: c.IsClosed,
					Persistent: c.Persistent,
					Presence: c.Presence,
					Private: c.Private,
				})
			}

			response.Channels = outChannels
		}
		// Otherwise there are no permissions for this request
	} else {
		return nil, &AuthError{reason: "Unauthorized"}
	}

	return &response, nil
}