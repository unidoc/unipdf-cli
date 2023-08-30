/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	"fmt"
	"io/ioutil"

	unicommon "github.com/unidoc/unipdf/v3/common"
	unilicense "github.com/unidoc/unipdf/v3/common/license"
)

// SetLicense sets the license for using the UniDOC library.
func SetLicense(licensePath string, customer string) error {
	// Read license file
	content, err := ioutil.ReadFile(licensePath)
	if err != nil {
		return err
	}

	return unilicense.SetLicenseKey(string(content), customer)
}

// SetMeteredKey sets the license key for using the UniDoc library with metered api key.
func SetMeteredKey(apiKey string) error {
	return unilicense.SetMeteredKey(apiKey)
}

// GetLicenseKey get information about user license key.
func GetLicenseKey() string {
	lk := unilicense.GetLicenseKey()
	if lk == nil {
		return "Failed retrieving license key"
	}
	return lk.ToString()
}

// GetMeteredState freshly checks the state, contacting the licensing server.
func GetMeteredState() {
	// GetMeteredState freshly checks the state, contacting the licensing server.
	state, err := unilicense.GetMeteredState()
	if err != nil {
		fmt.Printf("ERROR getting metered state: %+v\n", err)
		return
	}
	fmt.Printf("State: %+v\n", state)
	if state.OK {
		fmt.Printf("State is OK\n")
	} else {
		fmt.Printf("State is not OK\n")
	}
	fmt.Printf("Credits: %v\n", state.Credits)
	fmt.Printf("Used credits: %v\n", state.Used)
}

// SetLogLevel sets the verbosity of the output produced by the UniDOC library.
func SetLogLevel(level unicommon.LogLevel) {
	unicommon.SetLogger(unicommon.NewConsoleLogger(level))
}
