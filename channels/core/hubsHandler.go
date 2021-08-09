package core

import "sync"

func NewHubsHandler(hook HubsHandlerHook) *HubsHandler {
	return &HubsHandler{
		hubs: sync.Map{},
		hook: hook,
	}
}

// HubsHandler - Handle the hubs per application
type HubsHandler struct {
	hubs sync.Map // map[string]*Hub
	hook HubsHandlerHook
}

// ContainsHub - Return hub if exists
func (handler *HubsHandler) ContainsHub(AppID string) *Hub {
	data, isOK := handler.hubs.Load(AppID)

	if isOK {
		return data.(*Hub)
	}

	return nil
}

// GetHub - Get the Hub with the given AppID, if not found creates one
func (handler *HubsHandler) GetHub(AppID string) *Hub {
	data, isOK := handler.hubs.Load(AppID)

	if isOK {
		return data.(*Hub)
	}

	return handler.NewHub(AppID)
}

// NewHub - Create a new hub in this server add to the map
func (handler *HubsHandler) NewHub(AppID string) *Hub {
	hub := NewHub(AppID, nil)
	handler.hubs.Store(AppID, hub)

	if handler.hook != nil {
		hub.hook = handler.hook.OnNewHub(hub)
	}

	return hub
}

// RemoveHub - Remove hub from active hub and close all channels and connections
func (handler *HubsHandler) RemoveHub(AppID string) {
	data, isOK := handler.hubs.LoadAndDelete(AppID)

	if !isOK {
		return
	}

	hub := data.(*Hub)
	hub.Close()

	if handler.hook != nil {
		handler.hook.OnRemoveHub(hub)
	}
}
