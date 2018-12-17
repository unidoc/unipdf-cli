/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	"io/ioutil"

	unicommon "github.com/unidoc/unidoc/common"
	unilicense "github.com/unidoc/unidoc/common/license"
)

func SetLicense(licensePath string, customer string) error {
	// Read license file
	content, err := ioutil.ReadFile(licensePath)
	if err != nil {
		return err
	}

	return unilicense.SetLicenseKey(string(content), "")
}

func SetLogLevel(level unicommon.LogLevel) {
	unicommon.SetLogger(unicommon.NewConsoleLogger(level))
}
