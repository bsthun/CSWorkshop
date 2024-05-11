package text

import (
	"log"
	"runtime/debug"
)

var Commit = func() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		var hash string
		var modified string
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				hash = setting.Value[:7]
			}
			if setting.Key == "vcs.modified" {
				if setting.Value == "false" {
					modified = "/c" // Clean build
				} else {
					modified = "/d" // Dirty build
				}
			}
		}
		if hash == "" || modified == "" {
			log.Fatal("Failed to get build info", nil)
		}
		return hash + modified + "/" + Build
	}
	return ""
}()
