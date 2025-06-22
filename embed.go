package embed

//embeds migrations

import "embed"

//go:embed migrations/*.sql
var EmbedMigrations embed.FS
