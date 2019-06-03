/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	"io/ioutil"

	unicommon "github.com/unidoc/unipdf/v3/common"
	unilicense "github.com/unidoc/unipdf/v3/common/license"
)

// SetLicense sets the license for using the Unidoc library.
func SetLicense(licensePath string, customer string) error {
	// Read license file
	content, err := ioutil.ReadFile(licensePath)
	if err != nil {
		return err
	}

	return unilicense.SetLicenseKey(string(content), "")
}

// SetLogLevel sets the verbosity of the output produced by the Unidoc library.
func SetLogLevel(level unicommon.LogLevel) {
	unicommon.SetLogger(unicommon.NewConsoleLogger(level))
}
