//go:build !embed

package frontend

import (
	"embed"
)

var embed_static embed.FS

func init() {
	static_embedded = false
}
