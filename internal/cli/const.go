/*
 * This file is subject to the terms and conditions defined in
 * file 'LICENSE.md', which is part of this source code package.
 */

package cli

import (
	"errors"
	"strings"

	unicommon "github.com/unidoc/unipdf/v3/common"
	unisecurity "github.com/unidoc/unipdf/v3/core/security"
	unipdf "github.com/unidoc/unipdf/v3/model"
)

var encryptAlgoMap = map[string]unipdf.EncryptionAlgorithm{
	"rc4":    unipdf.RC4_128bit,
	"aes128": unipdf.AES_128bit,
	"aes256": unipdf.AES_256bit,
}

var logLevelMap = map[string]unicommon.LogLevel{
	"trace":   unicommon.LogLevelTrace,
	"debug":   unicommon.LogLevelDebug,
	"info":    unicommon.LogLevelInfo,
	"notice":  unicommon.LogLevelNotice,
	"warning": unicommon.LogLevelWarning,
	"error":   unicommon.LogLevelError,
}

var imageFormats = map[string]struct{}{
	"jpeg": struct{}{},
	"png":  struct{}{},
}

func parseEncryptionMode(mode string) (unipdf.EncryptionAlgorithm, error) {
	algo, ok := encryptAlgoMap[mode]
	if !ok {
		return 0, errors.New("invalid encryption mode")
	}

	return algo, nil
}

func parseLogLevel(levelStr string) (unicommon.LogLevel, error) {
	levelStr = strings.TrimSpace(levelStr)
	if levelStr == "" {
		return unicommon.LogLevelError, nil
	}

	level, ok := logLevelMap[levelStr]
	if !ok {
		return 0, errors.New("invalid log level")
	}

	return level, nil
}

func parsePermissionList(permStr string) (unisecurity.Permissions, error) {
	permStr = removeSpaces(permStr)
	if permStr == "" {
		return 0, nil
	}
	permList := strings.Split(permStr, ",")

	perms := unisecurity.Permissions(0)
	for _, perm := range permList {
		if perm == "" {
			continue
		}

		switch perm {
		case "all":
			perms = unisecurity.PermOwner
		case "none":
			perms = unisecurity.Permissions(0)
		case "print-low-res":
			perms |= unisecurity.PermPrinting
		case "print-high-res":
			perms |= unisecurity.PermFullPrintQuality
		case "modify":
			perms |= unisecurity.PermModify
		case "extract-graphics":
			perms |= unisecurity.PermExtractGraphics
		case "annotate":
			perms |= unisecurity.PermAnnotate
		case "fill-forms":
			perms |= unisecurity.PermFillForms
		case "rotate":
			perms |= unisecurity.PermRotateInsert

		default:
			return 0, errors.New("invalid permission")
		}
	}

	return perms, nil
}
