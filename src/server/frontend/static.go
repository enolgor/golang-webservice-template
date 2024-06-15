package frontend

import (
	"errors"
	"io/fs"
	"log/slog"
	"os"

	"github.com/enolgor/golang-webservice-template/config"
)

var static_embedded bool = false

func InitStatic(log *slog.Logger) error {
	var err error
	if static_embedded {
		log.Debug("serving static content from embedded filesystem")
		Static, err = fs.Sub(embed_static, "static")
	} else {
		log.Debug("serving static content from disk", "dir", config.SERVER_STATIC_DIR)
		if config.SERVER_STATIC_DIR == "" {
			return errors.New("server static dir env is not configured")
		}
		Static = os.DirFS(config.SERVER_STATIC_DIR)
	}
	return err
}
