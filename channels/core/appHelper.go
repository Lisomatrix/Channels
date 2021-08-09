package core

import (
	"errors"
	"log"
)

func GetApplications() ([]*App, error) {
	apps, err := GetEngine().GetAppRepository().GetApps()

	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Since we fetched the apps, let use the opportunity to cache them
	for _, app := range apps {
		GetEngine().GetCacheStorage().StoreApp(app.AppID, app.Name)
	}

	return apps, nil
}

func GetApplication(appID string) (*App, error) {
	// Check if app exists in cache
	app := GetEngine().GetCacheStorage().GetApp(appID)

	if app != nil {
		return app, nil
	}

	app, err := GetEngine().GetAppRepository().GetApp(appID)

	if err != nil {
		return nil, err
	}

	if app == nil {
		return nil, nil
	}

	return app, nil
}

func CreateApplication(appID, name string) error {

	app, err := GetApplication(appID)

	if err != nil {
		log.Println(err)
		return err
	}

	if app != nil {
		return errors.New("app with given ID already exists")
	}

	if err := GetEngine().GetAppRepository().CreateApp(appID, name); err != nil {
		log.Println(err)
		return err
	}

	// Store in cache
	GetEngine().GetCacheStorage().StoreApp(appID, name)

	return nil
}

func DeleteApplication(appID string) error {

	if err := GetEngine().GetAppRepository().DeleteApp(appID); err != nil {
		log.Println(err)
		return err
	}

	GetEngine().GetCacheStorage().RemoveApp(appID)

	return nil
}

func UpdateApplication(appID, name string) error {
	app, err := GetApplication(appID)

	if err != nil {
		log.Println(err)
		return err
	}

	if app == nil {
		return errors.New("app with given ID not found")
	}

	if err := GetEngine().GetAppRepository().UpdateApp(appID, name); err != nil {
		log.Println(err)
		return err
	}

	// Update cache
	GetEngine().GetCacheStorage().StoreApp(appID, name)

	return nil
}
