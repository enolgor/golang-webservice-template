//go:build embed

package frontend

import (
	"embed"
)

//go:embed static
var embed_static embed.FS

func init() {
	static_embedded = true
}
