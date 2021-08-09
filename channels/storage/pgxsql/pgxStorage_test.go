package pgxsql

import (
	"testing"
	"time"
)

const (
	testDBUsername = "ChannelsTest"
	testDBPassword = "test_password"
	testDBName     = "ChannelsTest"
	testDBPort     = "5432"
	testDBHost     = "127.0.0.1"
)

func TestPGXStorage(t *testing.T) {
	PGXSetConnectionParams(testDBUsername, testDBPassword, testDBHost, testDBPort, testDBName)
	storage := NewSQLStorageDatabase()

	if storage.GetAppRepository() == nil {
		t.Errorf("Failed to get app repository \n")
	}

	if storage.GetChannelRepository() == nil {
		t.Errorf("Failed to get channel repository \n")
	}

	if storage.GetClientRepository() == nil {
		t.Errorf("Failed to get client repository \n")
	}

	if storage.GetDeviceRepository() == nil {
		t.Errorf("Failed to get device repository \n")
	}
}

func TestPGXChannelStorage(t *testing.T) {
	PGXSetConnectionParams(testDBUsername, testDBPassword, testDBHost, testDBPort, testDBName)
	storage := NewSQLStorageDatabase()
	appRepo := storage.GetAppRepository()

	appID := "123"
	appName := "test_app"

	clientID := "123"
	username := "test_user"
	extra := "test_extra"

	channelID := "123"
	channelName := "test_channel"
	createdAt := time.Now().Unix()

	// Create app
	if err := appRepo.CreateApp(appID, appName); err != nil {
		t.Errorf("Failed to create app %s \n", err.Error())
		return
	}

	clientRepo := storage.GetClientRepository()

	// Create client
	if err := clientRepo.CreateClient(clientID, username, appID, extra); err != nil {
		t.Errorf("Failed to create client %s \n", err.Error())
	}

	// START HERE

	repo := storage.GetChannelRepository()

	// Create channel
	if err := repo.CreateChannel(channelID, appID, channelName, createdAt, false, extra, false, false, false, false); err != nil {
		t.Errorf("Failed to create channel %s \n", err.Error())
	}

	// App channel existence
	if exists, err := repo.ExistsAppChannel(appID, channelID); err != nil {
		t.Errorf("Failed to check channel existence %s \n", err.Error())
	} else if !exists {
		t.Errorf("Failed to check channel existence, returned false after one being created \n")
	}

	// Get channel
	if channel, err := repo.GetAppChannel(appID, channelID); err != nil {
		t.Errorf("Failed to get channel %s \n", err.Error())
	} else if channel == nil {
		t.Errorf("Failed to get channel, existent channel returned nil \n")
	}

	// Delete channel
	if err := repo.DeleteChannel(appID, channelID); err != nil {
		t.Errorf("Failed to delete channel %s \n", err.Error())
	}

	// Check channel existence
	if exists, err := repo.ExistsAppChannel(appID, channelID); err != nil {
		t.Errorf("Failed to check channel existence %s \n", err.Error())
	} else if exists {
		t.Errorf("Failed to check channel existence, return true for a deleted channel \n")
	}

	// END HERE

	// Delete client
	if err := clientRepo.DeleteClient(clientID); err != nil {
		t.Errorf("Failed to delete client after it being creted %s \n", err.Error())
	}

	// Delete app
	if err := appRepo.DeleteApp(appID); err != nil {
		t.Errorf("Failed to delete app %s \n", err.Error())
	}
}

func TestPGXClientStorage(t *testing.T) {
	PGXSetConnectionParams(testDBUsername, testDBPassword, testDBHost, testDBPort, testDBName)
	storage := NewSQLStorageDatabase()
	appRepo := storage.GetAppRepository()

	appID := "123"
	appName := "test_app"

	clientID := "123"
	username := "test_user"
	secondUsername := "test_user2"
	extra := "test_extra"
	secondExtra := "test_extra2"

	// Create app
	if err := appRepo.CreateApp(appID, appName); err != nil {
		t.Errorf("Failed to create app %s \n", err.Error())
		return
	}

	repo := storage.GetClientRepository()

	// Create client
	if err := repo.CreateClient(clientID, username, appID, extra); err != nil {
		t.Errorf("Failed to create client %s \n", err.Error())
	}

	// Update client
	if err := repo.UpdateClient(clientID, secondUsername, secondExtra); err != nil {
		t.Errorf("Failed to update client %s \n", err.Error())
	}

	// Client exists
	if exists, err := repo.ExistsAppClient(appID, clientID); err != nil {
		t.Errorf("Failed to verify client existence %s \n", err.Error())
	} else if !exists {
		t.Errorf("Failed to verify client existence after it being created \n")
	}

	// Get all clients
	if clients, err := repo.GetAllClients(); err != nil {
		t.Errorf("Failed to get all clients %s \n", err.Error())
	} else if len(clients) == 0 {
		t.Errorf("Failed to get all clients, the length after one being created \n")
	}

	// Get all app clients
	if clients, err := repo.GetAppClients(appID); err != nil {
		t.Errorf("Failed to get app clients %s \n", err.Error())
	} else if len(clients) == 0 {
		t.Errorf("Failed to get app clients, length is 0 after one being created \n")
	}

	// Get app client
	if client, err := repo.GetAppClient(appID, clientID); err != nil {
		t.Errorf("Failed to get client after it being creted %s \n", err.Error())
	} else if client == nil {
		t.Errorf("Failed to get client after it being creted \n")
	} else {
		if client.Extra != secondExtra {
			t.Errorf("Failed to update client \n")
		}

		if client.Username != secondUsername {
			t.Errorf("Failed to update client \n")
		}
	}

	// Delete client
	if err := repo.DeleteClient(clientID); err != nil {
		t.Errorf("Failed to delete client after it being creted %s \n", err.Error())
	}

	// Get client after being deleted
	if client, err := repo.GetAppClient(appID, clientID); err != nil {
		t.Errorf("Failed to get client after it being creted %s \n", err.Error())
	} else if client != nil {
		t.Errorf("Failed to delete client \n")
	}

	// Delete app
	if err := appRepo.DeleteApp(appID); err != nil {
		t.Errorf("Failed to delete app %s \n", err.Error())
	}
}

func TestPGXAppStorage(t *testing.T) {
	PGXSetConnectionParams(testDBUsername, testDBPassword, testDBHost, testDBPort, testDBName)
	storage := NewSQLStorageDatabase()
	repo := storage.GetAppRepository()

	appID := "123"
	appName := "test_app"
	secondAppName := "second_test_app"

	// Create app
	if err := repo.CreateApp(appID, appName); err != nil {
		t.Errorf("Failed to create app %s \n", err.Error())
		return
	}

	// Update app
	if err := repo.UpdateApp(appID, secondAppName); err != nil {
		t.Errorf("Failed to update app %s \n", err.Error())
	}

	// Get all apps
	if apps, err := repo.GetApps(); err != nil {
		t.Errorf("Failed to get all apps %s \n", err.Error())
	} else if len(apps) == 0 {
		t.Errorf("Failed to get all apps, the length is 0 after one being created \n")
	} else if apps[0].Name != secondAppName {
		t.Errorf("App name update failed \n")
	}

	// Delete app
	if err := repo.DeleteApp(appID); err != nil {
		t.Errorf("Failed to delete app %s \n", err.Error())
	}

	// Get all apps
	if apps, err := repo.GetApps(); err != nil {
		t.Errorf("Failed to get all apps %s \n", err.Error())
	} else if len(apps) != 0 {
		t.Errorf("Failed to delete app, the length is not 0 after one being created and delete \n")
	}
}
