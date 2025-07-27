// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package hexg

import "github.com/maloquacious/semver"

func Version() semver.Version {
	return semver.Version{
		Minor: 11,
		Patch: 0,
	}
}
