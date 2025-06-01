package app

import (
	"fmt"
	"runtime/debug"
)

// Build information variables (set by goreleaser)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// versionInfo holds version information from various sources
type versionInfo struct {
	version string
	commit  string
	date    string
}

// getVersionFromBuildInfo extracts version information from debug.ReadBuildInfo()
func getVersionFromBuildInfo() *versionInfo {
	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return nil
	}

	info := &versionInfo{
		version: buildInfo.Main.Version,
		commit:  "none",
		date:    "unknown",
	}

	// Normalize version
	if info.version == "(devel)" {
		info.version = "dev"
	}

	// Extract VCS information
	for _, setting := range buildInfo.Settings {
		switch setting.Key {
		case "vcs.revision":
			if len(setting.Value) >= 7 {
				info.commit = setting.Value[:7] // short commit hash
			} else {
				info.commit = setting.Value
			}
		case "vcs.time":
			info.date = setting.Value
		}
	}

	return info
}

// getLDFlagsVersionInfo returns version info set via ldflags (for go build)
func getLDFlagsVersionInfo() *versionInfo {
	if version != "" && version != "dev" {
		return &versionInfo{
			version: version,
			commit:  commit,
			date:    date,
		}
	}
	return nil
}

// getFallbackVersionInfo returns fallback version information
func getFallbackVersionInfo() *versionInfo {
	return &versionInfo{
		version: version,
		commit:  commit,
		date:    date,
	}
}

// PrintVersion prints version information to stdout
func PrintVersion() {
	var info *versionInfo

	// Try version sources in priority order:
	// 1. LDFLAGS (go build with custom flags)
	// 2. Build info (go install)
	// 3. Fallback (default values)
	if info = getLDFlagsVersionInfo(); info != nil {
		// go build with ldflags
	} else if info = getVersionFromBuildInfo(); info != nil {
		// go install or go build without ldflags
	} else {
		// fallback
		info = getFallbackVersionInfo()
	}

	fmt.Printf("todotui %s\n", info.version)
	fmt.Printf("Commit: %s\n", info.commit)
	fmt.Printf("Built: %s\n", info.date)
}

// GetVersion returns the current version string for logging purposes
func GetVersion() string {
	if info := getLDFlagsVersionInfo(); info != nil {
		return info.version
	}
	if info := getVersionFromBuildInfo(); info != nil {
		return info.version
	}
	return version
}

// GetCommit returns the current commit hash for logging purposes
func GetCommit() string {
	if info := getLDFlagsVersionInfo(); info != nil {
		return info.commit
	}
	if info := getVersionFromBuildInfo(); info != nil {
		return info.commit
	}
	return commit
}
