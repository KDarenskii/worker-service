package ehandler

import (
	"context"

	"github.com/KDarenskii/worker-service/internal/app/entity"
	"github.com/KDarenskii/worker-service/internal/pkg/broker"
)

type (
	OrderCreated interface {
		CallbackOrderCreated(
			ctx context.Context,
			ev *entity.EventOrderCreated, headers []broker.Header,
		) error
	}
)
