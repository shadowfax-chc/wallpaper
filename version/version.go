// Package version providers verstion information about the binary.
package version

import (
	"bytes"
	"fmt"
)

var (
	// BuildTime is supplied by the compiler as the time at which the binary was built
	BuildTime string
	// GitCommit is supplied by the compiler as the most recent commit the binary was built from
	GitCommit string
	// Version is supplied by the compiler as the most recent git tag the binary was built from
	// defaults to 0.0.1
	Version string
	// VersionDescription is a modifier to Version that describes the binary build
	VersionDescription = "dev"
)

// Description returns a string describing the binary build
// <version>(-<VersionDescription>) ( :: commit - <GitCommit> [ :: built @ <BuildTime> ] )
func Description() string {

	var versionString bytes.Buffer

	fmt.Fprintf(&versionString, "%s", Version)
	if VersionDescription != "" {
		fmt.Fprintf(&versionString, "-%s", VersionDescription)
	}

	if GitCommit != "" {
		fmt.Fprintf(&versionString, " :: commit - %s", GitCommit)
	}

	if BuildTime != "" {
		fmt.Fprintf(&versionString, " :: built @ %s", BuildTime)
	}

	return versionString.String()

}
