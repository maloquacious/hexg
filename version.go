// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg

import "github.com/maloquacious/semver"

var (
	version = semver.Version{
		Major: 0,
		Minor: 14,
		Patch: 3,
		Build: semver.Commit(),
	}
)

func Version() semver.Version {
	return version
}
