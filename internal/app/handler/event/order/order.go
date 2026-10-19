package eorder

import (
	"context"
	"fmt"

	"github.com/KDarenskii/worker-service/internal/app/entity"
	ehandler "github.com/KDarenskii/worker-service/internal/app/handler/event"
	"github.com/KDarenskii/worker-service/internal/pkg/broker"
	"github.com/rs/zerolog/log"
)

type handler struct{}

func NewHandler() ehandler.OrderCreated {
	return &handler{}
}

func (h *handler) CallbackOrderCreated(
	ctx context.Context, ev *entity.EventOrderCreated, _ []broker.Header,
) error {
	log.Info().Ctx(ctx).EmbedObject(ev).Msg("broker: received message")

	if ev.OrderGUID == "" {
		return broker.NotCriticalError(fmt.Errorf("broker: invalid order_guid: %s", ev.OrderGUID))
	}

	return nil
}
