// Package migrations embeds the .sql files into the Go binary, so the
// deployed API can run migrations without shipping the folder separately.
package migrations

import "embed"

//go:embed *.sql
var FS embed.FS
