package application

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/enolgor/golang-webservice-template/authenticator"
	"github.com/enolgor/golang-webservice-template/config"
	"github.com/enolgor/golang-webservice-template/db"
	"github.com/enolgor/golang-webservice-template/models"
)

type App struct {
	Log           *slog.Logger
	LogLevel      *slog.LevelVar
	Authenticator *authenticator.Authenticator
	DB            *db.DB
}

func New(log *slog.Logger, logLevel *slog.LevelVar) (*App, error) {
	var err error
	app := &App{}
	app.Log = log
	app.LogLevel = logLevel
	if app.Authenticator, err = authenticator.New(); err != nil {
		return nil, err
	}
	if app.DB, err = db.NewSQLiteDB(config.SQLITE_DB); err != nil {
		return nil, err
	}
	return app, nil
}

func (app *App) HelloWorld() string {
	return "hello world"
}

func (app *App) Shutdown(ctx context.Context) error {
	return nil
}

func (app *App) GetUserFromProfile(ctx context.Context, profile *models.Profile) (*models.User, error) {
	user, err := app.DB.ReadUser(ctx, profile.UserID)
	if errors.Is(err, sql.ErrNoRows) {
		roles := []models.Role{models.USER}
		if profile.Email == config.DEFAULT_ADMIN_EMAIL {
			roles = models.AllRoles
		}
		user = &models.User{
			ID:    profile.UserID,
			Email: profile.Email,
			Roles: roles,
		}
		err = app.DB.CreateUser(ctx, user)
	}
	return user, err
}
