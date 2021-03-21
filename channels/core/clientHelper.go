package core

import (
	"fmt"
	"os"
)

// GetClient - Get client from cache first, then try database and update cache if found
func GetClient(appID string, clientID string) (*Client, error) {

	client := GetEngine().GetCacheStorage().GetClient(appID, clientID)

	if client != nil {
		return client, nil
	}

	client, err := GetEngine().GetClientRepository().GetAppClient(appID, clientID)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Get client failed %v\n", err)
		return nil, err
	}

	if client == nil {
		return nil, nil
	}

	GetEngine().GetCacheStorage().StoreClient(appID, clientID, client)

	return client, nil
}

// CreateClient - Create client and also cache it
func CreateClient(appID string, clientID string, username string, extra string) (bool, error) {

	err := GetEngine().GetClientRepository().CreateClient(clientID, username, appID, extra)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Create Client failed %v\n", err)
		return false, err
	}

	GetEngine().GetCacheStorage().StoreClient(appID, clientID, &Client{
		ID:       clientID,
		Username: username,
		Extra:    extra,
		AppID:    appID,
	})

	return true, nil
}

// DeleteClient - Remove client from database and cache
func DeleteClient(appID string, clientID string) (bool, error) {
	exists, err := GetEngine().GetClientRepository().ExistsAppClient(appID, clientID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Delete Client: failed to check client existence %v\n", err)
		return false, err
	}

	if exists {
		return false, nil
	}

	err = GetEngine().GetClientRepository().DeleteClient(clientID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Delete Client failed %v\n", err)
		return false, err
	}

	// Remove client cache
	GetEngine().GetCacheStorage().RemoveClient(appID, clientID)
	GetEngine().GetCacheStorage().RemoveClientChannels(clientID)

	return true, nil
}

// UpdateClient - Update client on the database and cache
func UpdateClient(appID string, clientID string, username string, extra string) (bool, error) {
	exists, err := GetEngine().GetClientRepository().ExistsAppClient(appID, clientID)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Update Client: failed to check client existence %v\n", err)
		return false, err
	}

	if !exists {
		return false, nil
	}

	err = GetEngine().GetClientRepository().UpdateClient(clientID, username, extra)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Update Client failed %v\n", err)
		return false, err
	}

	GetEngine().GetCacheStorage().StoreClient(appID, clientID, &Client{
		ID:       clientID,
		AppID:    appID,
		Username: username,
		Extra:    extra,
	})

	return true, nil
}
