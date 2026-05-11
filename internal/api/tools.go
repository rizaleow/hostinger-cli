//go:build tools

// This file pins build-time tool dependencies so they survive `go mod tidy`.
// It is never compiled into the binary (guarded by the `tools` build tag).
package api

import _ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
