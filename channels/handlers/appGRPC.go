package handlers

import (
	"context"
	"fmt"
	"github.com/lisomatrix/channels/channels/auth"
	"github.com/lisomatrix/channels/channels/core"
	"os"
)

type AppServer struct{}

func (srv *AppServer) CreateApp(ctx context.Context, req *CreateAppRequest) (*DefaultResponse, error) {

	if req.GetToken() == "" {
		return &DefaultResponse{IsOK: false}, nil
	}

	// Check if is admin, and validate token
	identity, isOK := auth.AuthenticateAdmin(req.Token)

	// If not valid return
	if !isOK || !identity.IsSuperAdmin() {
		return &DefaultResponse{IsOK: false}, nil
	}

	// Check if app already exists in cache
	app := core.GetEngine().GetCacheStorage().GetApp(req.App.AppID)

	// If not found check from database
	if app == nil {
		exists, err := core.GetEngine().GetAppRepository().AppExists(req.App.AppID)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "GRPC Create App: failed to check app existence %v\n", err)
			return &DefaultResponse{IsOK: false}, err
		}

		if exists {
			return &DefaultResponse{IsOK: false}, nil
		}
	}

	// Store app in the database
	err := core.GetEngine().GetAppRepository().CreateApp(req.App.AppID, req.App.Name)

	// Store app in cache
	core.GetEngine().GetCacheStorage().StoreApp(req.App.AppID, req.App.Name)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GRPC Create App: failed to create app %v\n", err)
		return &DefaultResponse{IsOK: false}, err
	}

	return &DefaultResponse{IsOK: true}, nil
}

func (srv *AppServer) UpdateApp(ctx context.Context, req *UpdateAppRequest) (*DefaultResponse, error) {
	if req.GetToken() == "" {
		return nil, &AuthError{reason: "Invalid token"}
	}

	// Check if is admin, and validate token
	identity, isOK := auth.AuthenticateAdmin(req.Token)

	// If not valid return
	if !isOK || !identity.IsSuperAdmin() {
		return nil, &AuthError{reason: "Invalid token"}
	}

	// If client is just admin then check it is their app
	if !identity.CanUseAppID(req.Id) {
		return nil, &AuthError{reason: "Invalid token"}
	}

	// Check if app already exists in cache
	app := core.GetEngine().GetCacheStorage().GetApp(req.Id)

	// If app not found in cache
	if app == nil {
		// Check its exitence in the database
		exists, err := core.GetEngine().GetAppRepository().AppExists(req.Id)

		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "GRPC Update App: failed to check app existence %v\n", err)
			return nil, err
		}

		// If not found then the request is invalid
		if !exists {
			return &DefaultResponse{IsOK: false}, nil
		}
	}

	// Update app
	err := core.GetEngine().GetAppRepository().UpdateApp(req.Id, req.Name)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GRPC Update App: failed to update app %v\n", err)
		return nil, err
	}

	// Update cache entry
	core.GetEngine().GetCacheStorage().StoreApp(req.Id, req.Name)


	return &DefaultResponse{IsOK: true}, nil
}

func (srv *AppServer) GetApps(ctx context.Context, req *DefaultRequest) (*GetAppsResponse, error) {
	if req.GetToken() == "" {
		return nil, &AuthError{reason: "Invalid token"}
	}

	// Check if is admin, and validate token
	identity, isOK := auth.AuthenticateAdmin(req.Token)

	// If not valid return
	if !isOK || !identity.IsSuperAdmin() {
		return nil, &AuthError{reason: "Invalid token"}
	}

	// Get alls apps
	apps, err := core.GetEngine().GetAppRepository().GetApps()

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GRPC Get Apps: failed to get apps %v\n", err)
		return nil, &InternalError{reason: "Failed to get apps"}
	}

	// Since we fetched the apps, let use the opportunity to cache them
	go func() {
		for _, app := range apps {
			core.GetEngine().GetCacheStorage().StoreApp(app.AppID, app.Name)
		}
	}()

	outApps := make([]*App, len(apps))

	for _, p := range apps {
		outApps = append(outApps, &App{
			AppID:         p.AppID,
			Name:          p.Name,
		}	)
	}

	return &GetAppsResponse{Apps: outApps}, nil
}

func (srv *AppServer) DeleteApp(ctx context.Context, req *DeleteAppRequest) (*DefaultResponse, error) {
	if req.GetToken() == "" {
		return nil, &AuthError{reason: "Invalid token"}
	}

	// Check if is admin, and validate token
	identity, isOK := auth.AuthenticateAdmin(req.Token)

	// If not valid return
	if !isOK || !identity.IsSuperAdmin() {
		return nil, &AuthError{reason: "Invalid token"}
	}

	err := core.GetEngine().GetAppRepository().DeleteApp(req.Id)

	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "GRPC Delete App: failed to delete app %v\n", err)
		return nil, &InternalError{reason: "Failed to delete app"}
	}

	core.GetEngine().GetCacheStorage().RemoveApp(req.Id)

	return &DefaultResponse{IsOK: true}, nil
}

func (srv *AppServer) mustEmbedUnimplementedAppServiceServer() {

}
