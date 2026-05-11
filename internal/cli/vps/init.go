package vps

import (
	"github.com/rizaleow/hostinger-cli/internal/api"
	"github.com/rizaleow/hostinger-cli/internal/output"
)

// Register table renderers for the most common collection responses. Anything
// not registered here falls back to JSON automatically.
func init() {
	output.RegisterCollectionTable(api.VPSV1DataCenterDataCenterCollection{}, dataCentersTableCols)
}
