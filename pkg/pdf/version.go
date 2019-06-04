/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	unicommon "github.com/unidoc/unipdf/v3/common"
	unilicense "github.com/unidoc/unipdf/v3/common/license"
)

// VersionInfo contains version and license information
// about the Unidoc library.
type VersionInfo struct {
	Lib     string
	License string
}

// Version returns version and license information about the Unidoc library.
func Version() VersionInfo {
	var license string
	if key := unilicense.GetLicenseKey(); key != nil {
		license = key.ToString()
	}

	return VersionInfo{
		Lib:     unicommon.Version,
		License: license,
	}
}
