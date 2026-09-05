// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg

import "github.com/maloquacious/semver"

var (
	version = semver.Version{
		Major: 1,
		Minor: 1,
		Patch: 0,
		Build: semver.Commit(),
	}
)

// Version returns the current version of the hexg package.
func Version() semver.Version {
	return version
}
