/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	unicommon "github.com/unidoc/unidoc/common"
)

const appVersion = "0.1"

type VersionInfo struct {
	App string
	Lib string
}

func Version() VersionInfo {
	return VersionInfo{
		App: appVersion,
		Lib: unicommon.Version,
	}
}
