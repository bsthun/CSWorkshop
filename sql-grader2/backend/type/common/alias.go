package common

import (
	"io/fs"
)

type MigrationFS fs.FS

type FrontendFS fs.ReadFileFS
