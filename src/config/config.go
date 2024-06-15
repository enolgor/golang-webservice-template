package config

import (
	"log/slog"

	"github.com/enolgor/golang-webservice-template/utils/env"
)

const (
	DEVELOPMENT string = "development"
	PRODUCTION  string = "production"
)

var Mode string
var WaitForDebugger bool
var SQLITE_DB string
var DEFAULT_ADMIN_EMAIL string
var SESSION_SECRET string
var AUTH0_DOMAIN string
var AUTH0_CLIENT_ID string
var AUTH0_CLIENT_SECRET string
var AUTH0_CALLBACK_URL string
var SERVER_STATIC_DIR string

func Load(log *slog.Logger, level *slog.LevelVar) {
	env.ParseEnv(log, level,
		env.MustGetValid("MODE", &Mode, env.String, DEVELOPMENT, PRODUCTION),
		env.Get("WAIT_FOR_DEBUGGER", &WaitForDebugger, env.Bool, false),
		env.MustGet("DEFAULT_ADMIN_EMAIL", &DEFAULT_ADMIN_EMAIL, env.String),
		env.MustGet("SQLITE_DB", &SQLITE_DB, env.String),
		env.MustGet("AUTH0_DOMAIN", &AUTH0_DOMAIN, env.String),
		env.MustGet("AUTH0_CLIENT_ID", &AUTH0_CLIENT_ID, env.String),
		env.MustGet("AUTH0_CLIENT_SECRET", &AUTH0_CLIENT_SECRET, env.String),
		env.MustGet("AUTH0_CALLBACK_URL", &AUTH0_CALLBACK_URL, env.String),
		env.MustGet("SESSION_SECRET", &SESSION_SECRET, env.String),
		env.Get("SERVER_STATIC_DIR", &SERVER_STATIC_DIR, env.String, ""),
	)
}
