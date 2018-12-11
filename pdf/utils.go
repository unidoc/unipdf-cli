/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package pdf

import (
	unicommon "github.com/unidoc/unidoc/common"
	unicore "github.com/unidoc/unidoc/pdf/core"
)

func GetDict(obj unicore.PdfObject) *unicore.PdfObjectDictionary {
	if obj == nil {
		return nil
	}

	obj = unicore.TraceToDirectObject(obj)
	dict, ok := obj.(*unicore.PdfObjectDictionary)
	if !ok {
		unicommon.Log.Debug("Error type check error (got %T)", obj)
		return nil
	}

	return dict
}
