package monitor

import (
	"context"
)

type Repository interface {
	//Endpoint operation
	CreateEndpoint(ctx context.Context, endpoint *Endpoint) error
	GetEndpointByID(ctx context.Context, id int64) (*Endpoint, error)
	ListEndpoints(ctx context.Context) ([]*Endpoint, error)
	UpdateEndpoint(ctx context.Context, endpoint *Endpoint) error
	DeleteEndpoint(ctx context.Context, id int64) error

	// Monitoring results/history
	// SaveCheckResult(ctx context.Context, result *CheckResult) error
	// GetCheckHistory(ctx context.Context, endpointID int64, limit int) ([]*CheckResult, error)
}
