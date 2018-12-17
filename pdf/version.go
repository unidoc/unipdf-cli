/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	unicommon "github.com/unidoc/unidoc/common"
	unilicense "github.com/unidoc/unidoc/common/license"
)

const appVersion = "0.1"

type VersionInfo struct {
	App     string
	Lib     string
	License string
}

func Version() VersionInfo {
	var license string
	if key := unilicense.GetLicenseKey(); key != nil {
		license = key.ToString()
	}

	return VersionInfo{
		App:     appVersion,
		Lib:     unicommon.Version,
		License: license,
	}
}
